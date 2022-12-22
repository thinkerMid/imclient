package registerRequest

import (
	"bytes"
	"fmt"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/libsignal/ecc"
	"ws/framework/plugin/json"
	"ws/framework/utils"
	functionTools "ws/framework/utils/function_tools"
)

type phoneExistModel struct {
	Cc        string `json:"cc"`
	In        string `json:"in"`
	Rc        string `json:"rc"`
	Lg        string `json:"lg"`
	Lc        string `json:"lc"`
	AuthKey   string `json:"authkey"`
	Eregid    string `json:"e_regid"`
	Ekeytype  string `json:"e_keytype"`
	Eident    string `json:"e_ident"`
	EskeyId   string `json:"e_skey_id"`
	EskeyVal  string `json:"e_skey_val"`
	EskeySig  string `json:"e_skey_sig"`
	Fdid      string `json:"fdid"`
	Expid     string `json:"expid"`
	OfflineAb string `json:"offline_ab"`
	Id        string `json:"id"`
	VName     string `json:"vname"`
}

type metricsModel struct {
	ExpidCd int32 `json:"expid_cd"`
	ExpidMd int32 `json:"expid_md"`
	RcC     bool  `json:"rc_c"`
}

type offlineAbModel struct {
	Exposure []string      `json:"exposure"`
	Metrics  *metricsModel `json:"metrics"`
}

// PhoneExistRep .
type PhoneExistRep struct {
	Login       string `json:"login"`        //手机号
	Status      string `json:"status"`       //状态
	Reason      string `json:"reason"`       //描述
	SmsLength   int    `json:"sms_length"`   //短信验证码长度 不关注这个字段
	VoiceLength int    `json:"voice_length"` //语音验证码长度 不关注这个字段
	SmsWait     int    `json:"sms_wait"`     //不关注这个字段
	VoiceWait   int    `json:"voice_wait"`   //不关注这个字段
}

// HasError .
func (p *PhoneExistRep) HasError() error {
	if p.Status == "fail" {
		if p.Reason == "blocked" {
			return fmt.Errorf("手机号已封禁")
		} else if p.Reason == "length_long" || p.Reason == "length_short" {
			return fmt.Errorf("号码不正确")
		}
	}

	return nil
}

// MakePhoneExistBody .
func MakePhoneExistBody(appIocContainer containerInterface.IAppIocContainer) string {
	device := appIocContainer.ResolveDeviceService().Context()
	signedPreKeyKeyPair := appIocContainer.ResolveSignedPreKeyService().Context()
	identity := appIocContainer.ResolveIdentityService().Context()
	aesKey := appIocContainer.ResolveAesKeyService().Context()
	registrationToken := appIocContainer.ResolveRegistrationTokenService().Context()
	logger := appIocContainer.ResolveLogger().Named("PhoneExist")
	configuration := appIocContainer.ResolveWhatsappConfiguration()
	vname := appIocContainer.ResolveBusinessService().GenerateBusinessVerifiedName(false)

	phoneExistModel := phoneExistModel{}
	phoneExistModel.Cc = device.Area
	phoneExistModel.In = device.Phone
	phoneExistModel.Rc = "0"
	phoneExistModel.Lg = device.Language
	phoneExistModel.Lc = device.Country
	phoneExistModel.AuthKey = utils.Base64Encode(device.ClientStaticPubKey)
	phoneExistModel.Eregid = utils.Base64Encode(utils.IntToBigEndianBytes(int(device.RegistrationId)))
	phoneExistModel.Ekeytype = utils.Base64Encode(utils.IntToBigEndianBytes(ecc.DjbType))
	phoneExistModel.EskeyId = utils.Base64Encode(utils.IntToBigEndianBytes(int(signedPreKeyKeyPair.ID())))
	phoneExistModel.Eident = utils.Base64Encode(identity.Pub[:])

	eskeyval := signedPreKeyKeyPair.KeyPair().PublicKey().PublicKey()
	eskeysig := signedPreKeyKeyPair.Signature()

	phoneExistModel.EskeyVal = utils.Base64Encode(eskeyval[:])
	phoneExistModel.EskeySig = utils.Base64Encode(eskeysig[:])
	phoneExistModel.Fdid = device.FBUuid
	phoneExistModel.Expid = utils.Base64Encode(utils.ParseUUID4(device.Uuid))
	phoneExistModel.Id = utils.URLEncode(functionTools.B2S(registrationToken.RecoveryToken))
	phoneExistModel.VName = utils.Base64Encode(vname)

	offlineAbModel := &offlineAbModel{}
	offlineAbModel.Exposure = make([]string, 0)
	offlineAbModel.Metrics = &metricsModel{}
	offlineAbModel.Metrics.ExpidCd = device.FBUuidCreateTime
	offlineAbModel.Metrics.ExpidMd = device.FBUuidCreateTime
	offlineAbModel.Metrics.RcC = true //WAUserSessionPreferences  regRCCreated

	offlineAbModelJson, _ := json.Marshal(&offlineAbModel)
	offlineAb := utils.URLEncode(functionTools.B2S(offlineAbModelJson))

	phoneExistModel.OfflineAb = offlineAb

	phoneExistModelJson, _ := json.Marshal(&phoneExistModel)
	logger.Debug(functionTools.B2S(phoneExistModelJson))

	params := utils.GenURLParams(phoneExistModelJson)
	aesResult := utils.AesGcmEncrypt(aesKey.AesKey, []byte(params))

	var buffer bytes.Buffer
	buffer.Write(aesKey.PubKey)
	buffer.Write(aesResult)

	result := utils.Base64Encode(buffer.Bytes())

	return configuration.HttpUrl + "/v2/exist?ENC=" + result
}

// MakeEmptyExistBody .
func MakeEmptyExistBody(appIocContainer containerInterface.IAppIocContainer) string {
	configuration := appIocContainer.ResolveWhatsappConfiguration()

	return configuration.HttpUrl + "/v2/exist"
}
