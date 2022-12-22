package im

import (
	"fmt"
	"ws/business/im/control"
	"ws/framework/application"
	waProto "ws/framework/application/constant/binary/proto"
	appContainer "ws/framework/application/container"
	containerInterface "ws/framework/application/container/abstract_interface"
	"ws/framework/config"
	"ws/framework/external"
	networkConstant "ws/framework/plugin/network/constant"
)

// NewClient .
func NewClient(application application.IApplication) containerInterface.IIMClient {
	return &Client{application: application}
}

// Client .
type Client struct {
	application        application.IApplication
	imControl          control.IMControl
	messageMarkReadPID uint32
}

func (c *Client) checkAccountStatus() error {
	account := c.application.Container().ResolveAccountService().Context()

	if account == nil || !account.IsRegistered() || !account.BusinessAccount {
		return external.AccountNotFoundErr
	}

	// 如果返回不可用状态 不需要再走后面的逻辑了
	if !account.AvailableStatus() {
		return fmt.Errorf("%v", account.Status)
	}

	return nil
}

// Connect .
func (c *Client) Connect(optsFn ...config.OptionsFn) (err error) {
	err = c.checkAccountStatus()
	if err != nil {
		return err
	}

	opts := config.Options{}
	for _, optFn := range optsFn {
		optFn(&opts)
	}

	if len(opts.ConnectionConfig.ProxyAddress) == 0 {
		opts.ConnectionConfig.Type = networkConstant.Socket
	}

	c.imControl.AutoMessageMarkRead = opts.AutoMessageMarkRead

	c.application.Container().Inject(appContainer.ConnectionConfig, opts.ConnectionConfig)
	c.application.Container().Inject(appContainer.IMControl, &c.imControl)

	return c.application.Start()
}

// EnterScene .
func (c *Client) EnterScene(scene containerInterface.IScene, resultProcessor containerInterface.LocalResultProcessor) {
	channel := c.application.Container().ResolveMessageChannel()

	channel.AddProcessorAndAttach(scene.Build(), resultProcessor)
}

// AddGlobalResultProcessor .
func (c *Client) AddGlobalResultProcessor(p containerInterface.GlobalResultProcessor) {
	channel := c.application.Container().ResolveMessageChannel()
	channel.AddGlobalResultProcessor(p)
}

// EnableAutoMessageMarkRead .
func (c *Client) EnableAutoMessageMarkRead() {
	c.imControl.EnableAutoMessageMarkRead()
}

// DisableMessageMarkRead .
func (c *Client) DisableMessageMarkRead() {
	c.imControl.DisableMessageMarkRead()
}

// Logout .
func (c *Client) Logout() {
	c.application.Exit()
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
