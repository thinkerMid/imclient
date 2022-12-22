package contact

import (
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	"ws/framework/application/core/result/constant"
	"ws/framework/utils/xmpp"
)

// Delete .
type Delete struct {
	processor.BaseAction
	UserID string // 手机号
}

// Start .
func (m *Delete) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	usyncIQ := xmpp.UsyncIQTemplate(context, "delta", "interactive",
		[]waBinary.Node{
			{Tag: "query", Content: []waBinary.Node{
				{Tag: "business", Content: []waBinary.Node{
					{Tag: "verified_name"},
					{Tag: "profile",
						Attrs: waBinary.Attrs{
							"v": "116",
						},
					},
				}},
				{Tag: "contact"},
				{Tag: "disappearing_mode"},
				{Tag: "sidelist"},
				{Tag: "status"},
			}},
			{Tag: "list", Content: []waBinary.Node{{
				Tag:   "user",
				Attrs: waBinary.Attrs{"jid": types.NewJID(m.UserID, types.DefaultUserServer).String()},
				Content: []waBinary.Node{{Tag: "contact",
					Attrs: waBinary.Attrs{"type": "delete"},
				}},
			}}},
		},
	)

	m.SendMessageId, err = context.SendIQ(usyncIQ)

	return
}

// Receive .
func (m *Delete) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.DeleteContact,
	})

	// 只删除联系人关系
	_ = context.ResolveContactService().DeleteByJID(m.UserID)

	next()
	return nil
}

func (m *Delete) Error(context containerInterface.IMessageContext, err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.DeleteContact,
		Error:      err,
	})
}
