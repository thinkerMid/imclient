package events

import eventSerialize "ws/framework/plugin/event_serialize"

type WamEventCommunityTabAction struct {
	WAMessageEvent

	CommunityNoActionTabViews    float64
	CommunityTabGroupNavigations float64
	CommunityTabToHomeViews      float64
	CommunityTabViews            float64
}

func (event *WamEventCommunityTabAction) InitFields(option interface{}) {
	event.CommunityNoActionTabViews = 0
	event.CommunityTabGroupNavigations = 0
	event.CommunityTabToHomeViews = 0
	event.CommunityTabViews = 0
}

func (event *WamEventCommunityTabAction) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Body().
		SerializeNumber(0x4, event.CommunityNoActionTabViews).
		SerializeNumber(0x1, event.CommunityTabGroupNavigations).
		SerializeNumber(0x2, event.CommunityTabToHomeViews)

	buffer.Footer().
		SerializeNumber(0x3, event.CommunityTabViews)
}
