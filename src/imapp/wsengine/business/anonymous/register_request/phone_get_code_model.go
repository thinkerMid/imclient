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

// phoneGetCodeModel .
type phoneGetCodeModel struct {
	Cc       string `json:"cc"`
	In       string `json:"in"`
	Rc       string `json:"rc"`
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
	Method   string `json:"method"`
	SimMcc   string `json:"sim_mcc"`
	SimMnc   string `json:"sim_mnc"`
	Token    string `json:"token"`
	Id       string `json:"id"`
	VName    string `json:"vname"`
}

// PhoneGetCodeRep .
type PhoneGetCodeRep struct {
	Login       string `json:"login"`           //手机号
	Status      string `json:"status"`          //状态
	Reason      string `json:"reason"`          //描述
	NotifyAfter int    `json:"notify_after"`    //这个不知道是什么，不关注这个字段
	Length      int    `json:"length"`          //验证码长度
	Method      string `json:"method"`          //验证码方式
	RetryAfter  int    `json:"retry_after"`     //下次重试时间，如果接码失败，必须要在这个时间之后再次发起才有用
	SmsWait     int    `json:"sms_wait"`        //不关注这个字段
	VoiceWait   int    `json:"voice_wait"`      //不关注这个字段
	FlashType   int    `json:"flash_type"`      //这个不知道是什么，不关注这个字段
	Param       string `json:"param,omitempty"` //        //如果有错误，这里会提示错误的字段
}

// HasError .
func (p *PhoneGetCodeRep) HasError() error {
	if p.Status == "sent" {
		return nil
	}

	if p.Reason == "too_recent" {
		return fmt.Errorf(fmt.Sprintf("验证码已发送,下次重试时间:%vs", p.RetryAfter))
	} else if p.Reason == "no_routes" {
		return fmt.Errorf("无法发送验证码")
	} else if p.Reason == "blocked" {
		return fmt.Errorf("手机号已封禁")
	} else if p.Reason == "bad_param" {
		if p.Param == "number" {
			return fmt.Errorf("号码不正确")
		}

		return fmt.Errorf(p.Param + "参数错误")
	} else if p.Reason == "too_many" {
		return fmt.Errorf("请求验证码频繁")
	}

	return fmt.Errorf(p.Reason)
}

// MakePhoneGetSmsCodeBody .
func MakePhoneGetSmsCodeBody(appIocContainer containerInterface.IAppIocContainer) string {
	device := appIocContainer.ResolveDeviceService().Context()
	signedPreKeyKeyPair := appIocContainer.ResolveSignedPreKeyService().Context()
	identity := appIocContainer.ResolveIdentityService().Context()
	aesKey := appIocContainer.ResolveAesKeyService().Context()
	registrationToken := appIocContainer.ResolveRegistrationTokenService().Context()
	logger := appIocContainer.ResolveLogger().Named("GetSmsCode")
	configuration := appIocContainer.ResolveWhatsappConfiguration()
	vname := appIocContainer.ResolveBusinessService().GenerateBusinessVerifiedName(false)

	phoneGetCodeModel := phoneGetCodeModel{}
	phoneGetCodeModel.Cc = device.Area
	phoneGetCodeModel.In = device.Phone
	phoneGetCodeModel.Rc = "0"
	phoneGetCodeModel.Lg = device.Language
	phoneGetCodeModel.Lc = device.Country
	phoneGetCodeModel.AuthKey = utils.Base64Encode(device.ClientStaticPubKey)
	phoneGetCodeModel.Eregid = utils.Base64Encode(utils.IntToBigEndianBytes(int(device.RegistrationId)))
	phoneGetCodeModel.Ekeytype = utils.Base64Encode(utils.IntToBigEndianBytes(ecc.DjbType))
	phoneGetCodeModel.EskeyId = utils.Base64Encode(utils.IntToBigEndianBytes(int(signedPreKeyKeyPair.ID())))
	phoneGetCodeModel.Eident = utils.Base64Encode(identity.Pub[:])

	eskeyval := signedPreKeyKeyPair.KeyPair().PublicKey().PublicKey()
	eskeysig := signedPreKeyKeyPair.Signature()

	phoneGetCodeModel.EskeyVal = utils.Base64Encode(eskeyval[:])
	phoneGetCodeModel.EskeySig = utils.Base64Encode(eskeysig[:])
	phoneGetCodeModel.Fdid = device.FBUuid
	phoneGetCodeModel.Expid = utils.Base64Encode(utils.ParseUUID4(device.Uuid))
	phoneGetCodeModel.Method = "sms"
	phoneGetCodeModel.SimMcc = device.Mcc
	phoneGetCodeModel.SimMnc = device.Mnc
	phoneGetCodeModel.Token = utils.MD5Hex(configuration.AESPassword + configuration.BuildHash + device.Phone)
	phoneGetCodeModel.Id = utils.URLEncode(functionTools.B2S(registrationToken.RecoveryToken))
	phoneGetCodeModel.VName = utils.Base64Encode(vname)

	phoneGetCodeModelJson, _ := json.Marshal(&phoneGetCodeModel)
	logger.Debug(functionTools.B2S(phoneGetCodeModelJson))

	params := utils.GenURLParams(phoneGetCodeModelJson)
	aesResult := utils.AesGcmEncrypt(aesKey.AesKey, []byte(params))

	var buffer bytes.Buffer
	buffer.Write(aesKey.PubKey)
	buffer.Write(aesResult)

	result := utils.Base64Encode(buffer.Bytes())

	return configuration.HttpUrl + "/v2/code?ENC=" + result
}
