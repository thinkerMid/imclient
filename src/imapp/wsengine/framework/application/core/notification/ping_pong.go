package notification

import (
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
)

// PingPong .
type PingPong struct{}

// Receive .
func (m PingPong) Receive(context containerInterface.IMessageContext) (err error) {
	node := context.Message()
	if node.AttrGetter().String("xmlns") != "urn:xmpp:ping" {
		return
	}

	_, err = context.SendNode(waBinary.Node{
		Tag: "iq",
		Attrs: waBinary.Attrs{
			"to":   types.ServerJID,
			"type": "result",
		},
	})

	return
}
