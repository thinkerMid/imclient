package user

import (
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
)

// UploadRecoveryToken .
type UploadRecoveryToken struct {
	processor.BaseAction
}

// Start .
func (m *UploadRecoveryToken) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	token := context.ResolveRegistrationTokenService().Context()

	m.SendMessageId, err = context.SendIQ(message.InfoQuery{
		ID:        context.GenerateRequestID(),
		Namespace: "w:auth:token",
		Type:      "set",
		To:        types.ServerJID,
		Content: []waBinary.Node{{
			Tag:     "token",
			Content: token.RecoveryToken,
		}},
	})

	return
}

// Receive .
func (m *UploadRecoveryToken) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	next()

	return nil
}

// Error .
func (m *UploadRecoveryToken) Error(context containerInterface.IMessageContext, err error) {
}
