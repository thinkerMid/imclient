package deviceService

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"strings"
	"ws/framework/application/constant/binary/proto"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	deviceDB "ws/framework/application/data_storage/device/database"
	"ws/framework/application/libsignal/ecc"
	"ws/framework/application/libsignal/util/keyhelper"
	"ws/framework/lib/apple"
	"ws/framework/lib/nickname"
	"ws/framework/plugin/database"
	"ws/framework/plugin/database/database_tools"
	"ws/framework/utils"
)

var _ containerInterface.IDeviceService = &Device{}

// Device .
type Device struct {
	containerInterface.BaseService
	// 比较经常用 做个缓存
	context *deviceDB.Device
}

// Context .
func (d *Device) Context() *deviceDB.Device {
	if d.context != nil {
		return d.context
	}

	context := deviceDB.Device{JID: d.JID.User}

	err := databaseTools.FindByPrimaryKey(database.MasterDB(), &context)
	if err != nil {
		// TODO 兼容旧数据 一个月后或者一个版本后可以删除了
		context.JID = d.JID.String()
		err = databaseTools.FindByPrimaryKey(database.MasterDB(), &context)
	}

	if err != nil {
		d.Logger.Error(err)
		return nil
	}

	d.context = &context

	return d.context
}

// Create .
func (d *Device) Create(virtualDevice containerInterface.IVirtualDevice) (*deviceDB.Device, error) {
	identityKeyPair, err := keyhelper.GenerateIdentityKeyPair()
	if err != nil {
		return nil, err
	}

	number := virtualDevice.GetCC() + virtualDevice.GetPhoneNumber()

	context := deviceDB.Device{JID: number, UserName: number}

	djbECKeyPair, _ := ecc.GenerateKeyPair()
	clientStaticPriKey := djbECKeyPair.PrivateKey().Serialize()
	clientStaticPubKey := djbECKeyPair.PublicKey().PublicKey()

	context.Area = virtualDevice.GetCC()
	context.Phone = virtualDevice.GetPhoneNumber()
	context.Mnc = virtualDevice.GetMNC()
	context.Mcc = virtualDevice.GetMCC()
	context.PushName = nickname.New()
	context.ClientStaticPriKey = clientStaticPriKey[:]
	context.ClientStaticPubKey = clientStaticPubKey[:]
	context.FBUuid = utils.GenUUID4()
	context.FBUuidCreateTime = int32(utils.GetCurTime())
	context.Uuid = utils.GenUUID4()
	context.ServerStaticKey = make([]byte, 0)
	context.OsVersion = virtualDevice.GetOSVersion()
	context.Manufacturer = virtualDevice.GetManufacturer()
	context.Device = virtualDevice.GetProduction()
	context.Language = virtualDevice.GetLanguage()
	context.Country = virtualDevice.GetISO()
	context.BuildNumber = virtualDevice.GetBuildNumber()
	context.PrivateStatsId = utils.GenUUID4()

	identityPriKey := identityKeyPair.PrivateKey().Serialize()

	context.RegistrationId = keyhelper.GenerateRegistrationID()
	context.IdentityKey = identityPriKey[:]

	_, err = databaseTools.Create(database.MasterDB(), &context)
	if err != nil {
		d.Logger.Error(err)
		return nil, err
	}

	d.context = &context

	return d.context, nil
}

// Import .
func (d *Device) Import(device *deviceDB.Device) error {
	_, err := databaseTools.Create(database.MasterDB(), device)
	if err != nil {
		d.Logger.Error(err)
		return err
	}

	d.context = device

	return nil
}

// GetClientPayload .
func (d *Device) GetClientPayload() *waProto.ClientPayload {
	/**
	 username:6281334145057  passive:true  userAgent:{platform:IOS
	appVersion:{primary:2  secondary:21  tertiary:243  quaternary:1}
	mcc:"510"  mnc:"010"  osVersion:"14.6"  manufacturer:"Apple"
	device:"iPhone 6 Plus"  osBuildNumber:"16G183"  phoneId:"DE24EB5F-0EA7-4DFD-AA88-582980577E3B"
	releaseChannel:RELEASE  localeLanguageIso6391:"zh"  localeCountryIso31661Alpha2:"CN"}
	sessionId:134020439  shortConnect:true  connectType:WIFI_UNKNOWN  connectReason:USER_ACTIVATED
	dnsSource:{dnsMethod:SYSTEM  appCached:false}  connectAttemptCount:0  device:0
	*/

	context := d.Context()
	configuration := d.AppIocContainer.ResolveWhatsappConfiguration()

	// TODO(try) 版本更新需要核对这些成员是否修改
	payload := &waProto.ClientPayload{}
	payload.Username = proto.Uint64(d.JID.UserInt())
	payload.Passive = proto.Bool(false) // XMPP心跳设置  false:服务器主动ping true:客户端主动ping

	payload.UserAgent = &waProto.UserAgent{
		Platform: configuration.Platform.Enum(),
		AppVersion: &waProto.AppVersion{
			Primary:    &configuration.VersionCode[0],
			Secondary:  &configuration.VersionCode[1],
			Tertiary:   &configuration.VersionCode[2],
			Quaternary: &configuration.VersionCode[3],
		},
		Mcc:                         proto.String(context.Mcc),
		Mnc:                         proto.String(context.Mnc),
		OsVersion:                   proto.String(context.OsVersion),
		Manufacturer:                proto.String(context.Manufacturer),
		Device:                      proto.String(context.Device),
		OsBuildNumber:               proto.String(context.BuildNumber),
		PhoneId:                     proto.String(context.FBUuid),
		ReleaseChannel:              waProto.UserAgent_RELEASE.Enum(),
		LocaleLanguageIso6391:       proto.String(context.Language),
		LocaleCountryIso31661Alpha2: proto.String(context.Country),
	}

	if len(context.PushName) > 0 {
		payload.PushName = proto.String(context.PushName)
	}

	payload.SessionId = proto.Int32(int32(utils.RandInt64(1, 1000000000)))
	payload.ShortConnect = proto.Bool(true)
	payload.ConnectType = waProto.ClientPayload_CELLULAR_LTE.Enum()
	payload.ConnectReason = waProto.ClientPayload_USER_ACTIVATED.Enum()
	payload.DnsSource = &waProto.DNSSource{
		DnsMethod: waProto.DNSSource_SYSTEM.Enum(),
		AppCached: proto.Bool(false),
	}
	payload.ConnectAttemptCount = proto.Uint32(0)
	payload.Device = proto.Uint32(0)

	return payload
}

// DeviceAgent .
func (d *Device) DeviceAgent() string {
	configuration := d.AppIocContainer.ResolveWhatsappConfiguration()
	device := strings.ReplaceAll(d.Context().Device, " ", "_")

	return fmt.Sprintf(configuration.UserAgent, configuration.VersionString, d.Context().OsVersion, device)
}

// PrivateStatsAgent .
func (d *Device) PrivateStatsAgent() string {
	configuration := d.AppIocContainer.ResolveWhatsappConfiguration()
	cfnetwork, darwin := apple.DarwinSystemManagerInstance().GetCFNetworkAndDarwinVersion(d.Context().OsVersion)

	return fmt.Sprintf("WhatsApp/%s CFNetwork/%s Darwin/%s", configuration.VersionString, cfnetwork, darwin)
}

// ContextExecute .
func (d *Device) ContextExecute(f func(*deviceDB.Device)) {
	f(d.Context())

	_, err := databaseTools.Save(database.MasterDB(), d.Context())
	if err != nil {
		d.Logger.Errorf("update context error: %v", err)
	}
}

// OnJIDChangeWhenRegisterSuccess .
func (d *Device) OnJIDChangeWhenRegisterSuccess(newJID types.JID) {
	context := deviceDB.Device{JID: d.JID.User}
	context.UpdateJID(newJID.User)

	_, err := databaseTools.Save(database.MasterDB(), &context)
	if err != nil {
		d.Logger.Error(err)
	}
}

// CleanupAllData 只修改用户名不做删除操作
func (d *Device) CleanupAllData() {
	context := deviceDB.Device{JID: d.JID.User}
	context.UpdateJID("unavailable")

	_, _ = databaseTools.Save(database.MasterDB(), &context)
}
