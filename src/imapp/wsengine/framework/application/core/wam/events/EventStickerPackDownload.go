package events

import eventSerialize "ws/framework/plugin/event_serialize"

type WamEventStickerPackDownload struct {
	WAMessageEvent

	StickerPackDownloadOrigin float64
	StickerPackIsFirstParty   float64
}

func (event *WamEventStickerPackDownload) InitFields(option interface{}) {
	event.StickerPackDownloadOrigin = 4
	event.StickerPackIsFirstParty = 1.000000
}

func (event *WamEventStickerPackDownload) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Body().
		SerializeNumber(0x1, event.StickerPackDownloadOrigin)

	buffer.Footer().
		SerializeNumber(0x2, event.StickerPackIsFirstParty)
}
