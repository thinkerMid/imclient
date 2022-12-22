package events

import (
	eventSerialize "ws/framework/plugin/event_serialize"
)

type WamEventAppLaunch struct {
	WAMessageEvent

	LaunchType        float64
	LaunchMainPreTime float64
	LaunchMainRunTime float64
	LaunchTime        float64
}

func (event *WamEventAppLaunch) InitFields(option interface{}) {
	event.LaunchType = 1
	event.LaunchMainPreTime = 1949.3319988250732
	event.LaunchMainRunTime = 685.8890056610107
	event.LaunchTime = 2635.221004486084
}

func (event *WamEventAppLaunch) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Body().
		SerializeNumber(0x1, event.LaunchTime).
		SerializeNumber(0x3, event.LaunchMainPreTime).
		SerializeNumber(0x4, event.LaunchMainRunTime)

	buffer.Footer().
		SerializeNumber(0x5, event.LaunchType)
}
