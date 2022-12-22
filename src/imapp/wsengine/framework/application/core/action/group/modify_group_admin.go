package group

import (
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	"ws/framework/application/core/result/constant"
)

// ModifyGroupAdmin .
type ModifyGroupAdmin struct {
	processor.BaseAction
	AddOperate bool // 是否添加动作
	GroupID    string
	UserIDs    []string
}

// Start .
func (m *ModifyGroupAdmin) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	groupJID := types.NewJID(m.GroupID, types.GroupServer)
	content := []waBinary.Node{{Tag: "promote"}, {Tag: "demote"}}

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

	if m.AddOperate {
		content[0].Content = userNodes
	} else {
		content[1].Content = userNodes
	}

	iq := message.InfoQuery{
		ID:        context.GenerateRequestID(),
		Namespace: "w:g2",
		Type:      message.IqSet,
		To:        groupJID,
		Content: []waBinary.Node{
			{
				Tag:     "admin",
				Content: content,
			},
		},
	}

	m.SendMessageId, err = context.SendIQ(iq)

	return
}

// Receive .
func (m *ModifyGroupAdmin) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.ModifyGroupAdmin,
	})

	next()

	return nil
}

// Error .
func (m *ModifyGroupAdmin) Error(context containerInterface.IMessageContext, err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.ModifyGroupAdmin,
		Error:      err,
	})
}
