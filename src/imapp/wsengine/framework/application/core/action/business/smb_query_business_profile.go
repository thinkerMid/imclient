package business

import (
	waBinary "ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	containerInterface "ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	businessDB "ws/framework/application/data_storage/business/database"
)

// SMBQueryBusinessProfile .
type SMBQueryBusinessProfile struct {
	processor.BaseAction
	MagicV string
}

// Start .
func (m *SMBQueryBusinessProfile) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	if len(m.MagicV) == 0 {
		m.MagicV = "116"
	}

	profileAttrs := waBinary.Attrs{
		"jid": context.ResolveJID(),
	}

	profile := context.ResolveBusinessService().Context()
	if profile != nil && len(profile.ProfileTag) > 0 {
		profileAttrs["tag"] = profile.ProfileTag
	}

	m.SendMessageId, err = context.SendIQ(message.InfoQuery{
		ID:        context.GenerateRequestID(),
		Namespace: "w:biz",
		Type:      "get",
		To:        types.ServerJID,
		Content: []waBinary.Node{{
			Tag: "business_profile",
			Attrs: waBinary.Attrs{
				"v": m.MagicV,
			},
			Content: []waBinary.Node{{
				Tag:   "profile",
				Attrs: profileAttrs,
			}},
		}},
	})

	return
}

// Receive .
func (m *SMBQueryBusinessProfile) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	profileNode := context.Message().GetChildByTag("profile")

	context.ResolveBusinessService().ContextExecute(func(b *businessDB.BusinessProfile) {
		b.UpdateProfileTag(profileNode.AttrGetter().String("tag"))
	})

	next()

	return nil
}

// Error .
func (m *SMBQueryBusinessProfile) Error(context containerInterface.IMessageContext, err error) {
}
