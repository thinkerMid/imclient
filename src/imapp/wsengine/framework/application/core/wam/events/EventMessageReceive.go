package events

import (
	"time"
	eventSerialize "ws/framework/plugin/event_serialize"
	"ws/framework/utils"
)

type WamEventMessageReceive struct {
	WAMessageEvent

	MessageMediaType                    float64 //0x2
	MessageType                         float64 //0x1
	MessageReceiveT0                    float64 //0x6
	MessageReceiveT1                    float64 //0x7
	IsAReply                            float64 //0x13
	IsForwardedForward                  float64 //0x12
	IsViewOnce                          float64 //0x9
	MessageIsOffline                    float64 //0x5
	ReceiverDefaultDisappearingDuration float64 //0xc

}

type MessageReceiveOption struct {
	MediaType MediaType
}

func WithMessageReceiveOption(media MediaType) MessageReceiveOption {
	return MessageReceiveOption{
		MediaType: media,
	}
}

func (event *WamEventMessageReceive) InitFields(option interface{}) {
	// 注册时有构造该日志
	event.MessageType = 4 //?
	event.MessageReceiveT0 = utils.LogRandMillSecond(2*time.Second, 10*time.Second)
	event.MessageReceiveT1 = utils.LogRandMillSecond(0, time.Second)
	event.IsAReply = 0
	event.IsForwardedForward = 0
	event.IsViewOnce = 0
	event.MessageIsOffline = 0
	event.ReceiverDefaultDisappearingDuration = 0

	if opt, ok := option.(MessageReceiveOption); ok {
		event.MessageMediaType = float64(opt.MediaType)
	}
}

func (event *WamEventMessageReceive) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Body().
		SerializeNumber(0x2, event.MessageMediaType).
		SerializeNumber(0x1, event.MessageType).
		SerializeNumber(0x6, event.MessageReceiveT0).
		SerializeNumber(0x7, event.MessageReceiveT1).
		SerializeNumber(0x13, event.IsAReply).
		SerializeNumber(0x12, event.IsForwardedForward).
		SerializeNumber(0x9, event.IsViewOnce).
		SerializeNumber(0x5, event.MessageIsOffline)

	buffer.Footer().
		SerializeNumber(0xc, event.ReceiverDefaultDisappearingDuration)
}
