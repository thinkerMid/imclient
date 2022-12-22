package common

import (
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	"ws/framework/application/core/result/constant"
	"ws/framework/utils/xmpp"
)

// QueryUserDeviceList .
type QueryUserDeviceList struct {
	processor.BaseAction
	UserID               string
	CheckHaveMultiDevice bool // 需要检查有多设备才执行
}

// Start .
func (q *QueryUserDeviceList) Start(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	// 检查多设备特征
	if q.CheckHaveMultiDevice && !context.ResolveDeviceListService().HaveMultiDevice(q.UserID) {
		next()
		return
	}

	jid := types.NewJID(q.UserID, types.DefaultUserServer)

	q.SendMessageId, err = context.SendIQ(
		xmpp.UsyncIQTemplate(context, "query", "message",
			[]waBinary.Node{
				{Tag: "query", Content: []waBinary.Node{
					{Tag: "devices", Attrs: waBinary.Attrs{"version": "2"}},
				}},
				{Tag: "list", Content: []waBinary.Node{
					{
						Tag:   "user",
						Attrs: waBinary.Attrs{"jid": jid.String()},
						Content: []waBinary.Node{{Tag: "devices",
							Attrs: waBinary.Attrs{"device_hash": xmpp.GenerateDeviceHash(jid)},
						}},
					},
				}},
			},
		),
	)

	return
}

// Receive .
func (q *QueryUserDeviceList) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
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
func (q *QueryUserDeviceList) Error(context containerInterface.IMessageContext, err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.GetUserDeviceList,
		Error:      err,
	})
}
