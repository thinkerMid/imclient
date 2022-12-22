package processor

import (
	"ws/framework/application/constant/message"
	"ws/framework/application/container/abstract_interface"
)

// NewTimerProcessor .
func NewTimerProcessor(fn containerInterface.MakeActionFn, sets ...SetConfigFn) containerInterface.IProcessor {
	p := timerProcessor{}

	for _, set := range sets {
		set(&p)
	}

	p.makeActionFn = fn

	return &p
}

// timerProcessor 定时唤醒的处理器
type timerProcessor struct {
	BaseProcessor
	makeActionFn containerInterface.MakeActionFn
}

// ProcessMessage .
func (m *timerProcessor) ProcessMessage(context containerInterface.IProcessControl) {
	if !m.Timer.TimingEnd() {
		if context.MessageContext().Message().Tag != message.TickClockEvent {
			return
		}

		// 计时
		m.Timer.Timing()

		// 如果未计时结束 不放行
		if !m.Timer.TimingEnd() {
			return
		}
	}

	// 继续循环
	if m.loop {
		m.Timer.Reset()
	} else {
		context.Reject()
	}

	executeActionList := m.makeActionFn()
	process := NewOnceProcessor(executeActionList, AliasName("timer"))
	context.MessageChannel().AddMessageProcessor(process)
}

// ProcessorType .
func (m *timerProcessor) ProcessorType() containerInterface.ProcessorType {
	return containerInterface.TimerProcessorType
}

// DumpInfo .
func (m *timerProcessor) DumpInfo() containerInterface.ProcessorDump {
	dump := containerInterface.ProcessorDump{
		PID:                 m.pid,
		AliasName:           m.aliasName,
		ProcessorType:       m.ProcessorType().String(),
		CurrentActionName:   "/",
		CurrentActionStatus: "/",
		ActionQueue:         "/",
	}

	return dump
}
