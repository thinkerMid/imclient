package xmlStreamNotification

import (
	"ws/framework/application/container/abstract_interface"
)

// RoutingInfo .
type RoutingInfo struct{}

// Receive .
func (p RoutingInfo) Receive(context containerInterface.IMessageContext) (err error) {
	node, ok := context.Message().GetOptionalChildByTag("edge_routing")
	if !ok {
		return
	}

	node, ok = node.GetOptionalChildByTag("routing_info")
	if !ok {
		return
	}

	content := node.Content.([]byte)
	context.ResolveHandshakeHandler().SetEdgeRouting(content)

	contentList := node.ContentString()
	if len(contentList) == 0 {
		return
	}

	context.ResolveMemoryCache().AccountLoginData().RoutingInfo = contentList[0]

	return
}
