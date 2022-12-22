package groupNotification

import (
	"ws/framework/application/constant"
	"ws/framework/application/container/abstract_interface"
	groupDB "ws/framework/application/data_storage/group/database"
)

// UpdateDesc .
type UpdateDesc struct{}

// Receive .
func (m UpdateDesc) Receive(context containerInterface.IMessageContext) error {
	groupId, err := parseGroupNotification(context.Message())
	if err != nil {
		return nil
	}

	descriptionNode, ok := context.Message().GetOptionalChildByTag("description")
	if !ok {
		return nil
	}

	editDescKey := descriptionNode.AttrGetter().String("id")

	context.ResolveGroupService().ContextExecute(groupId, func(group *groupDB.Group) {
		group.UpdateEditDescKey(editDescKey)
	})

	return constant.AbortedError
}
