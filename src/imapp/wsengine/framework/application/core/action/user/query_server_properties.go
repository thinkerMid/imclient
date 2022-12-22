package user

import (
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
)

// QueryServerProperties .
type QueryServerProperties struct {
	processor.BaseAction
}

// Start .
func (m *QueryServerProperties) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	m.SendMessageId, err = context.SendIQ(message.InfoQuery{
		ID:        context.GenerateRequestID(),
		Namespace: "w",
		Type:      "get",
		To:        types.ServerJID,
		Content: []waBinary.Node{{
			Tag: "props",
			Attrs: waBinary.Attrs{
				"protocol": "2",
				"hash":     "", // TODO 这个值第一次是空的，往后是需要用回包返回的值，现在只有第一次注册才会发。所以没做回包解析处理
			},
		}},
	})

	return
}

// Receive .
func (m *QueryServerProperties) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	next()

	return nil
}

// Error .
func (m *QueryServerProperties) Error(context containerInterface.IMessageContext, err error) {
}
