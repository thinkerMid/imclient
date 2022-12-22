package processor

import (
	"go.uber.org/zap"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor/interval_timer"
	"ws/framework/application/core/processor/trigger"
)

// 处理器优先级级别类型 越小的级别越低
//
//	用于channel执行shutdown操作时判断的条件。目前分两种级别，当channel还可以进行重连时会销毁后台级别的处理器
const (
	// PriorityBackground 后台级别
	PriorityBackground uint8 = iota
	// PriorityForeground 前台级别
	PriorityForeground
)

// ----------------------------------------------------------------------------

// BaseAction .
type BaseAction struct {
	Name          string
	Index         uint8
	SendMessageId string
	Query         containerInterface.IAction
}

// SetActionIndex .
func (a *BaseAction) SetActionIndex(v uint8) {
	a.Index = v
}

// ActionIndex .
func (a *BaseAction) ActionIndex() uint8 {
	return a.Index
}

// SetActionName .
func (a *BaseAction) SetActionName(v string) {
	a.Name = v
}

// ActionName .
func (a *BaseAction) ActionName() string {
	return a.Name
}

// ReceiveID .
func (a *BaseAction) ReceiveID() string {
	if a.Query != nil {
		return a.Query.ReceiveID()
	}

	return a.SendMessageId
}

// RaiseErrorWhenNodeError 是否抛出node错误，如401，404等
func (a *BaseAction) RaiseErrorWhenNodeError() bool {
	return true
}

// AddFuture .
func (a *BaseAction) AddFuture(ctx containerInterface.IMessageContext, future containerInterface.IFuture) {
	if len(a.SendMessageId) == 0 {
		return
	}

	ctx.AddFutureProcessor(ctx.ProcessPID(), NewFutureProcessor(a.SendMessageId, future))
}

// Start .
func (a *BaseAction) Start(context containerInterface.IMessageContext, fn containerInterface.NextActionFn) error {
	panic("implement me")
}

// Receive .
func (a *BaseAction) Receive(context containerInterface.IMessageContext, fn containerInterface.NextActionFn) error {
	panic("implement me")
}

// Error .
func (a *BaseAction) Error(context containerInterface.IMessageContext, err error) {
	panic("implement me")
}

// ----------------------------------------------------------------------------

// BaseProcessor .
type BaseProcessor struct {
	pid               uint32
	aliasName         string
	processorPriority uint8
	loop              bool

	Trigger     containerInterface.ITrigger
	ActionQueue containerInterface.IActionQueue
	Monitors    []containerInterface.IMonitor
	Timer       *intervalTimer.Timer
	Logger      *zap.SugaredLogger
}

// Init .
func (m *BaseProcessor) Init(logger *zap.SugaredLogger) {
	m.Logger = logger

	if len(m.aliasName) == 0 {
		m.aliasName = "anonymous"
	}

	if m.Trigger == nil {
		m.Trigger = trigger.NewAutomaticTrigger()
	}

	if m.Monitors == nil {
		m.Monitors = []containerInterface.IMonitor{}
	}

	if m.Timer == nil {
		m.Timer = intervalTimer.New(0)
	}
}

// Start .
func (m *BaseProcessor) Start(context containerInterface.IProcessControl) {}

// ProcessMessage .
func (m *BaseProcessor) ProcessMessage(context containerInterface.IProcessControl) {}

// OnChannelError .
func (m *BaseProcessor) OnChannelError(context containerInterface.IProcessControl, err error) {}

// OnDestroy .
func (m *BaseProcessor) OnDestroy(context containerInterface.IProcessControl) {}

// ID .
func (m *BaseProcessor) ID() uint32 {
	return m.pid
}

// AliasName .
func (m *BaseProcessor) AliasName() string {
	return m.aliasName
}

// SetID .
func (m *BaseProcessor) SetID(i uint32) {
	m.pid = i
}

// SetAliasName .
func (m *BaseProcessor) SetAliasName(v string) {
	m.aliasName = v
}

// SetTrigger .
func (m *BaseProcessor) SetTrigger(t containerInterface.ITrigger) {
	m.Trigger = t
}

// SetActionQueue .
func (m *BaseProcessor) SetActionQueue(queue containerInterface.IActionQueue) {
	m.ActionQueue = queue
}

// SetMonitor .
func (m *BaseProcessor) SetMonitor(monitor containerInterface.IMonitor) {
	m.Monitors = append(m.Monitors, monitor)
}

// SetInterval uint8是限制不能小于0
func (m *BaseProcessor) SetInterval(second uint32) {
	m.Timer = intervalTimer.New(second)
}

// SetIntervalLoop .
func (m *BaseProcessor) SetIntervalLoop(loop bool) {
	m.loop = loop
}

// SetPriority .
func (m *BaseProcessor) SetPriority(v uint8) {
	m.processorPriority = v
}

// Priority .
func (m *BaseProcessor) Priority() uint8 {
	return m.processorPriority
}

// NeedAutoStart .
func (m *BaseProcessor) NeedAutoStart() bool {
	return m.Trigger.DefaultTrigger()
}

// SetResultProcessor .
func (m *BaseProcessor) SetResultProcessor(processor containerInterface.LocalResultProcessor) {}

// ProcessorType .
func (m *BaseProcessor) ProcessorType() containerInterface.ProcessorType {
	return containerInterface.UnknownProcessorType
}

// DumpInfo .
func (m *BaseProcessor) DumpInfo() containerInterface.ProcessorDump {
	return containerInterface.ProcessorDump{}
}

// ----------------------------------------------------------------------------

// SetConfigFn 设置函数
type SetConfigFn func(p containerInterface.IProcessor)

// TriggerTag 触发执行的Tag
func TriggerTag(key string) SetConfigFn {
	return func(p containerInterface.IProcessor) {
		p.SetTrigger(trigger.NewDefaultTrigger(key))
	}
}

// Priority 优先级
func Priority(v uint8) SetConfigFn {
	return func(p containerInterface.IProcessor) {
		p.SetPriority(v)
	}
}

// AttachMonitor 附加监控器
func AttachMonitor(monitor containerInterface.IMonitor) SetConfigFn {
	return func(p containerInterface.IProcessor) {
		p.SetMonitor(monitor)
	}
}

// Interval 执行间隔 只支持整秒的精度
func Interval(second uint32) SetConfigFn {
	return func(p containerInterface.IProcessor) {
		p.SetInterval(second)
	}
}

// IntervalLoop 循环间隔
func IntervalLoop(loop bool) SetConfigFn {
	return func(p containerInterface.IProcessor) {
		p.SetIntervalLoop(loop)
	}
}

// AliasName .
func AliasName(v string) SetConfigFn {
	return func(p containerInterface.IProcessor) {
		p.SetAliasName(v)
	}
}
