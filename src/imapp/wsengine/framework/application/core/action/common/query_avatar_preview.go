package common

import (
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
)

// QueryAvatarPreview .
type QueryAvatarPreview struct {
	processor.BaseAction
	UserID     string
	AvatarSize int
}

// RaiseErrorWhenNodeError 是否抛出异常
func (m *QueryAvatarPreview) RaiseErrorWhenNodeError() bool {
	return false
}

// Start .
func (m *QueryAvatarPreview) Start(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	targetJID := types.NewJID(m.UserID, types.DefaultUserServer)

	m.SendMessageId, err = context.SendIQ(message.InfoQuery{
		Namespace: "w:profile:picture",
		Type:      message.IqGet,
		To:        types.ServerJID,
		Target:    targetJID,
		ID:        context.GenerateRequestID(),
		Content: []waBinary.Node{{
			Tag: "picture",
			Attrs: waBinary.Attrs{
				"type": "preview",
			},
		}},
	})

	return
}

// Receive .
func (m *QueryAvatarPreview) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	defer next()

	picture, ok := context.Message().GetOptionalChildByTag("picture")
	if !ok {
		return
	}

	content := picture.Content.([]byte)
	m.AvatarSize = len(content) / 2
	return
}

// Error .
func (m *QueryAvatarPreview) Error(context containerInterface.IMessageContext, err error) {

}
