package application

import (
	"fmt"
	"sync/atomic"
	"ws/framework/application/constant/binary"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	panicDump "ws/framework/plugin/panic_dump"
	"ws/framework/plugin/webhook"
)

// AddProcessorAndAttach .
//
//	异步添加，可以在并发环境下调用，进来的是非异步事件的协程
//	给非内部环境操作的 通常用于操作场景
//	备注：这个里面的异常处理没对processor进行OnDestroy调用 因为没真正存在过于执行环境中
func (c *App) AddProcessorAndAttach(processor containerInterface.IProcessor, resultProcessor containerInterface.LocalResultProcessor) (pid uint32) {
	processor.Init(c.container.ResolveLogger().Named("Processor"))

	if c.shutdownErrorRemark.hasError() {
		context := NewProcessContext(c.container, nil, false)
		processor.OnChannelError(&context, c.shutdownErrorRemark.thrown())

		results := context.MessageResult()
		for i := range results {
			resultProcessor.ProcessResult(&context, &results[i])
		}

		resultProcessor.OnDestroy()
		return
	}

	pid = atomic.AddUint32(&c.messageProcessorId, 1)
	processor.SetID(pid)
	processor.SetResultProcessor(resultProcessor)

	c.workerPool.Invoke(newWakeUpProcessorEvent(processor), false)
	return
}

// AddMessageProcessor .
//
//	只有内部场景才会有添加消息处理操作， 进来的是异步事件的协程
//	备注：这个里面的异常处理没对processor进行OnDestroy调用 因为没真正存在过于执行环境中
func (c *App) AddMessageProcessor(processor containerInterface.IProcessor) (pid uint32) {
	processor.Init(c.container.ResolveLogger().Named("Processor"))

	if c.shutdownErrorRemark.hasError() {
		context := NewProcessContext(c.container, nil, false)
		processor.OnChannelError(&context, c.shutdownErrorRemark.thrown())
		return
	}

	pid = atomic.AddUint32(&c.messageProcessorId, 1)
	processor.SetID(pid)

	c.workerPool.Invoke(newWakeUpProcessorEvent(processor), true)
	return
}

// AddFutureProcessor .
//
//	同步添加，只能在action执行环境下使用，进来的是异步事件的协程
//	备注：这个里面的异常处理没对processor进行OnDestroy调用 因为没真正存在过于执行环境中
func (c *App) AddFutureProcessor(pid uint32, processor containerInterface.IProcessor) {
	processor.Init(c.container.ResolveLogger().Named("Processor"))

	// 关闭了
	if c.shutdownErrorRemark.hasError() {
		context := NewProcessContext(c.container, nil, false)
		processor.OnChannelError(&context, c.shutdownErrorRemark.thrown())
		return
	}

	processor.SetID(pid)

	c.messageProcessorQueue.PushBack(processor)

	return
}

// AddGlobalResultProcessor .
//
//	只允许在初始化成功后调用添加
func (c *App) AddGlobalResultProcessor(resultProcessor containerInterface.GlobalResultProcessor) {
	if atomic.LoadUint32(&c.status) == activate {
		c.logger.Warn("client already activated. can't add result processor when activated after. please add start before")
		return
	}

	c.globalResultProcessorQueue = append(c.globalResultProcessorQueue, resultProcessor)
}

// RemoveProcessor .
// 外部场景调用，进来的是非异步事件的协程
func (c *App) RemoveProcessor(pid uint32) {
	if c.shutdownErrorRemark.hasError() {
		return
	}

	c.workerPool.Invoke(newRemoveProcessorEvent(pid), false)
}

// ----------------------------------------------------------------------------

func (c *App) invokeMessage(node *waBinary.Node, isSignalMessage bool) {
	pCtx := NewProcessContext(c.container, node, isSignalMessage)
	nextElement := c.messageProcessorQueue.Front()
	var p containerInterface.IProcessor

	defer func() {
		if err := recover(); err != nil {
			stackScan := panicDump.Scan(3)
			stackScan.Print(c.container.ResolveLogger())

			processorInfo := p.DumpInfo()

			webhook.PanicPush(webhook.PanicPushTemplate{
				JID:                 c.container.ResolveJID().User,
				Message:             node.XMLString(),
				ProcessorID:         processorInfo.PID,
				ProcessorAlisaName:  processorInfo.AliasName,
				ProcessorType:       processorInfo.ProcessorType,
				CurrentActionName:   processorInfo.CurrentActionName,
				CurrentActionStatus: processorInfo.CurrentActionStatus,
				ActionQueue:         processorInfo.ActionQueue,
				PanicError:          fmt.Sprintf("%v", err),
				StackDumpID:         stackScan.DumpID,
				Time:                stackScan.Time,
			})

			panic(err)
		}
	}()

	// 遍历消息处理
	for nextElement != nil && nextElement.Value != nil {
		currentElement := nextElement
		nextElement = nextElement.Next()

		p = currentElement.Value.(containerInterface.IProcessor)
		pid := p.ID()

		pCtx.Reset()
		pCtx.SetProcessPID(pid)

		p.ProcessMessage(&pCtx)

		// 重新加入消息处理队列
		if !pCtx.AddToQueue() {
			p.OnDestroy(&pCtx)

			c.messageProcessorQueue.Remove(currentElement)
			c.watchDog.UnWatch(p)
		}

		c.dispatchGlobalResultProcessor(&pCtx)

		// 非信号才可中断
		if !pCtx.SignalMessage() && pCtx.MessageAborted() {
			break
		}
	}
}

func (c *App) dispatchGlobalResultProcessor(context *ProcessContext) {
	results := context.MessageResult()

	var find bool

	for i := range results {
		for j := range c.globalResultProcessorQueue {
			find = false

			p := c.globalResultProcessorQueue[j]
			if p == nil {
				continue
			}

			listenTagList := p.ListenTags()

			for _, tag := range listenTagList {
				if tag == results[i].ResultType {
					find = true
					break
				}
			}

			if find {
				p.ProcessResult(context, &results[i])

				if !p.Reside() {
					p.OnDestroy()
					c.globalResultProcessorQueue[j] = nil
				}
			}
		}
	}
}

func (c *App) wakeUpProcessor(processor containerInterface.IProcessor) {
	defer func() {
		if err := recover(); err != nil {
			stackScan := panicDump.Scan(3)
			stackScan.Print(c.container.ResolveLogger())

			processorInfo := processor.DumpInfo()

			webhook.PanicPush(webhook.PanicPushTemplate{
				JID:                 c.container.ResolveJID().User,
				Message:             "/",
				ProcessorID:         processorInfo.PID,
				ProcessorAlisaName:  processorInfo.AliasName,
				ProcessorType:       processorInfo.ProcessorType,
				CurrentActionName:   processorInfo.CurrentActionName,
				CurrentActionStatus: processorInfo.CurrentActionStatus,
				ActionQueue:         processorInfo.ActionQueue,
				PanicError:          fmt.Sprintf("%v", err),
				StackDumpID:         stackScan.DumpID,
				Time:                stackScan.Time,
			})

			panic(err)
		}
	}()

	pid := processor.ID()

	pCtx := NewProcessContext(c.container, nil, false)
	pCtx.SetProcessPID(pid)

	if c.shutdownErrorRemark.hasError() {
		processor.OnChannelError(&pCtx, c.shutdownErrorRemark.thrown())
		c.dispatchGlobalResultProcessor(&pCtx)
		processor.OnDestroy(&pCtx)

		return
	}

	if processor.NeedAutoStart() {
		processor.Start(&pCtx)
	}

	if pCtx.AddToQueue() {
		c.messageProcessorQueue.PushBack(processor)
		// 加入监控
		c.watchDog.Watch(processor)
	} else {
		processor.OnDestroy(&pCtx)
	}

	c.dispatchGlobalResultProcessor(&pCtx)
}

func (c *App) removeProcessor(pid uint32) {
	pCtx := NewProcessContext(c.container, nil, false)
	nextElement := c.messageProcessorQueue.Front()

	for nextElement != nil && nextElement.Value != nil {
		currentElement := nextElement
		nextElement = nextElement.Next()

		p := currentElement.Value.(containerInterface.IProcessor)
		if pid == p.ID() {
			p.OnDestroy(&pCtx)

			c.messageProcessorQueue.Remove(currentElement)

			// 移除监控
			c.watchDog.UnWatch(p)
			break
		}
	}
}

// 清理后台级别的处理器
//
//	为了解决某些情况下导致的包发送了但是连接断开没有响应的问题
//	这些载体的处理器一般是后台级别的
//	只清理 messageProcessorQueue
//	至于 globalResultProcessorQueue 都不处理。真正shutdown会清理
func (c *App) cleanupBackgroundProcessor() {
	pCtx := NewProcessContext(c.container, nil, false)

	nextElement := c.messageProcessorQueue.Front()

	for nextElement != nil && nextElement.Value != nil {
		currentElement := nextElement
		nextElement = nextElement.Next()

		p := currentElement.Value.(containerInterface.IProcessor)
		pid := p.ID()

		pCtx.Reset()
		pCtx.SetProcessPID(pid)

		// 移除非前台优先级的处理器
		if p.Priority() != processor.PriorityForeground {
			p.OnChannelError(&pCtx, c.shutdownErrorRemark.thrown())
			c.dispatchGlobalResultProcessor(&pCtx)
			p.OnDestroy(&pCtx)

			c.messageProcessorQueue.Remove(currentElement)
			// 移除监控
			c.watchDog.UnWatch(p)
		}
	}
}

// 生命周期结束了 不会再执行了
func (c *App) shutdown() {
	pCtx := NewProcessContext(c.container, nil, false)

	nextElement := c.messageProcessorQueue.Front()

	for nextElement != nil && nextElement.Value != nil {
		p := nextElement.Value.(containerInterface.IProcessor)

		pCtx.Reset()
		pCtx.SetProcessPID(p.ID())

		p.OnChannelError(&pCtx, c.shutdownErrorRemark.thrown())
		c.dispatchGlobalResultProcessor(&pCtx)
		p.OnDestroy(&pCtx)

		nextElement = nextElement.Next()
	}

	for i := range c.globalResultProcessorQueue {
		p := c.globalResultProcessorQueue[i]
		if p == nil {
			continue
		}

		c.globalResultProcessorQueue[i].OnDestroy()
	}

	c.messageProcessorQueue.Init()
	c.globalResultProcessorQueue = make([]containerInterface.GlobalResultProcessor, 0)
	c.watchDog.Stop()
}
