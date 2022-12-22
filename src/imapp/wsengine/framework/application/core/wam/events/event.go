package events

import (
	eventSerialize "ws/framework/plugin/event_serialize"
)

type WAMessageEvent struct {
	Channel uint8
	Code    int64
	Weight  float64
}

func (wae *WAMessageEvent) InitFields(option interface{}) {
	//TODO implement me
	panic("implement me")
}

func (wae *WAMessageEvent) Serialize(buffer eventSerialize.IEventBuffer) {
	//TODO implement me
	panic("implement me")
}

func (wae *WAMessageEvent) Init(channel uint8, code int64, weight float64) {
	wae.Channel = channel
	wae.Code = code
	wae.Weight = weight
}
