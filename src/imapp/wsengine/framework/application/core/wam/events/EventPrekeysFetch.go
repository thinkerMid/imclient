package events

import eventSerialize "ws/framework/plugin/event_serialize"

type WamEventPrekeysFetch struct {
	WAMessageEvent

	FetchContext   int32
	IdentityChange float64
	FetchCount     float64
}

func (event *WamEventPrekeysFetch) InitFields(option interface{}) {
	event.FetchContext = 1
	event.IdentityChange = 0
	event.FetchCount = 1
}

func (event *WamEventPrekeysFetch) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Body().
		SerializeNumber(0x3, float64(event.FetchContext)).
		SerializeNumber(0x1, event.IdentityChange)

	buffer.Footer().
		SerializeNumber(0x2, event.FetchCount)
}
