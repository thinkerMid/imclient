package notification

import (
	"ws/framework/application/constant/binary"
	"ws/framework/application/container/abstract_interface"
)

// Ack .
type Ack struct{}

// Receive .
func (r Ack) Receive(context containerInterface.IMessageContext) (err error) {
	node := context.Message()

	attrs := waBinary.Attrs{
		"class": node.Tag,
		"id":    node.Attrs["id"],
		"to":    node.Attrs["from"],
	}

	if participant, ok := node.Attrs["participant"]; ok {
		attrs["participant"] = participant
	}
	if recipient, ok := node.Attrs["recipient"]; ok {
		attrs["recipient"] = recipient
	}
	if receiptType, ok := node.Attrs["type"]; node.Tag != "message" && ok {
		attrs["type"] = receiptType
	}

	_, err = context.SendNode(waBinary.Node{
		Tag:   "ack",
		Attrs: attrs,
	})

	return
}
