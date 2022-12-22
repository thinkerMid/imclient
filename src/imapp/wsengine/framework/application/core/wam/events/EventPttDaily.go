package events

import eventSerialize "ws/framework/plugin/event_serialize"

type WamEventPttDaily struct {
	WAMessageEvent

	CancelBroadcast       float64
	CancelGroup           float64
	CancelIndividual      float64
	DraftReviewBroadcast  float64
	DraftReviewGroup      float64
	DraftReviewIndividual float64

	FastPlaybackBroadcast  float64
	FastPlaybackGroup      float64
	FastPlaybackIndividual float64
	LockBroadcast          float64
	LockGroup              float64
	LockIndividual         float64
	PlaybackBroadcast      float64
	PlaybackGroup          float64
	PlaybackIndividual     float64
	RecordBroadcast        float64
	RecordGroup            float64
	RecordIndividual       float64
	SendBroadcast          float64
	SendGroup              float64
	SendIndividual         float64
}

func (event *WamEventPttDaily) InitFields(option interface{}) {
	event.CancelBroadcast = 0
	event.CancelGroup = 0
	event.CancelIndividual = 0
	event.DraftReviewBroadcast = 0
	event.DraftReviewGroup = 0
	event.DraftReviewIndividual = 0
	event.FastPlaybackBroadcast = 0
	event.FastPlaybackGroup = 0
	event.FastPlaybackIndividual = 0
	event.LockBroadcast = 0
	event.LockGroup = 0
	event.LockIndividual = 0
	event.PlaybackBroadcast = 0
	event.PlaybackGroup = 0
	event.PlaybackIndividual = 0
	event.RecordBroadcast = 0
	event.RecordGroup = 0
	event.RecordIndividual = 0
	event.SendBroadcast = 0
	event.SendGroup = 0
	event.SendIndividual = 0
}

func (event *WamEventPttDaily) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	var idx int64
	for idx = 0x1; idx != 0x16; idx++ {
		if idx != 0x15 {
			buffer.Body().
				SerializeNumber(idx, 0.0)
		} else {
			buffer.Footer().
				SerializeNumber(idx, 0.0)
		}
	}
}
