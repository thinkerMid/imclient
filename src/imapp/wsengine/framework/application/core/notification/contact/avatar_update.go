package contactNotification

import (
	"ws/framework/application/constant"
	containerInterface "ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/action/common"
	"ws/framework/application/core/processor"
	messageResultType "ws/framework/application/core/result/constant"
	contactDB "ws/framework/application/data_storage/contact/database"
	"ws/framework/external"
)

// AvatarUpdate 头像更新
type AvatarUpdate struct {
	JID string
}

// Receive .
func (s *AvatarUpdate) Receive(context containerInterface.IMessageContext) (err error) {
	node := context.Message()

	ag := node.AttrGetter()
	if ag.String("type") != "picture" {
		return
	}

	// [S] <notification from="85256048573@s.whatsapp.net" id="530421206" notify="hhhhggffd" t="1661341815" type="picture"><set id="1661341815" jid="85256048573@s.whatsapp.net"/></notification>
	// [S] <notification from="85296475450@s.whatsapp.net" id="1961829908" notify="hhhhggffd" t="1661341845" type="picture"><set hash="Fg8T"/></notification>
	// [S] <notification from="85256048573@s.whatsapp.net" id="272688577" notify="hhhhggffd" t="1661341739" type="picture"><delete jid="85256048573@s.whatsapp.net"/></notification>

	child, ok := node.GetOptionalChildByTag("delete")
	if ok {
		childAttrGetter := child.AttrGetter()
		jid := childAttrGetter.JID("jid")
		s.JID = jid.User

		context.ResolveContactService().ContextExecute(jid.User, func(contact *contactDB.Contact) {
			contact.UpdateHaveAvatar(false)
		})

		context.AppendResult(containerInterface.MessageResult{
			ResultType: messageResultType.ContactAvatarUpdate,
			IContent: external.ProfileUpdate{
				JIDNumber: jid.User,
			},
		})

		return constant.AbortedError
	}

	child, ok = node.GetOptionalChildByTag("set")
	if ok {
		childAttrGetter := child.AttrGetter()

		if len(childAttrGetter.String("hash")) == 0 {
			jid := childAttrGetter.JID("jid")
			s.JID = jid.User

			// 查头像地址 这个真机是查Buffer 只有在看它资料的大图才有URL
			p := processor.NewOnceProcessor(
				[]containerInterface.IAction{
					&common.QueryAvatarUrl{UserID: jid.User},
				},
			)
			context.AddMessageProcessor(p)
		}

		return constant.AbortedError
	}

	return
}
