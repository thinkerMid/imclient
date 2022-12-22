package events

import eventSerialize "ws/framework/plugin/event_serialize"

type EventTestAnonymousDaily struct {
	WAMessageEvent
}

func (event *EventTestAnonymousDaily) InitFields(option interface{}) {

}

func (event *EventTestAnonymousDaily) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Footer().SerializeNumber(event.Code, event.Weight)
}
