package eventService

import (
	"time"
	"ws/framework/application/container/abstract_interface"
)

var _ containerInterface.IEventCache = &Channel2EventCache{}

// Channel2EventCache .
type Channel2EventCache struct {
	Channel0EventCache
}

// Init .
func (e *Channel2EventCache) Init() {
	e.serialNumber = time.Now().UnixNano()
	e.channelID = channel2
	e.dirtyCode = standbyWrite
}
