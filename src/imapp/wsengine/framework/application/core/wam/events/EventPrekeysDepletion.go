package events

import eventSerialize "ws/framework/plugin/event_serialize"

/*
 * 应该被弃用了
 */

// WamEventPrekeysDepletion 私聊只构造一个 群聊，邀请多少人构造多少次
type WamEventPrekeysDepletion struct {
	WAMessageEvent

	DeviceBuckets float64
	MessageType   float64 // 1:私聊 2:群聊
	FetchReason   float64
}

func (event *WamEventPrekeysDepletion) InitFields(option interface{}) {
	event.DeviceBuckets = 1
	event.MessageType = 2
	event.FetchReason = 1
}

func (event *WamEventPrekeysDepletion) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Body().
		SerializeNumber(0x3, event.DeviceBuckets).
		SerializeNumber(0x2, event.MessageType)

	buffer.Footer().
		SerializeNumber(0x1, event.FetchReason)
}
