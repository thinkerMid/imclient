package user

import (
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
)

// QueryParticipatingGroups .
type QueryParticipatingGroups struct {
	processor.BaseAction
}

// Start .
func (m *QueryParticipatingGroups) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	m.SendMessageId, err = context.SendIQ(message.InfoQuery{
		ID:        context.GenerateRequestID(),
		Namespace: "w:g2",
		Type:      "get",
		To:        types.GroupServerJID,
		Content: []waBinary.Node{{
			Tag: "participating",
			Content: []waBinary.Node{{
				Tag: "description",
			}, {
				Tag: "participants",
			}},
		}},
	})

	return
}

// Receive .
func (m *QueryParticipatingGroups) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	next()

	return nil
}

// Error .
func (m *QueryParticipatingGroups) Error(context containerInterface.IMessageContext, err error) {
}
