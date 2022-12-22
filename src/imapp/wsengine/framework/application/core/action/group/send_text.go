package group

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"sort"
	"ws/framework/application/constant/binary"
	waProto "ws/framework/application/constant/binary/proto"
	"ws/framework/application/constant/types"
	containerInterface "ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	messageResultType "ws/framework/application/core/result/constant"
	"ws/framework/utils"
	"ws/framework/utils/xmpp"
)

// SendText .
type SendText struct {
	processor.BaseAction
	GroupID     string
	MessageText string

	sendSKDMsgCount uint16 // max 512
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

	dstJID := types.NewJID(c.GroupID, types.GroupServer)

	node := xmpp.CreateMessageNode(dstJID, xmpp.TextMessageType)

	c.sendSKDMsgCount, err = encryptMessageNode(context, dstJID, &node, xmpp.NoneChild, &message)
	if err != nil {
		return err
	}

	c.SendMessageId, err = context.SendNode(node)

	return err
}

// Receive .
func (c *SendText) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	// 给群组内标记发送过消息
	if c.sendSKDMsgCount > 0 {
		context.ResolveSenderKeyService().SaveSentMessageByGroupID(c.GroupID)
	}

	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.SendGroupChatText,
		Content:    c.SendMessageId,
	})

	next()

	return
}

// Error .
func (c *SendText) Error(context containerInterface.IMessageContext, err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.SendGroupChatText,
		Error:      err,
	})
}

func encryptMessageNode(container containerInterface.IAppIocContainer, groupJID types.JID, node *waBinary.Node, mediaChildType xmpp.MessageChildType, msg *waProto.Message) (uint16, error) {
	spf := container.ResolveSignalProtocolFactory()

	// skmsg
	messagePlaintext, err := proto.Marshal(msg)

	// encrypt skmsg
	skmsg, skdmsg, err := spf.EncryptGroupMessage(groupJID.User, messagePlaintext)
	if err != nil {
		return 0, err
	}

	skdMessage := waProto.Message{
		SenderKeyDistributionMessage: &waProto.SenderKeyDistributionMessage{
			GroupId:                             proto.String(groupJID.String()),
			AxolotlSenderKeyDistributionMessage: skdmsg,
		},
	}

	// pkmsg senderkey
	skdPlaintext, err := proto.Marshal(&skdMessage)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal sender key distribution message to %s: %w", groupJID, err)
	}

	// all un sent member device
	senders, _ := container.ResolveSenderKeyService().FindUnSentMessageDeviceByGroupID(groupJID.User)

	tmpJID := types.NewJID("", types.DefaultUserServer)

	// pkmsg participants
	var participants []waBinary.Node
	var devices []string
	var content []waBinary.Node
	var encNode waBinary.Node
	var sendSKDMsgCount uint16

	for i := range senders {
		sendSKDMsgCount++

		tmpJID.User = senders[i].TheirJID
		tmpJID.Device = uint8(senders[i].DeviceID)
		tmpJID.AD = tmpJID.Device > 0

		// <enc type="?" v="2"><!-- 0 bytes --></enc>
		encNode, err = xmpp.EncryptMessageForJID(spf, tmpJID, mediaChildType, skdPlaintext)
		if err != nil {
			continue
		}

		jid := tmpJID.String()

		// <to jid="1000@s.whatsapp.net"><enc type="?" v="2"><!-- 0 bytes --></enc></to>
		participants = append(participants, waBinary.Node{
			Tag:     "to",
			Attrs:   waBinary.Attrs{"jid": jid},
			Content: []waBinary.Node{encNode},
		})

		devices = append(devices, jid)
	}

	if sendSKDMsgCount > 0 {
		//  <participants>
		//   <to jid="1000@s.whatsapp.net"><enc type="?" v="2"><!-- 0 bytes --></enc></to>
		//   <to jid="1000.0:1@s.whatsapp.net"><enc type="?" v="2"><!-- 0 bytes --></enc></to>
		//  </participants>
		content = append(content, waBinary.Node{
			Tag:     "participants",
			Content: participants,
		})

		sort.Strings(devices)

		node.Attrs["phash"] = xmpp.BatchGenerateDeviceHash(devices)
	}

	// <enc type="skmsg" v="2"><!-- 0 bytes --></enc>
	content = append(content, waBinary.Node{
		Tag:     "enc",
		Content: skmsg,
		Attrs:   waBinary.Attrs{"v": "2", "type": "skmsg"},
	})

	node.Content = content

	return sendSKDMsgCount, nil
}
