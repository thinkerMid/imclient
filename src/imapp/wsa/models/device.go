package models

import (
	"fmt"
	"labs/src/imapp/wsa/types/ecc"
	keyhelper "labs/src/imapp/wsa/types/helpers"
	"labs/src/lib/firmware"
	"labs/src/lib/msisdn"
	"labs/utils"
	"strings"
)

type Device struct {
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
	AdvSecretKey       []byte `gorm:"column:advSecretKey"`
	BusinessName       string `gorm:"column:businessName"`
	Platform           string `gorm:"column:platform"`
}

func (dev *Device) Agent(version string) string {
	name := strings.ReplaceAll(dev.OsVersion, " ", "_")
	return fmt.Sprintf("WhatsApp/%s iOS/%v Device/%v", version, name, dev.Device)
}

type VirtualDevice struct {
	msisdn.IMSIParser
	firmware.Apple
}

func CreateDevice(phoneNumber string) Device {
	imsi, _ := msisdn.Parse(phoneNumber, true)

	virtual := VirtualDevice{
		IMSIParser: imsi,
		Apple:      firmware.NewAppleFirmware(),
	}

	number := virtual.GetCC() + virtual.GetPhoneNumber()

	djbECKeyPair, _ := ecc.GenerateKeyPair()
	clientStaticPriKey := djbECKeyPair.PrivateKey().Serialize()
	clientStaticPubKey := djbECKeyPair.PublicKey().PublicKey()

	identityKeyPair, _ := keyhelper.GenerateIdentityKeyPair()
	identityPriKey := identityKeyPair.PrivateKey().Serialize()

	return Device{
		JID:                number,
		RegistrationId:     keyhelper.GenerateRegistrationID(),
		IdentityKey:        identityPriKey[:],
		PushName:           utils.RandNickName(),
		UserName:           number,
		FBUuid:             utils.GenUUID4(),
		Uuid:               utils.GenUUID4(),
		FBUuidCreateTime:   int32(utils.GetCurTime()),
		OsVersion:          virtual.GetOSVersion(),
		Mcc:                virtual.GetMCC(),
		Mnc:                virtual.GetMNC(),
		Manufacturer:       virtual.GetManufacturer(),
		Device:             virtual.GetProduction(),
		Language:           virtual.GetLanguage(),
		Country:            virtual.GetISO(),
		BuildNumber:        virtual.GetBuildNumber(),
		PrivateStatsId:     utils.GenUUID4(),
		Area:               virtual.GetCC(),
		Phone:              virtual.GetPhoneNumber(),
		SecurityCodeSet:    0,
		ClientStaticPriKey: clientStaticPriKey[:],
		ClientStaticPubKey: clientStaticPubKey[:],
		AdvSecretKey:       utils.RandBytes(32),
		BusinessName:       "",
		Platform:           "",
	}
}
