package group

import (
	"strings"
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	"ws/framework/application/core/result/constant"
)

// QueryIcon .
type QueryIcon struct {
	processor.BaseAction
	GroupID string
	Preview bool
}

// Start .
func (m *QueryIcon) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	groupJID := types.NewJID(m.GroupID, types.GroupServer)

	iq := message.InfoQuery{
		Target:    groupJID,
		ID:        context.GenerateRequestID(),
		Namespace: "w:profile:picture",
		Type:      message.IqGet,
		To:        types.ServerJID,
		Content: []waBinary.Node{
			{
				Tag: "picture",
				Attrs: waBinary.Attrs{
					"query": "url",
					"type":  "image",
				},
			},
		},
	}

	if m.Preview {
		iq.Content = []waBinary.Node{
			{
				Tag: "picture",
				Attrs: waBinary.Attrs{
					"type": "preview",
				},
			},
		}
	}

	m.SendMessageId, err = context.SendIQ(iq)

	return
}

// Receive .
func (m *QueryIcon) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	defer next()

	// preview类型是图片buffer
	if m.Preview {
		return nil
	}

	pictureNode := context.Message().GetChildByTag("picture")
	pictureUrl := pictureNode.AttrGetter().String("url")

	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.QueryGroupIcon,
		Content:    pictureUrl,
	})

	return nil
}

// Error .
func (m *QueryIcon) Error(context containerInterface.IMessageContext, err error) {
	errorStr := err.Error()

	// 404的异常当做空头像 并不是一个真正的异常
	if strings.Contains(errorStr, "404") {
		context.AppendResult(containerInterface.MessageResult{
			ResultType: messageResultType.QueryGroupIcon,
		})
		return
	}

	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.QueryGroupIcon,
		Error:      err,
	})
}
