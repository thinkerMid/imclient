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

type clientLogModel struct {
	Cc             string `json:"cc"`
	In             string `json:"in"`
	Rc             string `json:"rc"` // 0
	Lg             string `json:"lg"`
	Lc             string `json:"lc"`
	AuthKey        string `json:"authkey"`
	Eregid         string `json:"e_regid"`
	Ekeytype       string `json:"e_keytype"`
	Eident         string `json:"e_ident"`
	EskeyId        string `json:"e_skey_id"`
	EskeyVal       string `json:"e_skey_val"`
	EskeySig       string `json:"e_skey_sig"`
	Fdid           string `json:"fdid"`
	Expid          string `json:"expid"`
	CurrentScreen  string `json:"current_screen"`  //verify_sms
	PreviousScreen string `json:"previous_screen"` //enter_number
	ActionTaken    string `json:"action_taken"`    //continue
	Id             string `json:"id"`
	VName          string `json:"vname"`
}

// ClientSendLogResp .
type ClientSendLogResp struct {
	Login  string `json:"login"`  //手机号
	Status string `json:"status"` //状态
}

// HasError .
func (p *ClientSendLogResp) HasError() error {
	if p.Status == "ok" {
		return nil
	}

	return fmt.Errorf("发送日志失败")
}

// MakeClientLogBody .
func MakeClientLogBody(appIocContainer containerInterface.IAppIocContainer, currentScreen string, previousScreen string, actionTaken string) string {
	device := appIocContainer.ResolveDeviceService().Context()
	signedPreKeyKeyPair := appIocContainer.ResolveSignedPreKeyService().Context()
	identity := appIocContainer.ResolveIdentityService().Context()
	aesKey := appIocContainer.ResolveAesKeyService().Context()
	registrationToken := appIocContainer.ResolveRegistrationTokenService().Context()
	logger := appIocContainer.ResolveLogger().Named("ClientLog")
	configuration := appIocContainer.ResolveWhatsappConfiguration()
	vname := appIocContainer.ResolveBusinessService().GenerateBusinessVerifiedName(false)

	clientLog := clientLogModel{
		CurrentScreen:  currentScreen,
		PreviousScreen: previousScreen,
		ActionTaken:    actionTaken,
	}

	clientLog.Cc = device.Area
	clientLog.In = device.Phone
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
