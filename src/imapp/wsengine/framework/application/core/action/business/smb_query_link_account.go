package business

import (
	waBinary "ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	containerInterface "ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
)

// SMBQueryLinkedAccount .
type SMBQueryLinkedAccount struct {
	processor.BaseAction
}

// Start .
func (m *SMBQueryLinkedAccount) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	m.SendMessageId, err = context.SendNode(waBinary.Node{
		Tag: "iq",
		Attrs: waBinary.Attrs{
			"id":      context.GenerateRequestID(),
			"smax_id": "42",
			"from":    context.ResolveJID(),
			"to":      types.ServerJID,
			"type":    message.IqGet,
			"xmlns":   "fb:thrift_iq",
		},
		Content: []waBinary.Node{
			{Tag: "linked_accounts"},
		},
	})

	return
}

// Receive .
func (m *SMBQueryLinkedAccount) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	next()

	return nil
}

// Error .
func (m *SMBQueryLinkedAccount) Error(context containerInterface.IMessageContext, err error) {
}
