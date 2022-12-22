package contactNotification

import (
	"ws/framework/application/constant"
	containerInterface "ws/framework/application/container/abstract_interface"
)

// BeAddedOrDeleted .
type BeAddedOrDeleted struct {
	JID string
}

// Receive .
func (s *BeAddedOrDeleted) Receive(context containerInterface.IMessageContext) (err error) {
	// <notification from="79910515576@s.whatsapp.net" id="2798928688" t="1662428323" type="contacts"><update jid="84564844255@s.whatsapp.net"/></notification>
	node := context.Message()
	attr := node.AttrGetter()

	if attr.String("type") != "contacts" {
		return
	}

	s.JID = attr.JID("from").User

	return constant.AbortedError
}
