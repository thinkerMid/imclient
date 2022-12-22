package events

import (
	eventSerialize "ws/framework/plugin/event_serialize"
)

type EventGroupCreateInit struct {
	WAMessageEvent

	EntryPoint float64
}

func (event *EventGroupCreateInit) InitFields(option interface{}) {
	event.EntryPoint = 8 //?
}

func (event *EventGroupCreateInit) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Footer().
		SerializeNumber(0x1, event.EntryPoint)
}
