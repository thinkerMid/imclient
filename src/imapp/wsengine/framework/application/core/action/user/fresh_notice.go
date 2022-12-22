package user

import (
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
)

// FreshNotice .
type FreshNotice struct {
	processor.BaseAction
}

// Start .
func (m *FreshNotice) Start(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	configuration := context.ResolveWhatsappConfiguration()

	iq := message.InfoQuery{
		ID:        context.GenerateRequestID(),
		Namespace: "tos",
		Type:      message.IqGet,
		To:        types.ServerJID,
		Content: []waBinary.Node{{
			Tag: "request",
			Content: []waBinary.Node{{
				Tag: "notice",
				Attrs: waBinary.Attrs{
					"id": configuration.NoticeId,
				},
			}},
		}},
	}

	m.SendMessageId, err = context.SendIQ(iq)

	return
}

// Receive .
func (m *FreshNotice) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	next()

	return nil
}

// Error .
func (m *FreshNotice) Error(_ containerInterface.IMessageContext, _ error) {}
