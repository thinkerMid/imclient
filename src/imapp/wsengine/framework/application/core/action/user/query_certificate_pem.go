package user

import (
	"strconv"
	"time"
	waBinary "ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	containerInterface "ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
)

// QueryCertificatePEM .
type QueryCertificatePEM struct {
	processor.BaseAction
}

// Start .
func (m *QueryCertificatePEM) Start(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	m.SendMessageId, err = context.SendNode(waBinary.Node{
		Tag: "iq",
		Attrs: waBinary.Attrs{
			"id":      context.GenerateRequestID(),
			"smax_id": "99",
			"to":      types.ServerJID,
			"type":    message.IqGet,
			"xmlns":   "avatars",
		},
		Content: []waBinary.Node{
			{Tag: "password_pem"},
			{Tag: "timestamp", Content: strconv.Itoa(int(time.Now().Unix()))},
			{Tag: "payload_enc_certificates"},
		},
	})
	next()

	return
}

// Receive .
func (m *QueryCertificatePEM) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {

	return nil
}

// Error .
func (m *QueryCertificatePEM) Error(_ containerInterface.IMessageContext, _ error) {}
