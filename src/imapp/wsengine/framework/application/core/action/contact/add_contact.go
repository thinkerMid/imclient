package contact

import (
	"errors"
	"strings"
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	contactNotification "ws/framework/application/core/notification/contact"
	"ws/framework/application/core/processor"
	"ws/framework/application/core/result/constant"
	"ws/framework/external"
	functionTools "ws/framework/utils/function_tools"
	"ws/framework/utils/xmpp"
)

// Add .
type Add struct {
	processor.BaseAction
	UserID string // 手机号
}

// Start .
func (m *Add) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	//<iq id="1654710971-22" to="60174810922@s.whatsapp.net" type="get" xmlns="usync">
	//	<usync context="interactive" index="0" last="true" mode="delta" sid="1654711833-2066345802-5">
	//		<query><business><verified_name/><profile v="116"/></business><contact/><devices version="2"/><disappearing_mode/><sidelist/><status/></query>
	//			<list><user><contact>+5535999923962</contact></user></list><side_list/></usync></iq>

	//<iq id="1654715594-12" to="556499194737@s.whatsapp.net" type="get" xmlns="usync">
	//	<usync context="interactive" index="0" last="true" mode="delta" sid="1654715594-2779093937-2">
	//		<query><business><verified_name/><profile v="116"/></business><contact/><devices version="2"/><disappearing_mode/><sidelist/><status/></query>
	//			<list><user><contact>+5535999923962</contact></user></list><side_list/></usync></iq>

	m.SendMessageId, err = context.SendIQ(
		xmpp.UsyncIQTemplate(context, "delta", "interactive",
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
					{Tag: "devices", Attrs: waBinary.Attrs{"version": "2"}},
					{Tag: "disappearing_mode"},
					{Tag: "sidelist"},
					{Tag: "status"},
				}},
				{Tag: "list", Content: []waBinary.Node{{
					Tag: "user",
					Content: []waBinary.Node{{
						Tag:     "contact",
						Content: "+" + m.UserID,
					}},
				}}},
				{Tag: "side_list"},
			},
		),
	)

	return
}

/**
<iq from="85296475450@s.whatsapp.net" id="1661245651-15" type="result">
<usync context="interactive" index="0" last="true" mode="delta" sid="1661245651-741584589-4">
<result><status/><sidelist/><disappearing_mode/><devices/><business/><contact version="1661245652803150"/></result>
<list>
<user jid="85256048573@s.whatsapp.net">
<status t="1660562337">你好，我正在使用 WhatsApp</status>
<disappearing_mode duration="0" t="0"/>
<devices><device-list><device id="0"/></device-list></devices>
<contact type="in">+85256048573</contact>
</user>
</list></usync></iq>
*/

// Receive .
func (m *Add) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	nodes := context.Message().GetChildren()

	if len(nodes) > 0 {
		nodes = nodes[0].GetChildrenByTag("list")
	}

	if len(nodes) > 0 {
		nodes = nodes[0].GetChildren()
	}

	contactService := context.ResolveContactService()
	deviceListService := context.ResolveDeviceListService()

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

		// 先保存联系人
		if contactService.CreateAddressBookContactByJID(jid.User, m.UserID) != nil {
			return errors.New("add contact error")
		}

		// 设备信息
		devicesNode := userNode.GetChildByTag("devices")

		// 更新设备
		deviceIDList := xmpp.ParseDeviceIDList(&devicesNode)
		deviceListService.UpdateDeviceList(jid.User, deviceIDList)

		// 推新增联系人
		phoneNumber := contactNode.ContentString()[0]
		context.AppendResult(containerInterface.MessageResult{
			ResultType: messageResultType.AddContact,
			IContent: external.Contact{
				JIDNumber:   jid.User,
				PhoneNumber: strings.ReplaceAll(phoneNumber, "+", ""),
			},
		})

		// 签名
		statusNode := userNode.GetChildByTag("status")
		var signatureText string

		bStatus, ok := statusNode.Content.([]byte)
		if ok {
			if len(bStatus) > len(contactNotification.UnknownCharacter) && functionTools.SliceEqual(bStatus[:len(contactNotification.UnknownCharacter)], contactNotification.UnknownCharacter) {
				bStatus = bStatus[len(contactNotification.UnknownCharacter):]
			}

			signatureText = string(bStatus)
		}

		context.AppendResult(containerInterface.MessageResult{
			ResultType: messageResultType.ContactSignatureUpdate,
			IContent: external.ProfileUpdate{
				JIDNumber: jid.User,
				Content:   signatureText,
			},
		})

		if jid.User != m.UserID {
			// 更改这个UserID 如果有monitor监控着会更改后续的Action所持有的UserID
			m.UserID = jid.User
		}

		next()
		return nil
	}

	return nil
}

func (m *Add) Error(context containerInterface.IMessageContext, err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.AddContact,
		Error:      err,
	})
}
