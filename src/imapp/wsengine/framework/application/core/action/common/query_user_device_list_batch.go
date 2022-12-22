package common

import (
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	messageResultType "ws/framework/application/core/result/constant"
	"ws/framework/utils/xmpp"
)

// QueryUserDeviceListBatch .
type QueryUserDeviceListBatch struct {
	processor.BaseAction
	UserIDs []string
}

// Start .
func (q *QueryUserDeviceListBatch) Start(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	deviceListService := context.ResolveDeviceListService()

	jid := types.NewJID("", types.DefaultUserServer)
	nodeList := make([]waBinary.Node, 0)

	for i := range q.UserIDs {
		jid.User = q.UserIDs[i]

		// 是否有多设备特征
		if !deviceListService.HaveMultiDevice(jid.User) {
			continue
		}

		nodeList = append(nodeList, waBinary.Node{
			Tag:   "user",
			Attrs: waBinary.Attrs{"jid": jid.String()},
			Content: []waBinary.Node{{Tag: "devices",
				Attrs: waBinary.Attrs{"device_hash": xmpp.GenerateDeviceHash(jid)},
			}},
		})
	}

	if len(nodeList) == 0 {
		next()
		return
	}

	q.SendMessageId, err = context.SendIQ(
		xmpp.UsyncIQTemplate(context, "query", "message",
			[]waBinary.Node{
				{Tag: "query", Content: []waBinary.Node{
					{Tag: "devices", Attrs: waBinary.Attrs{"version": "2"}},
				}},
				{Tag: "list", Content: nodeList},
			},
		),
	)

	return
}

// Receive .
func (q *QueryUserDeviceListBatch) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	nodes := context.Message().GetChildren()

	if len(nodes) > 0 {
		nodes = nodes[0].GetChildrenByTag("list")
	}

	if len(nodes) > 0 {
		nodes = nodes[0].GetChildren()
	}

	deviceListService := context.ResolveDeviceListService()

	for _, user := range nodes {
		jid, ok := user.Attrs["jid"].(types.JID)

		if user.Tag != "user" || !ok {
			continue
		}

		deviceNode := user.GetChildByTag("devices")
		deviceIDList := xmpp.ParseDeviceIDList(&deviceNode)

		deviceListService.UpdateDeviceList(jid.User, deviceIDList)
	}

	next()
	return
}

// Error .
func (q *QueryUserDeviceListBatch) Error(context containerInterface.IMessageContext, err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.GetUserDeviceList,
		Error:      err,
	})
}
