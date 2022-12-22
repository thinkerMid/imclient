package event

import (
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	accountDB "ws/framework/application/data_storage/account/database"
	eventSerialize "ws/framework/plugin/event_serialize"
)

// UploadChannel0RecordEvent 渠道0
type UploadChannel0RecordEvent struct {
	processor.BaseAction
}

// Start .
func (u *UploadChannel0RecordEvent) Start(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	wec := context.ResolveChannel0EventCache()
	accountService := context.ResolveAccountService()

	if !accountService.NeedUploadRecordChannel0Event() {
		next()
		return
	}

	if wec.CacheBufferItem() == 0 {
		accountService.ContextExecute(func(account *accountDB.Account) {
			account.UpdateSendChannel0Time()
		})
		next()

		return
	}

	buffer := eventSerialize.AcquireEventBuffer()
	defer eventSerialize.ReleaseEventBuffer(buffer)

	wec.PackNotSentYetLog(accountService.Context().SendChannel0EventCount, buffer)

	u.SendMessageId, err = context.SendIQ(createEventIQ(context.GenerateRequestID(), buffer.Byte()))

	return
}

// Receive .
func (u *UploadChannel0RecordEvent) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	context.ResolveChannel0EventCache().ClearNotSentYetLog()

	context.ResolveAccountService().ContextExecute(func(account *accountDB.Account) {
		account.AddChannel0EventCount()
		account.UpdateSendChannel0Time()
	})

	next()
	return
}

// Error .
func (u *UploadChannel0RecordEvent) Error(_ containerInterface.IMessageContext, _ error) {}
