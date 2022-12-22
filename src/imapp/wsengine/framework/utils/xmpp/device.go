package xmpp

import (
	"ws/framework/application/constant/binary"
)

// ParseDeviceIDList .
func ParseDeviceIDList(node *waBinary.Node) []uint8 {
	deviceIDList := make([]uint8, 0)

	deviceNode, ok := node.GetOptionalChildByTag("device-list")
	if !ok {
		return deviceIDList
	}

	deviceList := deviceNode.GetChildren()

	for i := range deviceList {
		getter := deviceList[i].AttrGetter()

		deviceIDList = append(deviceIDList, uint8(getter.Uint64("id")))
	}

	return deviceIDList
}
