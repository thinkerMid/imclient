package user

import (
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
)

// QueryTwoFactorAuthentication .
type QueryTwoFactorAuthentication struct {
	processor.BaseAction
}

// Start .
func (m *QueryTwoFactorAuthentication) Start(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	m.SendMessageId, err = context.SendIQ(message.InfoQuery{
		ID:        context.GenerateRequestID(),
		Namespace: "urn:xmpp:whatsapp:account",
		Type:      "get",
		To:        types.ServerJID,
		Content: []waBinary.Node{{
			Tag: "2fa",
		}},
	})

	return
}

// Receive .
func (m *QueryTwoFactorAuthentication) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	// 2fa内容没有node的Tag为code的话，就是关闭状态
	next()

	return nil
}

// Error .
func (m *QueryTwoFactorAuthentication) Error(context containerInterface.IMessageContext, err error) {
}
