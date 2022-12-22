package processor

import (
	"errors"
	"ws/framework/application/constant"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor/queue"
	functionTools "ws/framework/utils/function_tools"
)

// NewNotificationProcessor .
func NewNotificationProcessor(fn containerInterface.MakeNotificationFn, sets ...SetConfigFn) containerInterface.IProcessor {
	p := notificationProcessor{}

	for _, set := range sets {
		set(&p)
	}

	p.notificationQueue = actionQueue.NewNotification(fn)

	return &p
}

// notificationProcessor 专门用于统一处理通知消息的处理器
type notificationProcessor struct {
	BaseProcessor
	notificationQueue containerInterface.INotificationQueue
}

// ProcessMessage .
func (m *notificationProcessor) ProcessMessage(context containerInterface.IProcessControl) {
	if m.Trigger == nil {
		context.Reject()
		return
	}

	msgCtx := context.MessageContext()

	if !m.Trigger.ContentedCondition(msgCtx.Message()) {
		return
	}

	for _, monitor := range m.Monitors {
		monitor.OnStart(msgCtx)
	}

	// 遍历处理 直到栈无内容为止
	for {
		action := m.notificationQueue.Pop()
		if action == nil {
			break
		}

		err := action.Receive(msgCtx)
		if err != nil {
			// 中断
			if errors.Is(err, constant.AbortedError) {
				break
			}
		}

		// 不管有没异常 接收到了就当执行成功的
		for _, monitor := range m.Monitors {
			monitor.ActionExecuteSuccess(action, msgCtx)
		}
	}

	m.Trigger.Reset()
	m.notificationQueue.Reload()
}

// OnDestroy .
func (m *notificationProcessor) OnDestroy(ctx containerInterface.IProcessControl) {
	for _, monitor := range m.Monitors {
		monitor.OnExit(ctx.MessageContext())
	}
}

// ProcessorType .
func (m *notificationProcessor) ProcessorType() containerInterface.ProcessorType {
	return containerInterface.NotificationProcessorType
}

// DumpInfo .
func (m *notificationProcessor) DumpInfo() containerInterface.ProcessorDump {
	dump := containerInterface.ProcessorDump{
		PID:           m.pid,
		AliasName:     m.aliasName,
		ProcessorType: m.ProcessorType().String(),
	}

	action := m.notificationQueue.Current()
	if action != nil {
		dump.CurrentActionName = functionTools.ReflectValueTypeName(action)
		dump.CurrentActionStatus = "/"
	} else {
		dump.CurrentActionName = "/"
		dump.CurrentActionStatus = "/"
	}

	dump.ActionQueue = m.notificationQueue.Dump()

	return dump
}
