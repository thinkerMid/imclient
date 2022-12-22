package events

import eventSerialize "ws/framework/plugin/event_serialize"

type WamEventContactInfo struct {
	WAMessageEvent

	ContactInfoVisit float64
	Mute             float64
}

func (event *WamEventContactInfo) InitFields(option interface{}) {
	event.ContactInfoVisit = 1
	event.Mute = 1
}

func (event *WamEventContactInfo) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Body().
		SerializeNumber(0x1, event.ContactInfoVisit)

	buffer.Footer().
		SerializeNumber(0xb, event.Mute)
}
