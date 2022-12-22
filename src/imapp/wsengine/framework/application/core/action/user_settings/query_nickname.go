package userSettings

import (
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	"ws/framework/application/core/result/constant"
)

// QueryNickname .
type QueryNickname struct {
	processor.BaseAction
}

// Start .
func (m *QueryNickname) Start(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.GetNickName,
		Content:    context.ResolveDeviceService().Context().PushName,
	})

	next()

	return
}

// Receive .
func (m *QueryNickname) Receive(_ containerInterface.IMessageContext, _ containerInterface.NextActionFn) error {
	return nil
}

// Error .
func (m *QueryNickname) Error(_ containerInterface.IMessageContext, _ error) {}
