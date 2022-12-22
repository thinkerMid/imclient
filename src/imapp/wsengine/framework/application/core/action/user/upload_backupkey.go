package user

import (
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
)

// UploadBackupKey .
// 一定会执行两次的，有两个backupKey要上传
type UploadBackupKey struct {
	processor.BaseAction
	UploadBackupKeyIndexOne bool
}

// Start .
func (m *UploadBackupKey) Start(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	token := context.ResolveRegistrationTokenService().Context()

	key := token.BackupKey
	if !m.UploadBackupKeyIndexOne {
		key = token.BackupKey2
	}

	m.SendMessageId, err = context.SendIQ(message.InfoQuery{
		ID:        context.GenerateRequestID(),
		Namespace: "urn:xmpp:whatsapp:account",
		Type:      "get",
		To:        types.ServerJID,
		Content: []waBinary.Node{{
			Tag: "crypto",
			Attrs: waBinary.Attrs{
				"action": "create",
			},
			Content: []waBinary.Node{{
				Tag:     "google",
				Content: key,
			}},
		}},
	})

	next()

	return
}

// Receive .
func (m *UploadBackupKey) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	return nil
}

// Error .
func (m *UploadBackupKey) Error(_ containerInterface.IMessageContext, _ error) {
}
