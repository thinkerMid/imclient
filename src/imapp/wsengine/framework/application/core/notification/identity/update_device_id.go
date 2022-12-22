package identityNotification

import (
	"ws/framework/application/constant"
	"ws/framework/application/container/abstract_interface"
)

// UpdateDeviceID 处理对方设备变更
type UpdateDeviceID struct{}

// Receive .
func (r UpdateDeviceID) Receive(context containerInterface.IMessageContext) (err error) {
	node := context.Message()

	ag := node.AttrGetter()
	if ag.String("type") != "devices" {
		return
	}

	var addTag bool

	child, ok := node.GetOptionalChildByTag("add")
	if ok {
		addTag = true
	} else {
		child, ok = node.GetOptionalChildByTag("remove")
		if !ok {
			return
		}
	}

	deviceNode, _ := child.GetOptionalChildByTag("device")
	deviceJID := deviceNode.AttrGetter().JID("jid")

	if addTag {
		// 设备列表新增
		context.ResolveDeviceListService().AddDevice(deviceJID.User, deviceJID.Device)
		// 群设备新增
		context.ResolveSenderKeyService().SearchSenderInGroupAndCreate(deviceJID)
	} else {
		// 设备列表移除
		context.ResolveDeviceListService().DeleteDevice(deviceJID.User, deviceJID.Device)
		// 群设备移除
		context.ResolveSenderKeyService().DeleteDevice(deviceJID)
	}

	return constant.AbortedError
}
