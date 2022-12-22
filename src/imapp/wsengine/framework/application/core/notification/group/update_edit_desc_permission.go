package groupNotification

import (
	"ws/framework/application/constant"
	"ws/framework/application/container/abstract_interface"
	groupDB "ws/framework/application/data_storage/group/database"
)

// UpdateEditDescPermission .
type UpdateEditDescPermission struct{}

// Receive .
func (m UpdateEditDescPermission) Receive(context containerInterface.IMessageContext) error {
	groupId, err := parseGroupNotification(context.Message())
	if err != nil {
		return nil
	}

	// 仅管理员可编辑
	_, ok := context.Message().GetOptionalChildByTag("locked")
	if ok {
		context.ResolveGroupService().ContextExecute(groupId, func(group *groupDB.Group) {
			group.UpdateEditDescPermission(false)
		})

		return constant.AbortedError
	}

	// 所有人可编辑
	_, ok = context.Message().GetOptionalChildByTag("unlocked")
	if ok {
		context.ResolveGroupService().ContextExecute(groupId, func(group *groupDB.Group) {
			group.UpdateEditDescPermission(true)
		})

		return constant.AbortedError
	}

	return nil
}
