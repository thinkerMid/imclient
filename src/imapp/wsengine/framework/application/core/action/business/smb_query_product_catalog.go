package business

import (
	waBinary "ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	containerInterface "ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
)

// SMBQueryProductCatalog .
type SMBQueryProductCatalog struct {
	processor.BaseAction

	ConsumerVisibleOnly bool
}

// RaiseErrorWhenNodeError 是否抛出node错误，如401，404等
func (m *SMBQueryProductCatalog) RaiseErrorWhenNodeError() bool {
	return false
}

// Start .
func (m *SMBQueryProductCatalog) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	attrs := waBinary.Attrs{
		"jid": context.ResolveJID(),
	}

	if m.ConsumerVisibleOnly {
		attrs["consumer_visible_only"] = "true"
	}

	m.SendMessageId, err = context.SendIQ(message.InfoQuery{
		ID:        context.GenerateRequestID(),
		Namespace: "w:biz:catalog",
		Type:      message.IqGet,
		To:        types.ServerJID,
		Content: []waBinary.Node{{
			Tag:   "product_catalog",
			Attrs: attrs,
			Content: []waBinary.Node{
				{Tag: "limit", Content: "10"},
				{Tag: "width", Content: "144"},
				{Tag: "height", Content: "144"},
			},
		}},
	})

	return
}

// Receive .
func (m *SMBQueryProductCatalog) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	next()

	return nil
}

// Error .
func (m *SMBQueryProductCatalog) Error(context containerInterface.IMessageContext, err error) {
}
