package events

import (
	"time"
	eventSerialize "ws/framework/plugin/event_serialize"
	"ws/framework/utils"
)

const (
	TMediaType_Text = 1
	TMediaType_Pic
	TMediaType_Video
	TMediaType_Voice
)

// MessageStage 125763 顺序发送
type MessageStage int

const (
	MsgStage_1 MessageStage = 1
	MsgStage_2
	MsgStage_3
	MsgStage_4
	MsgStage_5
	MsgStage_6
	MsgStage_7
)

type EventAndroidMessageSendPerf struct {
	WAMessageEvent

	DeviceSizeBucket                    float64     //0x1A
	DeviceCount                         float64     //0x24
	DurationAbs                         float64     //0xb
	DurationRelative                    float64     //0xc
	GroupSizeBucket                     float64     //0x11
	IsDirectedMessage                   float64     //0x21
	IsMessageFanout                     float64     //0x9
	IsMessageForward                    float64     //0x8
	IsRevokeMessage                     float64     //0x18
	MediaType                           MediaType   //0x3
	MessageIsFirstUserMessage           float64     //0x1e
	MessageIsInvisible                  float64     //0x1f
	MessageType                         MessageType //0x4
	ParticipantCount                    float64     //0x25
	PrekeysEligibleForPrallelProcessing float64     //0x1c
	SendCount                           float64     //0xd
	SendRetryCount                      float64     //0xa
	SendStage                           float64     //0x2
	TargetDeviceGroupSizeBucket         float64     //0x14
}

type MessageSendPerfOption struct {
	FirstMsg bool
	Stage    MessageStage
	Media    MediaType
	Message  MessageType
}

func WithMessageSendPerfOption(first bool, message MessageType, media MediaType, stage MessageStage) MessageSendPerfOption {
	return MessageSendPerfOption{
		FirstMsg: first,
		Stage:    stage,
		Media:    media,
		Message:  message,
	}
}

// InitFields .
// @media, @first, @stage
func (event *EventAndroidMessageSendPerf) InitFields(option interface{}) {
	event.DeviceSizeBucket = DefaultVal  // 1
	event.DeviceCount = DefaultVal       // 32
	event.IsDirectedMessage = DefaultVal // 0
	event.DurationAbs = float64(int(utils.LogRandSecond(10*time.Second, 300*time.Second)))
	event.DurationRelative = float64(int(utils.LogRandSecond(300*time.Second, 1000*time.Second)))
	event.GroupSizeBucket = 1
	event.IsMessageFanout = 1
	event.IsMessageForward = 0
	event.IsRevokeMessage = 0
	event.MessageIsFirstUserMessage = 0
	event.MessageIsInvisible = 0
	event.ParticipantCount = 32
	event.PrekeysEligibleForPrallelProcessing = 1
	event.SendCount = 0
	event.SendRetryCount = 0
	event.SendStage = 0
	event.TargetDeviceGroupSizeBucket = DefaultVal // 1

	if opt, ok := option.(MessageSendPerfOption); ok {
		event.MediaType = opt.Media
		event.MessageType = opt.Message

		if opt.FirstMsg {
			event.MessageIsFirstUserMessage = 1
		}
	}
}

func (event *EventAndroidMessageSendPerf) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Body().
		SerializeNumber(0x1A, event.DeviceSizeBucket).
		SerializeNumber(0x24, event.DeviceCount).
		SerializeNumber(0xb, event.DurationAbs).
		SerializeNumber(0xc, event.DurationRelative).
		SerializeNumber(0x11, event.GroupSizeBucket).
		SerializeNumber(0x21, event.IsDirectedMessage).
		SerializeNumber(0x9, event.IsMessageFanout).
		SerializeNumber(0x8, event.IsMessageForward).
		SerializeNumber(0x18, event.IsRevokeMessage).
		SerializeNumber(0x3, float64(event.MediaType)).
		SerializeNumber(0x1e, event.MessageIsFirstUserMessage).
		SerializeNumber(0x1f, event.MessageIsInvisible).
		SerializeNumber(0x4, float64(event.MessageType)).
		SerializeNumber(0x25, event.ParticipantCount).
		SerializeNumber(0x1c, event.PrekeysEligibleForPrallelProcessing).
		SerializeNumber(0xd, event.SendCount).
		SerializeNumber(0xa, event.SendRetryCount).
		SerializeNumber(0x2, event.SendStage)

	buffer.Footer().
		SerializeNumber(0x14, event.TargetDeviceGroupSizeBucket)
}
