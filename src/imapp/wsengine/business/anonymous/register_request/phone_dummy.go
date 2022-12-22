package registerRequest

import (
	"bytes"
	containerInterface "ws/framework/application/container/abstract_interface"
	"ws/framework/application/libsignal/ecc"
	"ws/framework/plugin/json"
	"ws/framework/utils"
	functionTools "ws/framework/utils/function_tools"
)

type SMBLogModel struct {
	Cc       string `json:"cc"`
	In       string `json:"in"`
	Rc       string `json:"rc"` // 0
	Lg       string `json:"lg"`
	Lc       string `json:"lc"`
	AuthKey  string `json:"authkey"`
	Eregid   string `json:"e_regid"`
	Ekeytype string `json:"e_keytype"`
	Eident   string `json:"e_ident"`
	EskeyId  string `json:"e_skey_id"`
	EskeyVal string `json:"e_skey_val"`
	EskeySig string `json:"e_skey_sig"`
	Fdid     string `json:"fdid"`
	Expid    string `json:"expid"`
	Id       string `json:"id"`
	//BackupToken    string `json:"backup_token"`
	VName          string `json:"vname"`
	EventName      string `json:"event_name"`
	HasConsumerApp int    `json:"has_consumer_app"`
	AppSource      string `json:"app_install_source"`
	OnBoardStep    int    `json:"smb_onboarding_step"`
	ConsumerLogin  int    `json:"is_logged_in_on_consumer_app"`
	SequenceNumber int    `json:"sequence_number"`
}

// MakeSMBLogBody .
func MakeSMBLogBody(appIocContainer containerInterface.IAppIocContainer, step, sequence int) string {
	device := appIocContainer.ResolveDeviceService().Context()
	signedPreKeyKeyPair := appIocContainer.ResolveSignedPreKeyService().Context()
	identity := appIocContainer.ResolveIdentityService().Context()
	aesKey := appIocContainer.ResolveAesKeyService().Context()
	registrationToken := appIocContainer.ResolveRegistrationTokenService().Context()
	logger := appIocContainer.ResolveLogger().Named("SMBLog")
	configuration := appIocContainer.ResolveWhatsappConfiguration()
	vname := appIocContainer.ResolveBusinessService().GenerateBusinessVerifiedName(false)

	clientLog := SMBLogModel{
		EventName:      "smb_client_onboarding_journey",
		HasConsumerApp: 0,
		AppSource:      "unknown%7Cunknown",
		OnBoardStep:    step,
		ConsumerLogin:  0,
		SequenceNumber: sequence,
	}

	clientLog.Cc = "1"          // 固定的
	clientLog.In = "2199990000" // 固定的
	clientLog.Rc = "0"
	clientLog.Lg = device.Language
	clientLog.Lc = device.Country
	clientLog.AuthKey = utils.Base64Encode(device.ClientStaticPubKey)
	clientLog.Eregid = utils.Base64Encode(utils.IntToBigEndianBytes(int(device.RegistrationId)))
	clientLog.Ekeytype = utils.Base64Encode(utils.IntToBigEndianBytes(ecc.DjbType))
	clientLog.EskeyId = utils.Base64Encode(utils.IntToBigEndianBytes(int(signedPreKeyKeyPair.ID())))
	clientLog.Eident = utils.Base64Encode(identity.Pub[:])

	eskeyval := signedPreKeyKeyPair.KeyPair().PublicKey().PublicKey()
	eskeysig := signedPreKeyKeyPair.Signature()

	clientLog.EskeyVal = utils.Base64Encode(eskeyval[:])
	clientLog.EskeySig = utils.Base64Encode(eskeysig[:])
	clientLog.Fdid = device.FBUuid
	clientLog.Expid = utils.Base64Encode(utils.ParseUUID4(device.Uuid))
	clientLog.Id = utils.URLEncode(functionTools.B2S(registrationToken.RecoveryToken))
	clientLog.VName = utils.Base64Encode(vname)

	clientLogModelJson, _ := json.Marshal(&clientLog)
	logger.Debug(functionTools.B2S(clientLogModelJson))

	params := utils.GenURLParams(clientLogModelJson)
	aesResult := utils.AesGcmEncrypt(aesKey.AesKey, []byte(params))

	var buffer bytes.Buffer
	buffer.Write(aesKey.PubKey)
	buffer.Write(aesResult)

	result := utils.Base64Encode(buffer.Bytes())

	return configuration.HttpUrl + "/v2/client_log?ENC=" + result
}
