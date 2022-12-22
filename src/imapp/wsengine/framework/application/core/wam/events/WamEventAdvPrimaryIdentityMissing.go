package events

import eventSerialize "ws/framework/plugin/event_serialize"

type WamEventAdvPrimaryIdentityMissing struct {
	WAMessageEvent

	PrimaryIdentityMissingProtoType float64
}

func (event *WamEventAdvPrimaryIdentityMissing) InitFields(option interface{}) {
	event.PrimaryIdentityMissingProtoType = 2
}

func (event *WamEventAdvPrimaryIdentityMissing) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Footer().
		SerializeNumber(0x1, event.PrimaryIdentityMissingProtoType)
}
