package user

import (
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
)

// QueryABProperties .
type QueryABProperties struct {
	processor.BaseAction
}

// Start .
func (m *QueryABProperties) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	m.SendMessageId, err = context.SendIQ(message.InfoQuery{
		ID:        context.GenerateRequestID(),
		Namespace: "abt",
		Type:      "get",
		To:        types.ServerJID,
		Content: []waBinary.Node{{
			Tag: "props",
			Attrs: waBinary.Attrs{
				"protocol": "1",
				//"hash":     "OFFLINE",
			},
		}},
	})

	return
}

// Receive .
func (m *QueryABProperties) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	next()

	props, ok := context.Message().GetOptionalChildByTag("props")
	if !ok {
		return nil
	}

	attrs := props.AttrGetter()
	abKey := attrs.String("ab_key")

	accountLoginData := context.ResolveMemoryCache().AccountLoginData()
	accountLoginData.ABKey2 = abKey

	_ = context.ResolveABKeyService().Create(abKey)

	return nil
}

// Error .
func (m *QueryABProperties) Error(_ containerInterface.IMessageContext, _ error) {
}
