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

	// ??????app??????????????????
	_ = registerRequest.SendAppLaunch(container)

	// ??????abprop
	_ = registerRequest.SendGetABProp(container)

	// ??????????????????????????????
	if err = registerRequest.CheckPhoneExist(container); err != nil {
		return
	}

	// ????????????????????????
	_ = registerRequest.SendClientLog(container, "verify_sms", "enter_number", "continue")

	// ???????????????
	return registerRequest.GetSmsCode(container)
}

// SendReceiveSmsCode .
func (c *Client) SendReceiveSmsCode(smsCode string, optsFn ...config.OptionsFn) (err error) {
	container := c.application.Container()

	c.initOptions(optsFn...)

	if container.ResolveAccountService().Context() == nil {
		c.newDevice()
	}

	// ???????????????
	resp, err := registerRequest.SendSmsCode(container, smsCode)
	if err != nil {
		return err
	}

	// ???????????????2?????????
	_ = registerRequest.SendClientLog(container, "no_backup_found", "verify_sms", "continue")
	_ = registerRequest.SendClientLog(container, "no_backup_found", "no_backup_found", "continue")

	// ??????????????????
	container.ResolveAccountService().ContextExecute(func(account *accountDB.Account) {
		account.UpdateAccountStatus(accountServiceConstant.Registered)
	})

	// ???????????????
	container.ResolveDeviceService().ContextExecute(func(device *deviceDB.Device) {
		device.UpdateUserName(resp.Login)
	})

	// ?????????JID?????????
	if resp.Login != container.ResolveJID().User {
		container.OnJIDChangeWhenRegisterSuccess(types.NewJID(resp.Login, types.DefaultUserServer))
	}

	return
}

// ImportDeviceData ????????????
func (c *Client) ImportDeviceData(data *importData.DeviceData) error {
	container := c.application.Container()

	if c.application.Container().ResolveAccountService().Context() != nil {
		return fmt.Errorf("?????????????????????")
	}

	err := importData.Do(container, data)
	if err != nil {
		container.CleanupDataStorage()
	}

	return err
}

// Unblock ??????
func (c *Client) Unblock(optsFn ...config.OptionsFn) (resp *unblock.UnblockContentResponse, err error) {
	container := c.application.Container()

	c.initOptions(optsFn...)

	account := container.ResolveAccountService().Context()

	if account == nil {
		// ??????????????????
		c.newDevice()
		// ??????????????????????????????
	} else if account.IsRegistered() {
		// ??????????????????
		if account.AvailableStatus() {
			return nil, fmt.Errorf("?????????????????????")
		}

		// ????????????
		container.CleanupDataStorage()
		// ??????????????????
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
