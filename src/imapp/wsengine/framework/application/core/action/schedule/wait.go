package scheduleAction

import (
	"ws/framework/application/constant/message"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
)

// Wait .
type Wait struct {
	processor.BaseAction
	nowWaitSecond uint32
	Second        uint32
}

// Start .
func (m *Wait) Start(_ containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	m.SendMessageId = message.TickClockEvent
	return
}

// Receive .
func (m *Wait) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	if m.nowWaitSecond == m.Second {
		next()
		return
	}

	m.nowWaitSecond++

	return
}

// Error .
func (m *Wait) Error(_ containerInterface.IMessageContext, _ error) {}
