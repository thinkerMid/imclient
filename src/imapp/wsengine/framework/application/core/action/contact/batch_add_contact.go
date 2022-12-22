package contact

import (
	"fmt"
	"strings"
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	"ws/framework/application/core/result/constant"
	"ws/framework/external"
	"ws/framework/utils/xmpp"
)

// BatchAdd .
type BatchAdd struct {
	processor.BaseAction
	UserIDs    []string // 手机号
	SaveResult bool     // 是否保存联系人结果
}

// Start .
func (m *BatchAdd) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	contacts := make([]waBinary.Node, len(m.UserIDs))

	for i := range contacts {
		// +852123456789
		contacts[i].Tag = "user"
		contacts[i].Content = []waBinary.Node{{
			Tag:     "contact",
			Content: "+" + m.UserIDs[i],
		}}
	}

	m.SendMessageId, err = context.SendIQ(
		xmpp.UsyncIQTemplate(context, "full", "interactive",
			[]waBinary.Node{
				{Tag: "query", Content: []waBinary.Node{
					{Tag: "business", Content: []waBinary.Node{
						{Tag: "verified_name"},
						{Tag: "profile",
							Attrs: waBinary.Attrs{
								"v": "372",
							},
						},
					}},
					{Tag: "contact"},
					{Tag: "devices", Attrs: waBinary.Attrs{"version": "2"}},
					{Tag: "disappearing_mode"},
					{Tag: "sidelist"},
					{Tag: "status"},
				}},
				{Tag: "list", Content: contacts},
				{Tag: "side_list"},
			},
		),
	)

	return
}

// Receive .
func (m *BatchAdd) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	defer next()

	deviceListService := context.ResolveDeviceListService()
	nodes := context.Message().GetChildren()

	if len(nodes) > 0 {
		nodes = nodes[0].GetChildrenByTag("list")
	}

	if len(nodes) > 0 {
		nodes = nodes[0].GetChildren()
	}

	// 没有返回任何联系人结果
	if len(nodes) == 0 {
		/**
		<iq from="14193015324@s.whatsapp.net" id="1667441868-17" type="result">
		<usync context="interactive" index="0" last="true" mode="full" sid="1667441868-3053743388-12">
		<result>
			<status refresh="592097"/>
			<sidelist refresh="577356"/><disappearing_mode refresh="549921"/>
			<devices refresh="546732"/>
			<business refresh="521302"/>
			<contact><error backoff="807" code="429" text="rate-overlimit"/></contact>
		</result>
		<list/>
		</usync>
		</iq>
		*/

		// 这个节点下面是不是有异常
		resultNode := context.Message().GetChildByTag("result")
		contactNode := resultNode.GetChildByTag("contact")
		errorNode, ok := contactNode.GetOptionalChildByTag("error")
		if ok {
			return fmt.Errorf(errorNode.AttrGetter().String("text"))
		}
	}

	result := BatchAddContact{}
	phoneNumbers := make([]string, 0)

	for i := range nodes {
		n := nodes[i]
		if n.Tag != "user" {
			continue
		}

		jid, ok := n.Attrs["jid"].(types.JID)
		if !ok {
			continue
		}

		// 没注册
		contactNode := n.GetChildByTag("contact")
		getter := contactNode.AttrGetter()
		if getter.String("type") != "in" {
			result.NotRegisterNumber = append(result.NotRegisterNumber, jid.User)
			continue
		}

		// 有异常
		statusNode := n.GetChildByTag("status")
		getter = statusNode.AttrGetter()
		if len(getter.OptionalString("code")) > 0 {
			result.ErrorStatusNumber = append(result.ErrorStatusNumber, jid.User)
			continue
		}

		// 可用的号码
		result.AvailableNumber = append(result.AvailableNumber, jid.User)

		// 设备信息
		devicesNode := n.GetChildByTag("devices")
		if _, ok := devicesNode.GetOptionalChildByTag("key-index-list"); ok {
			result.HaveKeyIndexNumber = append(result.HaveKeyIndexNumber, jid.User)
		}

		// 如果是需要保存结果的
		if m.SaveResult {
			// 保存设备数量
			deviceIDList := xmpp.ParseDeviceIDList(&devicesNode)
			deviceListService.UpdateDeviceList(jid.User, deviceIDList)
		}

		// 联系人号码
		phoneNumber := contactNode.ContentString()[0]
		phoneNumber = strings.ReplaceAll(phoneNumber, "+", "")
		phoneNumbers = append(phoneNumbers, phoneNumber)

		result.AvailableContact = append(result.AvailableContact, external.Contact{
			JIDNumber:   jid.User,
			PhoneNumber: phoneNumber,
		})
	}

	if m.SaveResult && len(result.AvailableNumber) > 0 {
		context.ResolveContactService().BatchCreateAddressBookContactByJIDList(result.AvailableNumber, phoneNumbers)
	}

	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.BatchAddContact,
		IContent:   result,
	})

	return
}

func (m *BatchAdd) Error(context containerInterface.IMessageContext, err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.BatchAddContact,
		Error:      err,
	})
}

// BatchAddContact .
type BatchAddContact struct {
	NotRegisterNumber  []string
	ErrorStatusNumber  []string
	HaveKeyIndexNumber []string
	AvailableNumber    []string
	AvailableContact   []external.Contact
}
