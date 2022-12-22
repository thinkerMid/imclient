package accountNotification

import (
	"ws/framework/application/container/abstract_interface"
)

// LoginData .
type LoginData struct{}

// Receive .
func (m LoginData) Receive(context containerInterface.IMessageContext) (err error) {
	attrs := context.Message().AttrGetter()
	val := attrs.String("location")

	accountLoginData := context.ResolveMemoryCache().AccountLoginData()
	accountLoginData.LastKnowDataCenter = val

	record := context.ResolveABKeyService().Context()
	if record == nil {
		return
	}

	accountLoginData.ABKey2 = record.Content
	return
}
