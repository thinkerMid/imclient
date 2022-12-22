package accountState

import (
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	"ws/framework/application/core/result/constant"
	accountDB "ws/framework/application/data_storage/account/database"
)

// Offline .
type Offline struct {
	processor.BaseAction
}

// Start .
func (m *Offline) Start(_ containerInterface.IMessageContext, _ containerInterface.NextActionFn) error {
	return nil
}

// Receive .
func (m *Offline) Receive(_ containerInterface.IMessageContext, _ containerInterface.NextActionFn) error {
	return nil
}

// Error .
func (m *Offline) Error(context containerInterface.IMessageContext, err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.Offline,
		Error:      err,
	})

	accountService := context.ResolveAccountService()
	event0 := context.ResolveChannel0EventCache()
	event2 := context.ResolveChannel2EventCache()
	account := accountService.Context()

	if account.Status > 400 && account.Status < 500 {
		event0.CleanupAllData()
		event2.CleanupAllData()
	} else {
		event0.FlushEventCache()
		event2.FlushEventCache()
	}

	accountService.ContextExecute(func(account *accountDB.Account) {
		account.UpdateLogoutTime()
	})
}
