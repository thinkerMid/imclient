package user

import (
	waBinary "ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	containerInterface "ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
)

// QueryPreKeyCount 查询服务器存储的密钥数量
type QueryPreKeyCount struct {
	processor.BaseAction
}

// Start .
func (s *QueryPreKeyCount) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	s.SendMessageId, err = context.SendIQ(message.InfoQuery{
		ID:        context.GenerateRequestID(),
		Namespace: "encrypt",
		Type:      "get",
		To:        types.ServerJID,
		Content: []waBinary.Node{
			{Tag: "count"},
		},
	})

	return
}

// Receive .
func (s *QueryPreKeyCount) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	next()

	return
}

// Error .
func (s *QueryPreKeyCount) Error(_ containerInterface.IMessageContext, _ error) {}
