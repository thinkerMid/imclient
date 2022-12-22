package events

import (
	"time"
	eventSerialize "ws/framework/plugin/event_serialize"
	"ws/framework/utils"
)

type WamEventLogin struct {
	WAMessageEvent

	ConnectionOrigin float64 //0x6

	LoginResult float64 //0x1
	LoginTime   float64 //0x3
	LongConnect float64 //0x4
	Passive     float64 //0x8
	RetryCount  float64 //0x2

}

func (event *WamEventLogin) InitFields(option interface{}) {
	event.ConnectionOrigin = 3
	event.LoginResult = 1
	event.LoginTime = utils.LogRandMillSecond(0, 2*time.Second)
	event.LongConnect = 0
	event.Passive = 1
	event.RetryCount = 0

}

func (event *WamEventLogin) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Body().
		SerializeNumber(0x6, event.ConnectionOrigin).
		SerializeNumber(0x1, event.LoginResult).
		SerializeNumber(0x3, event.LoginTime).
		SerializeNumber(0x4, event.LongConnect).
		SerializeNumber(0x8, event.Passive)

	buffer.Footer().
		SerializeNumber(0x2, event.RetryCount)
}
