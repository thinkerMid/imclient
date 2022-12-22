package userSettings

import (
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	"ws/framework/application/core/result/constant"
)

// QuerySignature .
type QuerySignature struct {
	processor.BaseAction
}

// Start .
func (m *QuerySignature) Start(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	account := context.ResolveAccountService().Context()

	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.GetSignature,
		Content:    account.Signature,
	})

	next()

	return
}

// Receive .
func (m *QuerySignature) Receive(_ containerInterface.IMessageContext, _ containerInterface.NextActionFn) error {
	return nil
}

// Error .
func (m *QuerySignature) Error(_ containerInterface.IMessageContext, _ error) {}
