package common

import (
	"fmt"
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	"ws/framework/application/core/result/constant"
)

// QueryDevicesIdentity .
type QueryDevicesIdentity struct {
	processor.BaseAction
	UserID               string
	CheckHaveMultiDevice bool // 需要检查有多设备才执行
}

// Start .
func (c *QueryDevicesIdentity) Start(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	// 是否有多设备特征
	if c.CheckHaveMultiDevice && !context.ResolveDeviceListService().HaveMultiDevice(c.UserID) {
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
				Tag: "identity",
				Content: []waBinary.Node{
					{
						Tag: "user",
						Attrs: waBinary.Attrs{
							"jid": fmt.Sprintf("%s@%s", c.UserID, types.DefaultUserServer),
						},
					},
				},
			}},
		},
	)

	return
}

// Receive .
func (c *QueryDevicesIdentity) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	defer next()

	// 不存了 代码里只有比对的作用

	//nodes := context.Message().GetChildrenByTag("list")
	//
	//if len(nodes) > 0 {
	//	nodes = nodes[0].GetChildren()
	//}
	//
	//identityService := context.ResolveIdentityService()
	//
	//for i := range nodes {
	//	child := nodes[i]
	//	if child.Tag != "user" {
	//		continue
	//	}
	//
	//	if identityNode, ok := child.GetOptionalChildByTag("identity"); ok {
	//		jid := child.AttrGetter().JID("jid")
	//
	//		byteContent := identityNode.Content.([]byte)
	//
	//		identityService.CreateIdentity(jid.SignalAddress(), byteContent)
	//	}
	//}
	return
}

func (c *QueryDevicesIdentity) Error(context containerInterface.IMessageContext, err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.GetUserDefaultDevice,
		Error:      err,
	})
}
