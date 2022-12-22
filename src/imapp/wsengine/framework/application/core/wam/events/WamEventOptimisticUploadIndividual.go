package events

import eventSerialize "ws/framework/plugin/event_serialize"

type WamEventOptimisticUploadIndividual struct {
	WAMessageEvent

	OptEndState    float64 //0x1
	OptOrigin      float64 //0x3
	OptUploadDelay float64 //0x4
	OptSizeDiff    float64 //0x7

}

func (event *WamEventOptimisticUploadIndividual) InitFields(option interface{}) {
	event.OptEndState = 7
	event.OptOrigin = 1
	event.OptUploadDelay = 2000
	event.OptSizeDiff = 0
}

func (event *WamEventOptimisticUploadIndividual) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Body().
		SerializeNumber(0x1, event.OptEndState).
		SerializeNumber(0x3, event.OptOrigin).
		SerializeNumber(0x4, event.OptUploadDelay)

	buffer.Footer().
		SerializeNumber(0x7, event.OptSizeDiff)
}
