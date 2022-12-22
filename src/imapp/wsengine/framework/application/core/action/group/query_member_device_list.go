package group

import (
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	"ws/framework/application/core/result/constant"
	"ws/framework/utils/xmpp"
)

// QueryMemberDeviceList
// 用于查询群成员的设备列表
type QueryMemberDeviceList struct {
	processor.BaseAction
	GroupID string
	UserIDs []string

	cacheSearchDeviceID map[string][]uint8
}

// Start .
func (q *QueryMemberDeviceList) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	deviceListService := context.ResolveDeviceListService()

	q.cacheSearchDeviceID = make(map[string][]uint8)

	jid := types.NewJID("", types.DefaultUserServer)
	nodeList := make([]waBinary.Node, 0)

	for i := range q.UserIDs {
		jid.User = q.UserIDs[i]

		idList := deviceListService.FindDeviceIDList(jid.User)

		// 是否有多设备特征
		if len(idList) > 1 {
			nodeList = append(nodeList, waBinary.Node{
				Tag:   "user",
				Attrs: waBinary.Attrs{"jid": jid.String()},
				Content: []waBinary.Node{{Tag: "devices",
					Attrs: waBinary.Attrs{"device_hash": xmpp.GenerateDeviceHash(jid)},
				}},
			})
		} else {
			nodeList = append(nodeList, waBinary.Node{
				Tag:   "user",
				Attrs: waBinary.Attrs{"jid": jid.String()},
			})
		}

		// 把查询的设备列表保存起来
		q.cacheSearchDeviceID[jid.User] = idList
	}

	// <iq id="1664331317-5" to="60175627234@s.whatsapp.net" type="get" xmlns="usync">
	//<usync context="notification" index="0" last="true" mode="query" sid="1664331318-792536448-2">
	//<query>
	//<devices version="2"/>
	//<disappearing_mode/>
	//</query>
	//<list>
	//<user jid="79052978547@s.whatsapp.net"/>
	//<user jid="60175627234@s.whatsapp.net"><devices device_hash="2:liWfapBZ"/></user>
	//<user jid="85295664081@s.whatsapp.net"/>
	//</list>
	//</usync>
	//</iq>

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
func (q *QueryMemberDeviceList) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	nodes := context.Message().GetChildren()

	if len(nodes) > 0 {
		nodes = nodes[0].GetChildrenByTag("list")
	}

	if len(nodes) > 0 {
		nodes = nodes[0].GetChildren()
	}

	deviceListService := context.ResolveDeviceListService()

	//<iq from="60175627234@s.whatsapp.net" id="1664331317-5" type="result">
	//<usync context="notification" index="0" last="true" mode="query" sid="1664331318-792536448-2">
	//<result>
	//<devices/>
	//<disappearing_mode/>
	//</result>
	//<list>
	//<user jid="85295664081@s.whatsapp.net"><devices><device-list><device id="0"/></device-list></devices><disappearing_mode duration="0" t="0"/></user>
	//<user jid="79052978547@s.whatsapp.net"><devices><device-list><device id="0"/></device-list></devices><disappearing_mode duration="0" t="0"/></user>
	//<user jid="60175627234@s.whatsapp.net"><disappearing_mode duration="0" t="0"/></user>
	//</list>
	//</usync>
	//</iq>
	var memberJIDs []types.JID

	for _, user := range nodes {
		if user.Tag != "user" {
			continue
		}

		jid, _ := user.Attrs["jid"].(types.JID)

		deviceNode, ok := user.GetOptionalChildByTag("devices")
		if !ok {
			continue
		}

		newDeviceIDList := xmpp.ParseDeviceIDList(&deviceNode)

		_, ok = q.cacheSearchDeviceID[jid.User]
		if ok {
			delete(q.cacheSearchDeviceID, jid.User)
		}

		deviceListService.UpdateDeviceList(jid.User, newDeviceIDList)

		// 用于群会话
		memberJIDs = append(memberJIDs, jid)
	}

	// 把没有匹配到的 作默认添加（有时候这个接口查询频率多了 一些号码似乎不返回）
	for jid, deviceIDList := range q.cacheSearchDeviceID {
		for i := 0; i < len(deviceIDList); i++ {
			jid := types.NewJID(jid, types.DefaultUserServer)
			jid.Device = deviceIDList[i]

			memberJIDs = append(memberJIDs, jid)
		}
	}

	// 批量创建群成员的设备
	context.ResolveSenderKeyService().BatchCreateDevice(q.GroupID, memberJIDs)

	next()
	return
}

// Error .
func (q *QueryMemberDeviceList) Error(context containerInterface.IMessageContext, err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.GetUserDeviceList,
		Error:      err,
	})
}
