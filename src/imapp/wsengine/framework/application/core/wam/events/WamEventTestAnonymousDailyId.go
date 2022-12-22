package events

import (
	"time"
	eventSerialize "ws/framework/plugin/event_serialize"
	"ws/framework/utils"
)

type WamEventTestAnonymousDailyId struct {
	WAMessageEvent

	PsTestEnumField  float64
	PsTestFloatField float64
}

func (event *WamEventTestAnonymousDailyId) InitFields(option interface{}) {
	event.PsTestEnumField = 1
	event.PsTestFloatField = utils.LogRandSecond(1*time.Second, 2*time.Second)
}

func (event *WamEventTestAnonymousDailyId) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Body().
		SerializeNumber(0x1, event.PsTestEnumField)

	buffer.Footer().
		SerializeNumber(0x2, event.PsTestFloatField)
}
