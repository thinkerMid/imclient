package messenger

import (
	"ws/framework/application"
	containerInterface "ws/framework/application/container/abstract_interface"
	"ws/framework/lib/msisdn"
	"ws/messenger/anonymous"
	"ws/messenger/im"
)

// NewAnonymousClient 账号注册功能
func NewAnonymousClient(phoneNumber string) (containerInterface.IAnonymousClient, error) {
	imsi, err := msisdn.Parse(phoneNumber, true)
	if err != nil {
		return nil, err
	}

	app := application.New(imsi.GetCC()+imsi.GetPhoneNumber(), configuration)

	return anonymous.NewClient(app), nil
}

// NewIMClient 即时通讯功能的客户端
func NewIMClient(jid string) containerInterface.IIMClient {
	return im.NewClient(application.New(jid, configuration))
}
