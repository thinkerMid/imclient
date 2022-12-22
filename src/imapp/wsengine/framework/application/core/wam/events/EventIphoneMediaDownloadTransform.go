package events

import (
	eventSerialize "ws/framework/plugin/event_serialize"
)

type WamEventIphoneMediaDownloadTransform struct {
	WAMessageEvent

	Amr  float64
	Edts float64
}

func (event *WamEventIphoneMediaDownloadTransform) InitFields(option interface{}) {
	event.Amr = 0
	event.Edts = 0
}

func (event *WamEventIphoneMediaDownloadTransform) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Body().
		SerializeNumber(0x1, event.Edts)

	buffer.Footer().
		SerializeNumber(0x2, event.Amr)
}
