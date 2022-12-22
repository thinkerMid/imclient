package control

import (
	"ws/framework/application/constant/message"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/action/event"
	eventNotification "ws/framework/application/core/notification/event"
	"ws/framework/application/core/processor"
)

// initTimerProcessor 以下定时只有实例退出了才会跟着移除 不需要再次注册
func initTimerProcessor(channel containerInterface.IMessageChannel) {
	// 日志写重置
	channel.AddMessageProcessor(processor.NewNotificationProcessor(
		func() []containerInterface.INotification {
			return []containerInterface.INotification{
				eventNotification.ResetEventLogState{},
			}
		},
		processor.TriggerTag(message.TickClockEvent),
		processor.Priority(processor.PriorityForeground),
	))

	// 上传0渠道日志 5分钟
	channel.AddMessageProcessor(processor.NewTimerProcessor(
		func() []containerInterface.IAction {
			return []containerInterface.IAction{&event.UploadChannel0Event{}}
		},
		processor.Interval(300),
		processor.IntervalLoop(true),
		processor.Priority(processor.PriorityForeground),
	))

	// 上传2渠道日志 10分钟
	channel.AddMessageProcessor(processor.NewTimerProcessor(
		func() []containerInterface.IAction {
			return []containerInterface.IAction{&event.UploadChannel2Event{}}
		},
		processor.Interval(600),
		processor.IntervalLoop(false),
		processor.Priority(processor.PriorityForeground),
	))
}
