package xmlStreamNotification

import (
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/result/constant"
)

// ReceiveStreamEnd .
type ReceiveStreamEnd struct{}

// Receive .
func (m ReceiveStreamEnd) Receive(context containerInterface.IMessageContext) (err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.StreamEnd,
	})

	return
}
