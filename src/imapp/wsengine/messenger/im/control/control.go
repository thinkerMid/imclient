package control

import (
	"ws/framework/application/constant/message"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/notification/private_chat"
	"ws/framework/application/core/processor"
)

var _ containerInterface.IIMControl = &IMControl{}

// IMControl .
type IMControl struct {
	containerInterface.BaseService

	AutoMessageMarkRead bool

	messageMarkReadPID uint32
	shortReconnect     bool
}

// OnApplicationStart .
func (c *IMControl) OnApplicationStart() {
	channel := c.AppIocContainer.ResolveMessageChannel()
	account := c.AppIocContainer.ResolveAccountService().Context()

	// 首次登录
	if account.FirstLogin() {
		_ = c.AppIocContainer.ResolveRegistrationTokenService().RefreshToken()

		initFirstLoginProcessor(c.AppIocContainer)
	} else {
		initLoginProcessor(c.AppIocContainer)
		c.EnableAutoMessageMarkRead()
	}

	// 通知型
	initNotificationProcessor(channel)
	// 定时型
	initTimerProcessor(channel)
}

// OnApplicationResume .
func (c *IMControl) OnApplicationResume() {
	// 只有非短重连的情况下 才会初始化登录逻辑
	// note: 短重连指的是 action主动发送 <xmlstreamend/> 然后等待连接重连后继续执行的场景
	if !c.shortReconnect {
		initLoginProcessor(c.AppIocContainer)
	}

	initNotificationProcessor(c.AppIocContainer.ResolveMessageChannel())

	// 重置
	c.shortReconnect = false
}

// ActiveDisconnectXMPP .
func (c *IMControl) ActiveDisconnectXMPP() {
	c.shortReconnect = true
}

// EnableAutoMessageMarkRead .
func (c *IMControl) EnableAutoMessageMarkRead() {
	if c.messageMarkReadPID > 0 || !c.AutoMessageMarkRead {
		return
	}

	channel := c.AppIocContainer.ResolveMessageChannel()

	pid := channel.AddMessageProcessor(processor.NewNotificationProcessor(
		func() []containerInterface.INotification {
			return []containerInterface.INotification{
				privateChatNotification.DelayReceiveMessageMarkRead{},
			}
		},
		processor.TriggerTag(message.ReceiveMessage),
		processor.Priority(processor.PriorityForeground),
	))

	c.messageMarkReadPID = pid
}

// DisableMessageMarkRead .
func (c *IMControl) DisableMessageMarkRead() {
	if c.messageMarkReadPID == 0 {
		return
	}

	c.messageMarkReadPID = 0
	c.AppIocContainer.ResolveMessageChannel().RemoveProcessor(c.messageMarkReadPID)
}
