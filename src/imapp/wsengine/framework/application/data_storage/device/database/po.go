package deviceDB

import (
	"ws/framework/plugin/database/database_tools"
)

// Device .
type Device struct {
	databaseTools.ChangeExtension

	JID                string `gorm:"column:jid;primaryKey"`
	RegistrationId     uint32 `gorm:"column:registrationId"`
	IdentityKey        []byte `gorm:"column:identityKey"`
	PushName           string `gorm:"column:pushName"`
	UserName           string `gorm:"column:userName"`
	FBUuid             string `gorm:"column:fbuuid"`
	Uuid               string `gorm:"column:uuid"`
	FBUuidCreateTime   int32  `gorm:"column:fbuuidCreateTime"`
	OsVersion          string `gorm:"column:osVersion"`
	Mcc                string `gorm:"column:mcc"`
	Mnc                string `gorm:"column:mnc"`
	Manufacturer       string `gorm:"column:manufacturer"`
	Device             string `gorm:"column:device"`
	Language           string `gorm:"column:language"`
	Country            string `gorm:"column:country"`
	BuildNumber        string `gorm:"column:buildNumber"`
	PrivateStatsId     string `gorm:"column:privateStatsId"`
	Area               string `gorm:"column:area"`
	Phone              string `gorm:"column:phone"`
	SecurityCodeSet    int8   `gorm:"column:securityCodeSet"`
	ClientStaticPriKey []byte `gorm:"column:clientStaticPriKey"`
	ClientStaticPubKey []byte `gorm:"column:clientStaticPubKey"`
	ServerStaticKey    []byte `gorm:"column:serverStaticKey"`
	BusinessName       string `gorm:"column:businessName"`
	Platform           string `gorm:"column:platform"`
}

// TableName .
func (m *Device) TableName() string {
	return "deviceinfo"
}

// UpdatePushName .
func (m *Device) UpdatePushName(name string) {
	if name == m.PushName {
		return
	}

	m.PushName = name
	m.Update("pushName", name)
}

// UpdateUserName .
func (m *Device) UpdateUserName(name string) {
	if name == m.UserName {
		return
	}

	m.UserName = name
	m.Update("userName", name)
}

// UpdateJID .
func (m *Device) UpdateJID(name string) {
	if name == m.JID {
		return
	}

	m.Update("jid", name)
}

// UpdateServerStaticKey .
func (m *Device) UpdateServerStaticKey(v []byte) {
	m.ServerStaticKey = v
	m.Update("serverStaticKey", v)
}
