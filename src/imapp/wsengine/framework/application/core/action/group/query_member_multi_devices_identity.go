package group

import (
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/action/common"
	"ws/framework/application/core/processor"
)

// QueryMemberMultiDevicesIdentity 批量查询群成员的设备信息
type QueryMemberMultiDevicesIdentity struct {
	processor.BaseAction
	GroupID string
}

// Start .
func (c *QueryMemberMultiDevicesIdentity) Start(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	jid := types.NewJID("", types.DefaultUserServer)
	deviceNode := make([]waBinary.Node, 0)

	// 找未发送过消息的群成员设备
	senders, _ := context.ResolveSenderKeyService().FindUnSentMessageDeviceByGroupID(c.GroupID)

	for i := range senders {
		jid.User = senders[i].TheirJID
		jid.Device = uint8(senders[i].DeviceID)
		jid.AD = jid.Device > 0

		/**
		85252448199@s.whatsapp.net
		85252448199.0:2@s.whatsapp.net
		85252448199.0:1@s.whatsapp.net
		*/

		// 看看有设备会话没有
		if context.ResolveDeviceListService().ContainsSession(jid.SignalAddress()) {
			continue
		}

		deviceNode = append(deviceNode, waBinary.Node{
			Tag:   "user",
			Attrs: waBinary.Attrs{"jid": jid.String()},
		})
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
func (c *QueryMemberMultiDevicesIdentity) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	q := common.QueryMultiDevicesIdentity{}
	err = q.Receive(context, next)
	return
}

func (c *QueryMemberMultiDevicesIdentity) Error(context containerInterface.IMessageContext, err error) {
	q := common.QueryMultiDevicesIdentity{}
	q.Error(context, err)
}
