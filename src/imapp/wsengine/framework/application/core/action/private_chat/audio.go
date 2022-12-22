package privateChat

import (
	"google.golang.org/protobuf/proto"
	waProto "ws/framework/application/constant/binary/proto"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	"ws/framework/application/core/result/constant"
	"ws/framework/lib/media_crypto"
	"ws/framework/plugin/media_decode/media_content"
	"ws/framework/utils"
	"ws/framework/utils/xmpp"
)

// SendAudio .
type SendAudio struct {
	processor.BaseAction
	UserID string             // 接收人
	Audio  mediaContent.Audio // 音频

	Parser mediaCrypto.Voice
}

// Start .
func (c *SendAudio) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	message := waProto.Message{
		AudioMessage: &waProto.AudioMessage{},
		MessageContextInfo: &waProto.MessageContextInfo{
			DeviceListMetadata: &waProto.DeviceListMetadata{
				SenderTimestamp: proto.Uint64(uint64(utils.GetCurTime())),
			},
			DeviceListMetadataVersion: proto.Int32(2),
		},
	}

	if err = c.fill(context, &message); err != nil {
		return
	}

	dstJID := types.NewJID(c.UserID, types.DefaultUserServer)
	node := xmpp.CreateMessageNode(dstJID, xmpp.MediaMessageType)

	err = encodeProtocolMessage(context, dstJID, &node, xmpp.AudioMedia, &message)
	if err != nil {
		return err
	}

	c.SendMessageId, err = context.SendNode(node)

	return err
}

// Receive .
func (c *SendAudio) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.SendPrivateChatAudio,
		Content:    c.SendMessageId,
	})

	next()

	return
}

// Error .
func (c *SendAudio) Error(context containerInterface.IMessageContext, err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.SendPrivateChatAudio,
		Error:      err,
	})
}

func (c *SendAudio) fill(container containerInterface.IAppIocContainer, message *waProto.Message) error {
	url, path, err := container.ResolveMultimediaMessagingService().UploadMediaFile(c.Parser.File)
	if err != nil {
		return err
	}

	message.AudioMessage.Url = proto.String(url)
	message.AudioMessage.DirectPath = proto.String(path)
	message.AudioMessage.Mimetype = proto.String(c.Audio.Mimetype)
	message.AudioMessage.FileSha256 = c.Parser.FileSHA256
	message.AudioMessage.FileLength = proto.Uint64(c.Parser.FileLength)
	message.AudioMessage.MediaKey = c.Parser.RandKey
	message.AudioMessage.FileEncSha256 = c.Parser.FileEncSHA256
	message.AudioMessage.MediaKeyTimestamp = proto.Int64(c.Parser.MediaKeyTimestamp)
	message.AudioMessage.Seconds = proto.Uint32(c.Audio.Duration) //音频时长，单位 秒
	message.AudioMessage.Ptt = proto.Bool(true)
	message.AudioMessage.StreamingSidecar = c.Parser.StreamSideCar
	message.MessageContextInfo.DeviceListMetadata = &waProto.DeviceListMetadata{}
	message.MessageContextInfo.DeviceListMetadata.SenderTimestamp = proto.Uint64(uint64(utils.GetCurTime())) //这个应该是创建会话的时间
	message.MessageContextInfo.DeviceListMetadataVersion = proto.Int32(2)

	return nil
}
