package privateChatNotification

import (
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	contactDB "ws/framework/application/data_storage/contact/database"
)

// HandleTrustedContact .
type HandleTrustedContact struct{}

// Receive .
func (r HandleTrustedContact) Receive(context containerInterface.IMessageContext) (err error) {
	node := context.Message()
	if node.Tag != "receipt" {
		return
	}

	from, ok := node.Attrs["from"].(types.JID)
	if !ok || from.Server != types.DefaultUserServer {
		return
	}

	jid := context.ResolveJID()
	if jid.User == from.User {
		return
	}

	context.ResolveContactService().ContextExecute(from.User, func(contact *contactDB.Contact) {
		contact.UpdateTrustedContact(true)
	})

	return
}
