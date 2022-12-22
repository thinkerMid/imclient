package monitor

import (
	containerInterface "ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/action/common"
	contactNotification "ws/framework/application/core/notification/contact"
	"ws/framework/application/core/processor"
	. "ws/framework/application/core/wam"
)

// NotificationMonitor .
type NotificationMonitor struct {
}

// OnStart .
func (p *NotificationMonitor) OnStart(ioc containerInterface.IAppIocContainer) {

}

// OnActionStartBefore .
func (p *NotificationMonitor) OnActionStartBefore(_ interface{}, _ containerInterface.IMessageContext) {
}

// OnActionStartAfter .
func (p *NotificationMonitor) OnActionStartAfter(_ interface{}, _ containerInterface.IMessageContext) {
}

func (p *NotificationMonitor) OnActionStartFail(action interface{}, context containerInterface.IMessageContext) {

}

func (p *NotificationMonitor) OnActionStartSuccess(action interface{}, context containerInterface.IMessageContext) {

}

// ActionExecuteSuccess .
func (p *NotificationMonitor) ActionExecuteSuccess(action interface{}, context containerInterface.IMessageContext) {
	switch action.(type) {
	case *contactNotification.BeAddedOrDeleted:
		query := action.(*contactNotification.BeAddedOrDeleted)
		if len(query.JID) == 0 {
			break
		}

		// 被添加为联系人 或者 被删除联系人 【真机也不是每次都推通知】
		LogManager().LogNotifyContactAddedOrDeleted(context)
	case *contactNotification.AvatarUpdate:
		// 接收联系人头像变化通知，查头像buffer
		query := action.(*contactNotification.AvatarUpdate)
		if len(query.JID) == 0 {
			break
		}

		context.AddMessageProcessor(processor.NewOnceProcessor(
			[]containerInterface.IAction{
				&common.QueryAvatarPreview{UserID: query.JID},
			},
			processor.AttachMonitor(&Session{}),
		))
	}
}

func (p *NotificationMonitor) ActionExecuteFailure(action interface{}, context containerInterface.IMessageContext) {

}

// OnExit .
func (p *NotificationMonitor) OnExit(containerInterface.IAppIocContainer) {

}
