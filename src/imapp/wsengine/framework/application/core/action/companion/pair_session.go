package companion

import (
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
)

type CompanionPairSession struct {
	processor.BaseAction
	JID types.JID
}

func (m *CompanionPairSession) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	m.SendMessageId, err = context.SendIQ(
		message.InfoQuery{
			ID:        context.GenerateRequestID(),
			Namespace: "encrypt",
			Type:      message.IqGet,
			To:        types.ServerJID,
			Content: []waBinary.Node{{
				Tag: "key",
				Content: []waBinary.Node{
					{
						Tag: "user",
						Attrs: waBinary.Attrs{
							"jid": "85268215347@s.whatsapp.net", //m.UserID.String(),
						},
					},
				},
			}},
		},
	)
	return nil
}

func (m *CompanionPairSession) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	next()
	return nil
}

func (m *CompanionPairSession) Error(context containerInterface.IMessageContext, err error) {

}
