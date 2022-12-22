package eventNotification

import (
	"ws/framework/application/container/abstract_interface"
)

// ResetEventLogState .
type ResetEventLogState struct{}

// Receive .
func (c ResetEventLogState) Receive(context containerInterface.IMessageContext) (err error) {
	wec := context.ResolveChannel0EventCache()
	wec2 := context.ResolveChannel2EventCache()

	wec.ResetAddLogState()
	wec2.ResetAddLogState()

	return
}
