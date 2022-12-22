package common

import (
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
)

// QueryMultiDevicesIdentityBatch 批量查询对方的设备信息
type QueryMultiDevicesIdentityBatch struct {
	processor.BaseAction
	UserIDs []string
}

// Start .
func (c *QueryMultiDevicesIdentityBatch) Start(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	jid := types.NewJID("", types.DefaultUserServer)
	deviceNode := make([]waBinary.Node, 0)

	for i := range c.UserIDs {
		jid.User = c.UserIDs[i]

		/**
		85252448199@s.whatsapp.net
		85252448199.0:2@s.whatsapp.net
		85252448199.0:1@s.whatsapp.net
		*/
		idList := context.ResolveDeviceListService().FindUnInitSessionDeviceIDList(jid.User)

		for _, id := range idList {
			jid.Device = id
			jid.AD = id > 0

			deviceNode = append(deviceNode, waBinary.Node{
				Tag:   "user",
				Attrs: waBinary.Attrs{"jid": jid.String()},
			})
		}
	}

	if len(deviceNode) == 0 {
		next()
		return
	}

	c.SendMessageId, err = context.SendIQ(
		message.InfoQuery{
			ID:        context.GenerateRequestID(),
			Namespace: "encrypt",
			Type:      message.IqGet,
			To:        types.ServerJID,
			Content: []waBinary.Node{{
				Tag:     "key",
				Content: deviceNode,
			}},
		},
	)

	return
}

// Receive .
func (c *QueryMultiDevicesIdentityBatch) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	q := QueryMultiDevicesIdentity{}
	err = q.Receive(context, next)
	return
}

func (c *QueryMultiDevicesIdentityBatch) Error(context containerInterface.IMessageContext, err error) {
	q := QueryMultiDevicesIdentity{}
	q.Error(context, err)
}
