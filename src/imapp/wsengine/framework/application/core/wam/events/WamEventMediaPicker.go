package events

import (
	"time"
	eventSerialize "ws/framework/plugin/event_serialize"
	"ws/framework/utils"
)

type WamEventMediaPicker struct {
	WAMessageEvent

	MediaPickerOrigin           float64
	MediaType                   MediaType
	MediaPickerTime             float64
	PhotoGalleryDurationTime    float64
	AudienceSelectorClicked     float64
	AudienceSelectorUpdated     float64
	ChatRecipients              float64
	MediaPickerChanged          float64
	MediaPickerCroppedRotated   float64
	MediaPickerDeleted          float64
	MediaPickerDrawing          float64
	MediaPickerFilter           float64
	MediaPickerLikeDoc          float64
	MediaPickerNotLikeDoc       float64
	MediaPickerOriginThirdParty float64
	MediaPickerSent             float64
	MediaPickerSentUnchanged    float64
	MediaPickerStickers         float64
	MediaPickerText             float64
	StatusRecipients            float64
}

type EventMediaPickerOption struct {
	MediaType MediaType
}

func WithEventMediaPicker(mt MediaType) EventMediaPickerOption {
	return EventMediaPickerOption{
		MediaType: mt,
	}
}

func (event *WamEventMediaPicker) InitFields(option interface{}) {
	if opt, ok := option.(EventMediaPickerOption); ok {
		event.MediaType = opt.MediaType
	}

	event.MediaPickerOrigin = 1
	//event.MediaType = MediaImage
	event.MediaPickerTime = utils.LogRandMillSecond(1*time.Second, 5*time.Second)
	event.PhotoGalleryDurationTime = utils.LogRandMillSecond(time.Second, 5*time.Second)
	event.AudienceSelectorClicked = 0
	event.AudienceSelectorUpdated = 0
	event.ChatRecipients = 1
	event.MediaPickerChanged = 0
	event.MediaPickerCroppedRotated = 0
	event.MediaPickerDeleted = 0
	event.MediaPickerDrawing = 0
	event.MediaPickerFilter = 0
	event.MediaPickerLikeDoc = 0
	event.MediaPickerNotLikeDoc = 1
	event.MediaPickerOriginThirdParty = 0
	event.MediaPickerSent = 1
	event.MediaPickerSentUnchanged = 1
	event.MediaPickerStickers = 0
	event.MediaPickerText = 0
	event.StatusRecipients = 0
}

func (event *WamEventMediaPicker) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Body().
		SerializeNumber(0xe, event.MediaPickerOrigin).
		SerializeNumber(0x1, float64(event.MediaType)).
		SerializeNumber(0xf, event.MediaPickerTime).
		SerializeNumber(0x17, event.PhotoGalleryDurationTime).
		SerializeNumber(0x18, event.AudienceSelectorClicked).
		SerializeNumber(0x19, event.AudienceSelectorUpdated).
		SerializeNumber(0x10, event.ChatRecipients).
		SerializeNumber(0x4, event.MediaPickerChanged).
		SerializeNumber(0xa, event.MediaPickerCroppedRotated).
		SerializeNumber(0x3, event.MediaPickerDeleted).
		SerializeNumber(0xb, event.MediaPickerDrawing).
		SerializeNumber(0x12, event.MediaPickerFilter).
		SerializeNumber(0x13, event.MediaPickerLikeDoc).
		SerializeNumber(0x14, event.MediaPickerNotLikeDoc).
		SerializeNumber(0x15, event.MediaPickerOriginThirdParty).
		SerializeNumber(0x2, event.MediaPickerSent).
		SerializeNumber(0x5, event.MediaPickerSentUnchanged).
		SerializeNumber(0xc, event.MediaPickerStickers).
		SerializeNumber(0xd, event.MediaPickerText)

	buffer.Footer().
		SerializeNumber(0x11, event.StatusRecipients)
}
