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

// SendImage .
type SendImage struct {
	processor.BaseAction
	UserID  string             // 接收人
	Image   mediaContent.Image // 图片
	Caption string             // 标题
	Parser  mediaCrypto.Image
}

// Start .
func (c *SendImage) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	message := waProto.Message{
		ImageMessage: &waProto.ImageMessage{},
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

	err = encodeProtocolMessage(context, dstJID, &node, xmpp.ImageMedia, &message)
	if err != nil {
		return err
	}

	c.SendMessageId, err = context.SendNode(node)

	return err
}

// Receive .
func (c *SendImage) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.SendPrivateChatImage,
		Content:    c.SendMessageId,
	})

	next()

	return
}

// Error .
func (c *SendImage) Error(context containerInterface.IMessageContext, err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.SendPrivateChatImage,
		Error:      err,
	})
}

func (c *SendImage) fill(container containerInterface.IAppIocContainer, message *waProto.Message) error {
	url, path, err := container.ResolveMultimediaMessagingService().UploadMediaFile(c.Parser.File)
	if err != nil {
		return err
	}

	message.ImageMessage.Url = proto.String(url)
	message.ImageMessage.Mimetype = proto.String(c.Image.Mimetype)
	message.ImageMessage.Height = proto.Uint32(c.Image.Height)
	message.ImageMessage.Width = proto.Uint32(c.Image.Width)
	message.ImageMessage.DirectPath = proto.String(path)
	message.ImageMessage.FileSha256 = c.Parser.FileSHA256
	message.ImageMessage.MediaKey = c.Parser.RandKey
	message.ImageMessage.FileEncSha256 = c.Parser.FileEncSHA256
	message.ImageMessage.FileLength = proto.Uint64(c.Parser.FileLength)
	message.ImageMessage.MediaKeyTimestamp = proto.Int64(c.Parser.MediaKeyTimestamp) // proto.Int64(time.Now().Unix())
	message.ImageMessage.ScanLengths = c.Parser.ScanLengths
	message.ImageMessage.FirstScanLength = &c.Parser.ScanLengths[0]
	message.ImageMessage.ScansSidecar = c.Parser.ScanSideCar
	message.ImageMessage.FirstScanSidecar = c.Parser.ScanSideCar[:10]
	message.ImageMessage.MidQualityFileSha256 = c.Parser.MidQualitySha256

	// Caption 非空则是图文消息
	if len(c.Caption) > 0 {
		message.ImageMessage.Caption = proto.String(c.Caption)
	}

	message.ImageMessage.JpegThumbnail = c.Image.ThumbnailJPEG

	return err
}
