package application

import (
	"sync/atomic"
	"time"
	"ws/framework/application/constant"
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/container/abstract_interface"
)

// 状态
const (
	standby  uint32 = iota // 待机
	activate               // 可用
	exit                   // 退出
	kill                   // 强制杀死
)

// ----------------------------------------------------------------------------

// 信号事件类型
const (
	disconnect uint8 = iota
	destroy
)

// 处理器事件子类型
const (
	wakeUp uint8 = iota
	remove
)

// ----------------------------------------------------------------------------

type event struct {
	EventType      eventType
	ChildEventType uint8
	Content        interface{}
}

type eventType uint8

const (
	// 信号事件类型
	signalEvent eventType = iota + 1
	// 处理器事件类型
	processorEvent
	// 消息事件类型
	messageEvent
	// 时钟事件类型
	tickClockEvent
)

func newMessageEvent(node *waBinary.Node) event {
	return event{EventType: messageEvent, Content: node}
}

func newSignalEvent(v uint8) event {
	return event{EventType: signalEvent, ChildEventType: v}
}

func newWakeUpProcessorEvent(p containerInterface.IProcessor) event {
	return event{EventType: processorEvent, ChildEventType: wakeUp, Content: p}
}

func newRemoveProcessorEvent(pid uint32) event {
	return event{EventType: processorEvent, ChildEventType: remove, Content: pid}
}

func newTickClockEvent() event {
	return event{EventType: tickClockEvent}
}

// ----------------------------------------------------------------------------

// 定时唤醒
// 内部场景调用，进来的是定时器的协程
func (c *App) timeTick() {
	c.workerPool.Invoke(newTickClockEvent(), false)

	// 下一次计数
	c.timeSchedule.Reset(time.Second)
}

// XMPP消息或连接事件的处理循环
func (c *App) eventLoop(i interface{}) {
	e := i.(event)

	switch e.EventType {
	case signalEvent:
		c.processSignalEvent(e.ChildEventType)
	case processorEvent:
		c.processProcessorEvent(e.ChildEventType, e.Content)
	case tickClockEvent:
		c.timerEvent()
	case messageEvent:
		node := e.Content.(*waBinary.Node)

		c.processMessage(node)
	}
}

// 定时器事件
func (c *App) timerEvent() {
	if atomic.LoadUint32(&c.status) != activate {
		return
	}

	tickNode := waBinary.Node{
		Tag:   message.TickClockEvent,
		Attrs: waBinary.Attrs{"id": message.TickClockEvent},
	}
	c.invokeMessage(&tickNode, true)

	err := c.watchDog.ScheduleCheck()
	if err != nil {
		c.logger.Errorf("KILLING Self: Blocked in action on processor. %v", err)
		c.changeStatus(kill)

		c.exit()
	}
}

// 处理器管理事件
func (c *App) processProcessorEvent(i uint8, content interface{}) {
	switch i {
	case wakeUp:
		c.wakeUpProcessor(content.(containerInterface.IProcessor))
	case remove:
		c.removeProcessor(content.(uint32))
	}
}

// 信号事件
func (c *App) processSignalEvent(i uint8) {
	switch i {
	// disconnect是连接关闭才有的 属于最底层的通知 尝试重连
	case disconnect:
		status := atomic.LoadUint32(&c.status)
		switch status {
		case exit:
			c.shutdownErrorRemark.streamErr = constant.LogoutError
		case kill:
			c.shutdownErrorRemark.streamErr = constant.ConnectionClosedError
		default:
			if c.shutdownErrorRemark.streamErr == nil {
				c.shutdownErrorRemark.streamErr = constant.ConnectionClosedError
			}
		}

		c.logger.Info("disconnect error: ", c.shutdownErrorRemark.thrown())

		if status != activate || c.shutdownErrorRemark.accountStateErr != nil {
			c.shutdown()
			return
		}

		err := c.reconnect()
		if err != nil {
			c.shutdown()
			return
		}

		c.resume()

		// 通知重连事件
		reconnectNode := waBinary.Node{
			Tag:   message.ReconnectEvent,
			Attrs: waBinary.Attrs{"id": message.ReconnectEvent},
		}

		c.invokeMessage(&reconnectNode, true)
	case destroy:
		c.exit()
	}
}

// 退出
func (c *App) exit() {
	logoutNode := waBinary.Node{
		Tag:   message.LogoutEvent,
		Attrs: waBinary.Attrs{"id": message.LogoutEvent},
	}

	c.invokeMessage(&logoutNode, true)

	c.container.ResolveConnection().Close()
}
