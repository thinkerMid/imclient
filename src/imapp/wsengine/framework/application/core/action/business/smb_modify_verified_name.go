package business

import (
	waBinary "ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	containerInterface "ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	messageResultType "ws/framework/application/core/result/constant"
	deviceDB "ws/framework/application/data_storage/device/database"
)

// SMBModifyVerifiedName .
type SMBModifyVerifiedName struct {
	processor.BaseAction
	ModifyName string
}

// Start .
func (m *SMBModifyVerifiedName) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	if len(m.ModifyName) > 0 {
		context.ResolveDeviceService().ContextExecute(func(device *deviceDB.Device) {
			device.UpdatePushName(m.ModifyName)
		})
	}

	vname := context.ResolveBusinessService().GenerateBusinessVerifiedName(true)

	m.SendMessageId, err = context.SendIQ(message.InfoQuery{
		ID:        context.GenerateRequestID(),
		Namespace: "w:biz",
		Type:      message.IqSet,
		To:        types.ServerJID,
		Content: []waBinary.Node{{
			Tag: "verified_name",
			Attrs: waBinary.Attrs{
				"v": "2",
			},
			Content: vname,
		}},
	})

	return
}

// Receive .
func (m *SMBModifyVerifiedName) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	next()

	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.SetNickName,
	})
	return nil
}

// Error .
func (m *SMBModifyVerifiedName) Error(context containerInterface.IMessageContext, err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.SetNickName,
		Error:      err,
	})
}
