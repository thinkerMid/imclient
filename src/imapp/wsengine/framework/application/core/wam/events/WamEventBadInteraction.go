package events

import (
	"time"
	eventSerialize "ws/framework/plugin/event_serialize"
	"ws/framework/utils"
)

type InteractionType string

const (
	InteractionLaunch   InteractionType = "app_launch"
	InteractionMessage  InteractionType = "message_send"
	InteractionChatOpen InteractionType = "chat_open" // 每次进入群聊
)

type WamEventBadInteraction struct {
	WAMessageEvent

	Actual    float64 //0x1
	Name      string  //0x2
	Threshold float64 //0x3
}

type BadInteractionOption struct {
	Name InteractionType
}

func WithBadInteractionOption(name InteractionType) BadInteractionOption {
	return BadInteractionOption{
		Name: name,
	}
}

func (event *WamEventBadInteraction) InitFields(option interface{}) {
	event.Actual = utils.LogRandMillSecond(500*time.Millisecond, 20*time.Second)
	event.Threshold = 0

	if opt, ok := option.(BadInteractionOption); ok {
		event.Name = string(opt.Name)
	}
}

func (event *WamEventBadInteraction) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Body().
		SerializeNumber(0x1, event.Actual).
		SerializeString(0x2, event.Name)

	buffer.Footer().
		SerializeNumber(0x3, event.Threshold)
}
