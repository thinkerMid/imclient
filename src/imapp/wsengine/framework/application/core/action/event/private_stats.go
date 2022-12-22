package event

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"github.com/chenzhuoyu/base64x"
	hertzConst "github.com/cloudwego/hertz/pkg/protocol/consts"
	"time"
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	eventSerialize "ws/framework/plugin/event_serialize"
	"ws/framework/plugin/json"
	httpApi "ws/framework/plugin/network/http_api"
	"ws/framework/utils"
)

type errorContent struct {
	ErrorUserTitle string `json:"error_user_title"`
}

type uploadResponse struct {
	LogReply string       `json:"log_reply"`
	Error    errorContent `json:"error"`
}

// PrivateStats .
type PrivateStats struct {
	processor.BaseAction
	random []byte
	scalar []byte
}

// Start .
func (m *PrivateStats) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	scalar, random, credential, err := utils.GenerateEd25519Credential()
	if err != nil {
		return err
	}

	m.scalar = scalar
	m.random = random

	iq := message.InfoQuery{
		ID:        context.GenerateRequestID(),
		Namespace: "privatestats",
		Type:      message.IqGet,
		To:        types.ServerJID,
		Content: []waBinary.Node{{
			Tag:   "sign_credential",
			Attrs: waBinary.Attrs{"version": "1"},
			Content: []waBinary.Node{{
				Tag:     "blinded_credential",
				Content: credential,
			}},
		}},
	}

	m.SendMessageId, err = context.SendIQ(iq)
	return
}

// Receive .
func (m *PrivateStats) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	defer next()

	buffer := eventSerialize.AcquireEventBuffer()
	defer eventSerialize.ReleaseEventBuffer(buffer)

	// 事件日志打包
	sendCount := context.ResolveAccountService().Context().SendChannel2EventCount
	context.ResolveChannel2EventCache().PackBuffer(sendCount, buffer)

	// 解析
	signCredential := context.Message().GetChildByTag("sign_credential")
	signedCredential := signCredential.GetChildByTag("signed_credential").Content.([]byte)
	acsPublicKey := signCredential.GetChildByTag("acs_public_key").Content.([]byte)

	var credential string

	// 签名
	{
		signature, err := utils.GenerateEd25519Signature(m.scalar, acsPublicKey, signedCredential)
		if err != nil {
			return err
		}

		sha512Digest := sha512.New()
		sha512Digest.Write(m.random)
		sha512Digest.Write(signature)

		sha512Buff := sha512Digest.Sum(nil)

		mac := hmac.New(sha256.New, sha512Buff)
		mac.Write(buffer.Byte())
		macB64Str := base64x.URLEncoding.EncodeToString(mac.Sum(nil))

		randomB64Str := base64x.URLEncoding.EncodeToString(m.random)
		credential = randomB64Str + "+" + macB64Str
	}

	configuration := context.ResolveWhatsappConfiguration()
	metaDataMapping := map[string]int64{"t": time.Now().Unix()}
	metaData, _ := json.Marshal(metaDataMapping)

	fields := []httpApi.MultipartField{
		{
			Name:        "credential",
			ContentType: httpApi.TextPlain,
			Body:        credential,
		},
		{
			Name:        "access_token",
			ContentType: httpApi.TextPlain,
			Body:        configuration.PrivateStatsAccessToken,
		},
		{
			Name:        "message",
			FileName:    "WAMEventBuffer.dat",
			ContentType: httpApi.ApplicationOctetStream,
			Body:        buffer.Byte(),
		},
		{
			Name:        "meta_data",
			ContentType: httpApi.ApplicationJson,
			Body:        metaData,
		},
	}

	resp := uploadResponse{}
	httpClient := context.ResolveHttpClient()

	err = httpApi.DoAndBind(
		httpClient, &resp,
		httpApi.Url(configuration.PrivateStatsURL),
		httpApi.Method(hertzConst.MethodPost),
		httpApi.UserAgent(context.ResolveDeviceService().PrivateStatsAgent()),
		httpApi.FormData(configuration.PrivateStatsBoundary, fields),
	)

	if err == nil && resp.LogReply != "forward" {
		return fmt.Errorf(resp.Error.ErrorUserTitle)
	}

	return
}

func (m *PrivateStats) Error(_ containerInterface.IMessageContext, _ error) {}
