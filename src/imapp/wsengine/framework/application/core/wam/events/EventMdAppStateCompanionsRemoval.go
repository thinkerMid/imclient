package events

import (
	eventSerialize "ws/framework/plugin/event_serialize"
)

type WamEventMdAppStateCompanionsRemoval struct {
	WAMessageEvent

	RetryCount float64
}

func (event *WamEventMdAppStateCompanionsRemoval) InitFields(option interface{}) {
	event.RetryCount = 0
}

func (event *WamEventMdAppStateCompanionsRemoval) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Footer().
		SerializeNumber(0x1, event.RetryCount)
}
