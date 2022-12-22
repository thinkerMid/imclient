package user

import (
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
)

// DisappearingMode .
type DisappearingMode struct {
	processor.BaseAction
}

// Start .
func (m *DisappearingMode) Start(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	iq := message.InfoQuery{
		ID:        context.GenerateRequestID(),
		Namespace: "disappearing_mode",
		Type:      message.IqGet,
		To:        types.ServerJID,
	}

	m.SendMessageId, err = context.SendIQ(iq)

	return
}

// Receive .
func (m *DisappearingMode) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	next()

	return nil
}

// Error .
func (m *DisappearingMode) Error(_ containerInterface.IMessageContext, _ error) {}
