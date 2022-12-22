package group

import (
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	"ws/framework/application/core/result/constant"
)

// ModifyName .
type ModifyName struct {
	processor.BaseAction
	GroupID string
	Name    string
}

// Start .
func (m *ModifyName) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	groupJID := types.NewJID(m.GroupID, types.GroupServer)

	iq := message.InfoQuery{
		ID:        context.GenerateRequestID(),
		Namespace: "w:g2",
		Type:      message.IqSet,
		To:        groupJID,
		Content: []waBinary.Node{
			{
				Tag:     "subject",
				Content: m.Name,
			},
		},
	}

	m.SendMessageId, err = context.SendIQ(iq)

	return
}

// Receive .
func (m *ModifyName) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.ModifyGroupName,
	})

	next()

	return nil
}

// Error .
func (m *ModifyName) Error(context containerInterface.IMessageContext, err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.ModifyGroupName,
		Error:      err,
	})
}
