package models

import (
	"github.com/golang/protobuf/proto"
	"labs/src/imapp/wsa/types"
	waProto "labs/src/imapp/wsa/types/binary/proto"
	"labs/src/lib/nickname"
	"labs/utils"
)

func CreateClientPayLoad(jid types.JID, dev Device) waProto.ClientPayload {
	return waProto.ClientPayload{
		Username: proto.Uint64(jid.UserInt()),
		Passive:  proto.Bool(false), // XMPP心跳设置  false:服务器发起 true:客户端发起,
		UserAgent: &waProto.UserAgent{
			Platform: waProto.UserAgent_SMB_IOS.Enum(),
			AppVersion: &waProto.AppVersion{
				Primary:    proto.Uint32(0x2),
				Secondary:  proto.Uint32(0x16),
				Tertiary:   proto.Uint32(0x15),
				Quaternary: proto.Uint32(0x4D),
			},
			Mcc:                         proto.String(dev.Mcc),
			Mnc:                         proto.String(dev.Mnc),
			OsVersion:                   proto.String(dev.OsVersion),
			Manufacturer:                proto.String(dev.Manufacturer),
			Device:                      proto.String(dev.Device),
			OsBuildNumber:               proto.String(dev.BuildNumber),
			PhoneId:                     proto.String(dev.FBUuid),
			ReleaseChannel:              waProto.UserAgent_RELEASE.Enum(),
			LocaleLanguageIso6391:       proto.String(dev.Language),
			LocaleCountryIso31661Alpha2: proto.String(dev.Country),
		},
		PushName:      proto.String(nickname.New()),
		SessionId:     proto.Int32(int32(utils.RandInt64(1, 1000000000))),
		ShortConnect:  proto.Bool(true),
		ConnectType:   waProto.ClientPayload_CELLULAR_LTE.Enum(),
		ConnectReason: waProto.ClientPayload_USER_ACTIVATED.Enum(),
		DnsSource: &waProto.DNSSource{
			DnsMethod: waProto.DNSSource_SYSTEM.Enum(),
			AppCached: proto.Bool(false),
		},
		ConnectAttemptCount: proto.Uint32(0),
		Device:              proto.Uint32(0),
		Oc:                  proto.Bool(false),
	}
}
