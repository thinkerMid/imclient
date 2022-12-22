package contactNotification

import (
	"fmt"
	"time"
	containerInterface "ws/framework/application/container/abstract_interface"
	messageResultType "ws/framework/application/core/result/constant"
	contactDB "ws/framework/application/data_storage/contact/database"
	"ws/framework/external"
	"ws/framework/utils/xmpp"
)

// TakeNicknameFromMessage .
type TakeNicknameFromMessage struct{}

// Receive .
func (t TakeNicknameFromMessage) Receive(context containerInterface.IMessageContext) (err error) {
	info, parseInfoErr := xmpp.ParseMessageInfo(context.ResolveJID(), context.Message())
	if parseInfoErr != nil || info.IsGroup {
		return
	}

	if len(info.PushName) > 0 {
		key := fmt.Sprintf("%s_%s_nickname", context.ResolveJID().User, info.Sender.User)

		if val, ok := context.ResolveMemoryCache().FindInCache(key); ok {
			cacheNickname := val.(string)

			if info.PushName == cacheNickname {
				return
			}
		}

		context.ResolveMemoryCache().CacheTTL(key, info.PushName, time.Hour)

		context.AppendResult(containerInterface.MessageResult{
			ResultType: messageResultType.ContactNicknameUpdate,
			IContent: external.ProfileUpdate{
				JIDNumber: info.Sender.User,
				Content:   info.PushName,
			},
		})
	}

	context.ResolveContactService().ContextExecute(info.Sender.User, func(contact *contactDB.Contact) {
		// 收到过消息
		contact.UpdateReceiveChat(true)
	})

	return
}
