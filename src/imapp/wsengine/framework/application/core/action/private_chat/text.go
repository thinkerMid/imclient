package privateChat

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"strconv"
	"ws/framework/application/constant/binary"
	waProto "ws/framework/application/constant/binary/proto"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	"ws/framework/application/core/result/constant"
	"ws/framework/utils"
	"ws/framework/utils/xmpp"
)

// SendText .
type SendText struct {
	processor.BaseAction
	UserID      string
	MessageText string
}

// Start .
func (c *SendText) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	message := waProto.Message{
		Conversation: proto.String(c.MessageText),
		MessageContextInfo: &waProto.MessageContextInfo{
			DeviceListMetadata: &waProto.DeviceListMetadata{
				SenderTimestamp: proto.Uint64(uint64(utils.GetCurTime())),
			},
			DeviceListMetadataVersion: proto.Int32(2),
		},
	}

	dstJID := types.NewJID(c.UserID, types.DefaultUserServer)

	node := xmpp.CreateMessageNode(dstJID, xmpp.TextMessageType)

	err = encodeProtocolMessage(context, dstJID, &node, xmpp.NoneChild, &message)
	if err != nil {
		return err
	}

	c.SendMessageId, err = context.SendNode(node)

	return err
}

// Receive .
func (c *SendText) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.SendPrivateChatText,
		Content:    c.SendMessageId,
	})

	next()

	return
}

// Error .
func (c *SendText) Error(context containerInterface.IMessageContext, err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.SendPrivateChatText,
		Error:      err,
	})
}

func encodeProtocolMessage(container containerInterface.IAppIocContainer, dstJID types.JID, node *waBinary.Node, mediaType xmpp.MessageChildType, message *waProto.Message) error {
	messagePlaintext, err := proto.Marshal(message)
	if err != nil {
		return err
	}

	var toNodes []waBinary.Node

	spf := container.ResolveSignalProtocolFactory()
	logger := container.ResolveLogger()

	idList := container.ResolveDeviceListService().FindDeviceIDList(dstJID.User)
	multipleDevices := len(idList) > 1

	for _, id := range idList {
		jid := types.NewJID(dstJID.User, types.DefaultUserServer)

		if id > 0 {
			jid.AD = true
			jid.Device = id
		}

		encNode, err := xmpp.EncryptMessageForJID(spf, jid, mediaType, messagePlaintext)
		if err != nil {
			logger.Warnf("%s %s", jid.SignalAddress().String(), err)
			continue
		}

		if multipleDevices {
			// <to jid="1000@s.whatsapp.net"><enc type="?" v="2"><!-- 0 bytes --></enc></to>
			toNodes = append(toNodes, waBinary.Node{
				Tag:     "to",
				Attrs:   waBinary.Attrs{"jid": jid.String()},
				Content: []waBinary.Node{encNode},
			})
			continue
		}

		// <enc type="?" v="2"><!-- 0 bytes --></enc>
		toNodes = append(toNodes, encNode)
	}

	// 没有消息内容 大概率就是sessoin的问题
	if len(toNodes) == 0 {
		return fmt.Errorf("cipher encryption failed: not find %s session", dstJID)
	}

	// 多设备消息
	//  <participants>
	//   <to jid="1000@s.whatsapp.net"><enc type="?" v="2"><!-- 0 bytes --></enc></to>
	//   <to jid="1000.0:1@s.whatsapp.net"><enc type="?" v="2"><!-- 0 bytes --></enc></to>
	//  </participants>
	if multipleDevices {
		toNodes = []waBinary.Node{{Tag: "participants", Content: toNodes}}
	}

	// 第一次给陌生人发送消息
	//  <url_number/>
	appendNodeTagWhenFirstSendMessageToStranger(container, dstJID.User, &toNodes)

	// 商业版本attr属性
	// <message id="?" to="?@s.whatsapp.net" type="?" verified_name="?"/>
	appendNodeAttrWhenBusinessPlatform(container, node)

	node.Content = toNodes
	return nil
}

// appendNodeTagWhenFirstSendMessageToStranger 第一次给陌生人发送消息
func appendNodeTagWhenFirstSendMessageToStranger(container containerInterface.IAppIocContainer, dstJID string, content *[]waBinary.Node) {
	contact := container.ResolveContactService().FindByJID(dstJID)

	// 陌生人 并且 未发送过一条消息 并且 未接收过消息
	// 符合这个条件的 是通过链接聊天 当第一条消息要有url_number
	if contact.InAddressBook == false && contact.ChatWith == false && contact.ReceiveChat == false {
		*content = append(*content, waBinary.Node{Tag: "url_number"})
	}
}

// appendNodeAttrWhenBusinessPlatform 商业版本attr属性
func appendNodeAttrWhenBusinessPlatform(container containerInterface.IAppIocContainer, messageNode *waBinary.Node) {
	account := container.ResolveAccountService().Context()
	if !account.BusinessAccount {
		return
	}

	businessProfile := container.ResolveBusinessService().Context()
	if businessProfile == nil || businessProfile.VNameSerial == 0 {
		return
	}

	// k=verified_name v=vname_serial
	messageNode.Attrs["verified_name"] = strconv.FormatInt(businessProfile.VNameSerial, 10)
}
