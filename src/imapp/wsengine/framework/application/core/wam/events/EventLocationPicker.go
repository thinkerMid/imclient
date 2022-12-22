package events

import (
	"math/rand"
	"time"
	eventSerialize "ws/framework/plugin/event_serialize"
	"ws/framework/utils"
)

// WamEventLocationPicker 发送位置信息
type WamEventLocationPicker struct {
	WAMessageEvent

	LocationPickerPlacesResponse     float64 //0x4
	LocationPickerPlacesSource       float64 //0x1
	LocationPickerResultType         float64 //0x3
	LocationPickerSpendTime          float64 //0xC
	LocationPickerFullScreen         float64 //0x7
	LocationPickerPlacesCount        float64 //0xA
	LocationPickerSelectedPlaceIndex float64 //0xB
}

func (event *WamEventLocationPicker) InitFields(option interface{}) {
	event.LocationPickerPlacesResponse = 1
	event.LocationPickerPlacesSource = 2
	event.LocationPickerResultType = 4
	event.LocationPickerSpendTime = utils.LogRandMillSecond(5*time.Second, 10*time.Second)
	event.LocationPickerFullScreen = 0
	event.LocationPickerPlacesCount = float64(int(utils.LogRandSecond(10*time.Second, 80*time.Second)))
	event.LocationPickerSelectedPlaceIndex = float64(rand.Intn(int(event.LocationPickerPlacesCount - 1)))
}

func (event *WamEventLocationPicker) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Body().
		SerializeNumber(0x4, event.LocationPickerPlacesResponse).
		SerializeNumber(0x1, event.LocationPickerPlacesSource).
		SerializeNumber(0x3, event.LocationPickerResultType).
		SerializeNumber(0xC, event.LocationPickerSpendTime).
		SerializeNumber(0x7, event.LocationPickerFullScreen).
		SerializeNumber(0xA, event.LocationPickerPlacesCount)

	buffer.Footer().
		SerializeNumber(0xB, event.LocationPickerSelectedPlaceIndex)
}
