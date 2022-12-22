package privateChatNotification

import (
	"encoding/hex"
	"google.golang.org/protobuf/proto"
	"ws/framework/application/constant"
	waProto "ws/framework/application/constant/binary/proto"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/result/constant"
	"ws/framework/external"
	"ws/framework/lib/media_crypto"
	"ws/framework/utils/xmpp"
)

// ReceiveMessage .
type ReceiveMessage struct{}

// Receive .
func (r ReceiveMessage) Receive(context containerInterface.IMessageContext) (err error) {
	info, parseInfoErr := xmpp.ParseMessageInfo(context.ResolveJID(), context.Message())
	// 不解析群组消息
	if parseInfoErr != nil || info.IsGroup {
		return
	}

	var decrypted []byte
	var decryptErr error

	spf := context.ResolveSignalProtocolFactory()
	children := context.Message().GetChildren()

	for _, child := range children {
		if child.Tag != "enc" {
			continue
		}

		encType, ok := child.Attrs["type"].(string)
		if !ok {
			continue
		}

		content, _ := child.Content.([]byte)

		if encType == "pkmsg" || encType == "msg" {
			decrypted, decryptErr = spf.DecryptPrivateChatMessage(info.Sender, content, encType == "pkmsg")
		} else {
			continue
		}

		if decryptErr != nil {
			context.ResolveLogger().Error(decryptErr)
			continue
		}

		var msg waProto.Message
		unmarshalErr := proto.Unmarshal(decrypted, &msg)
		if unmarshalErr != nil {
			context.ResolveLogger().Error(unmarshalErr)
			continue
		}

		pushed := r.pushMessageContent(context, &info, &msg)
		if !pushed {
			context.ResolveLogger().Warnf("received %s message not handle, proto hex:%s", info.Type, hex.EncodeToString(decrypted))
		}
	}

	return constant.AbortedError
}

func (r *ReceiveMessage) pushMessageContent(context containerInterface.IMessageContext, messageInfo *types.MessageInfo, message *waProto.Message) bool {
	senderID := messageInfo.Sender.User

	chatMsg := external.ChatMessage{
		JIDNumber: senderID,
		MessageID: messageInfo.ID,
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

		// 目前这个消息结构还是未知用途的 只知道触发源是因为打开对方的wa.me的一个链接 之后给他发送消息 会从这个proto里过来 不是常见的消息
	} else if message.ExtendedTextMessage != nil {
		chatMsg.Conversation = message.ExtendedTextMessage.GetText()
	} else {
		return false
	}

	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.ReceiveMessage,
		IContent:   chatMsg,
	})

	return true
}
