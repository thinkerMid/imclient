package events

import (
	"time"
	eventSerialize "ws/framework/plugin/event_serialize"
	"ws/framework/utils"
)

type EventMessageSend struct {
	WAMessageEvent

	DeviceCount                         float64     //0x1F
	DeviceSizeBucket                    float64     //0x19
	E2eBackFill                         float64     //0x17
	IsReply                             float64     //0x23
	IsViewOnce                          float64     //0x16
	MediaCaptionPresent                 float64     //0x8
	MessageIsFirstUserMessage           float64     //0x1a  1:会话第一条消息 	0:非第一条
	MessageIsForward                    float64     //0x4
	MessageIsInvisible                  float64     //0x1d
	MessageIsRevoke                     float64     //0x18
	MessageMediaType                    MediaType   //0x3
	MessageSendResult                   float64     //0x1
	MessageSendResultIsTerminal         float64     //0x11
	MessageSendTime                     float64     //0xb
	MessageType                         MessageType //0x2
	ParticipantCount                    float64     //0x20
	ReceiverDefaultDisappearingDuration float64     //0x1c
	ResendCount                         float64     //0x10
	RetryCount                          float64     //0x6
	SenderDefaultDisappearingDuration   float64     //0x1b
	ThumbSize                           float64     //0x14
}

type MessageSendOption struct {
	FirstMessage     bool
	MessageMediaType MediaType
	MessageType      MessageType
	ThumbSize        int32
}

func WithMessageSendOption(isFirst bool, message MessageType, media MediaType, thumbSize int32) MessageSendOption {
	return MessageSendOption{
		FirstMessage:     isFirst,
		MessageType:      message,
		MessageMediaType: media,
		ThumbSize:        thumbSize,
	}
}

// InitFields .
// @media, @first
func (event *EventMessageSend) InitFields(option interface{}) {
	event.DeviceCount = 32
	event.DeviceSizeBucket = 1
	event.E2eBackFill = 0
	event.IsReply = 0
	event.IsViewOnce = 0
	event.MediaCaptionPresent = 0
	event.MessageIsForward = 0
	event.MessageIsInvisible = 0
	event.MessageIsRevoke = 0
	event.MessageSendResult = 1
	event.MessageSendResultIsTerminal = 0
	event.MessageSendTime = utils.LogRandMillSecond(0, 2*time.Second)
	event.ParticipantCount = 32
	event.ReceiverDefaultDisappearingDuration = 0
	event.ResendCount = 0
	event.RetryCount = 0
	event.SenderDefaultDisappearingDuration = 0

	if opt, ok := option.(MessageSendOption); ok {
		if opt.FirstMessage {
			event.MessageIsFirstUserMessage = 1
		}
		event.MessageType = opt.MessageType
		event.MessageMediaType = opt.MessageMediaType
		event.ThumbSize = float64(opt.ThumbSize)
	}
}

func (event *EventMessageSend) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Body().
		SerializeNumber(0x1F, event.DeviceCount).
		SerializeNumber(0x19, event.DeviceSizeBucket).
		SerializeNumber(0x17, event.E2eBackFill).
		//SerializeNumber(0xf, event.FastForwardEnabled).
		SerializeNumber(0x23, event.IsReply).
		SerializeNumber(0x16, event.IsViewOnce).
		SerializeNumber(0x8, event.MediaCaptionPresent).
		SerializeNumber(0x1a, event.MessageIsFirstUserMessage).
		SerializeNumber(0x4, event.MessageIsForward).
		SerializeNumber(0x1d, event.MessageIsInvisible).
		SerializeNumber(0x18, event.MessageIsRevoke).
		SerializeNumber(0x3, float64(event.MessageMediaType)).
		SerializeNumber(0x1, event.MessageSendResult).
		SerializeNumber(0x11, event.MessageSendResultIsTerminal).
		SerializeNumber(0xb, event.MessageSendTime).
		SerializeNumber(0x2, float64(event.MessageType)).
		SerializeNumber(0x20, event.ParticipantCount).
		SerializeNumber(0x1c, event.ReceiverDefaultDisappearingDuration).
		SerializeNumber(0x10, event.ResendCount).
		SerializeNumber(0x6, event.RetryCount).
		SerializeNumber(0x1b, event.SenderDefaultDisappearingDuration).
		SerializeNumber(0x14, event.ThumbSize)
}
