package events

import (
	"time"
	eventSerialize "ws/framework/plugin/event_serialize"
	"ws/framework/utils"
)

type WamEventMediaPickerPerf struct {
	WAMessageEvent

	MediaPickerPerfOrigin           float64 //0x1
	MediaPickerPerfDismissTime      float64 //0x4
	MediaPickerPerfPresentationTime float64 //0x3
	MediaPickerPerfReadyTime        float64 //0x5
	MediaPickerPerfGifCount         float64 //0x9
	MediaPickerPerfImageCount       float64 //0x8
	MediaPickerPerfInvocationCount  float64 //0x2
	MediaPickerPerfNumAdded         float64 //0x6
	MediaPickerPerfNumRemoved       float64 //0x7
	MediaPickerPerfVideoCount       float64 //0xa
}

func (event *WamEventMediaPickerPerf) InitFields(option interface{}) {
	event.MediaPickerPerfOrigin = 1
	event.MediaPickerPerfDismissTime = utils.LogRandMillSecond(0, time.Second)
	event.MediaPickerPerfPresentationTime = utils.LogRandMillSecond(0, time.Second)
	event.MediaPickerPerfReadyTime = utils.LogRandSecond(0, 10*time.Second)
	event.MediaPickerPerfGifCount = 0
	event.MediaPickerPerfImageCount = 1
	event.MediaPickerPerfInvocationCount = 0
	event.MediaPickerPerfNumAdded = 1
	event.MediaPickerPerfNumRemoved = 0
	event.MediaPickerPerfVideoCount = 0

}

func (event *WamEventMediaPickerPerf) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Body().
		SerializeNumber(0x1, event.MediaPickerPerfOrigin).
		SerializeNumber(0x4, event.MediaPickerPerfDismissTime).
		SerializeNumber(0x3, event.MediaPickerPerfPresentationTime).
		SerializeNumber(0x5, event.MediaPickerPerfReadyTime).
		SerializeNumber(0x9, event.MediaPickerPerfGifCount).
		SerializeNumber(0x8, event.MediaPickerPerfImageCount).
		SerializeNumber(0x2, event.MediaPickerPerfInvocationCount).
		SerializeNumber(0x6, event.MediaPickerPerfNumAdded).
		SerializeNumber(0x7, event.MediaPickerPerfNumRemoved)

	buffer.Footer().
		SerializeNumber(0xa, event.MediaPickerPerfVideoCount)
}
