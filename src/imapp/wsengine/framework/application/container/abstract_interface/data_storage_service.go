package containerInterface

import (
	"time"
	"ws/framework/application/constant/binary/proto"
	"ws/framework/application/constant/types"
	"ws/framework/application/data_storage/ab_key/database"
	"ws/framework/application/data_storage/account/database"
	"ws/framework/application/data_storage/aes_key/database"
	businessDB "ws/framework/application/data_storage/business/database"
	"ws/framework/application/data_storage/cache/constant"
	"ws/framework/application/data_storage/contact/database"
	"ws/framework/application/data_storage/device/database"
	"ws/framework/application/data_storage/group/database"
	mmsConstant "ws/framework/application/data_storage/mms/constant"
	"ws/framework/application/data_storage/prekey/database"
	"ws/framework/application/data_storage/push_config/database"
	"ws/framework/application/data_storage/registration_token/database"
	"ws/framework/application/data_storage/routing_info/database"
	"ws/framework/application/data_storage/sender_key/database"
	"ws/framework/application/data_storage/signed_prekey/database"
	groupRecord "ws/framework/application/libsignal/groups/state/record"
	"ws/framework/application/libsignal/keys/identity"
	"ws/framework/application/libsignal/protocol"
	"ws/framework/application/libsignal/state/record"
	"ws/framework/application/libsignal/state/store"
	mediaCrypto "ws/framework/lib/media_crypto"
	eventSerialize "ws/framework/plugin/event_serialize"
	"ws/framework/utils/keys"
)

// IMemoryCacheService .
type IMemoryCacheService interface {
	Cache(key, value interface{})
	CacheTTL(key, value interface{}, ttl time.Duration)
	UnCache(key interface{})
	FindInCache(key interface{}) (interface{}, bool)
	AccountLoginData() *cacheConstant.AccountLoginData
}

// IMultimediaMessagingService .
type IMultimediaMessagingService interface {
	UpdateMMSEndpoints(v mmsConstant.MMSEndpoints)
	UploadMediaFile(file mediaCrypto.File) (url string, path string, err error)
}

// ISignalProtocolService .
type ISignalProtocolService interface {
	Context() store.SignalProtocol
	// EncryptGroupMessage 加密群组消息 skmsg
	EncryptGroupMessage(groupID string, plaintext []byte) ([]byte, []byte, error)
	// DecryptGroupSenderKey 解密群组 senderkey
	DecryptGroupSenderKey(senderKeyName *protocol.SenderKeyName, body []byte) error
	// DecryptGroupMessage 解密群组消息
	DecryptGroupMessage(senderKeyName *protocol.SenderKeyName, body []byte) ([]byte, error)
	// EncryptPrivateChatMessage 加密私信消息
	EncryptPrivateChatMessage(dstJID types.JID, plaintext []byte) (protocol.CiphertextMessage, error)
	// DecryptPrivateChatMessage 解密私信消息
	DecryptPrivateChatMessage(dstJID types.JID, plaintext []byte, pkmsg bool) ([]byte, error)
	// CreateGroupSession 创建自己的群组消息密钥
	CreateGroupSession(groupID string)
}

// IVirtualDevice .
type IVirtualDevice interface {
	// GetPhoneNumber 手机号
	GetPhoneNumber() string
	// GetCC 国际区号
	GetCC() string
	// GetMCC .
	GetMCC() string
	// GetMNC .
	GetMNC() string
	// GetISO .
	GetISO() string
	// GetLanguage .
	GetLanguage() string
	// GetOSVersion .
	GetOSVersion() string
	// GetManufacturer .
	GetManufacturer() string
	// GetProduction .
	GetProduction() string
	// GetBuildNumber .
	GetBuildNumber() string
}

// IDeviceService .
type IDeviceService interface {
	Context() *deviceDB.Device
	Create(virtualDevice IVirtualDevice) (*deviceDB.Device, error)
	Import(*deviceDB.Device) error
	GetClientPayload() *waProto.ClientPayload
	DeviceAgent() string
	PrivateStatsAgent() string
	ContextExecute(f func(*deviceDB.Device))
}

// IContactService .
type IContactService interface {
	CreateContact(dstJID string, aliasPhoneNumber string)
	DeleteByJID(dstJID string) error
	BatchCreateAddressBookContactByJIDList(jidList []string, aliasPhoneNumber []string)
	CreateAddressBookContactByJID(dstJID string, aliasPhoneNumber string) error
	FindByPhoneNumber(phoneNumber string) *contactDB.Contact
	FindByJID(dstJID string) *contactDB.Contact
	ContextExecute(dstJID string, f func(*contactDB.Contact))
}

// IRoutingInfoService .
type IRoutingInfoService interface {
	Context() *routingInfoDB.RoutingInfo
	Create(content []byte) error
	Save(content []byte)
}

// IABKeyService .
type IABKeyService interface {
	Context() *abKeyDB.QueryResult
	Create(content string) error
}

// IPushConfigService .
type IPushConfigService interface {
	Context() *pushConfigDB.PushConfig
	Create() (*pushConfigDB.PushConfig, error)
	Import(*pushConfigDB.PushConfig) error
	ContextExecute(f func(*pushConfigDB.PushConfig))
}

// IRegistrationTokenService .
type IRegistrationTokenService interface {
	Context() *registrationTokenDB.QueryResult
	Create() (*registrationTokenDB.RegistrationToken, error)
	Import(*registrationTokenDB.RegistrationToken) error
	RefreshToken() error
}

// ISessionService
//
//	Deprecated replace with IDeviceListService
type ISessionService interface{}

// IIdentityService .
type IIdentityService interface {
	Context() *keys.KeyPair
	GetIdentityKeyPair() *identity.KeyPair
	GetLocalRegistrationId() uint32
}

// ISenderKeyService .
type ISenderKeyService interface {
	CreateSenderKey(senderKeyName *protocol.SenderKeyName, keyRecord *groupRecord.SenderKey)
	UpdateSenderKey(senderKeyName *protocol.SenderKeyName, keyRecord *groupRecord.SenderKey)
	ResetSenderKey(senderKeyName *protocol.SenderKeyName, keyRecord *groupRecord.SenderKey)
	ContainsSenderKey(senderKeyName *protocol.SenderKeyName) bool
	FindSenderKey(senderKeyName *protocol.SenderKeyName) (*groupRecord.SenderKey, error)

	BatchCreateDevice(groupID string, deviceID []types.JID)                              // 批量创建群设备
	SearchSenderInGroupAndCreate(deviceID types.JID)                                     // 查询JID所在的群组批量创建群设备
	SaveSentMessageByGroupID(groupID string)                                             // 设置群组里的设备都已发过消息
	FindDevice(senderKeyName *protocol.SenderKeyName) *senderKeyDB.SenderDevice          // 查询群设备
	FindUnSentMessageDeviceByGroupID(groupID string) ([]senderKeyDB.SenderDevice, error) // 查找未发送过消息的群设备
	BatchDeleteSenderByGroupIDAndJID(groupID string, senderJID []string)                 // 批量删除群内群设备
	DeleteAllDeviceByGroupID(groupID string)                                             // 删除群组内所有群设备
	DeleteDevice(deviceID types.JID)                                                     // 从所有群内删除该设备
}

// IPreKeyService .
type IPreKeyService interface {
	StatementPreKeyCount() int
	SavePreKey(preKeyID uint32, preKeyRecord *record.PreKey)
	ContainsPreKey(preKeyID uint32) bool
	DeletePreKey(_ uint32)
	FindPreKey(preKeyID uint32) *record.PreKey
	InitPreKeys() ([]keys.PreKey, error)
	GeneratePreKeys(count int) ([]keys.PreKey, error)
	Import([]prekeyDB.PreKey) error
}

// ISignedPreKeyService .
type ISignedPreKeyService interface {
	Context() *record.SignedPreKey
	Create() (*signedPreKeyDB.SignedPreKey, error)
	FindSignedPreKey(signedPreKeyID uint32) *record.SignedPreKey
	SaveSignedPreKeyBuffer(signedPreKeyID uint32, signedPreKeyJsonBuffer []byte)
}

// IAccountService .
type IAccountService interface {
	Context() *accountDB.Account
	Create() (*accountDB.Account, error)
	NeedUploadRecordChannel0Event() bool
	NeedUploadRecordChannel2Event() bool
	ContextExecute(f func(*accountDB.Account))
	Import(*accountDB.Account) error
}

// IDeviceListService .
type IDeviceListService interface {
	// AddDevice 设备相关
	AddDevice(dstJID string, deviceID uint8)
	DeleteDevice(dstJID string, deviceID uint8)
	BatchCreateDevice(dst string, deviceIDList []uint8)
	BatchDeleteDevice(dst string, deviceIDList []uint8)
	FindDeviceIDList(dstJID string) []uint8
	FindUnInitSessionDeviceIDList(dstJID string) (idList []uint8)
	HaveMultiDevice(dstJID string) bool
	UpdateDeviceList(dstJID string, deviceIDList []uint8)

	// CreateSession 会话相关
	CreateSession(address *protocol.SignalAddress, record *record.Session)
	FindSession(address *protocol.SignalAddress) (*record.Session, error)
	SaveSession(remoteAddress *protocol.SignalAddress, record *record.Session)
	SaveEncryptSession(remoteAddress *protocol.SignalAddress, record *record.Session) // 不同场景下的保存操作
	SaveDecryptSession(remoteAddress *protocol.SignalAddress, record *record.Session) // 不同场景下的保存操作
	SaveRebuildSession(remoteAddress *protocol.SignalAddress, record *record.Session) // 不同场景下的保存操作
	ContainsSession(remoteAddress *protocol.SignalAddress) bool
	DeleteSession(remoteAddress *protocol.SignalAddress)
}

// IAesKeyService .
type IAesKeyService interface {
	Context() *aesKeyDB.QueryResult
	Create() (*aesKeyDB.AESKey, error)
}

// IGroupService .
type IGroupService interface {
	CreateGroup(groupID string, isAdmin bool)
	DeleteGroup(groupID string)
	Find(groupID string) *groupDB.Group
	ContextExecute(groupID string, f func(group *groupDB.Group))
}

// WaEvent .
type WaEvent interface {
	Init(uint8, int64, float64)
	InitFields(interface{})
	Serialize(buffer eventSerialize.IEventBuffer)
}

// IEventCache .
type IEventCache interface {
	// AddEvent 添加日志
	AddEvent(WaEvent)
	// ClearLog 清理之前和当前的日志
	ClearLog()
	// ClearNotSentYetLog 清理之前未发送的日志
	ClearNotSentYetLog()
	// PackNotSentYetLog 打包之前未发送的日志
	PackNotSentYetLog(sendEventCount int32, buffer eventSerialize.IEventBuffer)
	// PackBuffer 打包之前和当前的日志内容
	PackBuffer(sendEventCount int32, buffer eventSerialize.IEventBuffer)
	// CacheBufferItem 已储存的Buffer数量
	CacheBufferItem() int64
	// ResetAddLogState 重置添加状态
	ResetAddLogState()
	// FlushEventCache .
	FlushEventCache()
	// CleanupAllData .
	CleanupAllData()
}

// IBusinessService .
type IBusinessService interface {
	Create() (*businessDB.BusinessProfile, error)
	Context() *businessDB.BusinessProfile
	GenerateBusinessVerifiedName(appendPushName bool) []byte
	ContextExecute(f func(*businessDB.BusinessProfile))
}
