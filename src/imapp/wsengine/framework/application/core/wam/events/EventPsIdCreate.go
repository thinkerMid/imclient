package events

import eventSerialize "ws/framework/plugin/event_serialize"

type WamEventPsIdCreate struct {
	WAMessageEvent
}

func (event *WamEventPsIdCreate) InitFields(option interface{}) {
}

func (event *WamEventPsIdCreate) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Footer().SerializeNumber(event.Code, event.Weight)
}
