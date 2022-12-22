package groupNotification

import (
	"encoding/hex"
	"google.golang.org/protobuf/proto"
	"ws/framework/application/constant"
	waProto "ws/framework/application/constant/binary/proto"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/result/constant"
	"ws/framework/application/libsignal/protocol"
	"ws/framework/external"
	"ws/framework/lib/media_crypto"
	"ws/framework/utils/xmpp"
)

// ReceiveMessage .
type ReceiveMessage struct{}

// Receive .
func (r ReceiveMessage) Receive(context containerInterface.IMessageContext) (err error) {
	info, parseInfoErr := xmpp.ParseMessageInfo(context.ResolveJID(), context.Message())
	// 只解析群组消息
	if parseInfoErr != nil || !info.IsGroup {
		return
	}

	// 用设备会话解密出senderkey
	var pkmsgContent []byte
	// 需要使用senderkey解密的消息内容
	var skmsgContent []byte

	var decrypted []byte
	var decryptErr error

	spf := context.ResolveSignalProtocolFactory()
	senderKeyName := protocol.NewSenderKeyName(info.Chat.User, info.Sender.SignalAddress())

	children := context.Message().GetChildren()

	for _, child := range children {
		if child.Tag != "enc" {
			continue
		}

		encType, ok := child.Attrs["type"].(string)
		if !ok {
			continue
		}

		if encType == "pkmsg" || encType == "msg" {
			msgContent, _ := child.Content.([]byte)

			pkmsgContent, decryptErr = spf.DecryptPrivateChatMessage(info.Sender, msgContent, encType == "pkmsg")

			// save senderKey
			if len(pkmsgContent) > 0 {
				decryptErr = spf.DecryptGroupSenderKey(senderKeyName, pkmsgContent)
				if decryptErr != nil {
					context.ResolveLogger().Error(decryptErr)
				}
			}

		} else if encType == "skmsg" {
			skmsgContent, _ = child.Content.([]byte)

			// group msg
			decrypted, decryptErr = spf.DecryptGroupMessage(senderKeyName, skmsgContent)
			if decryptErr == nil {
				var msg waProto.Message
				decryptErr = proto.Unmarshal(decrypted, &msg)

				if decryptErr == nil {
					if !r.pushMessageContent(context, &info, &msg) {
						context.ResolveLogger().Warnf("received %s message not handle, proto hex:%s", info.Type, hex.EncodeToString(decrypted))
					}
				}
			}
		}
	}

	return constant.AbortedError
}

func (r *ReceiveMessage) pushMessageContent(context containerInterface.IMessageContext, messageInfo *types.MessageInfo, message *waProto.Message) bool {
	senderID := messageInfo.Sender.User
	groupID := messageInfo.Chat.User

	chatMsg := external.ChatMessage{
		GroupNumber: groupID,
		JIDNumber:   senderID,
		MessageID:   messageInfo.ID,
	}

	// 图片消息
	if message.ImageMessage != nil {
		cbcKeyPair := mediaCrypto.ParseMediaKey(message.ImageMessage.GetMediaKey(), mediaCrypto.TMediaImage)

		chatMsg.Mimetype = message.ImageMessage.GetMimetype()
		chatMsg.MediaUrl = message.ImageMessage.GetUrl()
		chatMsg.CBCIv = cbcKeyPair.Iv
		chatMsg.CBCKey = cbcKeyPair.Enc
		// 语音消息
	} else if message.AudioMessage != nil {
		cbcKeyPair := mediaCrypto.ParseMediaKey(message.AudioMessage.GetMediaKey(), mediaCrypto.TMediaAudio)

		chatMsg.Mimetype = message.AudioMessage.GetMimetype()
		chatMsg.MediaUrl = message.AudioMessage.GetUrl()
		chatMsg.CBCIv = cbcKeyPair.Iv
		chatMsg.CBCKey = cbcKeyPair.Enc
		chatMsg.Seconds = message.AudioMessage.GetSeconds()
		// 视频消息
	} else if message.VideoMessage != nil {
		cbcKeyPair := mediaCrypto.ParseMediaKey(message.VideoMessage.GetMediaKey(), mediaCrypto.TMediaVideo)

		chatMsg.Mimetype = message.VideoMessage.GetMimetype()
		chatMsg.MediaUrl = message.VideoMessage.GetUrl()
		chatMsg.CBCIv = cbcKeyPair.Iv
		chatMsg.CBCKey = cbcKeyPair.Enc
		chatMsg.Seconds = message.VideoMessage.GetSeconds()
		chatMsg.JpegThumbnail = message.VideoMessage.GetJpegThumbnail()

		// 文本消息
	} else if message.Conversation != nil {
		chatMsg.Conversation = message.GetConversation()

		// 转发的消息
	} else if message.ExtendedTextMessage != nil {
		chatMsg.Conversation = message.ExtendedTextMessage.GetText()
	} else {
		return false
	}

	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.ReceiveGroupChatMessage,
		IContent:   chatMsg,
	})

	return true
}
