package events

import (
	eventSerialize "ws/framework/plugin/event_serialize"
)

type WamEventArchiveStateDaily struct {
	WAMessageEvent

	SetKeepChatArchived           float64
	TotalGroupArchived            float64
	TotalIndividualArchived       float64
	TotalUnreadGroupArchived      float64
	TotalUnreadIndividualArchived float64
}

func (event *WamEventArchiveStateDaily) InitFields(option interface{}) {
	event.SetKeepChatArchived = 1
	event.TotalGroupArchived = 0.000000
	event.TotalIndividualArchived = 0.000000
	event.TotalUnreadGroupArchived = 0.000000
	event.TotalUnreadIndividualArchived = 0.000000
}

func (event *WamEventArchiveStateDaily) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Body().
		SerializeNumber(0x5, event.SetKeepChatArchived).
		SerializeNumber(0x2, event.TotalGroupArchived).
		SerializeNumber(0x1, event.TotalIndividualArchived).
		SerializeNumber(0x4, event.TotalUnreadGroupArchived)

	buffer.Footer().
		SerializeNumber(0x3, event.TotalUnreadIndividualArchived)
}
