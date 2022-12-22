package companion

import (
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
)

type CompanionRemove struct {
	processor.BaseAction
	JID types.JID
}

func (m *CompanionRemove) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	m.SendMessageId, err = context.SendIQ(
		message.InfoQuery{
			ID:        context.GenerateRequestID(),
			Namespace: "md",
			Type:      message.IqSet,
			To:        types.ServerJID,
			Content: []waBinary.Node{{
				Tag: "remove-companion-device",
				Attrs: waBinary.Attrs{
					"all":    "true",
					"reason": "user_initiated",
				},
			}},
		},
	)
	return nil
}

func (m *CompanionRemove) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	next()
	return nil
}

func (m *CompanionRemove) Error(context containerInterface.IMessageContext, err error) {

}
