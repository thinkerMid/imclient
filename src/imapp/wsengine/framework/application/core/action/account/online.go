package accountState

import (
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	"ws/framework/application/core/result/constant"
)

// Online .
type Online struct {
	processor.BaseAction
}

// Start .
func (m *Online) Start(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.Online,
	})

	next()

	return nil
}

// Receive .
func (m *Online) Receive(_ containerInterface.IMessageContext, _ containerInterface.NextActionFn) error {
	return nil
}

// Error .
func (m *Online) Error(_ containerInterface.IMessageContext, _ error) {}
