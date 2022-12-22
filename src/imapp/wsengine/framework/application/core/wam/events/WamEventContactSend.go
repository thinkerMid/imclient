package events

import (
	"time"
	eventSerialize "ws/framework/plugin/event_serialize"
	"ws/framework/utils"
)

type WamEventContactSend struct {
	WAMessageEvent

	Channel           float64 // 0x1
	IsMultiCard       float64 // 0x2
	MessageSendResult float64 // 0x3
	MessageSendTime   float64 // 0x4
	VCardSize         float64 // 0x5
}

type EventContactSendOption struct {
	MultiCard bool
	CardSize  int32
}

func WithEventContactSend(multi bool, size int32) EventContactSendOption {
	return EventContactSendOption{
		MultiCard: multi,
		CardSize:  size,
	}
}

func (event *WamEventContactSend) InitFields(option interface{}) {
	if opt, ok := option.(EventContactSendOption); ok {
		if opt.MultiCard {
			event.IsMultiCard = 1
			event.VCardSize = 1024
		}
	}
	event.Channel = 1
	event.MessageSendResult = 1
	event.MessageSendTime = utils.LogRandMillSecond(100*time.Millisecond, 1*time.Second)
}

func (event *WamEventContactSend) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Body().
		SerializeNumber(0x1, event.Channel).
		SerializeNumber(0x2, event.IsMultiCard).
		SerializeNumber(0x3, event.MessageSendResult).
		SerializeNumber(0x4, event.MessageSendTime)

	buffer.Footer().
		SerializeNumber(0x5, event.VCardSize)
}
