package userSettings

import (
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	"ws/framework/application/core/result/constant"
	accountDB "ws/framework/application/data_storage/account/database"
)

// ModifyAvatar .
type ModifyAvatar struct {
	processor.BaseAction
	Content []byte
}

// Start .
func (m *ModifyAvatar) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	m.SendMessageId, err = context.SendIQ(message.InfoQuery{
		Namespace: "w:profile:picture",
		Type:      message.IqSet,
		To:        types.ServerJID,
		ID:        context.GenerateRequestID(),
		Content: []waBinary.Node{{
			Tag: "picture",
			Attrs: waBinary.Attrs{
				"type": "image",
			},
			Content: m.Content,
		}},
	})

	return
}

// Receive .
func (m *ModifyAvatar) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	defer next()

	context.ResolveAccountService().ContextExecute(func(account *accountDB.Account) {
		account.UpdateHaveAvatar(len(m.Content) > 0)
	})

	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.SetAvatar,
	})

	return nil
}

// Error .
func (m *ModifyAvatar) Error(context containerInterface.IMessageContext, err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.SetAvatar,
		Error:      err,
	})
}
