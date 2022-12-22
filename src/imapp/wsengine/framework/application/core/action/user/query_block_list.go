package user

import (
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
)

// QueryBlockList .
type QueryBlockList struct {
	processor.BaseAction
}

// Start .
func (m *QueryBlockList) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	m.SendMessageId, err = context.SendIQ(message.InfoQuery{
		ID:        context.GenerateRequestID(),
		Namespace: "blocklist",
		Type:      "get",
		To:        types.ServerJID,
	})

	return
}

// Receive .
func (m *QueryBlockList) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	next()

	return nil
}

// Error .
func (m *QueryBlockList) Error(context containerInterface.IMessageContext, err error) {
}
