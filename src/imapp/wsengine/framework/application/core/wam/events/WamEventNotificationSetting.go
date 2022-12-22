package events

import eventSerialize "ws/framework/plugin/event_serialize"

type WamEventNotificationSetting struct {
	WAMessageEvent

	GroupSoundTone              float64 //0x3
	InAppNotificationAlertStyle float64 //0x4
	MessageSoundTone            float64 //0x9
	GroupReactionNotification   float64 //0x1
	GroupShowNotification       float64 //0x2
	InAppNotificationSound      float64 //0x5
	InAppNotificationVibrate    float64 //0x6
	MessageReactionNotification float64 //0x7
	MessageShowNotification     float64 //0x8
	ShowPreview                 float64 //0xa
}

func (event *WamEventNotificationSetting) InitFields(option interface{}) {
	event.GroupSoundTone = 1
	event.InAppNotificationAlertStyle = 2
	event.MessageSoundTone = 1
	event.GroupReactionNotification = 1
	event.GroupShowNotification = 1
	event.InAppNotificationSound = 1
	event.InAppNotificationVibrate = 1
	event.MessageReactionNotification = 1
	event.MessageShowNotification = 1
	event.ShowPreview = 1
}

func (event *WamEventNotificationSetting) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Body().
		SerializeNumber(0x3, event.GroupSoundTone).
		SerializeNumber(0x4, event.InAppNotificationAlertStyle).
		SerializeNumber(0x9, event.MessageSoundTone).
		SerializeNumber(0x1, event.GroupReactionNotification).
		SerializeNumber(0x2, event.GroupShowNotification).
		SerializeNumber(0x5, event.InAppNotificationSound).
		SerializeNumber(0x6, event.InAppNotificationVibrate).
		SerializeNumber(0x7, event.MessageReactionNotification).
		SerializeNumber(0x8, event.MessageShowNotification)

	buffer.Footer().
		SerializeNumber(0xa, event.ShowPreview)
}
