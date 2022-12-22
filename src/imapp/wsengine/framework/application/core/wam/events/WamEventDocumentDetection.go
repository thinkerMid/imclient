package events

import (
	"time"
	eventSerialize "ws/framework/plugin/event_serialize"
	"ws/framework/utils"
)

type WamEventDocumentDetection struct {
	WAMessageEvent

	AlgorithmType           float64
	AlgorithmTime           float64
	DecodeBitmapTime        float64
	DocumentImageMaxEdge    float64
	DocumentImageQuality    float64
	ImageHeightRorAlgorithm float64
	ImageHeightRorSending   float64
	ImageWidthRorAlgorithm  float64
	ImageWidthRorSending    float64
	IsImageLikeDocument     float64
	IsImageLikeDocumentByWA float64
	// marginX = ImageWidth * 0.25
	MarginX             float64
	MarginY             float64
	OriginalImageHeight float64
	OriginalImageWidth  float64
}

type DocumentDetectionOption struct {
	Width   int32
	Height  int32
	MarginX int32
	MarginY int32
}

func WithDocumentDetection(w, h int32) DocumentDetectionOption {
	//? 非1:1 图片有缩放
	return DocumentDetectionOption{
		Width:  w,
		Height: h,
	}
}

func (event *WamEventDocumentDetection) InitFields(option interface{}) {
	event.AlgorithmTime = float64(int(utils.LogRandSecond(3*time.Second, 250*time.Second)))
	event.DecodeBitmapTime = float64(int(utils.LogRandSecond(5*time.Second, 100*time.Second)))
	event.DocumentImageMaxEdge = 1280
	event.DocumentImageQuality = 0.75
	event.ImageHeightRorAlgorithm = 0
	event.ImageHeightRorSending = 0
	event.ImageWidthRorAlgorithm = 0
	event.ImageWidthRorSending = 0
	event.MarginX = 0 // w / 4
	event.MarginY = 0 // h / 4
	event.OriginalImageHeight = 0
	event.OriginalImageWidth = 0
	event.IsImageLikeDocument = 0
	event.IsImageLikeDocumentByWA = 0
	event.AlgorithmType = 2

	if opt, ok := option.(DocumentDetectionOption); ok {
		event.ImageWidthRorAlgorithm = float64(opt.Width)
		event.ImageWidthRorSending = float64(opt.Width)
		event.ImageHeightRorAlgorithm = float64(opt.Height)
		event.ImageHeightRorSending = float64(opt.Height)
		event.MarginX = float64(int(event.ImageWidthRorAlgorithm / 4))
		event.MarginY = float64(int(event.ImageHeightRorAlgorithm / 4))
	}
}

func (event *WamEventDocumentDetection) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Body().
		SerializeNumber(0x11, event.AlgorithmType).
		SerializeNumber(0x1, event.AlgorithmTime).
		SerializeNumber(0x2, event.DecodeBitmapTime).
		SerializeNumber(0x4, event.DocumentImageMaxEdge).
		SerializeNumber(0x5, event.DocumentImageQuality).
		SerializeNumber(0x7, event.ImageHeightRorAlgorithm).
		SerializeNumber(0x9, event.ImageHeightRorSending).
		SerializeNumber(0x6, event.ImageWidthRorAlgorithm).
		SerializeNumber(0x8, event.ImageWidthRorSending).
		SerializeNumber(0x3, event.IsImageLikeDocument).
		SerializeNumber(0x10, event.IsImageLikeDocumentByWA).
		SerializeNumber(0xc, event.MarginX).
		SerializeNumber(0xd, event.MarginY).
		SerializeNumber(0xb, event.OriginalImageHeight)

	buffer.Footer().
		SerializeNumber(0xa, event.OriginalImageWidth)
}
