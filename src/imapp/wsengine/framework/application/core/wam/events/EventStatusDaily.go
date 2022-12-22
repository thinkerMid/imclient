package events

import eventSerialize "ws/framework/plugin/event_serialize"

type WamEventStatusDaily struct {
	WAMessageEvent

	AvailableCount     float64
	AvailableRowsCount float64
	ViewCount          float64
	ViewRowsCount      float64
}

func (event *WamEventStatusDaily) InitFields(option interface{}) {
	event.AvailableCount = 0.000000
	event.AvailableRowsCount = 0.000000
	event.ViewCount = 0.000000
	event.ViewRowsCount = 0.000000
}

func (event *WamEventStatusDaily) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Body().
		SerializeNumber(0x3, event.AvailableCount).
		SerializeNumber(0x2, event.AvailableRowsCount).
		SerializeNumber(0x4, event.ViewCount)

	buffer.Footer().
		SerializeNumber(0x2, event.ViewRowsCount)
}
