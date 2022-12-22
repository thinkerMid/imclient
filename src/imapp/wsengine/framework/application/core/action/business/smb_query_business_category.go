package business

import (
	waBinary "ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	containerInterface "ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
)

// SMBQueryBusinessCategory .
type SMBQueryBusinessCategory struct {
	processor.BaseAction
}

// Start .
func (m *SMBQueryBusinessCategory) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	m.SendMessageId, err = context.SendIQ(message.InfoQuery{
		ID:        context.GenerateRequestID(),
		Namespace: "fb:thrift_iq",
		Type:      message.IqGet,
		To:        context.ResolveJID(),
		Content: []waBinary.Node{{
			Tag: "request",
			Attrs: waBinary.Attrs{
				"op":   "profile_typeahead",
				"type": "catkit",
				"v":    "1",
			},
			Content: []waBinary.Node{
				{Tag: "query"},
			},
		}},
	})

	return
}

// Receive .
func (m *SMBQueryBusinessCategory) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	next()

	return nil
}

// Error .
func (m *SMBQueryBusinessCategory) Error(context containerInterface.IMessageContext, err error) {
}
