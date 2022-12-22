package containerInterface

import (
	hertzClient "github.com/cloudwego/hertz/pkg/app/client"
	"go.uber.org/zap"
	"ws/framework/application/constant/types"
	networkConstant "ws/framework/plugin/network/constant"
)

// ProvideFn .
type ProvideFn func(IAppIocContainer) interface{}

// ClassType .
type ClassType uint8

// IAppIocContainer .
type IAppIocContainer interface {
	Inject(t ClassType, i interface{})
	Provide(t ClassType, fn ProvideFn)

	OnStart()
	OnResume()
	OnExit()

	OnJIDChangeWhenRegisterSuccess(types.JID)
	CleanupDataStorage()

	//region 基础组件

	// ResolveJID 全局上下文的号码
	ResolveJID() types.JID
	// ResolveLogger 日志
	ResolveLogger() *zap.SugaredLogger
	// ResolveSignalProtocolFactory 协议处理
	ResolveSignalProtocolFactory() ISignalProtocolService
	// ResolveConnectionConfig 连接设置
	ResolveConnectionConfig() networkConstant.ConnectionConfig
	// ResolveConnection TCP连接
	ResolveConnection() IConnection
	// ResolveHandshakeHandler 握手
	ResolveHandshakeHandler() IHandshakeHandler
	// ResolveMessageChannel XMPP消息处理
	ResolveMessageChannel() IMessageChannel
	// ResolveHttpClient HTTP客户端
	ResolveHttpClient() *hertzClient.Client
	// ResolveIMControl .
	ResolveIMControl() IIMControl

	//endregion

	//region 数据服务

	// ResolveDeviceService 设备
	ResolveDeviceService() IDeviceService
	// ResolveContactService 联系人
	ResolveContactService() IContactService
	// ResolveAccountService 账号
	ResolveAccountService() IAccountService
	// ResolveRoutingInfoService 区域标识
	ResolveRoutingInfoService() IRoutingInfoService
	// ResolveABKeyService ABKey
	ResolveABKeyService() IABKeyService
	// ResolvePushConfigService 推送设置
	ResolvePushConfigService() IPushConfigService
	// ResolveRegistrationTokenService 注册令牌
	ResolveRegistrationTokenService() IRegistrationTokenService
	// ResolveSessionService 聊天会话
	ResolveSessionService() ISessionService
	// ResolveIdentityService 设备身份
	ResolveIdentityService() IIdentityService
	// ResolveSenderKeyService 群组设备管理
	ResolveSenderKeyService() ISenderKeyService
	// ResolvePreKeyService 812个PreKey密钥
	ResolvePreKeyService() IPreKeyService
	// ResolveSignedPreKeyService SignedPreKey
	ResolveSignedPreKeyService() ISignedPreKeyService
	// ResolveDeviceListService 设备列表
	ResolveDeviceListService() IDeviceListService
	// ResolveAesKeyService AesKey
	ResolveAesKeyService() IAesKeyService
	// ResolveGroupService 群组信息管理
	ResolveGroupService() IGroupService
	// ResolveMemoryCache 数据缓存
	ResolveMemoryCache() IMemoryCacheService
	// ResolveWhatsappConfiguration whatsapp的配置信息
	ResolveWhatsappConfiguration() *WhatsappConfiguration
	// ResolveBusinessService 商业版数据管理
	ResolveBusinessService() IBusinessService
	// ResolveMultimediaMessagingService 多媒体管理
	ResolveMultimediaMessagingService() IMultimediaMessagingService
	//endregion

	//region 事件日志

	// ResolveChannel0EventCache 渠道0日志
	ResolveChannel0EventCache() IEventCache
	// ResolveChannel2EventCache 渠道2日志
	ResolveChannel2EventCache() IEventCache

	//endregion
}
