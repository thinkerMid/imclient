package companion

import (
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
)

type CompanionPresence struct {
	processor.BaseAction
	JID types.JID
}

func (m *CompanionPresence) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	node := waBinary.Node{
		Tag: "presence",
		Attrs: waBinary.Attrs{
			"to":   m.JID.String(),
			"type": "probe",
		},
	}
	m.SendMessageId, _ = context.SendNode(node)
	return nil
}

func (m *CompanionPresence) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	next()
	return nil
}

func (m *CompanionPresence) Error(context containerInterface.IMessageContext, err error) {

}
