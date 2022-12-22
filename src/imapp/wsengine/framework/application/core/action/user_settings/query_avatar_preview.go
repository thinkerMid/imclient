package userSettings

import (
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/action/common"
	"ws/framework/application/core/processor"
)

// QueryAvatarPreview .
type QueryAvatarPreview struct {
	processor.BaseAction
}

// Start .
func (m *QueryAvatarPreview) Start(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	account := context.ResolveAccountService().Context()

	// 有头像
	if account.HaveAvatar {
		next()
		return
	}

	m.Query = &common.QueryAvatarPreview{UserID: context.ResolveJID().User}
	return m.Query.Start(context, func() {})
}

// Receive .
func (m *QueryAvatarPreview) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	defer next()

	return m.Query.Receive(context, func() {})
}

// Error .
func (m *QueryAvatarPreview) Error(context containerInterface.IMessageContext, err error) {
	if m.Query != nil {
		m.Query.Error(context, err)
	}
}
