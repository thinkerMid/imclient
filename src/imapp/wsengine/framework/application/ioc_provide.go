package application

import (
	"ws/framework/application/connection"
	"ws/framework/application/constant/types"
	"ws/framework/application/container"
	containerInterface "ws/framework/application/container/abstract_interface"
	"ws/framework/application/data_storage/ab_key/service"
	"ws/framework/application/data_storage/account/service"
	"ws/framework/application/data_storage/aes_key/service"
	businessService "ws/framework/application/data_storage/business/service"
	"ws/framework/application/data_storage/cache"
	"ws/framework/application/data_storage/contact/service"
	"ws/framework/application/data_storage/device/service"
	"ws/framework/application/data_storage/device_list/service"
	"ws/framework/application/data_storage/event/service"
	"ws/framework/application/data_storage/group/service"
	"ws/framework/application/data_storage/identity/service"
	mmsService "ws/framework/application/data_storage/mms"
	"ws/framework/application/data_storage/prekey/service"
	"ws/framework/application/data_storage/push_config/service"
	"ws/framework/application/data_storage/registration_token/service"
	"ws/framework/application/data_storage/routing_info/service"
	"ws/framework/application/data_storage/sender_key/service"
	"ws/framework/application/data_storage/session/service"
	"ws/framework/application/data_storage/signal_protocol"
	"ws/framework/application/data_storage/signed_prekey/service"
	"ws/framework/plugin/logger"
	networkConstant "ws/framework/plugin/network/constant"
	"ws/framework/plugin/network/netpoll"
	tlsCert "ws/framework/plugin/ws_tls_cert"
)

var _ containerInterface.IMessageChannel = &App{}

// New .
func New(jid string, configuration *containerInterface.WhatsappConfiguration) IApplication {
	app := App{
		logger: logger.New(jid).Named("Application"),
	}

	container := provide(jid, configuration)
	// 把app注入成xmpp的通信载体
	container.Inject(appContainer.XMPPMessageChannel, &app)

	app.container = container

	return &app
}

func provide(jidNumber string, configuration *containerInterface.WhatsappConfiguration) containerInterface.IAppIocContainer {
	c := appContainer.NewAppIocContainer()

	jid := types.NewJID(jidNumber, types.DefaultUserServer)

	// 基础设置
	c.Inject(appContainer.JID, jid)
	c.Inject(appContainer.Logger, logger.New(jid.User))
	c.Inject(appContainer.ConnectionConfig, networkConstant.ConnectionConfig{Tls: tlsCert.TlsConfig()})
	c.Inject(appContainer.WhatsappConfiguration, configuration)

	// 全局数据缓存
	c.Provide(
		appContainer.MemoryCache,
		func(container containerInterface.IAppIocContainer) interface{} {
			return &memoryCacheService.MemoryCache{}
		},
	)

	// IM通信协议逻辑处理
	c.Provide(
		appContainer.SignalProtocol,
		func(container containerInterface.IAppIocContainer) interface{} {
			return &signalProtocolService.Factory{}
		},
	)

	// HTTP
	c.Provide(
		appContainer.HttpClient,
		func(container containerInterface.IAppIocContainer) interface{} {
			config := container.ResolveConnectionConfig()
			config.Tls = tlsCert.TlsConfig()

			return netpoll.HTTP(config)
		},
	)

	// XMPP握手
	c.Provide(
		appContainer.HandshakeHandler,
		func(container containerInterface.IAppIocContainer) interface{} {
			return &connection.Handshake{}
		},
	)

	// TCP
	c.Provide(
		appContainer.Connection,
		func(container containerInterface.IAppIocContainer) interface{} {
			return &connection.Connection{}
		},
	)

	// 数据存储服务
	provideDataStorageService(c)

	return c
}

func provideDataStorageService(c containerInterface.IAppIocContainer) {
	// ABKey
	c.Provide(
		appContainer.ABKey,
		func(container containerInterface.IAppIocContainer) interface{} {
			return &abKeyService.ABKey{}
		},
	)

	// 账号信息
	c.Provide(
		appContainer.Account,
		func(container containerInterface.IAppIocContainer) interface{} {
			return &accountService.Account{}
		},
	)

	// AesKey
	c.Provide(
		appContainer.AesKey,
		func(container containerInterface.IAppIocContainer) interface{} {
			return &aesKeyService.AesKey{}
		},
	)

	// 联系人
	c.Provide(
		appContainer.Contact,
		func(container containerInterface.IAppIocContainer) interface{} {
			return &contactService.Contact{}
		},
	)

	// 设备
	c.Provide(
		appContainer.Device,
		func(container containerInterface.IAppIocContainer) interface{} {
			return &deviceService.Device{}
		},
	)

	// 设备列表数量
	c.Provide(
		appContainer.DeviceList,
		func(container containerInterface.IAppIocContainer) interface{} {
			return &deviceListService.DeviceList{}
		},
	)

	// 设备标识
	c.Provide(
		appContainer.Identity,
		func(container containerInterface.IAppIocContainer) interface{} {
			return &identityService.Identity{}
		},
	)

	// 群信息
	c.Provide(
		appContainer.Group,
		func(container containerInterface.IAppIocContainer) interface{} {
			return &groupService.Group{}
		},
	)

	// 812私信加解密key
	c.Provide(
		appContainer.Prekey,
		func(container containerInterface.IAppIocContainer) interface{} {
			return &preKeyService.PreKey{}
		},
	)

	// 推送设置
	c.Provide(
		appContainer.PushConfig,
		func(container containerInterface.IAppIocContainer) interface{} {
			return &pushConfigService.PushConfig{}
		},
	)

	// RegistrationToken
	c.Provide(
		appContainer.ReverToken,
		func(container containerInterface.IAppIocContainer) interface{} {
			return &registrationTokenService.RegistrationToken{}
		},
	)

	// 登录区域信息
	c.Provide(
		appContainer.RoutingInfo,
		func(container containerInterface.IAppIocContainer) interface{} {
			return &routingInfoService.RoutingInfo{}
		},
	)

	// SenderKey
	c.Provide(
		appContainer.SenderKey,
		func(container containerInterface.IAppIocContainer) interface{} {
			return &senderKeyService.SenderKey{}
		},
	)

	// 会话管理
	c.Provide(
		appContainer.Session,
		func(container containerInterface.IAppIocContainer) interface{} {
			return &sessionService.Session{}
		},
	)

	// SignedPreKey
	c.Provide(
		appContainer.SignedPreKey,
		func(container containerInterface.IAppIocContainer) interface{} {
			return &signedPreKeyService.SignedPreKey{}
		},
	)

	// 商业版数据管理
	c.Provide(
		appContainer.Business,
		func(container containerInterface.IAppIocContainer) interface{} {
			return &businessService.Business{}
		},
	)

	// 多媒体管理
	c.Provide(
		appContainer.MMS,
		func(container containerInterface.IAppIocContainer) interface{} {
			return &mmsService.MultimediaMessaging{}
		},
	)

	// 渠道0
	c.Provide(
		appContainer.Channel0EventCache,
		func(container containerInterface.IAppIocContainer) interface{} {
			return &eventService.Channel0EventCache{}
		},
	)

	// 渠道2
	c.Provide(
		appContainer.Channel2EventCache,
		func(container containerInterface.IAppIocContainer) interface{} {
			return &eventService.Channel2EventCache{}
		},
	)
}
