package privateChatCommon

import (
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	contactDB "ws/framework/application/data_storage/contact/database"
)

// ChatWith 记录是否发送过消息
type ChatWith struct {
	processor.BaseAction
	UserID string
}

// Start .
func (m *ChatWith) Start(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	context.ResolveContactService().ContextExecute(m.UserID, func(contact *contactDB.Contact) {
		contact.UpdateChatWith(true)
	})

	next()
	return
}

// Receive .
func (m *ChatWith) Receive(_ containerInterface.IMessageContext, _ containerInterface.NextActionFn) error {
	return nil
}

// Error .
func (m *ChatWith) Error(_ containerInterface.IMessageContext, _ error) {}
