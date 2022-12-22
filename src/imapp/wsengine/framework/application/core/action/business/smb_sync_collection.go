package business

import (
	waBinary "ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	containerInterface "ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
)

// SMBSyncCollection .
type SMBSyncCollection struct {
	processor.BaseAction
}

// RaiseErrorWhenNodeError 是否抛出node错误，如401，404等
func (m *SMBSyncCollection) RaiseErrorWhenNodeError() bool {
	return false
}

// Start .
func (m *SMBSyncCollection) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	m.SendMessageId, err = context.SendIQ(message.InfoQuery{
		ID:        context.GenerateRequestID(),
		Namespace: "w:sync:app:state",
		Type:      message.IqSet,
		To:        types.ServerJID,
		Content: []waBinary.Node{{
			Tag: "sync",
			Content: []waBinary.Node{
				{
					Tag: "collection",
					Attrs: waBinary.Attrs{
						"name":    "critical_block",
						"version": "2", // TODO maybe autoincre?
					},
					Content: []waBinary.Node{
						{Tag: "patch", Content: []byte{}}, // TODO unknown bytes
					},
				},
			},
		}},
	})

	return
}

// Receive .
func (m *SMBSyncCollection) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	next()

	return nil
}

// Error .
func (m *SMBSyncCollection) Error(context containerInterface.IMessageContext, err error) {
	// TODO error code 409 发送需要置空操作
}
