package identityNotification

import (
	"ws/framework/application/constant"
	"ws/framework/application/container/abstract_interface"
)

// DeleteSession .
type DeleteSession struct{}

// Receive .
func (r DeleteSession) Receive(context containerInterface.IMessageContext) error {
	node := context.Message()

	ag := node.AttrGetter()
	if ag.String("type") != "encrypt" {
		return nil
	}

	_, ok := node.GetOptionalChildByTag("identity")
	if !ok {
		return nil
	}

	from := node.AttrGetter().JID("from")

	context.ResolveDeviceListService().DeleteDevice(from.User, from.Device)
	context.ResolveSenderKeyService().DeleteDevice(from)

	return constant.AbortedError
}
