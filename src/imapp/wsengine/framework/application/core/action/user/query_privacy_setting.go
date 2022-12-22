package user

import (
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
)

// QueryPrivacySetting .
type QueryPrivacySetting struct {
	processor.BaseAction
}

// Start .
func (m *QueryPrivacySetting) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	m.SendMessageId, err = context.SendIQ(message.InfoQuery{
		ID:        context.GenerateRequestID(),
		Namespace: "privacy",
		Type:      "get",
		To:        types.ServerJID,
		Content: []waBinary.Node{{
			Tag: "privacy",
		}},
	})

	return
}

// Receive .
func (m *QueryPrivacySetting) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	next()

	return nil
}

// Error .
func (m *QueryPrivacySetting) Error(context containerInterface.IMessageContext, err error) {
}
