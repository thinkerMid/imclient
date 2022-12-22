package monitor

import (
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/action/common"
	"ws/framework/application/core/action/contact"
)

// ReplaceUserIDWhenContact 用于从添加联系人流程场景，更正手机号码和whatsappJID不一致的情况
type ReplaceUserIDWhenContact struct {
	DstUserID string
}

// OnStart .
func (p *ReplaceUserIDWhenContact) OnStart(ioc containerInterface.IAppIocContainer) {}

// OnActionStartBefore .
func (p *ReplaceUserIDWhenContact) OnActionStartBefore(iAction interface{}, _ containerInterface.IMessageContext) {
	switch iAction.(type) {
	case *contact.Delete:
		action := iAction.(*contact.Delete)
		action.UserID = p.DstUserID
	case *common.DeleteDevices:
		action := iAction.(*common.DeleteDevices)
		action.UserID = p.DstUserID
	case *common.QueryAvatarUrl:
		action := iAction.(*common.QueryAvatarUrl)
		action.UserID = p.DstUserID
	case *common.QueryDevicesIdentity:
		action := iAction.(*common.QueryDevicesIdentity)
		action.UserID = p.DstUserID
	case *common.QueryUserDeviceList:
		action := iAction.(*common.QueryUserDeviceList)
		action.UserID = p.DstUserID
	}
}

// OnActionStartAfter .
func (p *ReplaceUserIDWhenContact) OnActionStartAfter(iAction interface{}, _ containerInterface.IMessageContext) {
	switch iAction.(type) {
	case *contact.Check:
		action := iAction.(*contact.Check)

		// 查询出来的ID是不一样的 用查询出来的
		if p.DstUserID != action.UserID {
			p.DstUserID = action.UserID
		}
	}
}

func (p *ReplaceUserIDWhenContact) OnActionStartFail(action interface{}, context containerInterface.IMessageContext) {

}

func (p *ReplaceUserIDWhenContact) OnActionStartSuccess(action interface{}, context containerInterface.IMessageContext) {

}

// ActionExecuteSuccess .
func (p *ReplaceUserIDWhenContact) ActionExecuteSuccess(iAction interface{}, _ containerInterface.IMessageContext) {
	switch iAction.(type) {
	case *contact.Add:
		action := iAction.(*contact.Add)

		// 添加成功后ID是不一样的 用添加完的出来的
		if p.DstUserID != action.UserID {
			p.DstUserID = action.UserID
		}
	}
}

// ActionExecuteFailure .
func (p *ReplaceUserIDWhenContact) ActionExecuteFailure(action interface{}, context containerInterface.IMessageContext) {
}

// OnExit .
func (p *ReplaceUserIDWhenContact) OnExit(ioc containerInterface.IAppIocContainer) {}
