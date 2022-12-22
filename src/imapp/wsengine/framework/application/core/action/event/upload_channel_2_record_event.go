package event

import (
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	accountDB "ws/framework/application/data_storage/account/database"
)

// UploadChannel2RecordEvent 渠道2
type UploadChannel2RecordEvent struct {
	processor.BaseAction
}

// Start .
func (u *UploadChannel2RecordEvent) Start(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	wec := context.ResolveChannel2EventCache()
	accountService := context.ResolveAccountService()

	if !accountService.NeedUploadRecordChannel2Event() {
		next()
		return
	}

	if wec.CacheBufferItem() == 0 {
		accountService.ContextExecute(func(account *accountDB.Account) {
			account.UpdateSendChannel2Time()
		})

		next()
		return
	}

	u.Query = &PrivateStats{}
	return u.Query.Start(context, next)
}

// Receive .
func (u *UploadChannel2RecordEvent) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	defer next()

	err = u.Query.Receive(context, next)
	if err != nil {
		return err
	}

	context.ResolveChannel2EventCache().ClearNotSentYetLog()

	context.ResolveAccountService().ContextExecute(func(account *accountDB.Account) {
		account.AddChannel2EventCount()
		account.UpdateSendChannel2Time()
	})

	return
}

// Error .
func (u *UploadChannel2RecordEvent) Error(_ containerInterface.IMessageContext, _ error) {}
