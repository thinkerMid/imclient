package xmlStream

import (
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
)

// SendStreamEnd .
type SendStreamEnd struct {
	processor.BaseAction
}

// Start .
func (m *SendStreamEnd) Start(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	m.SendMessageId, err = context.SendNode(waBinary.Node{Tag: message.StreamEnd})

	// 告知控制器是主动关闭的 不要做重连后的初始化逻辑
	context.ResolveIMControl().ActiveDisconnectXMPP()
	// 防止出现发送出去了包 但是连接没关闭的情况
	context.ResolveConnection().Close()

	next()
	return
}

// Receive .
func (m *SendStreamEnd) Receive(_ containerInterface.IMessageContext, _ containerInterface.NextActionFn) error {
	return nil
}

// Error .
func (m *SendStreamEnd) Error(_ containerInterface.IMessageContext, _ error) {}
