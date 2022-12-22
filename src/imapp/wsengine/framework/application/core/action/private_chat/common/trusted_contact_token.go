package privateChatCommon

import (
	"time"
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	"ws/framework/application/core/result/constant"
)

// TrustedContactToken .
type TrustedContactToken struct {
	processor.BaseAction
	UserID string
}

// Start .
func (m *TrustedContactToken) Start(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	contact := context.ResolveContactService().FindByJID(m.UserID)

	// 是否信任关系
	if contact != nil && contact.TrustedContact {
		next()
		return
	}

	iq := message.InfoQuery{
		ID:        context.GenerateRequestID(),
		Namespace: "privacy",
		Type:      message.IqSet,
		To:        types.ServerJID,
		Content: []waBinary.Node{
			{
				Tag: "tokens",
				Content: []waBinary.Node{
					{
						Tag: "token",
						Attrs: waBinary.Attrs{
							"t":    time.Now().Unix(),
							"type": "trusted_contact",
							"jid":  types.NewJID(m.UserID, types.DefaultUserServer),
						},
					},
				},
			},
		},
	}

	m.SendMessageId, err = context.SendIQ(iq)

	return
}

// Receive .
func (m *TrustedContactToken) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	next()

	return nil
}

// Error .
func (m *TrustedContactToken) Error(context containerInterface.IMessageContext, err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.TrustedContactToken,
		Error:      err,
	})
}
