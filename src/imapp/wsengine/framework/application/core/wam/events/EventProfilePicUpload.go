package events

import (
	"time"
	eventSerialize "ws/framework/plugin/event_serialize"
	"ws/framework/utils"
)

// WamEventProfilePicUpload .
type WamEventProfilePicUpload struct {
	WAMessageEvent

	PicSize         float64
	PicAvatar       float64
	PicTotalTime    float64
	PicUploadTime   float64
	PicUploadType   float64 //?
	PicUploadResult float64
}

type ProfilePicUploadOption struct {
	pictureSize int32
	isAvatar    int32
}

func WithProfilePicUploadOption(isAvatar, pictureSize uint32) ProfilePicUploadOption {
	return ProfilePicUploadOption{
		pictureSize: int32(pictureSize),
		isAvatar:    int32(isAvatar),
	}
}

func (event *WamEventProfilePicUpload) InitFields(option interface{}) {
	if opt, ok := option.(ProfilePicUploadOption); ok {
		event.PicAvatar = float64(opt.isAvatar)
		event.PicSize = float64(opt.pictureSize)
	}

	costTime := utils.LogRandMillSecond(time.Second, 10*time.Second)

	event.PicTotalTime = costTime
	event.PicUploadTime = costTime
	event.PicUploadType = 1
	event.PicUploadResult = 1
}

func (event *WamEventProfilePicUpload) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Body().
		SerializeNumber(0x4, event.PicSize).
		SerializeNumber(0x6, event.PicTotalTime).
		SerializeNumber(0x3, event.PicUploadTime).
		SerializeNumber(0x1, event.PicUploadResult)

	buffer.Footer().
		SerializeNumber(0x5, event.PicUploadType)
}
