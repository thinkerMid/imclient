package registerRequest

import (
	"fmt"
	hertzConst "github.com/cloudwego/hertz/pkg/protocol/consts"
	"ws/framework/application/container/abstract_interface"
	httpApi "ws/framework/plugin/network/http_api"
)

func SendAppLaunch(appIocContainer containerInterface.IAppIocContainer) error {
	var resp PhoneExistRep

	client := appIocContainer.ResolveHttpClient()

	err := httpApi.DoAndBind(
		client, &resp,
		httpApi.Url(MakeEmptyExistBody(appIocContainer)),
		httpApi.Method(hertzConst.MethodGet),
		httpApi.Header("Content-Type", "text/json;charset=utf-8"),
		httpApi.UserAgent(appIocContainer.ResolveDeviceService().DeviceAgent()),
	)

	if err != nil {
		return fmt.Errorf("网络异常")
	}

	return resp.HasError()
}

func SendGetABProp(appIocContainer containerInterface.IAppIocContainer) error {
	var resp PhoneABPropResp

	client := appIocContainer.ResolveHttpClient()

	err := httpApi.DoAndBind(
		client, &resp,
		httpApi.Url(MakePhoneABPropBody(appIocContainer)),
		httpApi.Method(hertzConst.MethodGet),
		httpApi.Header("Content-Type", "text/json;charset=utf-8"),
		httpApi.UserAgent(appIocContainer.ResolveDeviceService().DeviceAgent()),
	)

	if err != nil {
		return fmt.Errorf("网络异常")
	}

	return resp.HasError()
}

// CheckPhoneExist 检测手机号是否存在
func CheckPhoneExist(appIocContainer containerInterface.IAppIocContainer) error {
	var resp PhoneExistRep

	client := appIocContainer.ResolveHttpClient()

	err := httpApi.DoAndBind(
		client, &resp,
		httpApi.Url(MakePhoneExistBody(appIocContainer)),
		httpApi.Method(hertzConst.MethodGet),
		httpApi.Header("Content-Type", "text/json;charset=utf-8"),
		httpApi.UserAgent(appIocContainer.ResolveDeviceService().DeviceAgent()),
	)

	if err != nil {
		return fmt.Errorf("网络异常")
	}

	return resp.HasError()
}

// SendClientLog .
func SendClientLog(appIocContainer containerInterface.IAppIocContainer, currentScreen string, previousScreen string, actionTaken string) error {
	var clientSendLog ClientSendLogResp

	client := appIocContainer.ResolveHttpClient()

	// 发送日志
	err := httpApi.DoAndBind(
		client,
		&clientSendLog,
		httpApi.Url(MakeClientLogBody(appIocContainer, currentScreen, previousScreen, actionTaken)),
		httpApi.Method(hertzConst.MethodGet),
		httpApi.Header("Content-Type", "text/json;charset=utf-8"),
		httpApi.UserAgent(appIocContainer.ResolveDeviceService().DeviceAgent()),
	)

	if err != nil {
		return fmt.Errorf("网络异常")
	}

	return clientSendLog.HasError()
}

// GetSmsCode 获取验证码
func GetSmsCode(appIocContainer containerInterface.IAppIocContainer) error {
	var resp PhoneGetCodeRep

	client := appIocContainer.ResolveHttpClient()

	err := httpApi.DoAndBind(
		client, &resp,
		httpApi.Url(MakePhoneGetSmsCodeBody(appIocContainer)),
		httpApi.Method(hertzConst.MethodGet),
		httpApi.Header("Content-Type", "text/json;charset=utf-8"),
		httpApi.UserAgent(appIocContainer.ResolveDeviceService().DeviceAgent()),
	)

	if err != nil {
		return fmt.Errorf("网络异常")
	}

	return resp.HasError()
}

// SendSmsCode 发送验证码
func SendSmsCode(appIocContainer containerInterface.IAppIocContainer, smsCode string) (*PhoneSendCodeRep, error) {
	var resp PhoneSendCodeRep

	client := appIocContainer.ResolveHttpClient()

	err := httpApi.DoAndBind(
		client,
		&resp,
		httpApi.Url(MakePhoneSendSmsCodeBody(appIocContainer, smsCode)),
		httpApi.Method(hertzConst.MethodGet),
		httpApi.Header("Content-Type", "text/json;charset=utf-8"),
		httpApi.UserAgent(appIocContainer.ResolveDeviceService().DeviceAgent()),
	)

	if err != nil {
		return &resp, fmt.Errorf("网络异常")
	}

	return &resp, resp.HasError()
}
