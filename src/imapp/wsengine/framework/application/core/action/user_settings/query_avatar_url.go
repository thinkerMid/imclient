package userSettings

import (
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/action/common"
	"ws/framework/application/core/processor"
)

// QueryAvatarUrl .
type QueryAvatarUrl struct {
	processor.BaseAction
}

// RaiseErrorWhenNodeError 是否抛出node错误，如401，404等
func (m *QueryAvatarUrl) RaiseErrorWhenNodeError() bool {
	return false
}

// Start .
func (m *QueryAvatarUrl) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	m.Query = &common.QueryAvatarUrl{UserID: context.ResolveJID().User}
	return m.Query.Start(context, func() {})
}

// Receive .
func (m *QueryAvatarUrl) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	defer next()

	err := m.Query.Receive(context, func() {})
	if err != nil {
		return err
	}

	return nil
}

// Error .
func (m *QueryAvatarUrl) Error(context containerInterface.IMessageContext, err error) {
	if m.Query != nil {
		m.Query.Error(context, err)
	}
}
