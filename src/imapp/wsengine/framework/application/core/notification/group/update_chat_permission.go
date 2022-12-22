package groupNotification

import (
	"ws/framework/application/constant"
	"ws/framework/application/container/abstract_interface"
	groupDB "ws/framework/application/data_storage/group/database"
)

// UpdateChatPermission .
type UpdateChatPermission struct{}

// Receive .
func (m UpdateChatPermission) Receive(context containerInterface.IMessageContext) error {
	groupId, err := parseGroupNotification(context.Message())
	if err != nil {
		return nil
	}

	// 仅管理员可发言
	_, ok := context.Message().GetOptionalChildByTag("announcement")
	if ok {
		context.ResolveGroupService().ContextExecute(groupId, func(group *groupDB.Group) {
			group.UpdateChatPermission(false)
		})

		return constant.AbortedError
	}

	// 所有人可发言
	_, ok = context.Message().GetOptionalChildByTag("not_announcement")
	if ok {
		context.ResolveGroupService().ContextExecute(groupId, func(group *groupDB.Group) {
			group.UpdateChatPermission(true)
		})

		return constant.AbortedError
	}

	return nil
}
