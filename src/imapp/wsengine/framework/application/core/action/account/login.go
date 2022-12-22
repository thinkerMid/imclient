package accountState

import (
	containerInterface "ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	"ws/framework/application/core/wam"
	accountDB "ws/framework/application/data_storage/account/database"
)

// LoginSuccess .
type LoginSuccess struct {
	processor.BaseAction
	FirstLogin bool
}

// Start .
func (m *LoginSuccess) Start(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	// 日志 注册后首登/其他正常登录
	if m.FirstLogin {
		wam.LogManager().LogRegisterLaunch(context)
	} else {
		wam.LogManager().LogLogin(context)
	}

	context.ResolveAccountService().ContextExecute(func(account *accountDB.Account) {
		account.AddLoginCount()
		account.UpdateAccountStatus(0)
	})

	next()

	return nil
}

// Receive .
func (m *LoginSuccess) Receive(_ containerInterface.IMessageContext, _ containerInterface.NextActionFn) error {
	return nil
}

// Error .
func (m *LoginSuccess) Error(_ containerInterface.IMessageContext, _ error) {}
