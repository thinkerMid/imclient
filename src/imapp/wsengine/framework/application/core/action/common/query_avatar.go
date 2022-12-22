package common

import (
	"strings"
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	messageResultType "ws/framework/application/core/result/constant"
	accountDB "ws/framework/application/data_storage/account/database"
	contactDB "ws/framework/application/data_storage/contact/database"
	"ws/framework/external"
)

// QueryAvatarUrl 查头像的URL
type QueryAvatarUrl struct {
	processor.BaseAction
	UserID string
}

// RaiseErrorWhenNodeError 是否抛出异常
func (m *QueryAvatarUrl) RaiseErrorWhenNodeError() bool {
	return false
}

// Start .
func (m *QueryAvatarUrl) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	targetJID := types.NewJID(m.UserID, types.DefaultUserServer)

	m.SendMessageId, err = context.SendIQ(message.InfoQuery{
		Namespace: "w:profile:picture",
		Type:      message.IqGet,
		To:        types.ServerJID,
		Target:    targetJID,
		ID:        context.GenerateRequestID(),
		Content: []waBinary.Node{{
			Tag: "picture",
			Attrs: waBinary.Attrs{
				"query": "url",
				"type":  "image",
			},
		}},
	})

	return
}

// Receive .
func (m *QueryAvatarUrl) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	defer next()

	picture, _ := context.Message().GetOptionalChildByTag("picture")
	pictureUrl := picture.AttrGetter().String("url")

	// 如果是自己的
	if context.ResolveJID().User == m.UserID {
		context.ResolveAccountService().ContextExecute(func(account *accountDB.Account) {
			account.UpdateHaveAvatar(len(pictureUrl) > 0)
		})

		context.AppendResult(containerInterface.MessageResult{
			ResultType: messageResultType.GetAvatar,
			Content:    pictureUrl,
		})

		return
	}

	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.ContactAvatarUpdate,
		IContent: external.ProfileUpdate{
			JIDNumber: m.UserID,
			Content:   pictureUrl,
		},
	})

	context.ResolveContactService().ContextExecute(m.UserID, func(contact *contactDB.Contact) {
		contact.UpdateHaveAvatar(len(pictureUrl) > 0)
	})

	return
}

// Error .
func (m *QueryAvatarUrl) Error(context containerInterface.IMessageContext, err error) {
	// 如果是自己的
	if context.ResolveJID().User == m.UserID {
		context.ResolveAccountService().ContextExecute(func(account *accountDB.Account) {
			account.UpdateHaveAvatar(false)
		})

		context.AppendResult(containerInterface.MessageResult{
			ResultType: messageResultType.GetAvatar,
		})
		return
	}

	errorStr := err.Error()

	// 404的异常当做空头像 并不是一个真正的异常
	if strings.Contains(errorStr, "404") {
		context.AppendResult(containerInterface.MessageResult{
			ResultType: messageResultType.ContactAvatarUpdate,
			IContent: external.ProfileUpdate{
				JIDNumber: m.UserID,
			},
		})

		context.ResolveContactService().ContextExecute(m.UserID, func(contact *contactDB.Contact) {
			contact.UpdateHaveAvatar(false)
		})
	}
}
