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

// SendVideo .
type SendVideo struct {
	processor.BaseAction
	UserID string             // 接收人
	Video  mediaContent.Video // 视频

	Parser mediaCrypto.Video
}

// Start .
func (c *SendVideo) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	message := waProto.Message{
		VideoMessage: &waProto.VideoMessage{},
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

	err = encodeProtocolMessage(context, dstJID, &node, xmpp.VideoMedia, &message)
	if err != nil {
		return err
	}

	c.SendMessageId, err = context.SendNode(node)

	return err
}

// Receive .
func (c *SendVideo) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.SendPrivateChatVideo,
		Content:    c.SendMessageId,
	})

	next()

	return
}

// Error .
func (c *SendVideo) Error(context containerInterface.IMessageContext, err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.SendPrivateChatVideo,
		Error:      err,
	})
}

func (c *SendVideo) fill(container containerInterface.IAppIocContainer, message *waProto.Message) error {
	url, path, err := container.ResolveMultimediaMessagingService().UploadMediaFile(c.Parser.File)
	if err != nil {
		return err
	}

	message.VideoMessage.Url = proto.String(url)
	message.VideoMessage.DirectPath = proto.String(path)
	message.VideoMessage.Mimetype = proto.String(c.Video.Mimetype)
	message.VideoMessage.FileSha256 = c.Parser.FileSHA256
	message.VideoMessage.FileLength = proto.Uint64(c.Parser.FileLength)
	message.VideoMessage.MediaKey = c.Parser.RandKey
	message.VideoMessage.FileEncSha256 = c.Parser.FileEncSHA256
	message.VideoMessage.MediaKeyTimestamp = proto.Int64(c.Parser.MediaKeyTimestamp)
	message.VideoMessage.Seconds = proto.Uint32(c.Video.Duration) //视频时长，单位 秒
	message.VideoMessage.Width = proto.Uint32(c.Video.Width)
	message.VideoMessage.Height = proto.Uint32(c.Video.Height)
	message.VideoMessage.StreamingSidecar = c.Parser.StreamSideCar
	message.VideoMessage.JpegThumbnail = c.Video.ThumbnailJPEG

	return nil
}
