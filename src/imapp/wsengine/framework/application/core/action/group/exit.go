package group

import (
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	"ws/framework/application/core/result/constant"
)

// Exit .
type Exit struct {
	processor.BaseAction
	GroupID string
}

// Start .
func (m *Exit) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	groupJID := types.NewJID(m.GroupID, types.GroupServer)

	iq := message.InfoQuery{
		ID:        context.GenerateRequestID(),
		Namespace: "w:g2",
		Type:      message.IqSet,
		To:        types.GroupServerJID,
		Content: []waBinary.Node{
			{
				Tag: "leave",
				Content: []waBinary.Node{
					{
						Tag: "group",
						Attrs: waBinary.Attrs{
							"id": groupJID,
						},
					},
				},
			},
		},
	}

	m.SendMessageId, err = context.SendIQ(iq)

	return
}

// Receive .
func (m *Exit) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	next()

	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.ExitGroup,
	})

	return nil
}

// Error .
func (m *Exit) Error(context containerInterface.IMessageContext, err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.ExitGroup,
		Error:      err,
	})
}
