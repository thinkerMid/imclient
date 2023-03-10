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
	// EncryptGroupMessage ?????????????????? skmsg
	EncryptGroupMessage(groupID string, plaintext []byte) ([]byte, []byte, error)
	// DecryptGroupSenderKey ???????????? senderkey
	DecryptGroupSenderKey(senderKeyName *protocol.SenderKeyName, body []byte) error
	// DecryptGroupMessage ??????????????????
	DecryptGroupMessage(senderKeyName *protocol.SenderKeyName, body []byte) ([]byte, error)
	// EncryptPrivateChatMessage ??????????????????
	EncryptPrivateChatMessage(dstJID types.JID, plaintext []byte) (protocol.CiphertextMessage, error)
	// DecryptPrivateChatMessage ??????????????????
	DecryptPrivateChatMessage(dstJID types.JID, plaintext []byte, pkmsg bool) ([]byte, error)
	// CreateGroupSession ?????????????????????????????????
	CreateGroupSession(groupID string)
}

// IVirtualDevice .
type IVirtualDevice interface {
	// GetPhoneNumber ?????????
	GetPhoneNumber() string
	// GetCC ????????????
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

	BatchCreateDevice(groupID string, deviceID []types.JID)                              // ?????????????????????
	SearchSenderInGroupAndCreate(deviceID types.JID)                                     // ??????JID????????????????????????????????????
	SaveSentMessageByGroupID(groupID string)                                             // ??????????????????????????????????????????
	FindDevice(senderKeyName *protocol.SenderKeyName) *senderKeyDB.SenderDevice          // ???????????????
	FindUnSentMessageDeviceByGroupID(groupID string) ([]senderKeyDB.SenderDevice, error) // ????????????????????????????????????
	BatchDeleteSenderByGroupIDAndJID(groupID string, senderJID []string)                 // ???????????????????????????
	DeleteAllDeviceByGroupID(groupID string)                                             // ??????????????????????????????
	DeleteDevice(deviceID types.JID)                                                     // ??????????????????????????????
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
	// AddDevice ????????????
	AddDevice(dstJID string, deviceID uint8)
	DeleteDevice(dstJID string, deviceID uint8)
	BatchCreateDevice(dst string, deviceIDList []uint8)
	BatchDeleteDevice(dst string, deviceIDList []uint8)
	FindDeviceIDList(dstJID string) []uint8
	FindUnInitSessionDeviceIDList(dstJID string) (idList []uint8)
	HaveMultiDevice(dstJID string) bool
	UpdateDeviceList(dstJID string, deviceIDList []uint8)

	// CreateSession ????????????
	CreateSession(address *protocol.SignalAddress, record *record.Session)
	FindSession(address *protocol.SignalAddress) (*record.Session, error)
	SaveSession(remoteAddress *protocol.SignalAddress, record *record.Session)
	SaveEncryptSession(remoteAddress *protocol.SignalAddress, record *record.Session) // ??????????????????????????????
	SaveDecryptSession(remoteAddress *protocol.SignalAddress, record *record.Session) // ??????????????????????????????
	SaveRebuildSession(remoteAddress *protocol.SignalAddress, record *record.Session) // ??????????????????????????????
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
	// AddEvent ????????????
	AddEvent(WaEvent)
	// ClearLog ??????????????????????????????
	ClearLog()
	// ClearNotSentYetLog ??????????????????????????????
	ClearNotSentYetLog()
	// PackNotSentYetLog ??????????????????????????????
	PackNotSentYetLog(sendEventCount int32, buffer eventSerialize.IEventBuffer)
	// PackBuffer ????????????????????????????????????
	PackBuffer(sendEventCount int32, buffer eventSerialize.IEventBuffer)
	// CacheBufferItem ????????????Buffer??????
	CacheBufferItem() int64
	// ResetAddLogState ??????????????????
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
