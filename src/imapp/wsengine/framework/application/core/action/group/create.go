package group

import (
	"time"
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	"ws/framework/application/core/result/constant"
	"ws/framework/external"
)

// Create .
type Create struct {
	processor.BaseAction
	GroupName   string
	HaveIcon    bool
	JoinUserIDs []string
}

// Start .
func (m *Create) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	participants := make([]waBinary.Node, len(m.JoinUserIDs))

	for i, id := range m.JoinUserIDs {
		participants[i].Tag = "participant"
		participants[i].Attrs = waBinary.Attrs{
			"jid": types.NewJID(id, types.DefaultUserServer),
		}
	}

	iq := message.InfoQuery{
		ID:        context.GenerateRequestID(),
		Namespace: "w:g2",
		Type:      message.IqSet,
		To:        types.GroupServerJID,
		Content: []waBinary.Node{{
			Tag: "create",
			Attrs: waBinary.Attrs{
				"subject": m.GroupName,
				"key":     time.Now().Unix(),
			},
			Content: participants,
		}},
	}

	m.SendMessageId, err = context.SendIQ(iq)
	return
}

// Receive .
func (m *Create) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	defer next()

	groupInfo := m.parseGroupInfo(context.Message())

	context.ResolveGroupService().CreateGroup(groupInfo.GroupNumber, true)

	// TODO 优化
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.CreateGroup,
		Content:    groupInfo.GroupNumber,
		IContent:   groupInfo,
	})

	return nil
}

// Error .
func (m *Create) Error(context containerInterface.IMessageContext, err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.CreateGroup,
		Error:      err,
	})
}

func (m *Create) parseGroupInfo(node *waBinary.Node) (info external.GroupInfo) {
	right := func(role string) int32 {
		switch role {
		case "admin":
			return 1
		case "superadmin":
			return 2
		}
		return 0
	}

	grp := node.GetChildByTag("group")
	attrs := grp.AttrGetter()
	// 群
	info.GroupNumber = attrs.String("id")
	info.CreateTime = attrs.Int64("creation")
	// 群昵称
	info.Title.Text = attrs.String("subject")
	//info.Title.EditTime = attrs.Int64("s_t")
	//info.Title.Editor = attrs.OptionalJIDOrEmpty("s_o").String()

	// 群描述
	//desc := grp.GetChildByTag("description")
	//attrs = desc.AttrGetter()
	//
	//info.Description.Editor = attrs.jid("participant").String()
	//info.Description.EditTime = attrs.Int64("t")
	//info.Description.EditId = attrs.String("id")
	//
	//var description string
	//body := desc.GetChildByTag("body")
	//descriptions := body.ContentString()
	//if len(descriptions) > 0 {
	//	description = descriptions[0]
	//}
	//info.Description.Text = description

	// 成员
	var members []external.GroupMember
	children := grp.GetChildrenByTag("participant")
	if len(children) > 0 {
		for _, child := range children {
			attrs = child.AttrGetter()

			jid := attrs.JID("jid")
			role, _ := attrs.GetString("type", false)

			member := external.GroupMember{
				MemberNumber: jid.String(),
				Right:        right(role),
			}
			members = append(members, member)
		}
	}
	info.Members = members
	return
}
