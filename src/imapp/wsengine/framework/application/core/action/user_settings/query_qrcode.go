package userSettings

import (
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	"ws/framework/application/core/result/constant"
)

var qrCodeUrlPrefix = "https://wa.me/qr/"

// QueryQrCode .
type QueryQrCode struct {
	processor.BaseAction
}

// Start .
func (m *QueryQrCode) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	m.SendMessageId, err = context.SendIQ(message.InfoQuery{
		ID:        context.GenerateRequestID(),
		Namespace: "w:qr",
		Type:      "set",
		To:        types.ServerJID,
		Content: []waBinary.Node{{
			Tag: "qr",
			Attrs: waBinary.Attrs{
				"action": "get",
				"type":   "contact",
			},
		}},
	})

	return
}

// Receive .
func (m *QueryQrCode) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	node := context.Message().GetChildByTag("qr")

	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.GetQrCode,
		Content:    qrCodeUrlPrefix + node.AttrGetter().String("code"),
	})

	next()
	return nil
}

// Error .
func (m *QueryQrCode) Error(context containerInterface.IMessageContext, err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.GetQrCode,
		Error:      err,
	})
}
