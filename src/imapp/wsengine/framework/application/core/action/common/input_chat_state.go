package common

import (
	"fmt"
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	"ws/framework/application/core/result/constant"
)

// InputChatState 输入状态
type InputChatState struct {
	processor.BaseAction
	UserID string // 对方
	Input  bool   // 是否输入

	ToGroup bool // 发送给群组的
}

// Start .
func (m *InputChatState) Start(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	err = m.sendInputNode(context)
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.InputChatState,
	})

	next()
	return
}

// Receive .
func (m *InputChatState) Receive(_ containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	return
}

// Error .
func (m *InputChatState) Error(context containerInterface.IMessageContext, err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.InputChatState,
		Error:      err,
	})
}

func (m *InputChatState) sendInputNode(context containerInterface.IMessageContext) (err error) {
	stateTag := "composing"
	if !m.Input {
		stateTag = "paused"
	}

	toServer := types.DefaultUserServer
	if m.ToGroup {
		toServer = types.GroupServer
	}

	_, err = context.SendNode(waBinary.Node{
		Tag: "chatstate",
		Attrs: waBinary.Attrs{
			"to": fmt.Sprintf("%s@%s", m.UserID, toServer),
		},
		Content: []waBinary.Node{{Tag: stateTag}},
	})

	return
}
