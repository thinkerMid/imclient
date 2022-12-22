package event

import (
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	"ws/framework/application/data_storage/account/database"
)

// UploadChannel2Event 渠道2
type UploadChannel2Event struct {
	processor.BaseAction
}

// Start .
func (u *UploadChannel2Event) Start(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	wec := context.ResolveChannel2EventCache()
	accountService := context.ResolveAccountService()

	if wec.CacheBufferItem() == 0 {
		accountService.ContextExecute(func(account *accountDB.Account) {
			account.UpdateSendChannel2Time()
		})

		next()
		return
	}

	u.Query = &PrivateStats{}
	return u.Query.Start(context, func() {})
}

// Receive .
func (u *UploadChannel2Event) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	defer next()

	err = u.Query.Receive(context, func() {})
	if err != nil {
		return err
	}

	context.ResolveChannel2EventCache().ClearLog()

	context.ResolveAccountService().ContextExecute(func(account *accountDB.Account) {
		account.AddChannel2EventCount()
		account.UpdateSendChannel2Time()
	})
	return
}

// Error .
func (u *UploadChannel2Event) Error(_ containerInterface.IMessageContext, _ error) {}
