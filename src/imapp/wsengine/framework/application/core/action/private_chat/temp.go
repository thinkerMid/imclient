package privateChat

import (
	"google.golang.org/protobuf/proto"
	"ws/framework/application/constant/binary"
	waProto "ws/framework/application/constant/binary/proto"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	"ws/framework/application/core/result/constant"
	"ws/framework/utils"
	"ws/framework/utils/xmpp"
)

type MessageButton struct {
	ButtonUrl   string //按钮的跳转
	ButtonTitle string //按钮的标题
}

// SendTemp .
type SendTemp struct {
	processor.BaseAction
	UserID      string
	Title       string //标题
	MessageText string //内容
	Footer      string //页脚
	Button      []MessageButton
}

// Start .
func (c *SendTemp) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	tempLateMessage := &waProto.TemplateMessage{}
	tempLateMessage.HydratedTemplate = &waProto.HydratedFourRowTemplate{}
	tempLateMessage.HydratedTemplate.HydratedContentText = proto.String(c.MessageText)
	tempLateMessage.HydratedTemplate.HydratedFooterText = proto.String(c.Footer)
	tempLateMessage.HydratedTemplate.Title = &waProto.HydratedFourRowTemplate_HydratedTitleText{HydratedTitleText: c.Title}

	for i, v := range c.Button {
		button := &waProto.HydratedTemplateButton{}
		button.Index = proto.Uint32(uint32(i))
		obj := &waProto.HydratedURLButton{Url: proto.String(v.ButtonUrl), DisplayText: proto.String(v.ButtonTitle)}
		button.HydratedButton = &waProto.HydratedTemplateButton_UrlButton{UrlButton: obj}

		tempLateMessage.HydratedTemplate.HydratedButtons = append(tempLateMessage.HydratedTemplate.HydratedButtons, button)
	}

	message := waProto.Message{
		TemplateMessage: tempLateMessage,
		MessageContextInfo: &waProto.MessageContextInfo{
			DeviceListMetadata: &waProto.DeviceListMetadata{
				SenderTimestamp: proto.Uint64(uint64(utils.GetCurTime())),
			},
			DeviceListMetadataVersion: proto.Int32(2),
		},
	}

	dstJID := types.NewJID(c.UserID, types.DefaultUserServer)
	node := xmpp.CreateMessageNode(dstJID, xmpp.TextMessageType)

	err = encodeTemplateProtocolMessage(context, dstJID, &node, &message)
	if err != nil {
		return err
	}

	c.SendMessageId, err = context.SendNode(node)

	return err
}

// Receive .
func (c *SendTemp) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.SendPrivateChatTemp,
		Content:    c.SendMessageId,
	})

	next()

	return
}

// Error .
func (c *SendTemp) Error(context containerInterface.IMessageContext, err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.SendPrivateChatTemp,
		Error:      err,
	})
}

func encodeTemplateProtocolMessage(container containerInterface.IAppIocContainer, dstJid types.JID, node *waBinary.Node, message *waProto.Message) error {
	var tmpNode waBinary.Node

	err := encodeProtocolMessage(container, dstJid, &tmpNode, xmpp.NoneChild, message)
	if err != nil {
		return err
	}

	content := make([]waBinary.Node, 3)

	content[0].Tag = "hsm"
	content[1].Tag = "biz"
	content[2] = tmpNode.Content.([]waBinary.Node)[0]

	node.Content = content

	return nil
}
