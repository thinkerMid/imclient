package processor

import (
	gormLogger "gorm.io/gorm/logger"
	"runtime"
	"ws/framework/application/constant"
	"ws/framework/application/constant/message"
	"ws/framework/application/container/abstract_interface"
)

const (
	// 等待激活
	waitingActive uint8 = iota
	// 等待计时
	waitingTiming
	// 等待执行结果
	waitingExecuteResult
)

var (
	actionErrorTmp    = "pid=%v alias=%s queue=%s progress=%v/%v %s action err: `%s`"
	startupSuccessTmp = "pid=%v alias=%s queue=%s progress=%v/%v startup [%s] action `success`"
	startupFailureTmp = "pid=%v alias=%s queue=%s progress=%v/%v startup [%s] action `failure`"
	executeSuccessTmp = "pid=%v alias=%s queue=%s progress=%v/%v execute [%s] action `success`"
	executeFailureTmp = "pid=%v alias=%s queue=%s progress=%v/%v execute [%s] action `failure`"
	executeDoneTmp    = "pid=%v alias=%s queue=%s progress=%v/%v execute done"
)

func init() {
	if runtime.GOOS == "windows" {
		actionErrorTmp = gormLogger.BlueBold + actionErrorTmp + gormLogger.Reset
		startupSuccessTmp = gormLogger.BlueBold + startupSuccessTmp + gormLogger.Reset
		startupFailureTmp = gormLogger.BlueBold + startupFailureTmp + gormLogger.Reset
		executeSuccessTmp = gormLogger.BlueBold + executeSuccessTmp + gormLogger.Reset
		executeFailureTmp = gormLogger.BlueBold + executeFailureTmp + gormLogger.Reset
		executeDoneTmp = gormLogger.BlueBold + executeDoneTmp + gormLogger.Reset
	}
}

// LogicProcessor .
type LogicProcessor struct {
	BaseProcessor

	next            bool
	state           uint8
	markActionIndex uint8

	resultProcessor containerInterface.LocalResultProcessor
}

func newLogicProcessor() LogicProcessor {
	return LogicProcessor{}
}

// Start .
func (m *LogicProcessor) Start(context containerInterface.IProcessControl) {
	if m.Trigger.DefaultTrigger() {
		m.onProcessorStart(context)
		m.startAction(context)
		m.onNext(context)

		m.Trigger.Disable()
	}

	if m.ActionQueue.Empty() {
		context.Reject()
	}
}

// ProcessMessage .
func (m *LogicProcessor) ProcessMessage(context containerInterface.IProcessControl) {
	msgCtx := context.MessageContext()

	switch m.state {
	case waitingActive:
		if !m.Trigger.ContentedCondition(msgCtx.Message()) {
			return
		}

		m.onProcessorStart(context)
		m.startAction(context)
	case waitingTiming:
		if msgCtx.Message().Tag != message.TickClockEvent {
			return
		}

		// 计时
		m.Timer.Timing()

		// 如果未计时结束 不放行
		if !m.Timer.TimingEnd() {
			return
		}

		m.startAction(context)
	case waitingExecuteResult:
		m.onReceive(context)
	}

	m.onNext(context)

	if m.ActionQueue.Empty() {
		context.Reject()
	}
}

func (m *LogicProcessor) switchNext() {
	m.next = true
}

func (m *LogicProcessor) onNext(context containerInterface.IProcessControl) {
	for m.next {
		m.next = false

		// 是否还有
		action := m.ActionQueue.Pop()
		if action == nil {
			break
		}

		// 重置计时器，等待下一次间隔
		m.Timer.Reset()

		// 如果未计时结束
		if !m.Timer.TimingEnd() {
			// 等计时器满足了再来执行
			m.state = waitingTiming
			return
		}

		m.startAction(context)
	}

	// 默认等待结果状态
	m.state = waitingExecuteResult
}

// 处理action接收的结果
func (m *LogicProcessor) onReceive(context containerInterface.IProcessControl) {
	msgCtx := context.MessageContext()
	action := m.ActionQueue.Current()

	if len(action.ReceiveID()) == 0 || action.ReceiveID() != msgCtx.Message().ID() {
		return
	}

	// 向上告知这个包不需要再继续传递
	context.Aborted()

	err := msgCtx.Message().HasError()
	if err == nil {
		err = action.Receive(msgCtx, m.switchNext)
	}

	if err != nil {
		m.Logger.Infof(
			executeFailureTmp,
			m.ID(), m.AliasName(), m.ActionQueue.Name(), action.ActionIndex(), m.ActionQueue.Size(), action.ActionName(),
		)

		for _, monitor := range m.Monitors {
			monitor.ActionExecuteFailure(action, msgCtx)
		}

		m.onActionError(context, err)
	} else {
		m.Logger.Debugf(
			executeSuccessTmp,
			m.ID(), m.AliasName(), m.ActionQueue.Name(), action.ActionIndex(), m.ActionQueue.Size(), action.ActionName(),
		)

		for _, monitor := range m.Monitors {
			monitor.ActionExecuteSuccess(action, msgCtx)
		}

		m.dispatchMessageResult(context.MessageContext())
	}
}

// 执行启动新的action流程
func (m *LogicProcessor) startAction(context containerInterface.IProcessControl) {
	action := m.ActionQueue.Current()
	msgCtx := context.MessageContext()

	m.markActionIndex = action.ActionIndex()

	for _, monitor := range m.Monitors {
		monitor.OnActionStartBefore(action, msgCtx)
	}

	err := action.Start(msgCtx, m.switchNext)

	for _, monitor := range m.Monitors {
		monitor.OnActionStartAfter(action, msgCtx)
	}

	if err != nil {
		m.Logger.Infof(
			startupFailureTmp,
			m.ID(), m.AliasName(), m.ActionQueue.Name(), action.ActionIndex(), m.ActionQueue.Size(), action.ActionName(),
		)

		for _, monitor := range m.Monitors {
			monitor.OnActionStartFail(action, msgCtx)
		}

		m.onActionError(context, err)
	} else {
		m.Logger.Debugf(
			startupSuccessTmp,
			m.ID(), m.AliasName(), m.ActionQueue.Name(), action.ActionIndex(), m.ActionQueue.Size(), action.ActionName(),
		)

		for _, monitor := range m.Monitors {
			monitor.OnActionStartSuccess(action, msgCtx)
		}

		m.dispatchMessageResult(context.MessageContext())
	}
}

// action异常
func (m *LogicProcessor) onActionError(context containerInterface.IProcessControl, err error) {
	_, ignoreError := err.(*constant.IgnoreError)

	action := m.ActionQueue.Current()
	action.Error(context.MessageContext(), err)

	if action.RaiseErrorWhenNodeError() {
		m.ActionQueue.OnError(err)
	}

	if !ignoreError {
		m.Logger.Infof(
			actionErrorTmp,
			m.ID(), m.AliasName(), m.ActionQueue.Name(), action.ActionIndex(), m.ActionQueue.Size(), action.ActionName(), err.Error(),
		)
	}

	m.dispatchMessageResult(context.MessageContext())
	m.next = true
}

// 派发消息结果
func (m *LogicProcessor) dispatchMessageResult(context containerInterface.IMessageContext) {
	// TODO 应该只调度新的消息结果 不应该从头遍历
	if m.resultProcessor != nil {
		results := context.MessageResult()
		for i := range results {
			m.resultProcessor.ProcessResult(context, &results[i])
		}
	}
}

// 处理器启动
func (m *LogicProcessor) onProcessorStart(context containerInterface.IProcessControl) {
	for _, monitor := range m.Monitors {
		monitor.OnStart(context.AppIocContainer())
	}
}

// OnChannelError .
func (m *LogicProcessor) OnChannelError(context containerInterface.IProcessControl, err error) {
	// 指派给当前的action
	m.onActionError(context, err)
}

// OnDestroy .
func (m *LogicProcessor) OnDestroy(context containerInterface.IProcessControl) {
	for _, monitor := range m.Monitors {
		monitor.OnExit(context.AppIocContainer())
	}

	if m.resultProcessor != nil {
		m.resultProcessor.OnDestroy()
	}

	m.Logger.Infof(
		executeDoneTmp,
		m.ID(), m.AliasName(), m.ActionQueue.Name(), m.markActionIndex, m.ActionQueue.Size(),
	)
}

// SetResultProcessor .
func (m *LogicProcessor) SetResultProcessor(processor containerInterface.LocalResultProcessor) {
	m.resultProcessor = processor
}

// ProcessorType .
func (m *LogicProcessor) ProcessorType() containerInterface.ProcessorType {
	return containerInterface.LogicProcessorType
}

// DumpInfo .
func (m *LogicProcessor) DumpInfo() containerInterface.ProcessorDump {
	dump := containerInterface.ProcessorDump{
		PID:           m.pid,
		AliasName:     m.aliasName,
		ProcessorType: m.ProcessorType().String(),
	}

	action := m.ActionQueue.Current()
	if action != nil {
		dump.CurrentActionName = action.ActionName()
		if len(action.ReceiveID()) == 0 {
			dump.CurrentActionStatus = "Start"
		} else {
			dump.CurrentActionStatus = "Receive"
		}
	} else {
		dump.CurrentActionName = "/"
		dump.CurrentActionStatus = "/"
	}

	dump.ActionQueue = m.ActionQueue.Dump()

	return dump
}
