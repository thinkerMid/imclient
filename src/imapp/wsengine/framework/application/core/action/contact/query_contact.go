package contact

import (
	"errors"
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	"ws/framework/application/core/result/constant"
	"ws/framework/utils/xmpp"
)

// Query .
type Query struct {
	processor.BaseAction
	UserID string // 手机号
}

// Start .
func (m *Query) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	m.SendMessageId, err = context.SendIQ(
		xmpp.UsyncIQTemplate(context, "query", "add",
			[]waBinary.Node{
				{Tag: "query", Content: []waBinary.Node{
					{Tag: "business",
						Content: []waBinary.Node{
							{Tag: "verified_name"},
							{Tag: "profile",
								Attrs: waBinary.Attrs{
									"v": "372",
								},
							},
						},
					},
					{Tag: "contact"},
					{Tag: "disappearing_mode"},
				}},
				{Tag: "list", Content: []waBinary.Node{{
					Tag: "user",
					Content: []waBinary.Node{{
						Tag:     "contact",
						Content: "+" + m.UserID,
					}},
				}}},
			},
		),
	)

	return
}

// Receive .
func (m *Query) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	nodes := context.Message().GetChildren()

	if len(nodes) > 0 {
		nodes = nodes[0].GetChildrenByTag("list")
	}

	if len(nodes) > 0 {
		nodes = nodes[0].GetChildren()
	}

	for i := range nodes {
		n := nodes[i]
		if n.Tag != "user" {
			continue
		}

		_, ok := n.Attrs["jid"].(types.JID)
		if !ok {
			continue
		}

		contactNode := n.GetChildByTag("contact")
		if contactNode.AttrGetter().String("type") != "in" {
			break
		}

		next()
		return nil
	}

	return errors.New("not find dstUserID in whatsapp")
}

func (m *Query) Error(context containerInterface.IMessageContext, err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.QueryContact,
		Error:      err,
	})
}
