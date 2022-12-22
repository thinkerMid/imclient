package events

import (
	"time"
	eventSerialize "ws/framework/plugin/event_serialize"
	"ws/framework/utils"
)

type WamEventIphoneIcloudUbt struct {
	WAMessageEvent

	LoadTime float64
}

func (event *WamEventIphoneIcloudUbt) InitFields(option interface{}) {
	event.LoadTime = utils.LogRandSecond(time.Second, 20*time.Second)
}

func (event *WamEventIphoneIcloudUbt) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Footer().
		SerializeNumber(1, event.LoadTime)
}
