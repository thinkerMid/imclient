package groupNotification

import (
	"ws/framework/application/constant"
	containerInterface "ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/action/group"
	"ws/framework/application/core/processor"
)

// JoinGroup .
type JoinGroup struct{}

// Receive .
func (m JoinGroup) Receive(context containerInterface.IMessageContext) error {
	node := context.Message()

	groupID, err := parseGroupNotification(node)
	if err != nil {
		return nil
	}

	add, ok := node.GetOptionalChildByTag("add")
	if !ok {
		return nil
	}

	var memberJIDs []string

	participantList := add.GetChildrenByTag("participant")

	for i := range participantList {
		getter := participantList[i].AttrGetter()
		jid := getter.JID("jid").User

		memberJIDs = append(memberJIDs, jid)
	}

	if len(memberJIDs) > 0 {
		// 查询群成员的设备列表
		context.AddMessageProcessor(processor.NewOnceProcessor(
			[]containerInterface.IAction{
				&group.QueryMemberDeviceList{GroupID: groupID, UserIDs: memberJIDs},
			},
		))
	}

	return constant.AbortedError
}
