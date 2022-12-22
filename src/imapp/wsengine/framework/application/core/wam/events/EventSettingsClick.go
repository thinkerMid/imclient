package events

import eventSerialize "ws/framework/plugin/event_serialize"

type WamEventSettingsClick struct {
	WAMessageEvent

	ItemNumber float64 //?
}

func (event *WamEventSettingsClick) InitFields(option interface{}) {
	event.ItemNumber = 0
}

func (event *WamEventSettingsClick) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Footer().
		SerializeNumber(0x1, event.ItemNumber)
}
