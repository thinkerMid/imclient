package sms

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	. "labs/utils"
	"sync"
)

type SMS struct {
}

var (
	sms       *SMS
	once      sync.Once
	tlsConfig *tls.Config
)

func Instance() *SMS {
	once.Do(func() {
		sms = &SMS{}
		tlsConfig = &tls.Config{ServerName: service}
	})
	return sms
}

func (s *SMS) QueryCountries() {
	opts := []RequestOptions{
		WithHeaders(map[string]string{
			"Content-Type": "application/json",
		}),
		WithTLSConfig(tlsConfig),
	}

	type RequestBody struct {
		ApiKey string `json:"api_key"`
		Action string `json:"action"`
	}

	req, _ := json.Marshal(RequestBody{
		ApiKey: apiKey,
		Action: "getCountries",
	})
	body := bytes.NewBuffer(req)

	buff, err := HttpRequest(apiUrl, "POST", body, opts...)
	if err != nil {
		return
	}

	countries := make(map[string]Country, 0)
	if err = json.Unmarshal(buff, &countries); err != nil {
		return
	}
}
