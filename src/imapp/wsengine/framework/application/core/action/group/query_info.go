package group

import (
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	"ws/framework/application/core/result/constant"
	groupDB "ws/framework/application/data_storage/group/database"
	"ws/framework/external"
)

// QueryInfo .
type QueryInfo struct {
	processor.BaseAction
	GroupID string
}

// Start .
func (m *QueryInfo) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	groupJID := types.NewJID(m.GroupID, types.GroupServer)

	iq := message.InfoQuery{
		ID:        context.GenerateRequestID(),
		Namespace: "w:g2",
		Type:      message.IqGet,
		To:        groupJID,
		Content: []waBinary.Node{
			{
				Tag: "query",
			},
		},
	}

	m.SendMessageId, err = context.SendIQ(iq)

	return
}

// Receive .
func (m *QueryInfo) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	groupInfo := m.parseGroupInfo(context.Message())

	// 更新编辑的Key
	context.ResolveGroupService().ContextExecute(m.GroupID, func(group *groupDB.Group) {
		group.UpdateEditDescKey(groupInfo.Description.EditId)
	})

	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.QueryGroupInfo,
		IContent:   groupInfo,
	})

	next()
	return nil
}

// Error .
func (m *QueryInfo) Error(context containerInterface.IMessageContext, err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.QueryGroupInfo,
		Error:      err,
	})
}

func (m *QueryInfo) parseGroupInfo(node *waBinary.Node) (info external.GroupInfo) {
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
	info.GroupNumber = types.NewJID(attrs.String("id"), types.GroupServer).String()
	info.CreateTime = attrs.Int64("creation")
	// 群昵称
	info.Title.Text = attrs.String("subject")
	info.Title.EditTime = attrs.Int64("s_t")
	info.Title.Editor = attrs.OptionalJIDOrEmpty("s_o").String()

	// 群描述
	desc := grp.GetChildByTag("description")
	attrs = desc.AttrGetter()

	info.Description.Editor = attrs.JID("participant").String()
	info.Description.EditTime = attrs.Int64("t")
	info.Description.EditId = attrs.String("id")

	var description string
	body := desc.GetChildByTag("body")
	descriptions := body.ContentString()
	if len(descriptions) > 0 {
		description = descriptions[0]
	}
	info.Description.Text = description

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
