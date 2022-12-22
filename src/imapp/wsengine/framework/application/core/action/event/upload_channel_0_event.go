package event

import (
	"strconv"
	"time"
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	accountDB "ws/framework/application/data_storage/account/database"
	eventSerialize "ws/framework/plugin/event_serialize"
)

// UploadChannel0Event 渠道0
type UploadChannel0Event struct {
	processor.BaseAction
}

// Start .
func (u *UploadChannel0Event) Start(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	wec := context.ResolveChannel0EventCache()
	accountService := context.ResolveAccountService()

	if wec.CacheBufferItem() == 0 {
		accountService.ContextExecute(func(account *accountDB.Account) {
			account.UpdateSendChannel0Time()
		})

		next()
		return
	}

	buffer := eventSerialize.AcquireEventBuffer()
	defer eventSerialize.ReleaseEventBuffer(buffer)

	wec.PackBuffer(accountService.Context().SendChannel0EventCount, buffer)

	// [iq id=<'1643273074-9'> xmlns=<'w:stats'> type=<'set'> to=<s.whatsapp.net> [add t=<'1643273378'> {327b} ] ]
	// <iq id="1643872228-1" to="s.whatsapp.net" type="set" xmlns="w:stats"><add t="1643872528"><!-- 1352 bytes --></add></iq>
	u.SendMessageId, err = context.SendIQ(createEventIQ(context.GenerateRequestID(), buffer.Byte()))

	return
}

// Receive .
func (u *UploadChannel0Event) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	context.ResolveChannel0EventCache().ClearLog()

	context.ResolveAccountService().ContextExecute(func(account *accountDB.Account) {
		account.AddChannel0EventCount()
		account.UpdateSendChannel0Time()
	})

	next()
	return
}

// Error .
func (u *UploadChannel0Event) Error(_ containerInterface.IMessageContext, _ error) {}

func createEventIQ(iqId string, content []byte) message.InfoQuery {
	return message.InfoQuery{
		ID:        iqId,
		Namespace: "w:stats",
		Type:      message.IqSet,
		To:        types.ServerJID,
		Content: []waBinary.Node{
			{
				Tag: "add",
				Attrs: waBinary.Attrs{
					"t": strconv.FormatInt(time.Now().Unix(), 10),
				},
				Content: content,
			},
		},
	}
}
