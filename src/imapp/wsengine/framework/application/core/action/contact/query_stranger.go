package contact

import (
	"errors"
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	"ws/framework/application/core/result/constant"
	"ws/framework/utils/xmpp"
)

// QueryStranger .
type QueryStranger struct {
	processor.BaseAction
	UserID string // 手机号
}

// Start .
func (m *QueryStranger) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	//<iq id="1663074247-76" to="84564844255@s.whatsapp.net" type="get" xmlns="usync">
	//<usync context="interactive" index="0" last="true" mode="query" sid="1663080416-2261066114-40">
	//<query><business><verified_name/><profile v="116"/></business><contact/></query>
	//<list><user><contact>+85256048573</contact></user></list>
	//</usync>
	//</iq>
	m.SendMessageId, err = context.SendIQ(
		xmpp.UsyncIQTemplate(context, "query", "interactive",
			[]waBinary.Node{
				{Tag: "query", Content: []waBinary.Node{
					{Tag: "business", Content: []waBinary.Node{
						{Tag: "verified_name"},
						{Tag: "profile",
							Attrs: waBinary.Attrs{
								"v": "116",
							},
						},
					}},
					{Tag: "contact"},
				}},
				{Tag: "list", Content: []waBinary.Node{{
					Tag: "user",
					Content: []waBinary.Node{{
						Tag:     "contact",
						Content: "+" + m.UserID,
					}},
				}}},
			},
		),
	)

	return
}

/**
<iq from="84564844255@s.whatsapp.net" id="1663074247-76" type="result">
<usync context="interactive" index="0" last="true" mode="query" sid="1663080416-2261066114-40">
<result><business/><contact/></result>
<list><user jid="85256048573@s.whatsapp.net"><contact type="in">+85256048573</contact></user></list>
</usync>
</iq>
*/

// Receive .
func (m *QueryStranger) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	usyncNode := context.Message().GetChildByTag("usync")
	listNode := usyncNode.GetChildByTag("list")
	nodes := listNode.GetChildren()

	for i := range nodes {
		userNode := nodes[i]
		if userNode.Tag != "user" {
			continue
		}

		jid, ok := userNode.Attrs["jid"].(types.JID)
		if !ok {
			continue
		}

		// 是否使用whatsapp
		contactNode := userNode.GetChildByTag("contact")
		getter := contactNode.AttrGetter()

		if getter.String("type") != "in" {
			return errors.New("not in whatsapp")
		}

		// 创建联系人
		context.ResolveContactService().CreateContact(jid.User, m.UserID)
		// 默认给它创建一个设备
		context.ResolveDeviceListService().AddDevice(jid.User, jid.Device)

		next()
		return nil
	}

	return errors.New("create stranger error")
}

func (m *QueryStranger) Error(context containerInterface.IMessageContext, err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.QueryStranger,
		Error:      err,
	})
}
