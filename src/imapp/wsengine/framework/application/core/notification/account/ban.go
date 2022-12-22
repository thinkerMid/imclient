package accountNotification

import (
	"strconv"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/result/constant"
	accountDB "ws/framework/application/data_storage/account/database"
)

// Ban .
type Ban struct{}

// Receive .
func (m Ban) Receive(context containerInterface.IMessageContext) (err error) {
	reasonCode := context.Message().AttrGetter().String("reason")

	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.Ban,
		Content:    reasonCode,
	})

	parseInt, _ := strconv.ParseInt(reasonCode, 10, 16)

	context.ResolveAccountService().ContextExecute(func(account *accountDB.Account) {
		account.UpdateAccountStatus(int16(parseInt))
	})

	return
}
