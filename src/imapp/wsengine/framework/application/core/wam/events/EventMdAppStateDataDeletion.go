package events

import (
	eventSerialize "ws/framework/plugin/event_serialize"
)

type WamEventMdAppStateDataDeletion struct {
	WAMessageEvent

	RetryCount   float64
	ReasonNumber float64
}

func (event *WamEventMdAppStateDataDeletion) InitFields(option interface{}) {
	event.ReasonNumber = 2
	event.RetryCount = 0
}

func (event *WamEventMdAppStateDataDeletion) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Body().
		SerializeNumber(0x1, event.ReasonNumber)

	buffer.Footer().
		SerializeNumber(0x2, event.RetryCount)
}
