package user

import (
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
)

// CleanDirtyType .
type CleanDirtyType struct {
	processor.BaseAction
	Type string
}

// Start .
func (m *CleanDirtyType) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	m.SendMessageId, err = context.SendIQ(message.InfoQuery{
		ID:        context.GenerateRequestID(),
		Namespace: "urn:xmpp:whatsapp:dirty",
		Type:      "set",
		To:        types.ServerJID,
		Content: []waBinary.Node{{
			Tag: "clean",
			Attrs: waBinary.Attrs{
				"type": m.Type,
			},
		}},
	})

	return
}

// Receive .
func (m *CleanDirtyType) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	next()

	return nil
}

// Error .
func (m *CleanDirtyType) Error(context containerInterface.IMessageContext, err error) {
}
