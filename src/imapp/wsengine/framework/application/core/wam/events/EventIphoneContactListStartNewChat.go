package events

import (
	eventSerialize "ws/framework/plugin/event_serialize"
)

type WamEventIphoneContactListStartNewChat struct {
	WAMessageEvent

	FrequentContacted float64
	Search            float64
	ChatType          int32
}

func (event *WamEventIphoneContactListStartNewChat) InitFields(option interface{}) {
	event.FrequentContacted = 0
	event.Search = 0
	event.ChatType = 1 //? 群组类型
}

func (event *WamEventIphoneContactListStartNewChat) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Body().
		SerializeNumber(0x1, event.FrequentContacted).
		SerializeNumber(0x2, event.Search)

	buffer.Footer().
		SerializeNumber(0x3, float64(event.ChatType))
}
