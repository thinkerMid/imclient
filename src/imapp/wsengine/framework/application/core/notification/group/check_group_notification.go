package groupNotification

import (
	"fmt"
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/types"
)

func parseGroupNotification(node *waBinary.Node) (string, error) {
	attrs := node.AttrGetter()
	groupJID := attrs.OptionalJIDOrEmpty("from")

	if groupJID.Server != types.GroupServer {
		return "", fmt.Errorf("not gourp notification")
	}

	return groupJID.User, nil
}
