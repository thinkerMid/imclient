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

// SMBUploadPushConfig .
type SMBUploadPushConfig struct {
	processor.BaseAction
}

// Start .
func (m *SMBUploadPushConfig) Start(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
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

	/**
	<iq id="1669434694-5" to="s.whatsapp.net" type="set" xmlns="urn:xmpp:whatsapp:push">
	<config background_location="1" call="Opening.m4r" default="note.m4r" groups="note.m4r" id="be272c052202a14f46c8337122b4d655ba94744a13e766b25512dbb752f94046"
	lc="CN" lg="zh" nse_call="0" nse_read="0" nse_retry="0" nse_ver="1" pkey="0cB03p0yeaGLruy2xw8NhICIDvd41LahFRvQyfyuhZc=" platform="smbi" preview="1" reg_push="1" version="2"
	voip="9d930fd0217e5d0a7070c9d21e2acb423596893acda07b175763fc903ed9fe39" voip_payload_type="0"/>
	</iq>
	*/

	m.SendMessageId, err = context.SendIQ(message.InfoQuery{
		ID:        context.GenerateRequestID(),
		Namespace: "urn:xmpp:whatsapp:push",
		Type:      "set",
		To:        types.ServerJID,
		Content: []waBinary.Node{
			{
				Tag: "config",
				Attrs: waBinary.Attrs{
					"background_location": "1",
					"call":                "Opening.m4r",
					"default":             "note.m4r",
					"groups":              "note.m4r",
					"id":                  config.ApnsToken,
					"lc":                  device.Country,
					"lg":                  device.Language,
					"nse_call":            "0",
					"nse_read":            "0",
					"nse_retry":           "0",
					"nse_ver":             "1",
					"pkey":                base64x.URLEncoding.EncodeToString(config.Pkey),
					"platform":            "smbi",
					"preview":             "1",
					"reg_push":            "1",
					"version":             "2",
					"voip":                config.VoipToken,
					"voip_payload_type":   "0",
				},
			},
		},
	})
	next()
	return
}

// Receive .
func (m *SMBUploadPushConfig) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	next()

	return nil
}

// Error .
func (m *SMBUploadPushConfig) Error(context containerInterface.IMessageContext, err error) {
}
