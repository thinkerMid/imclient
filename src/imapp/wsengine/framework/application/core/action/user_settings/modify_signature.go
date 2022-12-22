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

// ModifySignature .
type ModifySignature struct {
	processor.BaseAction
	Content string
}

// Start .
func (m *ModifySignature) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	m.SendMessageId, err = context.SendIQ(message.InfoQuery{
		ID:        context.GenerateRequestID(),
		Namespace: "status",
		Type:      "set",
		To:        types.ServerJID,
		Content: []waBinary.Node{{
			Tag:     "status",
			Content: m.Content,
		}},
	})

	return
}

// Receive .
func (m *ModifySignature) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	defer next()

	context.ResolveAccountService().ContextExecute(func(account *accountDB.Account) {
		account.UpdateSignature(m.Content)
	})

	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.SetSignature,
	})

	return nil
}

// Error .
func (m *ModifySignature) Error(context containerInterface.IMessageContext, err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.SetSignature,
		Error:      err,
	})
}
