package monitor

import (
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/action/common"
	"ws/framework/application/core/action/contact"
	privateChat "ws/framework/application/core/action/private_chat"
	"ws/framework/application/core/action/private_chat/common"
)

// ReplaceUserIDWhenSession 用于从添加到发起私信流程场景，更正手机号码和whatsappJID不一致的情况
type ReplaceUserIDWhenSession struct {
	DstUserID string
}

// OnStart .
func (p *ReplaceUserIDWhenSession) OnStart(ioc containerInterface.IAppIocContainer) {}

// OnActionStartBefore .
func (p *ReplaceUserIDWhenSession) OnActionStartBefore(iAction interface{}, _ containerInterface.IMessageContext) {
	switch iAction.(type) {
	case *common.QueryAvatarUrl:
		action := iAction.(*common.QueryAvatarUrl)
		action.UserID = p.DstUserID
	case *common.QueryDevicesIdentity:
		action := iAction.(*common.QueryDevicesIdentity)
		action.UserID = p.DstUserID
	case *common.QueryUserDeviceList:
		action := iAction.(*common.QueryUserDeviceList)
		action.UserID = p.DstUserID
	case *common.SubscribeStatus:
		action := iAction.(*common.SubscribeStatus)
		action.UserID = p.DstUserID
	case *common.QueryMultiDevicesIdentity:
		action := iAction.(*common.QueryMultiDevicesIdentity)
		action.UserID = p.DstUserID
	case *privateChatCommon.TrustedContactToken:
		action := iAction.(*privateChatCommon.TrustedContactToken)
		action.UserID = p.DstUserID
	case *privateChat.SendImage:
		action := iAction.(*privateChat.SendImage)
		action.UserID = p.DstUserID
	case *privateChat.SendText:
		action := iAction.(*privateChat.SendText)
		action.UserID = p.DstUserID
	case *privateChat.SendAudio:
		action := iAction.(*privateChat.SendAudio)
		action.UserID = p.DstUserID
	case *privateChat.SendTemp:
		action := iAction.(*privateChat.SendTemp)
		action.UserID = p.DstUserID
	case *privateChat.SendVideo:
		action := iAction.(*privateChat.SendVideo)
		action.UserID = p.DstUserID
	}
}

// OnActionStartAfter .
func (p *ReplaceUserIDWhenSession) OnActionStartAfter(iAction interface{}, _ containerInterface.IMessageContext) {
	switch iAction.(type) {
	case *contact.Check:
		action := iAction.(*contact.Check)

		// 查询出来的ID是不一样的 用查询出来的
		if p.DstUserID != action.UserID {
			p.DstUserID = action.UserID
		}
	}
}

func (p *ReplaceUserIDWhenSession) OnActionStartFail(action interface{}, context containerInterface.IMessageContext) {

}

func (p *ReplaceUserIDWhenSession) OnActionStartSuccess(action interface{}, context containerInterface.IMessageContext) {

}

// ActionExecuteSuccess .
func (p *ReplaceUserIDWhenSession) ActionExecuteSuccess(iAction interface{}, _ containerInterface.IMessageContext) {
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
func (p *ReplaceUserIDWhenSession) ActionExecuteFailure(action interface{}, context containerInterface.IMessageContext) {
}

// OnExit .
func (p *ReplaceUserIDWhenSession) OnExit(ioc containerInterface.IAppIocContainer) {}
