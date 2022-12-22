package importData

import (
	"google.golang.org/protobuf/proto"
	"strconv"
	"time"
	waProto "ws/framework/application/constant/binary/proto"
	containerInterface "ws/framework/application/container/abstract_interface"
	accountDB "ws/framework/application/data_storage/account/database"
	deviceDB "ws/framework/application/data_storage/device/database"
	prekeyDB "ws/framework/application/data_storage/prekey/database"
	pushConfigDB "ws/framework/application/data_storage/push_config/database"
	registrationTokenDB "ws/framework/application/data_storage/registration_token/database"
	"ws/framework/application/libsignal/ecc"
	"ws/framework/application/libsignal/serialize"
	"ws/framework/application/libsignal/state/record"
	"ws/framework/lib/msisdn"
	"ws/framework/utils"
)

// DeviceData .
type DeviceData struct {
	JIDNumber   string `json:"jidNumber" validate:"required,min=1,max=15"`
	PhoneNumber string `json:"phoneNumber" validate:"required,min=1,max=15"`

	FBUuid              string `json:"fbuuid"`
	Uuid                string `json:"uuid"`
	IdentityKeypairData []byte `json:"identityKeypairData"`
	LoginPayload        []byte `json:"loginPayload"`
	PushName            string `json:"pushName"`
	UuidTime            string `json:"uuidTime"`
	RegistrationId      []byte `json:"registrationId"`
	ClientStaticPriKey  []byte `json:"clientStaticPriKey"`

	PushId   string `json:"pushId"`
	PushVoip string `json:"pushvoip"`
	PushApns string `json:"pushpkey"`

	EdgeRoutingInfo []byte `json:"edgeRoutingInfo"`

	PreKeysList [][]byte `json:"preKeysList"`

	SignedPreKeyRecord []byte `json:"signedPreKeyRecord"`

	RegTokenPBData []byte   `json:"regTokenPBData"`
	BackupKeyList  [][]byte `json:"backupKeyList"`
}

func parseDevice(data *DeviceData) (*deviceDB.Device, error) {
	payLoad := waProto.ClientPayload{}
	err := proto.Unmarshal(data.LoginPayload, &payLoad)
	if err != nil {
		return nil, err
	}

	IdentityKeypair := waProto.IdentityKeypair{}
	err = proto.Unmarshal(data.IdentityKeypairData, &IdentityKeypair)
	if err != nil {
		return nil, err
	}

	fbUuidCreateTime, _ := strconv.Atoi(data.UuidTime)
	userAgent := payLoad.GetUserAgent()

	imsi, _ := msisdn.Parse(data.PhoneNumber, true)

	djbECKeyPair := ecc.CreateKeyPair(data.ClientStaticPriKey)
	clientStaticPubKey := djbECKeyPair.PublicKey().PublicKey()

	return &deviceDB.Device{
		JID:              data.JIDNumber,
		RegistrationId:   uint32(utils.BigEndianBytesToInt(data.RegistrationId)),
		IdentityKey:      IdentityKeypair.GetPrivateKey(),
		PushName:         data.PushName,
		UserName:         strconv.FormatUint(payLoad.GetUsername(), 10),
		FBUuid:           data.FBUuid,
		Uuid:             data.Uuid,
		FBUuidCreateTime: int32(fbUuidCreateTime),
		OsVersion:        userAgent.GetOsVersion(),
		Mcc:              imsi.GetMCC(),
		Mnc:              imsi.GetMNC(),
		//Mcc:                userAgent.GetMcc(),
		//Mnc:                userAgent.GetMnc(),
		Manufacturer:       userAgent.GetManufacturer(),
		Device:             userAgent.GetDevice(),
		Language:           userAgent.GetLocaleLanguageIso6391(),
		Country:            userAgent.GetLocaleCountryIso31661Alpha2(),
		BuildNumber:        userAgent.GetOsBuildNumber(),
		PrivateStatsId:     userAgent.GetPhoneId(),
		Area:               imsi.GetCC(),
		Phone:              imsi.GetPhoneNumber(),
		ClientStaticPriKey: data.ClientStaticPriKey,
		ClientStaticPubKey: clientStaticPubKey[:],
		ServerStaticKey:    make([]byte, 0),
	}, nil
}

func parseRegistrationToken(data *DeviceData) (*registrationTokenDB.RegistrationToken, error) {
	regData := waProto.RegistrationToken{}
	err := proto.Unmarshal(data.RegTokenPBData, &regData)
	if err != nil {
		return nil, err
	}

	return &registrationTokenDB.RegistrationToken{
		JID:           data.JIDNumber,
		RecoveryToken: regData.GetRecoveryToken(),
		BackupToken:   regData.GetBackupToken(),
		BackupKey:     data.BackupKeyList[0],
		BackupKey2:    data.BackupKeyList[1],
	}, nil
}

func parsePreKeyList(data *DeviceData) ([]prekeyDB.PreKey, error) {
	preKeyRecords := make([]prekeyDB.PreKey, len(data.PreKeysList))
	for i, keyBuffer := range data.PreKeysList {
		keyPair := waProto.PreKeypair{}
		err := proto.Unmarshal(keyBuffer, &keyPair)
		if err != nil {
			return nil, err
		}

		recordPreKeyStructure := record.PreKeyStructure{
			ID:         keyPair.GetId(),
			PublicKey:  keyPair.GetPublicKey(),
			PrivateKey: keyPair.GetPrivateKey(),
		}

		// 去掉 ecc.DjbType 0x05
		recordPreKeyStructure.PublicKey = recordPreKeyStructure.PublicKey[1:]

		buffer := serialize.Proto.PreKeyRecord.Serialize(&recordPreKeyStructure)

		preKeyRecords[i].JID = data.JIDNumber
		preKeyRecords[i].KeyId = recordPreKeyStructure.ID
		preKeyRecords[i].KeyBuff = buffer
	}

	return preKeyRecords, nil
}

func parseSignedPreKey(data *DeviceData) (uint32, []byte, error) {
	signedPreKeyRecord := waProto.SignedPreKeyRecord{}
	err := proto.Unmarshal(data.SignedPreKeyRecord, &signedPreKeyRecord)
	if err != nil {
		return 0, nil, err
	}

	recordSignedPreKeyStructure := record.SignedPreKeyStructure{
		ID:         signedPreKeyRecord.GetId(),
		PublicKey:  signedPreKeyRecord.GetPublicKey(),
		PrivateKey: signedPreKeyRecord.GetPrivateKey(),
		Signature:  signedPreKeyRecord.GetSignature(),
		Timestamp:  time.Now().Unix(),
	}

	// 去掉 ecc.DjbType 0x05
	recordSignedPreKeyStructure.PublicKey = recordSignedPreKeyStructure.PublicKey[1:]

	buffer := serialize.Proto.SignedPreKeyRecord.Serialize(&recordSignedPreKeyStructure)

	return recordSignedPreKeyStructure.ID, buffer, nil
}

// Do .
func Do(container containerInterface.IAppIocContainer, data *DeviceData) error {
	// Device
	device, err := parseDevice(data)
	if err != nil {
		return err
	}
	err = container.ResolveDeviceService().Import(device)
	if err != nil {
		return err
	}

	// RegistrationToken
	registrationToken, err := parseRegistrationToken(data)
	if err != nil {
		return err
	}
	err = container.ResolveRegistrationTokenService().Import(registrationToken)
	if err != nil {
		return err
	}

	// PushConfig
	if len(data.PushVoip) > 0 && len(data.PushApns) > 0 {
		pushConfigData := pushConfigDB.PushConfig{
			JID:       data.JIDNumber,
			VoipToken: data.PushVoip,
			ApnsToken: data.PushApns,
		}

		err = container.ResolvePushConfigService().Import(&pushConfigData)
		if err != nil {
			return err
		}
	} else {
		_, err = container.ResolvePushConfigService().Create()
		if err != nil {
			return err
		}
	}

	// PreKey
	preKeyList, err := parsePreKeyList(data)
	if err != nil {
		return err
	}
	err = container.ResolvePreKeyService().Import(preKeyList)
	if err != nil {
		return err
	}

	// RoutingInfo
	err = container.ResolveRoutingInfoService().Create(data.EdgeRoutingInfo)
	if err != nil {
		return err
	}

	// SignedPreKey
	signedPreKeyID, signedPreKeyBuffer, err := parseSignedPreKey(data)
	if err != nil {
		return err
	}

	container.ResolveSignedPreKeyService().SaveSignedPreKeyBuffer(signedPreKeyID, signedPreKeyBuffer)

	// Account
	account := accountDB.Account{
		JID:        data.JIDNumber,
		LoginCount: 1,
	}

	err = container.ResolveAccountService().Import(&account)
	if err != nil {
		return err
	}

	return nil
}
