package msisdn

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	hertzConst "github.com/cloudwego/hertz/pkg/protocol/consts"
	"strconv"
	"time"
	msisdnDB "ws/framework/application/data_storage/msisdn/database"
	"ws/framework/env"
	"ws/framework/plugin/database"
	databaseTools "ws/framework/plugin/database/database_tools"
	"ws/framework/plugin/json"
	networkConstant "ws/framework/plugin/network/constant"
	httpApi "ws/framework/plugin/network/http_api"
	"ws/framework/plugin/network/netpoll"
)

type hlrLookupSearchResp struct {
	MCC string `json:"mcc,omitempty"`
	MNC string `json:"mnc,omitempty"`
}

// 国家码+前三位查询数据库
func hlrLookupSearchInCache(number string, imsi *IMSIParser) bool {
	where := msisdnDB.Msisdn{PhoneNumber: "+" + number[:len(imsi.GetCC())+3]}
	result := msisdnDB.QueryResult{}

	err := databaseTools.Find(database.MasterDB(), &where, &result)
	if err == nil {
		imsi.MCC = result.MCC
		imsi.MNC = result.MNC
	}

	return err == nil
}

func hlrLookupSearch(imsi *IMSIParser) {
	number := imsi.GetCC() + imsi.GetPhoneNumber()
	if hlrLookupSearchInCache(number, imsi) {
		return
	}

	timestamp := time.Now().Unix()
	path := "/hlr-lookup"
	url := "https://www.hlr-lookups.com/api/v2/hlr-lookup"

	secret := env.NacosConfig.HlrLookup.Secret
	apiKey := env.NacosConfig.HlrLookup.ApiKey
	body := map[string]string{"msisdn": "+" + number, "route": "PTX"}
	bodyJson, _ := json.Marshal(body)
	dataString := string(bodyJson)
	message := path + strconv.Itoa(int(timestamp)) + "POST" + dataString
	signature := computeHmac256(message, secret)

	client := netpoll.HTTP(networkConstant.ConnectionConfig{
		Tls: &tls.Config{InsecureSkipVerify: true},
	})

	var resp hlrLookupSearchResp

	err := httpApi.DoAndBind(
		client, &resp,
		httpApi.Url(url),
		httpApi.Method(hertzConst.MethodPost),
		httpApi.Header("X-Digest-Key", apiKey),
		httpApi.Header("X-Digest-Signature", signature),
		httpApi.Header("X-Digest-Timestamp", strconv.Itoa(int(timestamp))),
		httpApi.Body(bodyJson),
	)
	if err != nil {
		return
	}

	if resp.MCC != "" && resp.MNC != "" {
		data := msisdnDB.Msisdn{PhoneNumber: "+" + number[:len(imsi.GetCC())+3], MCC: resp.MCC, MNC: resp.MNC}
		_, _ = databaseTools.Create(database.MasterDB(), &data)

		imsi.MCC = resp.MCC
		imsi.MNC = resp.MNC
	}
}

func computeHmac256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}
