package user

import (
	waBinary "ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
)

// ModifyTwoFactorAuthentication .
type ModifyTwoFactorAuthentication struct {
	processor.BaseAction
	Code string // 空就是关闭两步验证
}

// Start .
func (m *ModifyTwoFactorAuthentication) Start(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	codeNode := waBinary.Node{
		Tag: "code",
	}

	if len(m.Code) > 0 {
		codeNode.Content = m.Code
	}

	m.SendMessageId, err = context.SendIQ(message.InfoQuery{
		ID:        context.GenerateRequestID(),
		Namespace: "urn:xmpp:whatsapp:account",
		Type:      message.IqSet,
		To:        types.ServerJID,
		Content: []waBinary.Node{{
			Tag:     "2fa",
			Content: []waBinary.Node{codeNode},
		}},
	})

	return
}

// Receive .
func (m *ModifyTwoFactorAuthentication) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	next()

	return nil
}

// Error .
func (m *ModifyTwoFactorAuthentication) Error(context containerInterface.IMessageContext, err error) {
}
