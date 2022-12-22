package events

import (
	"fmt"
	"time"
	eventSerialize "ws/framework/plugin/event_serialize"
)

type WamEventGroupProfilePicture struct {
	WAMessageEvent

	GroupProfileAction     float64 // 0x7
	PreciseGroupSizeBucket float64 // 0x5
	ProfilePictureType     float64 // 0x6
	GroupCreationDs        string  // 0x1
	HasProfilePicture      float64 // 0x3
	IsAdmin                float64 // 0x4
}

type GroupProfilePictureOption struct {
	Action    float64
	HasAvatar int8
	IsAdmin   int8
}

func WithGroupProfilePictureOption(action float64, hasAvatar, isAdmin int8) GroupProfilePictureOption {
	return GroupProfilePictureOption{
		Action:    action,
		HasAvatar: hasAvatar,
		IsAdmin:   isAdmin,
	}
}

func (event *WamEventGroupProfilePicture) InitFields(option interface{}) {

	if opt, ok := option.(GroupProfilePictureOption); ok {
		event.GroupProfileAction = opt.Action

		if opt.HasAvatar == 1 {
			event.HasProfilePicture = 1
		}

		if opt.IsAdmin == 1 {
			event.IsAdmin = 1
		}
	}

	y, m, d := time.Now().Date()

	event.PreciseGroupSizeBucket = 1
	event.GroupCreationDs = fmt.Sprintf("%v-%v-%v", y, m, d)
	if event.GroupProfileAction == 10 {
		event.ProfilePictureType = 2
	}
}

func (event *WamEventGroupProfilePicture) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Body().
		SerializeNumber(0x7, event.GroupProfileAction).
		SerializeNumber(0x5, event.PreciseGroupSizeBucket).
		SerializeNumber(0x6, event.ProfilePictureType).
		SerializeString(0x1, event.GroupCreationDs).
		SerializeNumber(0x3, event.HasProfilePicture)

	buffer.Footer().
		SerializeNumber(0x4, event.IsAdmin)
}
