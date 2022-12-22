package group

import (
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	"ws/framework/application/core/result/constant"
)

// ModifyDescription .
type ModifyDescription struct {
	processor.BaseAction
	GroupID     string
	DescContent string
}

// Start .
func (m *ModifyDescription) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	group := context.ResolveGroupService().Find(m.GroupID)

	var editDescKey string
	if len(group.EditDescKey) == 0 {
		editDescKey = "none"
	} else {
		editDescKey = group.EditDescKey
	}

	groupJID := types.NewJID(m.GroupID, types.GroupServer)

	iq := message.InfoQuery{
		ID:        context.GenerateRequestID(),
		Namespace: "w:g2",
		Type:      message.IqSet,
		To:        groupJID,
		Content: []waBinary.Node{
			{
				Tag: "description",
				Attrs: waBinary.Attrs{
					"id":   context.GenerateRequestID(),
					"prev": editDescKey,
				},
				Content: []waBinary.Node{
					{
						Tag:     "body",
						Content: m.DescContent,
					},
				},
			},
		},
	}

	m.SendMessageId, err = context.SendIQ(iq)

	return
}

// Receive .
func (m *ModifyDescription) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.ModifyGroupDesc,
	})

	next()

	return nil
}

// Error .
func (m *ModifyDescription) Error(context containerInterface.IMessageContext, err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.ModifyGroupDesc,
		Error:      err,
	})
}
