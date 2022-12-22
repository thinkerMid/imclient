package xmlStream

import (
	"ws/framework/application/constant/message"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
)

// ConnectStream .
type ConnectStream struct {
	processor.BaseAction
}

// Start .
func (m *ConnectStream) Start(_ containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	m.SendMessageId = message.ReconnectEvent

	return
}

// Receive .
func (m *ConnectStream) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	// 连接成功了才会进来的

	next()

	return nil
}

// Error .
func (m *ConnectStream) Error(_ containerInterface.IMessageContext, _ error) {}
