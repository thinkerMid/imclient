package events

import eventSerialize "ws/framework/plugin/event_serialize"

type WamEventMp4Repair struct {
	WAMessageEvent

	OldOriginator     float64 //0x4
	NewOriginator     float64 //0x8
	NewMajorVersion   float64 //0x5
	NewMinorVersion   float64 //0x6
	NewReleaseVersion float64 //0x7
	RepairRequired    float64 //0x9
	RepairSuccessful  float64 //0xa

}

func (event *WamEventMp4Repair) InitFields(option interface{}) {
	event.OldOriginator = 0
	event.NewOriginator = 5
	event.NewMajorVersion = 1
	event.NewMinorVersion = 1
	event.NewReleaseVersion = 0
	event.RepairRequired = 1
	event.RepairSuccessful = 1
}

func (event *WamEventMp4Repair) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Body().
		SerializeNumber(0x4, event.OldOriginator).
		SerializeNumber(0x8, event.NewOriginator).
		SerializeNumber(0x5, event.NewMajorVersion).
		SerializeNumber(0x6, event.NewMinorVersion).
		SerializeNumber(0x7, event.NewReleaseVersion).
		SerializeNumber(0x9, event.RepairRequired)

	buffer.Footer().
		SerializeNumber(0xa, event.RepairSuccessful)
}
