package containerInterface

import (
	"go.uber.org/zap"
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
)

// ProcessorType 处理器类型
type ProcessorType uint8

const (
	// UnknownProcessorType .
	UnknownProcessorType ProcessorType = iota
	// LogicProcessorType 逻辑处理器
	LogicProcessorType
	// TimerProcessorType 定时处理器
	TimerProcessorType
	// NotificationProcessorType 通知处理器
	NotificationProcessorType
	// FutureProcessorType 结果接收处理器
	FutureProcessorType
)

func (t ProcessorType) String() string {
	switch t {
	case LogicProcessorType:
		return "LogicProcessor"
	case TimerProcessorType:
		return "TimerProcessor"
	case NotificationProcessorType:
		return "NotificationProcessor"
	case FutureProcessorType:
		return "FutureProcessor"
	default:
		return "UnknownProcessor"
	}
}

// ----------------------------------------------------------------------------

// IProcessControl .
type IProcessControl interface {
	Aborted()
	Reject()
	AppIocContainer() IAppIocContainer
	MessageContext() IMessageContext
	MessageChannel() IMessageChannel
}

// IMessageContext .
type IMessageContext interface {
	Message() *waBinary.Node
	AppendResult(result MessageResult)
	VisitResult(index int) MessageResult
	MessageResult() []MessageResult
	Reset()
	CleanResult()
	ProcessPID() uint32
	IAppIocContainer
	IMessageChannel
}

// MessageResult .
type MessageResult struct {
	ResultType uint8
	Content    string
	IContent   interface{}
	Error      error
}

// ProcessorDump .
type ProcessorDump struct {
	PID                 uint32
	AliasName           string
	ProcessorType       string
	CurrentActionName   string
	CurrentActionStatus string
	ActionQueue         string
}

// ----------------------------------------------------------------------------

// IScene .
type IScene interface {
	Build() IProcessor
}

// ----------------------------------------------------------------------------

// MakeActionFn .
type MakeActionFn func() []IAction

// MakeNotificationFn .
type MakeNotificationFn func() []INotification

// NextActionFn .
type NextActionFn func()

// ITrigger 触发器
type ITrigger interface {
	DefaultTrigger() bool
	WaitActive() bool
	Disable()
	Reset()
	ContentedCondition(*waBinary.Node) bool
}

// IMonitor 监控器
type IMonitor interface {
	OnStart(IAppIocContainer)
	OnActionStartBefore(action interface{}, context IMessageContext)
	OnActionStartAfter(action interface{}, context IMessageContext)
	OnActionStartFail(action interface{}, context IMessageContext)
	OnActionStartSuccess(action interface{}, context IMessageContext)
	ActionExecuteSuccess(action interface{}, context IMessageContext)
	ActionExecuteFailure(action interface{}, context IMessageContext)
	OnExit(IAppIocContainer)
}

// IAction 基础类型的action
//
//	具有发送消息和处理消息的特性
//	start后可以调用next，会丢失接收消息回复的机会
//	receive后并处理完毕一定要调用next，下一个action才会得到执行机会
type IAction interface {
	SetActionIndex(uint8)
	ActionIndex() uint8
	SetActionName(string)
	ActionName() string
	ReceiveID() string
	RaiseErrorWhenNodeError() bool
	Start(IMessageContext, NextActionFn) error
	Receive(IMessageContext, NextActionFn) error
	Error(IMessageContext, error)
}

// INotification 通知类型的action
type INotification interface {
	Receive(IMessageContext) error
}

// ITimer 定时器类型的action
type ITimer interface {
	ReceiveID() string
	Start(IMessageContext) error
	Receive(IMessageContext) error
	Error(IMessageContext, error)
}

// IFuture 处理特定内容的一个载体
type IFuture interface {
	Receive(IMessageContext) error
	Error(IMessageContext, error)
}

// IProcessor 消息处理器
type IProcessor interface {
	ID() uint32
	AliasName() string
	Priority() uint8
	NeedAutoStart() bool
	ProcessorType() ProcessorType

	SetID(uint32)
	SetAliasName(string)
	SetPriority(uint8)
	SetInterval(uint32)
	SetIntervalLoop(bool)
	SetMonitor(IMonitor)
	SetTrigger(ITrigger)
	SetResultProcessor(LocalResultProcessor)

	Init(*zap.SugaredLogger)
	Start(IProcessControl)
	ProcessMessage(IProcessControl)
	OnChannelError(IProcessControl, error)
	OnDestroy(IProcessControl)

	DumpInfo() ProcessorDump
}

// IActionQueue .
type IActionQueue interface {
	Name() string
	Size() int
	Empty() bool
	OnError(_ error)
	Current() IAction
	Pop() IAction
	Dump() string
}

// INotificationQueue .
type INotificationQueue interface {
	Reload()
	Current() INotification
	Pop() INotification
	Dump() string
}

// LocalResultProcessor 本地结果处理器
//
//	和MessageProcessor是绑定的 只能处理已绑定MessageProcessor输出的结果
type LocalResultProcessor interface {
	ProcessResult(iocProvider IAppIocContainer, result *MessageResult)
	OnDestroy()
}

// GlobalResultProcessor 全局结果处理器
//
//	处于全局消费结果队列中,可以处理其他消息处理器未消费的结果
type GlobalResultProcessor interface {
	ListenTags() []uint8 // 监听的数据类型
	Reside() bool        // 是否常驻
	ProcessResult(iocProvider IAppIocContainer, result *MessageResult)
	OnDestroy()
}

// IMessageChannel .
type IMessageChannel interface {
	SendNode(waBinary.Node) (string, error)
	SendIQ(message.InfoQuery) (string, error)
	GenerateRequestID() string
	GenerateSID() string
	AddMessageProcessor(IProcessor) uint32
	AddFutureProcessor(uint32, IProcessor)
	AddProcessorAndAttach(IProcessor, LocalResultProcessor) uint32
	AddGlobalResultProcessor(GlobalResultProcessor)
	RemoveProcessor(uint32)
}

// IIMControl .
type IIMControl interface {
	ActiveDisconnectXMPP()
}
