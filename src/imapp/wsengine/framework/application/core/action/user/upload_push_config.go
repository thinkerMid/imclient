package user

import (
	"github.com/chenzhuoyu/base64x"
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	"ws/framework/application/data_storage/push_config/database"
)

// UploadPushConfig .
type UploadPushConfig struct {
	processor.BaseAction
}

// Start .
func (m *UploadPushConfig) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	device := context.ResolveDeviceService().Context()

	var config *pushConfigDB.PushConfig

	// 查询 顺便 刷新pkey
	context.ResolvePushConfigService().ContextExecute(func(p *pushConfigDB.PushConfig) {
		config = p

		p.RefreshPkey()
	})

	// 没有则创建
	if config == nil {
		config, _ = context.ResolvePushConfigService().Create()
	}

	m.SendMessageId, err = context.SendIQ(message.InfoQuery{
		ID:        context.GenerateRequestID(),
		Namespace: "urn:xmpp:whatsapp:push",
		Type:      "set",
		To:        types.ServerJID,
		Content: []waBinary.Node{
			{
				Tag: "config",
				Attrs: waBinary.Attrs{
					"id":                  config.VoipToken,
					"voip":                config.ApnsToken,
					"pkey":                base64x.URLEncoding.EncodeToString(config.Pkey),
					"nse_call":            "0",
					"groups":              "note.m4r",
					"preview":             "1",
					"call":                "Opening.m4r",
					"version":             "2",
					"lg":                  device.Language,
					"reg_push":            "1",
					"lc":                  device.Country,
					"default":             "note.m4r",
					"background_location": "1",
					"platform":            "apple",
					"nse_ver":             "1",
					"voip_payload_type":   "0",
				},
			},
		},
	})

	return
}

// Receive .
func (m *UploadPushConfig) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	next()

	return nil
}

// Error .
func (m *UploadPushConfig) Error(context containerInterface.IMessageContext, err error) {
}
