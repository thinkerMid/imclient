package groupNotification

import (
	"ws/framework/application/constant"
	"ws/framework/application/container/abstract_interface"
	messageResultType "ws/framework/application/core/result/constant"
)

// LeftGroup .
type LeftGroup struct{}

// Receive .
func (m LeftGroup) Receive(context containerInterface.IMessageContext) error {
	groupId, err := parseGroupNotification(context.Message())
	if err != nil {
		return nil
	}

	removeParticipant, ok := context.Message().GetOptionalChildByTag("remove")
	if !ok {
		return nil
	}

	var foundMe bool
	var memberJIDs []string

	mJID := context.ResolveJID().User
	participantList := removeParticipant.GetChildren()

	for i := range participantList {
		getter := participantList[i].AttrGetter()
		jid := getter.JID("jid").User

		if jid == mJID {
			foundMe = true
			break
		}

		memberJIDs = append(memberJIDs, jid)
	}

	if foundMe {
		context.ResolveGroupService().DeleteGroup(groupId)
		context.ResolveSenderKeyService().DeleteAllDeviceByGroupID(groupId)

		context.AppendResult(containerInterface.MessageResult{
			ResultType: messageResultType.LeftGroup,
			Content:    groupId,
		})
	} else {
		context.ResolveSenderKeyService().BatchDeleteSenderByGroupIDAndJID(groupId, memberJIDs)
	}

	return constant.AbortedError
}
