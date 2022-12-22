package events

import (
	eventSerialize "ws/framework/plugin/event_serialize"
)

type EventIphoneGroupCreate struct {
	WAMessageEvent

	EntryPoint float64
}

func (event *EventIphoneGroupCreate) InitFields(option interface{}) {
	event.EntryPoint = 8 //? 同EventGroupCreateInit
}

func (event *EventIphoneGroupCreate) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Footer().
		SerializeNumber(0x1, event.EntryPoint)
}
