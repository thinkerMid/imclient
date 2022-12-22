package events

import (
	eventSerialize "ws/framework/plugin/event_serialize"
)

// WamEventMalformedMediaResponse .
type WamEventMalformedMediaResponse struct {
	WAMessageEvent

	RequestType   float64 // 0x1
	ResponseError float64 // 0x2
}

func (event *WamEventMalformedMediaResponse) InitFields(option interface{}) {
	event.RequestType = 2
	event.ResponseError = 2
}

func (event *WamEventMalformedMediaResponse) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Body().
		SerializeNumber(0x1, event.RequestType)

	buffer.Footer().
		SerializeNumber(0x2, event.ResponseError)
}
