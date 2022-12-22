package events

import (
	eventSerialize "ws/framework/plugin/event_serialize"
)

type EventE2eMessageSend struct {
	WAMessageEvent

	E2eCiphertextType    float64
	E2eDestination       float64
	E2eReceiverType      float64
	E2eCiphertextVersion float64
	E2eSuccessful        float64
	MessageMediaType     MediaType
	RetryCount           float64
}

type E2eMessageSendOption struct {
	Media MediaType
}

func WithE2eMessageSendOption(media MediaType) E2eMessageSendOption {
	return E2eMessageSendOption{
		Media: media,
	}
}

func (event *EventE2eMessageSend) InitFields(option interface{}) {
	event.E2eCiphertextType = 1
	event.E2eDestination = 0
	event.E2eReceiverType = 1 // ?
	event.E2eCiphertextVersion = 2
	event.E2eSuccessful = 1
	event.MessageMediaType = MediaText
	event.RetryCount = 0

	if opt, ok := option.(E2eMessageSendOption); ok {
		event.MessageMediaType = opt.Media
	}
}

func (event *EventE2eMessageSend) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Body().
		SerializeNumber(0x5, event.E2eCiphertextType).
		SerializeNumber(0x6, event.E2eCiphertextVersion).
		SerializeNumber(0x4, event.E2eDestination).
		SerializeNumber(0x8, event.E2eReceiverType).
		SerializeNumber(0x1, event.E2eSuccessful).
		SerializeNumber(0x7, float64(event.MessageMediaType))

	buffer.Footer().
		SerializeNumber(0x3, event.RetryCount)
}
