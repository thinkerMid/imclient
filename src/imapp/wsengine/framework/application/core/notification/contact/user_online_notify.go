package contactNotification

import (
	"strconv"
	"time"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/result/constant"
	"ws/framework/external"
)

// UserOnlineNotify .
type UserOnlineNotify struct{}

// [presence from=<8618898739725@s.whatsapp.net> last=<10 bytes> type=<'unavailable'> ]

// Receive .
func (m UserOnlineNotify) Receive(context containerInterface.IMessageContext) (err error) {
	attrs := context.Message().AttrGetter()

	JID := attrs.JID("from")
	_, typeOk := attrs.GetString("type", true)
	timeStr, lastTimeOk := attrs.GetString("last", true)

	var timeNumber int64
	// 离线通知
	if typeOk && lastTimeOk {
		parseInt, parseErr := strconv.ParseInt(timeStr, 10, 64)
		if parseErr != nil {
			parseInt = time.Now().Unix()
		}

		timeNumber = parseInt
	} else {
		// 在线通知
		timeNumber = time.Now().Unix()
	}

	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.UserLastOnlineNotify,
		IContent: external.OnlineNotify{
			JIDNumber: JID.User,
			Time:      timeNumber,
		},
	})

	return
}
