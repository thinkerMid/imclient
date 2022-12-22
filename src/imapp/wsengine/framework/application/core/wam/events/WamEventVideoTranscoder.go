package events

import (
	"math/rand"
	"time"
	eventSerialize "ws/framework/plugin/event_serialize"
	"ws/framework/utils"
)

type WamEventVideoTranscoder struct {
	WAMessageEvent

	SourceFormat                       float64 //0xe
	SourceDuration                     float64 //0x8
	SourceAudioBitRate                 float64 //0xc
	SourceFileSize                     float64 //0x7
	SourceFrameRate                    float64 //0xd
	SourceHeight                       float64 //0xa
	SourceVideoBitRate                 float64 //0xb
	SourceWidth                        float64 //0x9
	TargetFormat                       float64 //0x16
	TranscoderAlgorithm                float64 //0x1
	TranscoderResult                   float64 //0x2
	TargetDuration                     float64 //0x10
	TranscoderTime                     float64 //0x3
	TargetAudioBitRate                 float64 //0x14
	TargetFileSize                     float64 //0xf
	TargetFrameRate                    float64 //0x15
	TargetHeight                       float64 //0x12
	TargetVideoBitRate                 float64 //0x13
	TargetWidth                        float64 //0x11
	TranscoderContainsVideoComposition float64 //0x5
	TranscoderHasEdits                 float64 //0x6
	TranscoderIsPassThrough            float64 //0x4
}

type VideoTranscoderOption struct {
	Width    int32
	Height   int32
	FileSize int32
}

func WithVideoTranscoderOption(w, h, size int32) VideoTranscoderOption {
	return VideoTranscoderOption{
		Width:    w,
		Height:   h,
		FileSize: size,
	}
}

func (event *WamEventVideoTranscoder) InitFields(option interface{}) {
	event.SourceFormat = 1
	event.SourceDuration = utils.LogRandMillSecond(3*time.Second, 20*time.Second)
	event.SourceAudioBitRate = utils.LogRandSecond(150000*time.Second, 170000*time.Second)
	event.SourceVideoBitRate = utils.LogRandSecond(300000*time.Second, 400000*time.Second)
	event.SourceFrameRate = utils.LogRandSecond(5*time.Second, 60*time.Second) + rand.Float64()

	event.TargetAudioBitRate = 64000
	event.TargetFormat = 1
	event.TargetDuration = event.SourceDuration - utils.LogRandMillSecond(time.Second, 3*time.Second)
	event.TargetHeight = 480
	event.TargetWidth = 848

	event.TranscoderAlgorithm = 0
	event.TranscoderResult = 1
	event.TranscoderTime = utils.LogRandMillSecond(0, 1*time.Second)
	event.TranscoderContainsVideoComposition = 0
	event.TranscoderHasEdits = 0
	event.TranscoderIsPassThrough = 1

	if opt, ok := option.(VideoTranscoderOption); ok {
		event.SourceFileSize = float64(opt.FileSize)
		event.SourceHeight = float64(opt.Width)
		event.SourceWidth = float64(opt.Height)

		event.TargetFileSize = event.SourceFileSize / 10
	}
}

func (event *WamEventVideoTranscoder) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Body().
		SerializeNumber(0x1, event.TranscoderAlgorithm).
		SerializeNumber(0x2, event.TranscoderResult).
		SerializeNumber(0x3, event.TranscoderTime).
		SerializeNumber(0x4, event.TranscoderIsPassThrough).
		SerializeNumber(0x5, event.TranscoderContainsVideoComposition).
		SerializeNumber(0x6, event.TranscoderHasEdits).
		SerializeNumber(0x7, event.SourceFileSize).
		SerializeNumber(0x8, event.SourceDuration).
		SerializeNumber(0x9, event.SourceWidth).
		SerializeNumber(0xa, event.SourceHeight).
		SerializeNumber(0xb, event.SourceVideoBitRate).
		SerializeNumber(0xc, event.SourceAudioBitRate).
		SerializeNumber(0xd, event.SourceFrameRate).
		SerializeNumber(0xe, event.SourceFormat).
		SerializeNumber(0xf, event.TargetFileSize).
		SerializeNumber(0x10, event.TargetDuration).
		SerializeNumber(0x11, event.TargetWidth).
		SerializeNumber(0x12, event.TargetHeight).
		SerializeNumber(0x13, event.TargetVideoBitRate).
		SerializeNumber(0x14, event.TargetAudioBitRate).
		SerializeNumber(0x15, event.TargetFrameRate)

	buffer.Footer().
		SerializeNumber(0x16, event.TargetFormat)
}
