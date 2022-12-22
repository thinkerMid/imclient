package appContainer

import (
	hertzClient "github.com/cloudwego/hertz/pkg/app/client"
	"go.uber.org/zap"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	networkConstant "ws/framework/plugin/network/constant"
	functionTools "ws/framework/utils/function_tools"
)

const (
	//region 基础设施声明
	_infraStart containerInterface.ClassType = iota
	JID
	Logger
	SignalProtocol
	HttpClient
	XMPPMessageChannel
	HandshakeHandler
	ConnectionConfig
	Connection
	WhatsappConfiguration
	_infraEnd
	//endregion

	//region 数据服务声明
	_dataStorageServiceStart
	MemoryCache
	ABKey
	Account
	AesKey
	Contact
	Device
	DeviceList
	Identity
	Prekey
	PushConfig
	ReverToken
	RoutingInfo
	SenderKey
	Session
	SignedPreKey
	Channel0EventCache
	Channel2EventCache
	Group
	Business
	MMS
	_dataStorageServiceEnd
	//endregion

	IMControl

	// ----------------------------------------------
	_limitNumber
)

// ----------------------------------------------------------------------------

// API

// NewAppIocContainer .
func NewAppIocContainer() containerInterface.IAppIocContainer {
	return &AppIocContainer{
		store:   make([]interface{}, _limitNumber),
		provide: make([]containerInterface.ProvideFn, _limitNumber),
	}
}

// ----------------------------------------------------------------------------

// AppIocContainer .
type AppIocContainer struct {
	store   []interface{}
	provide []containerInterface.ProvideFn
}

// Provide 懒加载的函数
func (c *AppIocContainer) Provide(t containerInterface.ClassType, fn containerInterface.ProvideFn) {
	c.provide[t] = fn
}

// Inject 注入变量
func (c *AppIocContainer) Inject(t containerInterface.ClassType, instance interface{}) {
	c.store[t] = instance

	service, ok := instance.(containerInterface.IService)
	if ok {
		serviceName := functionTools.ReflectValueTypeName(instance)

		service.SetJID(c.ResolveJID())
		service.SetLogger(c.ResolveLogger().Named(serviceName))
		service.SetAppIocContainer(c)

		service.Init()
	}
}

func (c *AppIocContainer) initialize(t containerInterface.ClassType) interface{} {
	if c.store[t] == nil && c.provide[t] != nil {
		initializeFn := c.provide[t]
		c.provide[t] = nil

		instance := initializeFn(c)
		c.store[t] = instance

		service, ok := instance.(containerInterface.IService)
		if ok {
			serviceName := functionTools.ReflectValueTypeName(instance)

			service.SetJID(c.ResolveJID())
			service.SetLogger(c.ResolveLogger().Named(serviceName + "Service"))
			service.SetAppIocContainer(c)

			service.Init()
		}

		return instance
	}

	return c.store[t]
}

// OnJIDChangeWhenRegisterSuccess .
func (c *AppIocContainer) OnJIDChangeWhenRegisterSuccess(jid types.JID) {
	start := int(_dataStorageServiceStart)
	end := int(_dataStorageServiceEnd)

	for i := start; i < end; i++ {
		instance := c.initialize(containerInterface.ClassType(i))

		service, ok := instance.(containerInterface.IDataStorageService)
		if ok {
			service.OnJIDChangeWhenRegisterSuccess(jid)
			service.SetJID(jid)
		}
	}

	c.store[JID] = jid
}

// CleanupDataStorage .
func (c *AppIocContainer) CleanupDataStorage() {
	start := int(_dataStorageServiceStart)
	end := int(_dataStorageServiceEnd)

	for i := start; i < end; i++ {
		instance := c.initialize(containerInterface.ClassType(i))

		service, ok := instance.(containerInterface.IDataStorageService)
		if ok {
			service.CleanupAllData()
		}
	}

	c.ResolveLogger().Named("Application").Info("cleanup data storage")
}

// OnStart .
func (c *AppIocContainer) OnStart() {
	// 遍历实例化过的
	for i := range c.store {
		instance := c.store[i]
		if instance == nil {
			continue
		}

		service, ok := instance.(containerInterface.IService)
		if !ok {
			continue
		}

		service.OnApplicationStart()
	}
}

// OnResume .
func (c *AppIocContainer) OnResume() {
	// 遍历实例化过的
	for i := range c.store {
		instance := c.store[i]
		if instance == nil {
			continue
		}

		service, ok := instance.(containerInterface.IService)
		if !ok {
			continue
		}

		service.OnApplicationResume()
	}
}

// OnExit .
func (c *AppIocContainer) OnExit() {
	// 遍历实例化过的
	for i := range c.store {
		instance := c.store[i]
		if instance == nil {
			continue
		}

		service, ok := instance.(containerInterface.IService)
		if !ok {
			continue
		}

		service.OnApplicationExit()
	}
}

// ResolveAesKeyService .
func (c *AppIocContainer) ResolveAesKeyService() containerInterface.IAesKeyService {
	return c.initialize(AesKey).(containerInterface.IAesKeyService)
}

// ResolveDeviceListService .
func (c *AppIocContainer) ResolveDeviceListService() containerInterface.IDeviceListService {
	return c.initialize(DeviceList).(containerInterface.IDeviceListService)
}

// ResolveSignalProtocolFactory .
func (c *AppIocContainer) ResolveSignalProtocolFactory() containerInterface.ISignalProtocolService {
	return c.initialize(SignalProtocol).(containerInterface.ISignalProtocolService)
}

// ResolveSessionService .
func (c *AppIocContainer) ResolveSessionService() containerInterface.ISessionService {
	return c.initialize(Session).(containerInterface.ISessionService)
}

// ResolveIdentityService .
func (c *AppIocContainer) ResolveIdentityService() containerInterface.IIdentityService {
	return c.initialize(Identity).(containerInterface.IIdentityService)
}

// ResolveSenderKeyService .
func (c *AppIocContainer) ResolveSenderKeyService() containerInterface.ISenderKeyService {
	return c.initialize(SenderKey).(containerInterface.ISenderKeyService)
}

// ResolvePreKeyService .
func (c *AppIocContainer) ResolvePreKeyService() containerInterface.IPreKeyService {
	return c.initialize(Prekey).(containerInterface.IPreKeyService)
}

// ResolveSignedPreKeyService .
func (c *AppIocContainer) ResolveSignedPreKeyService() containerInterface.ISignedPreKeyService {
	return c.initialize(SignedPreKey).(containerInterface.ISignedPreKeyService)
}

// ResolvePushConfigService .
func (c *AppIocContainer) ResolvePushConfigService() containerInterface.IPushConfigService {
	return c.initialize(PushConfig).(containerInterface.IPushConfigService)
}

// ResolveRegistrationTokenService .
func (c *AppIocContainer) ResolveRegistrationTokenService() containerInterface.IRegistrationTokenService {
	return c.initialize(ReverToken).(containerInterface.IRegistrationTokenService)
}

// ResolveMemoryCache .
func (c *AppIocContainer) ResolveMemoryCache() containerInterface.IMemoryCacheService {
	return c.initialize(MemoryCache).(containerInterface.IMemoryCacheService)
}

// ResolveDeviceService .
func (c *AppIocContainer) ResolveDeviceService() containerInterface.IDeviceService {
	return c.initialize(Device).(containerInterface.IDeviceService)
}

// ResolveABKeyService .
func (c *AppIocContainer) ResolveABKeyService() containerInterface.IABKeyService {
	return c.initialize(ABKey).(containerInterface.IABKeyService)
}

// ResolveRoutingInfoService .
func (c *AppIocContainer) ResolveRoutingInfoService() containerInterface.IRoutingInfoService {
	return c.initialize(RoutingInfo).(containerInterface.IRoutingInfoService)
}

// ResolveJID .
func (c *AppIocContainer) ResolveJID() types.JID {
	return c.initialize(JID).(types.JID)
}

// ResolveChannel0EventCache .
func (c *AppIocContainer) ResolveChannel0EventCache() containerInterface.IEventCache {
	return c.initialize(Channel0EventCache).(containerInterface.IEventCache)
}

// ResolveChannel2EventCache .
func (c *AppIocContainer) ResolveChannel2EventCache() containerInterface.IEventCache {
	return c.initialize(Channel2EventCache).(containerInterface.IEventCache)
}

// ResolveAccountService .
func (c *AppIocContainer) ResolveAccountService() containerInterface.IAccountService {
	return c.initialize(Account).(containerInterface.IAccountService)
}

// ResolveGroupService .
func (c *AppIocContainer) ResolveGroupService() containerInterface.IGroupService {
	return c.initialize(Group).(containerInterface.IGroupService)
}

// ResolveConnection .
func (c *AppIocContainer) ResolveConnection() containerInterface.IConnection {
	return c.initialize(Connection).(containerInterface.IConnection)
}

// ResolveMessageChannel .
func (c *AppIocContainer) ResolveMessageChannel() containerInterface.IMessageChannel {
	return c.initialize(XMPPMessageChannel).(containerInterface.IMessageChannel)
}

// ResolveLogger .
func (c *AppIocContainer) ResolveLogger() *zap.SugaredLogger {
	return c.initialize(Logger).(*zap.SugaredLogger)
}

// ResolveConnectionConfig .
func (c *AppIocContainer) ResolveConnectionConfig() networkConstant.ConnectionConfig {
	return c.initialize(ConnectionConfig).(networkConstant.ConnectionConfig)
}

// ResolveHandshakeHandler .
func (c *AppIocContainer) ResolveHandshakeHandler() containerInterface.IHandshakeHandler {
	return c.initialize(HandshakeHandler).(containerInterface.IHandshakeHandler)
}

// ResolveHttpClient .
func (c *AppIocContainer) ResolveHttpClient() *hertzClient.Client {
	return c.initialize(HttpClient).(*hertzClient.Client)
}

// ResolveContactService .
func (c *AppIocContainer) ResolveContactService() containerInterface.IContactService {
	return c.initialize(Contact).(containerInterface.IContactService)
}

// ResolveWhatsappConfiguration .
func (c *AppIocContainer) ResolveWhatsappConfiguration() *containerInterface.WhatsappConfiguration {
	return c.initialize(WhatsappConfiguration).(*containerInterface.WhatsappConfiguration)
}

// ResolveBusinessService .
func (c *AppIocContainer) ResolveBusinessService() containerInterface.IBusinessService {
	return c.initialize(Business).(containerInterface.IBusinessService)
}

// ResolveMultimediaMessagingService .
func (c *AppIocContainer) ResolveMultimediaMessagingService() containerInterface.IMultimediaMessagingService {
	return c.initialize(MMS).(containerInterface.IMultimediaMessagingService)
}

// ResolveIMControl .
func (c *AppIocContainer) ResolveIMControl() containerInterface.IIMControl {
	return c.initialize(IMControl).(containerInterface.IIMControl)
}
