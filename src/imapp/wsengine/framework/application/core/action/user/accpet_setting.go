package user

import (
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
)

// QueryAcceptSetting .
type QueryAcceptSetting struct {
	processor.BaseAction
}

// Start .
func (m *QueryAcceptSetting) Start(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	iq := message.InfoQuery{
		ID:        context.GenerateRequestID(),
		Namespace: "urn:xmpp:whatsapp:account",
		Type:      message.IqGet,
		To:        types.ServerJID,
		Content: []waBinary.Node{{
			Tag: "accept",
		}},
	}

	m.SendMessageId, err = context.SendIQ(iq)

	return
}

// Receive .
func (m *QueryAcceptSetting) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	next()

	return nil
}

// Error .
func (m *QueryAcceptSetting) Error(_ containerInterface.IMessageContext, _ error) {}
