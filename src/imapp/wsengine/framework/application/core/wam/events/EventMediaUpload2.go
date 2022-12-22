package events

import (
	"time"
	eventSerialize "ws/framework/plugin/event_serialize"
	"ws/framework/utils"
)

// WamEventMediaUpload2 .
type WamEventMediaUpload2 struct {
	WAMessageEvent

	OverallMediaKeyReuse      float64 //0x28
	OverallMediaType          float64 //0x1
	OverallOptimisticFlag     float64 //0xc
	OverallUploadMode         float64 //0x27
	OverallUploadOrigin       float64 //0x2c
	OverallUploadResult       float64 //0x23
	OverallConnBlockFetchTime float64 //0xa
	OverallCumTime            float64 //0x25
	OverallCumUserVisibleTime float64 //0x26
	OverallQueueTime          float64 //0x9
	OverallTime               float64 //0x8
	OverallTranscodeTime      float64 //0xf
	OverallUserVisibleTime    float64 //0xe
	OverallAttemptCount       float64 //0x4
	OverallIsFinal            float64 //0x24
	OverallIsForward          float64 //0x10
	OverallIsManual           float64 //0xd
	OverallMediaSize          float64 //0x7
	OverallMmsVersion         float64 //0x6
	OverallRetryCount         float64 //0x3
	UploadBytesTransferred    float64 //0x1b
	UploadHttpCode            float64 //0x19
	UploadIsReuse             float64 //0x18
	UploadIsStreaming         float64 //0x1a
	UploadResumePoint         float64 //0x15
	UploadConnectTime         float64 //0x16
	UploadNetworkTime         float64 //0x17
	ResumeHttpCode            float64 //0x14
	ResumeNetworkTime         float64 //0x12
	IsViewOnce                float64 //0x31
}

type MediaUpload2Option struct {
	MediaSize int32
	MediaType MediaType
}

func WithMediaUpload2Option(mediaSize int32, mediaType MediaType) MediaUpload2Option {
	return MediaUpload2Option{
		MediaSize: mediaSize,
		MediaType: mediaType,
	}
}

func (event *WamEventMediaUpload2) InitFields(option interface{}) {
	event.OverallMediaKeyReuse = 1
	event.OverallOptimisticFlag = 0
	event.OverallUploadMode = 1
	event.OverallUploadOrigin = 2
	event.OverallUploadResult = 1
	event.OverallConnBlockFetchTime = 0
	event.OverallCumTime = utils.LogRandMillSecond(1*time.Second, 5*time.Second)
	event.OverallCumUserVisibleTime = event.OverallCumTime + utils.LogRandMillSecond(time.Millisecond, 10*time.Millisecond)
	event.OverallQueueTime = utils.LogRandSecond(0, time.Second)
	event.OverallTime = event.OverallCumTime + utils.LogRandMillSecond(time.Millisecond, 10*time.Millisecond)
	event.OverallTranscodeTime = -1
	event.OverallUserVisibleTime = event.OverallCumTime + utils.LogRandMillSecond(time.Millisecond, 10*time.Millisecond)
	event.OverallAttemptCount = 0
	event.OverallIsFinal = 1
	event.OverallIsForward = 0
	event.OverallIsManual = 1
	event.OverallMmsVersion = 4
	event.OverallRetryCount = 0
	event.UploadHttpCode = 200
	event.UploadIsReuse = 1
	event.UploadIsStreaming = 0
	event.UploadResumePoint = 0
	event.UploadConnectTime = 0
	event.UploadNetworkTime = utils.LogRandMillSecond(0, time.Second)
	event.ResumeHttpCode = 200
	event.ResumeNetworkTime = utils.LogRandMillSecond(1*time.Second, 3*time.Second)
	event.IsViewOnce = 0

	if opt, ok := option.(MediaUpload2Option); ok {
		event.OverallMediaSize = float64(opt.MediaSize)
		event.UploadBytesTransferred = event.OverallMediaSize - 17 //传输大小与实际大小的关系 ?
		event.OverallMediaType = float64(opt.MediaType)
	}
}

func (event *WamEventMediaUpload2) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Body().
		SerializeNumber(0x28, event.OverallMediaKeyReuse).
		SerializeNumber(0x1, event.OverallMediaType).
		SerializeNumber(0xc, event.OverallOptimisticFlag).
		SerializeNumber(0x27, event.OverallUploadMode).
		SerializeNumber(0x2c, event.OverallUploadOrigin).
		SerializeNumber(0x23, event.OverallUploadResult).
		SerializeNumber(0xa, event.OverallConnBlockFetchTime).
		SerializeNumber(0x25, event.OverallCumTime).
		SerializeNumber(0x26, event.OverallCumUserVisibleTime).
		SerializeNumber(0x9, event.OverallQueueTime).
		SerializeNumber(0x8, event.OverallTime).
		SerializeNumber(0xf, event.OverallTranscodeTime).
		SerializeNumber(0xe, event.OverallUserVisibleTime).
		SerializeNumber(0x4, event.OverallAttemptCount).
		SerializeNumber(0x24, event.OverallIsFinal).
		SerializeNumber(0x10, event.OverallIsForward).
		SerializeNumber(0xd, event.OverallIsManual).
		SerializeNumber(0x7, event.OverallMediaSize).
		SerializeNumber(0x6, event.OverallMmsVersion).
		SerializeNumber(0x3, event.OverallRetryCount).
		SerializeNumber(0x1b, event.UploadBytesTransferred).
		SerializeNumber(0x19, event.UploadHttpCode).
		SerializeNumber(0x18, event.UploadIsReuse).
		SerializeNumber(0x1a, event.UploadIsStreaming).
		SerializeNumber(0x15, event.UploadResumePoint).
		SerializeNumber(0x16, event.UploadConnectTime).
		SerializeNumber(0x17, event.UploadNetworkTime).
		SerializeNumber(0x14, event.ResumeHttpCode).
		SerializeNumber(0x12, event.ResumeNetworkTime)

	buffer.Footer().
		SerializeNumber(0x31, event.IsViewOnce)
}
