package groupNotification

import (
	"ws/framework/application/constant"
	"ws/framework/application/constant/binary"
	"ws/framework/application/container/abstract_interface"
	groupDB "ws/framework/application/data_storage/group/database"
)

// UpdateAdmin .
type UpdateAdmin struct{}

// Receive .
func (m UpdateAdmin) Receive(context containerInterface.IMessageContext) error {
	groupId, err := parseGroupNotification(context.Message())
	if err != nil {
		return nil
	}

	mUserID := context.ResolveJID().User

	// 设置管理员
	promoteNode, ok := context.Message().GetOptionalChildByTag("promote")
	if ok {
		if m.findUserID(promoteNode, mUserID) {
			context.ResolveGroupService().ContextExecute(groupId, func(group *groupDB.Group) {
				group.UpdateAdmin(true)
			})
		}

		return constant.AbortedError
	}

	// 取消管理员
	demoteNode, ok := context.Message().GetOptionalChildByTag("demote")
	if ok {
		if m.findUserID(demoteNode, mUserID) {
			context.ResolveGroupService().ContextExecute(groupId, func(group *groupDB.Group) {
				group.UpdateAdmin(false)
			})
		}

		return constant.AbortedError
	}

	return nil
}

func (m *UpdateAdmin) findUserID(node waBinary.Node, dstUserID string) bool {
	participantList := node.GetChildren()
	for i := range participantList {
		participant := participantList[i]

		getter := participant.AttrGetter()
		userID := getter.JID("jid").User

		if userID == dstUserID {
			return true
		}
	}

	return false
}
