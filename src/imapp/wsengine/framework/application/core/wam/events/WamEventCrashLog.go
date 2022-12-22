package events

import eventSerialize "ws/framework/plugin/event_serialize"

/*
* WhatsApp杀进程后重启，会产生该日志
 */

type WamEventCrashLog struct {
	WAMessageEvent

	CrashType                              float64
	IphoneAppState                         float64
	IphoneLastMemoryPressureLevel          float64
	IphoneTimeSinceLastMemoryPressureEvent float64
	CrashContext                           string
	CrashCount                             float64
	CrashReason                            string
	IphoneMemoryPressureEventCount         float64
	//IphonePossibleFalsePositive            float64
}

func (event *WamEventCrashLog) InitFields(option interface{}) {
	event.CrashType = 11
	event.IphoneAppState = 1
	event.IphoneLastMemoryPressureLevel = 0
	event.IphoneTimeSinceLastMemoryPressureEvent = 0
	event.CrashContext = ""
	event.CrashCount = 1
	event.CrashReason = "FOOM"
	event.IphoneMemoryPressureEventCount = 0
}

func (event *WamEventCrashLog) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Body().
		SerializeNumber(0x6, event.CrashType).
		SerializeNumber(0x4, event.IphoneAppState).
		SerializeNumber(0xc, event.IphoneLastMemoryPressureLevel).
		SerializeNumber(0xa, event.IphoneTimeSinceLastMemoryPressureEvent).
		SerializeString(0x3, event.CrashContext).
		SerializeNumber(0x5, event.CrashCount).
		SerializeString(0x2, event.CrashReason)

	buffer.Footer().
		SerializeNumber(0xb, event.IphoneMemoryPressureEventCount)
}
