package group

import (
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	"ws/framework/application/core/result/constant"
)

// ModifyGroupMember .
type ModifyGroupMember struct {
	processor.BaseAction
	GroupID    string
	UserIDs    []string
	AddOperate bool // 是否添加动作
}

// Start .
func (m *ModifyGroupMember) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	groupJID := types.NewJID(m.GroupID, types.GroupServer)

	userNodes := make([]waBinary.Node, 0)
	for _, id := range m.UserIDs {
		node := waBinary.Node{
			Tag: "participant",
			Attrs: waBinary.Attrs{
				"jid": types.NewJID(id, types.DefaultUserServer),
			},
		}

		userNodes = append(userNodes, node)
	}

	tag := "add"
	if !m.AddOperate {
		tag = "remove"
	}

	iq := message.InfoQuery{
		ID:        context.GenerateRequestID(),
		Namespace: "w:g2",
		Type:      message.IqSet,
		To:        groupJID,
		Content: []waBinary.Node{
			{
				Tag:     tag,
				Content: userNodes,
			},
		},
	}

	m.SendMessageId, err = context.SendIQ(iq)

	return
}

// Receive .
func (m *ModifyGroupMember) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.GroupMemberChange,
	})

	next()

	return nil
}

// Error .
func (m *ModifyGroupMember) Error(context containerInterface.IMessageContext, err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.GroupMemberChange,
		Error:      err,
	})
}
