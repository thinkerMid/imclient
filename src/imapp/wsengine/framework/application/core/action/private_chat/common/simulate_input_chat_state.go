package privateChatCommon

import (
	"fmt"
	"math"
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	"ws/framework/application/core/processor/interval_timer"
	"ws/framework/application/core/result/constant"
	"ws/framework/utils"
)

// SimulateInputChatState 欺骗对方正在输入状态的action
type SimulateInputChatState struct {
	processor.BaseAction
	UserID            string // 对方
	TotalInputLatency uint32 // 输入总耗时 多少字则多少秒
	TotalCount        int32  // 发送次数

	inputStopCount   uint32 // 输入停止次数
	onceInputLatency uint32 // 单次输入耗时

	inputState bool // 输入状态
	restState  bool // 休息状态

	timer *intervalTimer.Timer
}

// Start .
func (m *SimulateInputChatState) Start(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	contact := context.ResolveContactService().FindByJID(m.UserID)
	// 没聊过天不处理
	if contact == nil || contact.ChatWith == false {
		next()
		return
	}

	// 每输入5秒暂停一次  停顿次数 = 总时长 / 5秒
	m.inputStopCount = uint32(math.Ceil(float64(m.TotalInputLatency / 5)))
	// 不足则默认一次
	if m.inputStopCount == 0 {
		m.inputStopCount = 1
	} else if m.inputStopCount > 5 {
		// 最多5次
		m.inputStopCount = 5
	}

	m.TotalCount = int32(m.inputStopCount)

	// 每次输入多少时长后停顿  输入时长 = 总时长 / 停止次数
	m.onceInputLatency = uint32(math.Ceil(float64(m.TotalInputLatency / m.inputStopCount)))
	if m.onceInputLatency > 24 {
		// 输入时长最多24秒 减少1秒确保通信延时也能续上输入状态
		m.onceInputLatency = 24
	}

	m.timer = intervalTimer.New(m.onceInputLatency)

	// 默认输入中
	m.inputState = true
	err = m.sendInputNode(context)

	// 用于接收时钟通知包
	m.SendMessageId = message.TickClockEvent

	return
}

// Receive .
func (m *SimulateInputChatState) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	// 计时
	m.timer.Timing()
	if !m.timer.TimingEnd() {
		return
	}

	// 从休息状态恢复
	if m.restState {
		m.restState = false

		// 继续输入
		m.inputState = !m.inputState
		err = m.sendInputNode(context)

		// 开始下一个计时
		m.timer.ChangeWaitSecond(m.onceInputLatency)
		return
	}

	// 变换输入状态
	m.inputState = !m.inputState
	err = m.sendInputNode(context)

	// 减去一次停顿
	m.inputStopCount--

	// 不需要再停顿了
	if m.inputStopCount == 0 {
		next()
		return
	}

	// 休息片刻的时间
	waitNext := uint32(utils.RandInt64(1, 2))
	if waitNext == 0 {
		waitNext = 1
	}

	m.restState = true
	m.timer.ChangeWaitSecond(waitNext)

	return
}

// Error .
func (m *SimulateInputChatState) Error(context containerInterface.IMessageContext, err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.InputChatState,
		Error:      err,
	})
}

func (m *SimulateInputChatState) sendInputNode(context containerInterface.IMessageContext) (err error) {
	stateTag := "composing"
	if !m.inputState {
		stateTag = "paused"
	}

	_, err = context.SendNode(waBinary.Node{
		Tag: "chatstate",
		Attrs: waBinary.Attrs{
			"to": fmt.Sprintf("%s@%s", m.UserID, types.DefaultUserServer),
		},
		Content: []waBinary.Node{{Tag: stateTag}},
	})

	return
}
