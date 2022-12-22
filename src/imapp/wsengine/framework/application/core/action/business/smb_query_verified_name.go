package business

import (
	waBinary "ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	containerInterface "ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
)

// SMBQueryVerifiedName .
type SMBQueryVerifiedName struct {
	processor.BaseAction
}

// Start .
func (m *SMBQueryVerifiedName) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	m.SendMessageId, err = context.SendIQ(message.InfoQuery{
		ID:        context.GenerateRequestID(),
		Namespace: "w:biz",
		Type:      message.IqGet,
		To:        types.ServerJID,
		Content: []waBinary.Node{{
			Tag: "verified_name",
			Attrs: waBinary.Attrs{
				"jid": context.ResolveJID(),
			},
		}},
	})

	return
}

// Receive .
func (m *SMBQueryVerifiedName) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	next()

	return nil
}

// Error .
func (m *SMBQueryVerifiedName) Error(context containerInterface.IMessageContext, err error) {
}
