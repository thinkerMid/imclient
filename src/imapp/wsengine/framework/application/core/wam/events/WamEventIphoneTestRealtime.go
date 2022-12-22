package events

import eventSerialize "ws/framework/plugin/event_serialize"

type WamEventIphoneTestRealtime struct {
	WAMessageEvent

	IsConnectedToChat float64
}

func (event *WamEventIphoneTestRealtime) InitFields(option interface{}) {
	event.IsConnectedToChat = 1
}

func (event *WamEventIphoneTestRealtime) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Footer().
		SerializeNumber(0x1, event.IsConnectedToChat)
}
