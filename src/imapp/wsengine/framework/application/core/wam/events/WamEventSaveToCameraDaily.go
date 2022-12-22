package events

import eventSerialize "ws/framework/plugin/event_serialize"

type WamEventSaveToCameraDaily struct {
	WAMessageEvent

	ChatsAlways   float64 // 0x1
	ChatsDefault  float64 // 0x2
	ChatsDisabled float64 // 0x5
	ChatsNever    float64 // 0x3
	SettingCount  float64 // 0x4
}

func (event *WamEventSaveToCameraDaily) InitFields(option interface{}) {
	event.ChatsAlways = 0
	event.ChatsDefault = 1
	event.ChatsDisabled = 0
	event.ChatsNever = 0
	event.SettingCount = 0
}

func (event *WamEventSaveToCameraDaily) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Body().
		SerializeNumber(0x1, event.ChatsAlways).
		SerializeNumber(0x2, event.ChatsDefault).
		SerializeNumber(0x5, event.ChatsDisabled).
		SerializeNumber(0x3, event.ChatsNever)

	buffer.Footer().
		SerializeNumber(0x4, event.SettingCount)
}
