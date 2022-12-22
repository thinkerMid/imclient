package registerRequest

import (
	"fmt"
	containerInterface "ws/framework/application/container/abstract_interface"
)

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
func MakePhoneABPropBody(appIocContainer containerInterface.IAppIocContainer) string {
	device := appIocContainer.ResolveDeviceService().Context()
	configuration := appIocContainer.ResolveWhatsappConfiguration()

	abProp := fmt.Sprintf("cc=%v&in=%v&rc=0", device.Area, device.Phone)
	return configuration.HttpUrl + "/v2/reg_onboard_abprop?" + abProp
}
