package other

import (
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
)

// RemoveCompanionDevice .
type RemoveCompanionDevice struct {
	processor.BaseAction
}

// Start .
func (m *RemoveCompanionDevice) Start(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	iq := message.InfoQuery{
		ID:        context.GenerateRequestID(),
		Namespace: "md",
		Type:      message.IqSet,
		To:        types.ServerJID,
		Content: []waBinary.Node{{
			Tag: "remove-companion-device",
			Attrs: waBinary.Attrs{
				"all":    "true",
				"reason": "md_opt_out",
			},
		}},
	}

	m.SendMessageId, err = context.SendIQ(iq)

	return
}

// Receive .
func (m *RemoveCompanionDevice) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	next()

	return nil
}

// Error .
func (m *RemoveCompanionDevice) Error(_ containerInterface.IMessageContext, _ error) {}
