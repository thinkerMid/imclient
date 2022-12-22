package events

import eventSerialize "ws/framework/plugin/event_serialize"

type WamEventIphoneVideoCaching struct {
	WAMessageEvent

	VideoCachingResult float64
	VideoSize          float64
}

func (event *WamEventIphoneVideoCaching) InitFields(option interface{}) {
	event.VideoCachingResult = 1
	event.VideoSize = 0
}

func (event *WamEventIphoneVideoCaching) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Body().
		SerializeNumber(0x1, event.VideoCachingResult)

	buffer.Footer().
		SerializeNumber(0x2, event.VideoSize)
}
