package user

import (
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
)

// DeleteAllData .
type DeleteAllData struct {
	processor.BaseAction
}

// Start .
func (m *DeleteAllData) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	m.SendMessageId, err = context.SendIQ(message.InfoQuery{
		ID:        context.GenerateRequestID(),
		Namespace: "w:sync:app:state",
		Type:      "set",
		To:        types.ServerJID,
		Content: []waBinary.Node{{
			Tag: "delete_all_data",
		}},
	})
	return
}

// Receive .
func (m *DeleteAllData) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	next()

	return nil
}

// Error .
func (m *DeleteAllData) Error(_ containerInterface.IMessageContext, _ error) {
}
