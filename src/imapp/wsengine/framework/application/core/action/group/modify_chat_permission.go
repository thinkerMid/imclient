package group

import (
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	"ws/framework/application/core/result/constant"
)

// ChatPermission .
type ChatPermission struct {
	processor.BaseAction
	GroupID string
	Enabled bool // true 所有人  false 只有管理员
}

// Start .
func (m *ChatPermission) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	groupJID := types.NewJID(m.GroupID, types.GroupServer)

	code := "not_announcement"
	if !m.Enabled {
		code = "announcement"
	}

	iq := message.InfoQuery{
		ID:        context.GenerateRequestID(),
		Namespace: "w:g2",
		Type:      message.IqSet,
		To:        groupJID,
		Content:   []waBinary.Node{{Tag: code}},
	}

	m.SendMessageId, err = context.SendIQ(iq)

	return
}

// Receive .
func (m *ChatPermission) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.GroupEditDescPermissionChange,
	})

	next()

	return nil
}

// Error .
func (m *ChatPermission) Error(context containerInterface.IMessageContext, err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.GroupEditDescPermissionChange,
		Error:      err,
	})
}
