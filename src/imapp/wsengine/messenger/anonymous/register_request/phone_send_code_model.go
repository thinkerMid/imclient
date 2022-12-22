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

type phoneSendCodeModel struct {
	Cc          string `json:"cc"`
	In          string `json:"in"`
	Rc          string `json:"rc"`
	Lg          string `json:"lg"`
	Lc          string `json:"lc"`
	AuthKey     string `json:"authkey"`
	Eregid      string `json:"e_regid"`
	Ekeytype    string `json:"e_keytype"`
	Eident      string `json:"e_ident"`
	EskeyId     string `json:"e_skey_id"`
	EskeyVal    string `json:"e_skey_val"`
	EskeySig    string `json:"e_skey_sig"`
	Fdid        string `json:"fdid"`
	Expid       string `json:"expid"`
	Code        string `json:"code"`
	Entered     string `json:"entered"`
	Id          string `json:"id"`
	BackupToken string `json:"backup_token"`
}

// PhoneSendCodeRep .
type PhoneSendCodeRep struct {
	Login           string `json:"login"`             //手机号
	Status          string `json:"status"`            //状态
	Reason          string `json:"reason"`            //描述
	RetryAfter      int    `json:"retry_after"`       //下次重试时间
	EdgeRoutingInfo string `json:"edge_routing_info"` //dns
	SecurityCodeSet bool   `json:"security_code_set"` //是否设置安全码
	Type            string `json:"type"`              //登陆类型
	Param           string `json:"param,omitempty"`   //        //如果有错误，这里会提示错误的字段
}

// HasError .
func (p *PhoneSendCodeRep) HasError() error {
	if p.Status == "ok" {
		return nil
	}

	if p.Reason == "bad_param" {
		if p.Param == "number" {
			return fmt.Errorf("号码不正确")
		}

		return fmt.Errorf(p.Param + "参数错误")
	} else if p.Reason == "guessed_too_fast" {
		return fmt.Errorf(fmt.Sprintf("操作频繁,下次重试时间:%vs", p.RetryAfter))
	} else if p.Reason == "mismatch" {
		return fmt.Errorf("验证码错误")
		// 绕过了请求验证码那步骤就发送验证码的问题
	} else if p.Reason == "stale" || p.Reason == "missing" {
		return fmt.Errorf("请先发送验证码")
	}

	return fmt.Errorf(p.Reason)
}

// MakePhoneSendSmsCodeBody .
func MakePhoneSendSmsCodeBody(appIocContainer containerInterface.IAppIocContainer, smsCode string) string {
	device := appIocContainer.ResolveDeviceService().Context()
	signedPreKeyKeyPair := appIocContainer.ResolveSignedPreKeyService().Context()
	identity := appIocContainer.ResolveIdentityService().Context()
	aesKey := appIocContainer.ResolveAesKeyService().Context()
	reverToken := appIocContainer.ResolveRegistrationTokenService().Context()
	logger := appIocContainer.ResolveLogger().Named("SendSmsCode")
	configuration := appIocContainer.ResolveWhatsappConfiguration()

	phoneSendCodeModel := phoneSendCodeModel{}
	phoneSendCodeModel.Cc = device.Area
	phoneSendCodeModel.In = device.Phone
	phoneSendCodeModel.Rc = "0"
	phoneSendCodeModel.Lg = device.Language
	phoneSendCodeModel.Lc = device.Country
	phoneSendCodeModel.AuthKey = utils.Base64Encode(device.ClientStaticPubKey)
	phoneSendCodeModel.Eregid = utils.Base64Encode(utils.IntToBigEndianBytes(int(device.RegistrationId)))
	phoneSendCodeModel.Ekeytype = utils.Base64Encode(utils.IntToBigEndianBytes(ecc.DjbType))
	phoneSendCodeModel.EskeyId = utils.Base64Encode(utils.IntToBigEndianBytes(int(signedPreKeyKeyPair.ID())))
	phoneSendCodeModel.Eident = utils.Base64Encode(identity.Pub[:])

	eskeyval := signedPreKeyKeyPair.KeyPair().PublicKey().PublicKey()
	eskeysig := signedPreKeyKeyPair.Signature()
	phoneSendCodeModel.EskeyVal = utils.Base64Encode(eskeyval[:])
	phoneSendCodeModel.EskeySig = utils.Base64Encode(eskeysig[:])
	phoneSendCodeModel.Fdid = device.FBUuid
	phoneSendCodeModel.Expid = utils.Base64Encode(utils.ParseUUID4(device.Uuid))
	phoneSendCodeModel.Code = smsCode
	phoneSendCodeModel.Entered = "1"
	phoneSendCodeModel.Id = utils.URLEncode(functionTools.B2S(reverToken.RecoveryToken))
	phoneSendCodeModel.BackupToken = utils.URLEncode(functionTools.B2S(reverToken.BackupToken))

	phoneSendCodeModelJson, _ := json.Marshal(&phoneSendCodeModel)
	logger.Debug(functionTools.B2S(phoneSendCodeModelJson))

	params := utils.GenURLParams(phoneSendCodeModelJson)
	aesResult := utils.AesGcmEncrypt(aesKey.AesKey, []byte(params))

	var buffer bytes.Buffer
	buffer.Write(aesKey.PubKey)
	buffer.Write(aesResult)

	result := utils.Base64Encode(buffer.Bytes())

	return configuration.HttpUrl + "/v2/register?ENC=" + result
}
