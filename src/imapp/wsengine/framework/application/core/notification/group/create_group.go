package groupNotification

import (
	"ws/framework/application/constant"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/action/group"
	"ws/framework/application/core/processor"
	"ws/framework/application/core/result/constant"
)

// CreateGroup .
type CreateGroup struct{}

// Receive .
func (m CreateGroup) Receive(context containerInterface.IMessageContext) error {
	groupID, err := parseGroupNotification(context.Message())
	if err != nil {
		return nil
	}

	node, ok := context.Message().GetOptionalChildByTag("create")
	if !ok {
		return nil
	}

	node, ok = node.GetOptionalChildByTag("group")
	if !ok {
		return nil
	}

	var memberJIDs []string
	var isAdmin bool

	mJID := context.ResolveJID().User
	participantList := node.GetChildren()

	for i := range participantList {
		if participantList[i].Tag != "participant" {
			continue
		}

		getter := participantList[i].AttrGetter()
		jid := getter.JID("jid").User

		if jid == mJID {
			adminType := getter.OptionalString("type")
			isAdmin = len(adminType) > 0
			continue
		}

		memberJIDs = append(memberJIDs, jid)
	}

	context.ResolveGroupService().CreateGroup(groupID, isAdmin)
	context.ResolveSignalProtocolFactory().CreateGroupSession(groupID)

	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.JoinGroup,
		Content:    groupID,
	})

	// 查询群成员的设备
	context.AddMessageProcessor(processor.NewOnceProcessor(
		[]containerInterface.IAction{
			&group.QueryMemberDeviceList{GroupID: groupID, UserIDs: memberJIDs},
		},
	))

	return constant.AbortedError
}
