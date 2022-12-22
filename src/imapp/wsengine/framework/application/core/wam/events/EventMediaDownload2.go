package events

import (
	"time"
	eventSerialize "ws/framework/plugin/event_serialize"
	"ws/framework/utils"
)

/*
 * 1. index:0x12 field:3212c8  真机:4212c800
 * 2. ida最后一个字段为不参与序列化字段
 * 需要具体跟踪序列化内部实现
 */
const (
	DownloadMode1 float64 = 1
	DownloadMode2 float64 = 2

	DefaultVal float64 = -1 // 值为DefaultVal不参与序列化
)

type WamEventMediaDownload2 struct {
	WAMessageEvent

	OverallBackendStore         float64
	OverallDownloadMode         float64
	OverallDownloadOrigin       float64
	OverallDownloadResult       float64
	OverallMediaType            float64
	OverallConnBlockFetchTime   float64
	OverallCumTime              float64
	OverallDecryptTime          float64
	OverallFileValidationTime   float64
	OverallQueueTime            float64
	OverallTime                 float64
	OverallIsEncrypted          float64
	OverallIsFinal              float64
	OverallMediaSize            float64
	OverallMmsVersion           float64
	OverallAttemptCount         float64
	DownloadBytesTransferred    float64
	DownloadHttpCode            float64
	DownloadIsReuse             float64
	DownloadIsStreaming         float64
	DownloadResumePoint         float64
	DownloadConnectTime         float64
	DownloadNetworkTime         float64
	DownloadTimeToFirstByteTime float64
	OverallRetryCount           float64
	IsViewOnce                  float64
}

type MediaDownload2Option struct {
	Mode int
}

func (event *WamEventMediaDownload2) InitFields(option interface{}) {
	if opt, ok := option.(MediaDownload2Option); ok {
		event.OverallDownloadMode = float64(opt.Mode)
	}

	if event.OverallDownloadMode == DownloadMode1 {
		event.OverallBackendStore = 2
		event.OverallDownloadMode = 1
		event.OverallDownloadOrigin = 7
		event.OverallDownloadResult = 1
		event.OverallMediaType = 16
		event.OverallConnBlockFetchTime = 0
		event.OverallCumTime = utils.LogRandMillSecond(time.Second, 5*time.Second)
		event.OverallDecryptTime = utils.LogRandMillSecond(time.Second, 200*time.Second)
		event.OverallFileValidationTime = utils.LogRandMillSecond(time.Second, 30*time.Second)
		event.OverallQueueTime = utils.LogRandMillSecond(time.Second, 30*time.Second)
		event.OverallTime = utils.LogRandMillSecond(time.Second, 5*time.Second)
		event.OverallIsEncrypted = 1
		event.OverallIsFinal = 1
		event.OverallMediaSize = float64(int(utils.LogRandMillSecond(5*time.Second, 20*time.Second)))
		event.OverallMmsVersion = 4
		event.OverallRetryCount = 0
		event.OverallAttemptCount = DefaultVal
		event.DownloadBytesTransferred = float64(int(utils.LogRandMillSecond(time.Second, 100*time.Second)))
		event.DownloadHttpCode = 200
		event.DownloadIsReuse = 1
		event.DownloadIsStreaming = 0
		event.DownloadResumePoint = 0
		event.DownloadConnectTime = 0
		event.DownloadNetworkTime = utils.LogRandSecond(50*time.Second, 200*time.Second)
		event.DownloadTimeToFirstByteTime = event.DownloadNetworkTime - utils.LogRandSecond(time.Second, 5*time.Second)
		event.IsViewOnce = DefaultVal
	} else if event.OverallDownloadMode == DownloadMode2 {
		event.OverallBackendStore = DefaultVal
		event.OverallDownloadOrigin = 3
		event.OverallConnBlockFetchTime = DefaultVal
		event.OverallCumTime = utils.LogRandMillSecond(time.Second, 10*time.Second)
		event.OverallDecryptTime = utils.LogRandMillSecond(time.Second, 200*time.Second)
		event.OverallFileValidationTime = utils.LogRandSecond(0, time.Second)
		event.OverallQueueTime = DefaultVal
		event.OverallTime = utils.LogRandMillSecond(time.Second, 2*time.Second)
		event.OverallAttemptCount = 0
		event.DownloadBytesTransferred = float64(int(utils.LogRandMillSecond(time.Second, 100*time.Second)))
		event.DownloadHttpCode = 200
		event.DownloadIsReuse = DefaultVal
		event.DownloadResumePoint = DefaultVal
		event.DownloadConnectTime = DefaultVal
		event.DownloadNetworkTime = DefaultVal
		event.DownloadTimeToFirstByteTime = DefaultVal
		event.IsViewOnce = 0

		event.OverallDownloadMode = 1
		event.OverallDownloadResult = 1
		event.OverallIsEncrypted = 1
		event.OverallIsFinal = 1
		event.OverallMmsVersion = 4
		event.OverallRetryCount = 0
		event.DownloadIsStreaming = 0

		event.OverallMediaType = 2 //?
		event.OverallMediaSize = float64(int(utils.LogRandMillSecond(5*time.Second, 20*time.Second)))
	}
}

func (event *WamEventMediaDownload2) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Body().
		SerializeNumber(0x27, event.OverallBackendStore).
		SerializeNumber(0xb, event.OverallDownloadMode).
		SerializeNumber(0x23, event.OverallDownloadOrigin).
		SerializeNumber(0x19, event.OverallDownloadResult).
		SerializeNumber(0x1, event.OverallMediaType).
		SerializeNumber(0xa, event.OverallConnBlockFetchTime).
		SerializeNumber(0x1b, event.OverallCumTime).
		SerializeNumber(0xc, event.OverallDecryptTime).
		SerializeNumber(0xd, event.OverallFileValidationTime).
		SerializeNumber(0x9, event.OverallQueueTime).
		SerializeNumber(0x8, event.OverallTime).
		SerializeNumber(0x1c, event.OverallIsEncrypted).
		SerializeNumber(0x1a, event.OverallIsFinal).
		SerializeNumber(0x7, event.OverallMediaSize).
		SerializeNumber(0x6, event.OverallMmsVersion).
		SerializeNumber(0x3, event.OverallRetryCount).
		SerializeNumber(0x4, event.OverallAttemptCount).
		SerializeNumber(0x14, event.DownloadBytesTransferred).
		SerializeNumber(0x12, event.DownloadHttpCode).
		SerializeNumber(0x11, event.DownloadIsReuse).
		SerializeNumber(0x13, event.DownloadIsStreaming).
		SerializeNumber(0xe, event.DownloadResumePoint).
		SerializeNumber(0xf, event.DownloadConnectTime).
		SerializeNumber(0x10, event.DownloadNetworkTime).
		SerializeNumber(0x15, event.DownloadTimeToFirstByteTime)

	buffer.Footer().
		SerializeNumber(0x29, event.IsViewOnce)
}
