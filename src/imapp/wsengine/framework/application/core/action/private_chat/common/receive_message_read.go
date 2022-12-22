package privateChatCommon

import (
	"time"
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	messageResultType "ws/framework/application/core/result/constant"
)

// ReceiveMessageMarkRead .
type ReceiveMessageMarkRead struct {
	UserID     string
	MessageIDs []string
	processor.BaseAction
}

// Start .
func (r *ReceiveMessageMarkRead) Start(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	node := waBinary.Node{
		Tag: "receipt",
		Attrs: waBinary.Attrs{
			"id":   r.MessageIDs[0],
			"type": "read",
			"to":   types.NewJID(r.UserID, types.DefaultUserServer),
			"t":    time.Now().Unix(),
		},
	}

	if len(r.MessageIDs) > 1 {
		var messageIDNodes []waBinary.Node

		size := len(r.MessageIDs)
		for i := 1; i < size; i++ {
			messageIDNodes = append(messageIDNodes, waBinary.Node{
				Tag: "item",
				Attrs: waBinary.Attrs{
					"id": r.MessageIDs[i],
				},
			})
		}

		node.Content = []waBinary.Node{{Tag: "list", Content: messageIDNodes}}
	}

	_, err = context.SendNode(node)

	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.ReceiveMessageMarkRead,
	})

	next()
	return err
}

// Receive .
func (r *ReceiveMessageMarkRead) Receive(context containerInterface.IMessageContext, fn containerInterface.NextActionFn) error {
	return nil
}

// Error .
func (r *ReceiveMessageMarkRead) Error(context containerInterface.IMessageContext, err error) {

}
