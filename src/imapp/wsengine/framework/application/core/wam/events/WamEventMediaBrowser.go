package events

import (
	"time"
	eventSerialize "ws/framework/plugin/event_serialize"
	"ws/framework/utils"
)

type WamEventMediaBrowser struct {
	WAMessageEvent

	MediaBrowserPresentationTime float64
}

func (event *WamEventMediaBrowser) InitFields(option interface{}) {
	event.MediaBrowserPresentationTime = utils.LogRandMillSecond(0, time.Second)
}

func (event *WamEventMediaBrowser) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Footer().
		SerializeNumber(0x1, event.MediaBrowserPresentationTime)
}
