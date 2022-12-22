package monitor

import (
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/action/common"
	. "ws/framework/application/core/wam"
)

// Session .
type Session struct {
}

// OnStart .
func (p *Session) OnStart(ioc containerInterface.IAppIocContainer) {
	LogManager().SwitchAppMenu(ioc, PageSession)
}

// OnActionStartBefore .
func (p *Session) OnActionStartBefore(_ interface{}, _ containerInterface.IMessageContext) {}

// OnActionStartAfter .
func (p *Session) OnActionStartAfter(_ interface{}, _ containerInterface.IMessageContext) {}

func (p *Session) OnActionStartFail(action interface{}, context containerInterface.IMessageContext) {

}

func (p *Session) OnActionStartSuccess(action interface{}, context containerInterface.IMessageContext) {

}

// ActionExecuteSuccess .
func (p *Session) ActionExecuteSuccess(action interface{}, context containerInterface.IMessageContext) {
	switch action.(type) {
	case *common.QueryMultiDevicesIdentity:
		// 创建会话
		LogManager().LogSessionNew(context)
	case *common.QueryAvatarPreview:
		// 联系人头像变化 记录日志
		query := action.(*common.QueryAvatarPreview)

		contact := context.ResolveContactService().FindByJID(query.UserID)
		if contact != nil {
			LogManager().LogNotifyContactAvatar(context, int32(query.AvatarSize))
		}
	}
}

func (p *Session) ActionExecuteFailure(action interface{}, context containerInterface.IMessageContext) {
	switch action.(type) {
	case *common.QueryAvatarPreview:
		// 联系人头像变化 记录日志
		query := action.(*common.QueryAvatarPreview)

		contact := context.ResolveContactService().FindByJID(query.UserID)
		if contact != nil {
			LogManager().LogNotifyContactAvatar(context, 0)
		}
	}
}

// OnExit .
func (p *Session) OnExit(containerInterface.IAppIocContainer) {

}
