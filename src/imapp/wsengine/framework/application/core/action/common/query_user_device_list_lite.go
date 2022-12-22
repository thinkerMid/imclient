package common

import (
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	"ws/framework/utils/xmpp"
)

// QueryUserDeviceListLite
// context是notification的usync包 跟其他查询设备列表请求相比 这个似乎是数据权限比较低的查询接口 有时候频率高了还会存在空结果的情况
//
//	目前用于查自己设备 没什么用
type QueryUserDeviceListLite struct {
	processor.BaseAction
	UserIDs []string
}

// Start .
func (q *QueryUserDeviceListLite) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	jid := types.NewJID("", types.DefaultUserServer)
	nodeList := make([]waBinary.Node, 0)

	for i := range q.UserIDs {
		jid.User = q.UserIDs[i]

		nodeList = append(nodeList, waBinary.Node{
			Tag:   "user",
			Attrs: waBinary.Attrs{"jid": jid.String()},
			Content: []waBinary.Node{{Tag: "devices",
				Attrs: waBinary.Attrs{"device_hash": xmpp.GenerateDeviceHash(jid)},
			}},
		})
	}

	q.SendMessageId, err = context.SendIQ(
		xmpp.UsyncIQTemplate(context, "query", "notification",
			[]waBinary.Node{
				{Tag: "query", Content: []waBinary.Node{
					{Tag: "devices", Attrs: waBinary.Attrs{"version": "2"}},
					{Tag: "disappearing_mode"},
				}},
				{Tag: "list", Content: nodeList},
			},
		),
	)

	return
}

// Receive .
func (q *QueryUserDeviceListLite) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	next()
	return
}

// Error .
func (q *QueryUserDeviceListLite) Error(context containerInterface.IMessageContext, err error) {}
