package events

import (
	"encoding/hex"
	"github.com/chenzhuoyu/base64x"
	"strings"
	"time"
	eventSerialize "ws/framework/plugin/event_serialize"
	"ws/framework/utils"
)

type WamEventRegistrationComplete struct {
	WAMessageEvent

	RegistrationTime float64
	DeviceIdentifier string
}

type RegistrationCompleteOption struct {
	Identify string
}

func (event *WamEventRegistrationComplete) InitFields(option interface{}) {
	if opt, ok := option.(RegistrationCompleteOption); ok {
		opt.Identify = strings.ReplaceAll(opt.Identify, "-", "")
		buff, err := hex.DecodeString(opt.Identify)
		if err == nil {
			event.DeviceIdentifier = base64x.StdEncoding.EncodeToString(buff)
		}
	}

	event.RegistrationTime = utils.LogRandMillSecond(2*60*time.Second, 10*60*time.Second)
}

func (event *WamEventRegistrationComplete) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Body().
		SerializeString(0x9, event.DeviceIdentifier)

	buffer.Footer().
		SerializeNumber(0x1, event.RegistrationTime)
}
