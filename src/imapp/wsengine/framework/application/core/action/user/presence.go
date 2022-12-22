package user

import (
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/types"
	containerInterface "ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
)

// Presence 告知服务器上下线的东西
type Presence struct {
	processor.BaseAction
	PresenceState types.Presence
}

// Start .
func (m *Presence) Start(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	node := waBinary.Node{
		Tag: "presence",
		Attrs: waBinary.Attrs{
			"type": string(m.PresenceState),
		},
	}

	if m.PresenceState == types.PresenceAvailable {
		device := context.ResolveDeviceService().Context()
		if len(device.PushName) > 0 {
			node.Attrs["name"] = device.PushName
		}
	}

	m.SendMessageId, err = context.SendNode(node)

	next()

	return
}

// Receive .
func (m *Presence) Receive(_ containerInterface.IMessageContext, _ containerInterface.NextActionFn) error {
	return nil
}

// Error .
func (m *Presence) Error(_ containerInterface.IMessageContext, _ error) {}
