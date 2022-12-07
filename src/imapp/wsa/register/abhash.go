package register

import "fmt"

type PhoneABPropResp struct {
	AbHash string `json:"ab_hash"`
	AbKey  string `json:"ab_key"`
	ExpCfg string `json:"exp_cfg"`
	Login  string `json:"login"`
	Status string `json:"status"`
}

// HasError .
func (p *PhoneABPropResp) HasError() error {
	return nil
}

// MakePhoneABPropBody .
func MakePhoneABPropBody() string {
	device := appIocContainer.ResolveDeviceService().Context()
	configuration := appIocContainer.ResolveWhatsappConfiguration()

	abProp := fmt.Sprintf("cc=%v&in=%v&rc=0", device.Area, device.Phone)
	return configuration.HttpUrl + "/v2/reg_onboard_abprop?" + abProp
}
