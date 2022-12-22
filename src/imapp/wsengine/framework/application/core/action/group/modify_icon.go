package group

import (
	goContext "context"
	"time"
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	"ws/framework/application/core/result/constant"
	"ws/framework/plugin/media_decode/vips"
)

// ModifyIcon .
type ModifyIcon struct {
	processor.BaseAction
	GroupID string
	Icon    []byte
}

// Start .
func (m *ModifyIcon) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	const (
		IconImageQuality  = 50
		IconImageProgress = false
	)

	ctx, cancel := goContext.WithTimeout(goContext.Background(), time.Second*30)
	defer cancel()

	groupJID := types.NewJID(m.GroupID, types.GroupServer)

	var iconBuff []byte
	iconBuff, err = vips.CompressImageWithProgress(ctx, m.Icon, IconImageQuality, IconImageProgress)
	if err != nil {
		return err
	}

	iq := message.InfoQuery{
		Target:    groupJID,
		ID:        context.GenerateRequestID(),
		Namespace: "w:profile:picture",
		Type:      message.IqSet,
		To:        types.ServerJID,
		Content: []waBinary.Node{
			{
				Tag: "picture",
				Attrs: waBinary.Attrs{
					"type": "image",
				},
				Content: iconBuff,
			},
		},
	}

	m.SendMessageId, err = context.SendIQ(iq)

	return
}

// Receive .
func (m *ModifyIcon) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	// 推出操作结果
	pictureNode := context.Message().GetChildByTag("picture")
	pictureUrl := pictureNode.AttrGetter().String("url")

	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.ModifyGroupIcon,
		Content:    pictureUrl,
	})

	next()

	return nil
}

// Error .
func (m *ModifyIcon) Error(context containerInterface.IMessageContext, err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.ModifyGroupIcon,
		Error:      err,
	})
}
