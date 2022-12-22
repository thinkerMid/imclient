package privateChatNotification

import (
	"ws/framework/application/constant"
	containerInterface "ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/action/common"
	"ws/framework/application/core/processor"
)

// PrivacyToken 对方第一次发送消息时通过通知给自己的一个验证token
type PrivacyToken struct{}

// Receive .
func (r PrivacyToken) Receive(context containerInterface.IMessageContext) (err error) {
	node := context.Message()
	attrGetter := node.AttrGetter()

	privacyTokenType := attrGetter.String("type")
	if privacyTokenType != "privacy_token" {
		return
	}

	jid := attrGetter.JID("from")

	tokens := node.GetChildByTag("tokens")
	child := tokens.GetChildByTag("token")
	tcToToken := child.ContentString()[0]

	// 查找是否有联系人
	contact := context.ResolveContactService().FindByJID(jid.User)
	if contact == nil {
		// 默认创建联系人
		context.ResolveContactService().CreateContact(jid.User, jid.User)

		// 请求他的信息
		context.ResolveMessageChannel().AddMessageProcessor(
			processor.NewOnceIgnoreErrorProcessor(
				[]containerInterface.IAction{
					&common.SubscribeStatus{UserID: jid.User, TcToToken: tcToToken},
					&common.QueryAvatarPreview{UserID: jid.User},
					&common.QueryUserDeviceList{UserID: jid.User},
					&common.QueryDisappearingMode{UserID: jid.User},

					// 真机是不会查询这个的
					&common.QueryAvatarUrl{UserID: jid.User},
				},
			),
		)
	}

	return constant.AbortedError
}
