package events

import eventSerialize "ws/framework/plugin/event_serialize"

type WamEventUiUsage struct {
	WAMessageEvent

	EntryPoint float64
	UIUsage    float64
}

func (event *WamEventUiUsage) InitFields(option interface{}) {
	event.EntryPoint = 10 //?
	event.UIUsage = 1
}

func (event *WamEventUiUsage) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Body().
		SerializeNumber(0x2, event.UIUsage)

	buffer.Footer().
		SerializeNumber(0x1, event.EntryPoint)
}
