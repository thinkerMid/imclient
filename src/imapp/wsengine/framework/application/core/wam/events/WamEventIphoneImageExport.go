package events

import (
	"math/rand"
	"time"
	eventSerialize "ws/framework/plugin/event_serialize"
	"ws/framework/utils"
)

type WamEventIphoneImageExport struct {
	WAMessageEvent

	IphoneExportAlgorithmTime       float64
	IphoneExportWholeFlowTime       float64
	IphoneImageFirstScanLength      float64
	IphoneImageHeightFoExport       float64
	IphoneImageLowQualityScanLength float64
	IphoneImageMidQualityScanLength float64
	IphoneImageQuality              float64
	IphoneImageSize                 float64
	IphoneImageWidthForExport       float64
	IphoneIsHighResImage            float64
}

type IphoneImageExportOption struct {
	FirstScanLength      uint32
	HeightExport         int32
	WidthExport          int32
	LowQualityScanLength uint32
	MidQualityScanLength uint32
	Size                 int32
}

func WithIphoneImageExportOption(w, h, size int32, first, low, mid uint32) IphoneImageExportOption {
	return IphoneImageExportOption{
		WidthExport:          w,
		HeightExport:         h,
		Size:                 size,
		FirstScanLength:      first,
		LowQualityScanLength: low,
		MidQualityScanLength: mid,
	}
}

func (event *WamEventIphoneImageExport) InitFields(option interface{}) {
	event.IphoneExportAlgorithmTime = float64(int(utils.LogRandSecond(30*time.Second, 100*time.Second)))
	event.IphoneExportWholeFlowTime = event.IphoneExportAlgorithmTime + float64(rand.Intn(3))
	event.IphoneImageQuality = 75
	event.IphoneIsHighResImage = 0

	if opt, ok := option.(IphoneImageExportOption); ok {
		event.IphoneImageSize = float64(opt.Size)
		event.IphoneImageFirstScanLength = float64(opt.FirstScanLength)
		event.IphoneImageWidthForExport = float64(opt.WidthExport)
		event.IphoneImageHeightFoExport = float64(opt.HeightExport)
		event.IphoneImageLowQualityScanLength = float64(opt.LowQualityScanLength)
		event.IphoneImageMidQualityScanLength = float64(opt.MidQualityScanLength)
	}
}

func (event *WamEventIphoneImageExport) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Body().
		SerializeNumber(0x6, event.IphoneExportAlgorithmTime).
		SerializeNumber(0x5, event.IphoneExportWholeFlowTime).
		SerializeNumber(0x9, event.IphoneImageFirstScanLength).
		SerializeNumber(0x4, event.IphoneImageHeightFoExport).
		SerializeNumber(0xa, event.IphoneImageLowQualityScanLength).
		SerializeNumber(0xb, event.IphoneImageMidQualityScanLength).
		SerializeNumber(0x2, event.IphoneImageQuality).
		SerializeNumber(0x8, event.IphoneImageSize).
		SerializeNumber(0x3, event.IphoneImageWidthForExport)

	buffer.Footer().
		SerializeNumber(0x7, event.IphoneIsHighResImage)
}
