package events

import (
	"time"
	eventSerialize "ws/framework/plugin/event_serialize"
	"ws/framework/utils"
)

type PicType float64

const (
	PTGroup     PicType = 1 // 群头像
	PTMine      PicType = 2 // 自己头像
	PTThumbnail PicType = 1 // 缩略图
	PTFull      PicType = 2 // 大图
)

// WamEventProfilePicDownload .
type WamEventProfilePicDownload struct {
	WAMessageEvent

	DownloadResult float64
	DownloadTime   float64
	DownloadSize   float64
	PicType        float64
}

type ProfilePicDownloadOption struct {
	PictureType    PicType
	RequestCode    int
	PictureChanged bool
	PictureSize    int32
}

func WithProfilePicDownloadOption(pic PicType, code int, changed bool, size int32) ProfilePicDownloadOption {
	if pic == PTMine || pic == PTGroup {
		changed = true
	}

	return ProfilePicDownloadOption{
		PictureType:    pic,
		RequestCode:    code,
		PictureChanged: changed,
		PictureSize:    size,
	}
}

// InitFields .
func (event *WamEventProfilePicDownload) InitFields(option interface{}) {
	if opt, ok := option.(ProfilePicDownloadOption); ok {
		event.PicType = float64(opt.PictureType)
		event.DownloadResult = code2Result(opt.RequestCode, opt.PictureChanged)
		event.DownloadSize = float64(opt.PictureSize)
	}
	event.DownloadTime = utils.LogRandMillSecond(0, 3*time.Second)
}

func (event *WamEventProfilePicDownload) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Body().
		SerializeNumber(0x4, event.DownloadSize).
		SerializeNumber(0x1, event.DownloadResult).
		SerializeNumber(0x3, event.DownloadTime)

	buffer.Footer().
		SerializeNumber(0x2, event.PicType)
}

func code2Result(code int, changed bool) float64 {
	if code == 0 {
		if changed {
			return 1
		} else {
			return 2
		}
	}
	if code == 401 {
		return 5
	} else if code == 404 {
		return 3
	} else if code < 500 {
		return 4
	} else {
		return 6
	}
}
