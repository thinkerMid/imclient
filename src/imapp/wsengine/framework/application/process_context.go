package application

import (
	"ws/framework/application/constant/binary"
	containerInterface "ws/framework/application/container/abstract_interface"
)

// NewProcessContext .
func NewProcessContext(ioc containerInterface.IAppIocContainer, node *waBinary.Node, signalMessage bool) ProcessContext {
	return ProcessContext{
		IAppIocContainer: ioc,
		IMessageChannel:  ioc.ResolveMessageChannel(),
		signalMessage:    signalMessage,
		reAppendQueue:    true, // 默认true
		node:             node,
		results:          []containerInterface.MessageResult{},
	}
}

// ProcessContext .
type ProcessContext struct {
	containerInterface.IAppIocContainer
	containerInterface.IMessageChannel

	pid            uint32 // 处理当前消息的处理器ID
	messageAborted bool   // 中断消息传播
	reAppendQueue  bool   // 重新加入处理队列 默认true
	signalMessage  bool   // 是否由信号产生的消息

	next bool

	node    *waBinary.Node                     // xmpp消息内容
	results []containerInterface.MessageResult // 产生的消息结果
}

// Message .
func (p *ProcessContext) Message() *waBinary.Node {
	return p.node
}

// VisitResult .
func (p *ProcessContext) VisitResult(index int) containerInterface.MessageResult {
	return p.results[index]
}

// MessageResult .
func (p *ProcessContext) MessageResult() []containerInterface.MessageResult {
	return p.results
}

// AppendResult .
func (p *ProcessContext) AppendResult(result containerInterface.MessageResult) {
	p.results = append(p.results, result)
}

// Reset .
func (p *ProcessContext) Reset() {
	p.results = p.results[0:0]
	p.pid = 0
	p.messageAborted = false
	p.reAppendQueue = true
}

// CleanResult .
func (p *ProcessContext) CleanResult() {
	p.results = make([]containerInterface.MessageResult, 0)
}

// SetProcessPID .
func (p *ProcessContext) SetProcessPID(pid uint32) {
	p.pid = pid
}

// ProcessPID .
func (p *ProcessContext) ProcessPID() uint32 {
	return p.pid
}

// Aborted .
func (p *ProcessContext) Aborted() {
	p.messageAborted = true
}

// MessageAborted .
func (p *ProcessContext) MessageAborted() bool {
	return p.messageAborted
}

// Reject .
func (p *ProcessContext) Reject() {
	p.reAppendQueue = false
}

// AddToQueue .
func (p *ProcessContext) AddToQueue() bool {
	return p.reAppendQueue
}

// SignalMessage .
func (p *ProcessContext) SignalMessage() bool {
	return p.signalMessage
}

// AppIocContainer .
func (p *ProcessContext) AppIocContainer() containerInterface.IAppIocContainer {
	return p.IAppIocContainer
}

// MessageChannel .
func (p *ProcessContext) MessageChannel() containerInterface.IMessageChannel {
	return p.IMessageChannel
}

// MessageContext .
func (p *ProcessContext) MessageContext() containerInterface.IMessageContext {
	return p
}
