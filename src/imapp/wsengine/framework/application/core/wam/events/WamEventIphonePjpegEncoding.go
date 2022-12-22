package events

import eventSerialize "ws/framework/plugin/event_serialize"

type WamEventIphonePjpegEncoding struct {
	WAMessageEvent

	IphonePjpegEncodingResultType float64
}

func (event *WamEventIphonePjpegEncoding) InitFields(option interface{}) {
	event.IphonePjpegEncodingResultType = 0
}

func (event *WamEventIphonePjpegEncoding) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Footer().
		SerializeNumber(0x1, event.IphonePjpegEncodingResultType)
}
