package anonymous

import (
	"fmt"
	"ws/framework/application"
	waProto "ws/framework/application/constant/binary/proto"
	"ws/framework/application/constant/types"
	appContainer "ws/framework/application/container"
	containerInterface "ws/framework/application/container/abstract_interface"
	accountServiceConstant "ws/framework/application/data_storage/account/constant"
	accountDB "ws/framework/application/data_storage/account/database"
	deviceDB "ws/framework/application/data_storage/device/database"
	"ws/framework/config"
	"ws/framework/lib/firmware"
	"ws/framework/lib/msisdn"
	networkConstant "ws/framework/plugin/network/constant"
	"ws/messenger/anonymous/import_device"
	"ws/messenger/anonymous/register_request"
	"ws/messenger/anonymous/unblock"
)

type virtualDevice struct {
	msisdn.IMSIParser
	firmware.Apple
}

// NewClient .
func NewClient(application application.IApplication) containerInterface.IAnonymousClient {
	return &Client{application: application}
}

// Client .
type Client struct {
	application application.IApplication
}

func (c *Client) initOptions(optsFn ...config.OptionsFn) {
	opts := config.Options{}
	for _, optFn := range optsFn {
		optFn(&opts)
	}

	if len(opts.ConnectionConfig.ProxyAddress) == 0 {
		opts.ConnectionConfig.Type = networkConstant.Socket
	}

	c.application.Container().Inject(appContainer.ConnectionConfig, opts.ConnectionConfig)
}

func (c *Client) newDevice() {
	container := c.application.Container()

	imsi, _ := msisdn.Parse(container.ResolveJID().User, false)

	virtualDevice := virtualDevice{IMSIParser: imsi, Apple: firmware.NewAppleFirmware()}
	_, _ = container.ResolveDeviceService().Create(&virtualDevice)
	_, _ = container.ResolveSignedPreKeyService().Create()
	_, _ = container.ResolveRegistrationTokenService().Create()
	_, _ = container.ResolveAesKeyService().Create()
	_, _ = container.ResolveAccountService().Create()
}

// JID .
func (c *Client) JID() string {
	return c.application.Container().ResolveJID().User
}

// GetSmsCode .
func (c *Client) GetSmsCode(optsFn ...config.OptionsFn) (err error) {
	container := c.application.Container()

	c.initOptions(optsFn...)

	if container.ResolveAccountService().Context() == nil {
		c.newDevice()
	}

	// 启动app发送一个空包
	_ = registerRequest.SendAppLaunch(container)

	// 发送abprop
	_ = registerRequest.SendGetABProp(container)

	// 检查手机号是否被注册
	if err = registerRequest.CheckPhoneExist(container); err != nil {
		return
	}

	// 发送检查后的日志
	_ = registerRequest.SendClientLog(container, "verify_sms", "enter_number", "continue")

	// 获取验证码
	return registerRequest.GetSmsCode(container)
}

// SendReceiveSmsCode .
func (c *Client) SendReceiveSmsCode(smsCode string, optsFn ...config.OptionsFn) (err error) {
	container := c.application.Container()

	c.initOptions(optsFn...)

	if container.ResolveAccountService().Context() == nil {
		c.newDevice()
	}

	// 发送验证码
	resp, err := registerRequest.SendSmsCode(container, smsCode)
	if err != nil {
		return err
	}

	// 注册成功发2次日志
	_ = registerRequest.SendClientLog(container, "no_backup_found", "verify_sms", "continue")
	_ = registerRequest.SendClientLog(container, "no_backup_found", "no_backup_found", "continue")

	// 重置账号状态
	container.ResolveAccountService().ContextExecute(func(account *accountDB.Account) {
		account.UpdateAccountStatus(accountServiceConstant.Registered)
	})

	// 修改用户名
	container.ResolveDeviceService().ContextExecute(func(device *deviceDB.Device) {
		device.UpdateUserName(resp.Login)
	})

	// 返回的JID不一致
	if resp.Login != container.ResolveJID().User {
		container.OnJIDChangeWhenRegisterSuccess(types.NewJID(resp.Login, types.DefaultUserServer))
	}

	return
}

// ImportDeviceData 导入数据
func (c *Client) ImportDeviceData(data *importData.DeviceData) error {
	container := c.application.Container()

	if c.application.Container().ResolveAccountService().Context() != nil {
		return fmt.Errorf("设备信息已存在")
	}

	err := importData.Do(container, data)
	if err != nil {
		container.CleanupDataStorage()
	}

	return err
}

// Unblock 解封
func (c *Client) Unblock(optsFn ...config.OptionsFn) (resp *unblock.UnblockContentResponse, err error) {
	container := c.application.Container()

	c.initOptions(optsFn...)

	account := container.ResolveAccountService().Context()

	if account == nil {
		// 新建一个设备
		c.newDevice()
		// 存在已注册的账号数据
	} else if account.IsRegistered() {
		// 仍然可用状态
		if account.AvailableStatus() {
			return nil, fmt.Errorf("账号不需要解封")
		}

		// 清理数据
		container.CleanupDataStorage()
		// 重新创建设备
		c.newDevice()
	}

	resp, err = unblock.Do(container)

	return
}

// CleanupDataStorage .
func (c *Client) CleanupDataStorage() {
	c.application.Container().CleanupDataStorage()
}

// Version .
func (c *Client) Version() string {
	return c.application.Container().ResolveWhatsappConfiguration().VersionString
}

// Platform .
func (c *Client) Platform() waProto.UserAgent_UserAgentPlatform {
	return c.application.Container().ResolveWhatsappConfiguration().Platform
}
