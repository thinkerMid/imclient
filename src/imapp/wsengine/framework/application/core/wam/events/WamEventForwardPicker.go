package events

import (
	"time"
	eventSerialize "ws/framework/plugin/event_serialize"
	"ws/framework/utils"
)

type WamEventForwardPicker struct {
	WAMessageEvent

	ForwardPickerResult                float64 //0x1
	ForwardPickerSpendTime             float64 //0xB
	ForwardPickerContactsSelected      float64 //0x3
	ForwardPickerFrequentsDisplayed    float64 //0x6
	ForwardPickerFrequentsLimit        float64 //0x5
	ForwardPickerFrequentsNumberOfDays float64 //0x4
	ForwardPickerFrequentsSelected     float64 //0x7
	ForwardPickerMulticastEnabled      float64 //0x2
	ForwardPickerRecentsSelected       float64 //0x8
	ForwardPickerSearchResultsSelected float64 //0x9
	ForwardPickerSearchUsed            float64 //0xA

}

func (event *WamEventForwardPicker) InitFields(option interface{}) {
	event.ForwardPickerResult = 2
	event.ForwardPickerSpendTime = utils.LogRandMillSecond(2*time.Second, 5*time.Second)
	event.ForwardPickerContactsSelected = 1
	event.ForwardPickerFrequentsDisplayed = 1
	event.ForwardPickerFrequentsLimit = 3
	event.ForwardPickerFrequentsNumberOfDays = 8
	event.ForwardPickerFrequentsSelected = 1
	event.ForwardPickerMulticastEnabled = 1
	event.ForwardPickerRecentsSelected = 0
	event.ForwardPickerSearchResultsSelected = 0
	event.ForwardPickerSearchUsed = 0
}

func (event *WamEventForwardPicker) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Body().
		SerializeNumber(0x1, event.ForwardPickerResult).
		SerializeNumber(0xB, event.ForwardPickerSpendTime).
		SerializeNumber(0x3, event.ForwardPickerContactsSelected).
		SerializeNumber(0x6, event.ForwardPickerFrequentsDisplayed).
		SerializeNumber(0x5, event.ForwardPickerFrequentsLimit).
		SerializeNumber(0x4, event.ForwardPickerFrequentsNumberOfDays).
		SerializeNumber(0x7, event.ForwardPickerFrequentsSelected).
		SerializeNumber(0x2, event.ForwardPickerMulticastEnabled).
		SerializeNumber(0x8, event.ForwardPickerRecentsSelected).
		SerializeNumber(0x9, event.ForwardPickerSearchResultsSelected)

	buffer.Footer().
		SerializeNumber(0xA, event.ForwardPickerSearchUsed)
}
