package userSettings

import (
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	"ws/framework/application/core/result/constant"
	deviceDB "ws/framework/application/data_storage/device/database"
)

// ModifyNickName .
//
//	这个修改其实是本地数据库备注，如果要查询的时候去查本地数据库
type ModifyNickName struct {
	processor.BaseAction
	DstName string
}

// Start .
func (m *ModifyNickName) Start(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	node := waBinary.Node{
		Tag: "presence",
		Attrs: waBinary.Attrs{
			"type": string(types.PresenceAvailable),
		},
	}

	if len(m.DstName) > 0 {
		node.Attrs["name"] = m.DstName
	}

	m.SendMessageId, err = context.SendNode(node)

	if err == nil {
		context.ResolveDeviceService().ContextExecute(func(device *deviceDB.Device) {
			device.UpdatePushName(m.DstName)
		})

		context.AppendResult(containerInterface.MessageResult{
			ResultType: messageResultType.SetNickName,
		})
	}

	next()

	return
}

// Receive .
func (m *ModifyNickName) Receive(_ containerInterface.IMessageContext, _ containerInterface.NextActionFn) error {
	return nil
}

// Error .
func (m *ModifyNickName) Error(context containerInterface.IMessageContext, err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.SetNickName,
		Error:      err,
	})
}
