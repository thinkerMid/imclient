package events

import (
	eventSerialize "ws/framework/plugin/event_serialize"
)

type WamEventICloudRestore struct {
	WAMessageEvent

	RestoreResult  float64
	StartReason    float64
	RestoreVersion float64
}

func (event *WamEventICloudRestore) InitFields(option interface{}) {
	event.RestoreResult = 1
	event.StartReason = 1
	event.RestoreVersion = 0
}

func (event *WamEventICloudRestore) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Body().
		SerializeNumber(0x1, event.RestoreResult).
		SerializeNumber(0x2, event.StartReason)

	buffer.Footer().
		SerializeNumber(0x7, event.RestoreVersion)
}
