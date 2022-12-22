package events

import (
	eventSerialize "ws/framework/plugin/event_serialize"
)

type WamEventPrivacyHighlightDaily struct {
	WAMessageEvent

	PrivacyHighlightCategory float64
	PrivacyHighlightSurface  float64
	DialogAppearCount        float64
	DialogSelectCount        float64
	NarrativeAppearCount     float64
}

func (event *WamEventPrivacyHighlightDaily) InitFields(option interface{}) {
	event.PrivacyHighlightCategory = 0
	event.PrivacyHighlightSurface = 8
	event.DialogAppearCount = 0
	event.DialogSelectCount = 0
	event.NarrativeAppearCount = 1
}

func (event *WamEventPrivacyHighlightDaily) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Body().
		SerializeNumber(0x1, event.DialogAppearCount)
	buffer.Footer().
		SerializeNumber(0x2, event.DialogSelectCount)
	buffer.Footer().
		SerializeNumber(0x3, event.NarrativeAppearCount)
	buffer.Footer().
		SerializeNumber(0x4, event.PrivacyHighlightCategory)
	buffer.Footer().
		SerializeNumber(0x5, event.PrivacyHighlightSurface)
}
