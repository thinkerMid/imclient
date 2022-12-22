package xmlStreamNotification

import (
	"strconv"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/result/constant"
)

// StreamErrorStatus .
type StreamErrorStatus struct{}

// Receive .
func (m StreamErrorStatus) Receive(context containerInterface.IMessageContext) (err error) {
	code := context.Message().AttrGetter().String("code")

	if numCode, err := strconv.Atoi(code); err == nil && numCode > 0 {
		context.AppendResult(containerInterface.MessageResult{
			ResultType: messageResultType.StreamErrorCode,
			Content:    code,
		})
	}

	return
}
