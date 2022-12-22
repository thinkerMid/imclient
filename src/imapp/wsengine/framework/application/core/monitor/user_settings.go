package monitor

import (
	"ws/framework/application/container/abstract_interface"
	userSettings "ws/framework/application/core/action/user_settings"
	. "ws/framework/application/core/wam"
)

// UserSettings .
type UserSettings struct{}

// OnStart .
func (p *UserSettings) OnStart(ioc containerInterface.IAppIocContainer) {
	LogManager().SwitchAppMenu(ioc, PageProfile)
}

// OnActionStartBefore .
func (p *UserSettings) OnActionStartBefore(action interface{}, context containerInterface.IMessageContext) {
	switch action.(type) {
	case *userSettings.QueryAvatarPreview:
		LogManager().LogProfileView(context)
	}
}

// OnActionStartAfter .
func (p *UserSettings) OnActionStartAfter(_ interface{}, _ containerInterface.IMessageContext) {}

func (p *UserSettings) OnActionStartFail(action interface{}, context containerInterface.IMessageContext) {

}

func (p *UserSettings) OnActionStartSuccess(action interface{}, context containerInterface.IMessageContext) {

}

// ActionExecuteSuccess .
func (p *UserSettings) ActionExecuteSuccess(action interface{}, context containerInterface.IMessageContext) {
	switch action.(type) {
	case *userSettings.ModifyAvatar:
		ma := action.(*userSettings.ModifyAvatar)
		LogManager().LogProfileSetAvatar(context, uint32(len(ma.Content)))
	case *userSettings.ModifyNickName:
	case *userSettings.ModifySignature:
	}
}

// ActionExecuteFailure .
func (p *UserSettings) ActionExecuteFailure(action interface{}, context containerInterface.IMessageContext) {
}

// OnExit .
func (p *UserSettings) OnExit(ioc containerInterface.IAppIocContainer) {}
