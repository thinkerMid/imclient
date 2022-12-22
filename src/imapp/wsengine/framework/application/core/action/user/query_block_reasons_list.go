package user

import (
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
)

// QueryBlockReasonsList .
type QueryBlockReasonsList struct {
	processor.BaseAction
}

// Start .
func (m *QueryBlockReasonsList) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	m.SendMessageId, err = context.SendIQ(message.InfoQuery{
		ID:        context.GenerateRequestID(),
		Namespace: "w:biz",
		Type:      "get",
		To:        types.ServerJID,
		Content: []waBinary.Node{{
			Tag: "mobile_config",
			Attrs: waBinary.Attrs{
				"name": "biz_block_reasons",
				"v":    "2",
			},
		}},
	})

	return
}

// Receive .
func (m *QueryBlockReasonsList) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	next()

	return nil
}

// Error .
func (m *QueryBlockReasonsList) Error(context containerInterface.IMessageContext, err error) {
}
