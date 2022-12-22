package wam

import (
	containerInterface "ws/framework/application/container/abstract_interface"
	. "ws/framework/application/core/wam/events"
)

type Media struct {
	*Image
	*Video
	*Voice
}

type Image struct {
	Width     int32
	Height    int32
	Size      int32
	FirstScan uint32
	LowScan   uint32
	MidScan   uint32
}

type Video struct {
	Width  int32
	Height int32
	Size   int32
}

type Voice struct {
	Size int32
}

func (m *manager) LogSendText(container containerInterface.IAppIocContainer, stateCount int32, fstMsg bool) {
	cache := container.ResolveChannel0EventCache()
	cache2 := container.ResolveChannel2EventCache()
	_ = cache2

	var evt containerInterface.WaEvent
	// 点击输入框触发一次
	evt = NewWAMEvent(ET_WamEventChatComposerAction, WithChatComposerActionOption(TargetText))
	cache2.AddEvent(evt)

	// 输入文本过程中触发
	if stateCount > 1 {
		for idx := 0; idx != int(stateCount-1); idx++ {
			evt = NewWAMEvent(ET_WamEventChatComposerAction, WithChatComposerActionOption(TargetText))
			cache2.AddEvent(evt)
		}
	}

	if fstMsg {
		evt = NewWAMEvent(ET_EventMessageSend, WithMessageSendOption(fstMsg, MessagePrivate, MediaText, 0))
		cache.AddEvent(evt)
	}
}

func (m *manager) LogSendMedia(container containerInterface.IAppIocContainer, mediaType MediaType, media Media) {
	cache := container.ResolveChannel0EventCache()
	cache2 := container.ResolveChannel2EventCache()
	_ = cache2

	var evt containerInterface.WaEvent

	if mediaType != MediaVoice {
		// 点击加号打开媒体选择框 触发一次
		evt = NewWAMEvent(ET_WamEventChatComposerAction, WithChatComposerActionOption(TargetMedia))
		cache2.AddEvent(evt)
		// 选择相册触发一次
		evt = NewWAMEvent(ET_WamEventMediaBrowser, nil)
		cache.AddEvent(evt)
	}

	switch mediaType {
	case MediaImage:
		// 选中图片
		evt = NewWAMEvent(ET_WamEventIphonePjpegEncoding, nil)
		cache.AddEvent(evt)
		// 点发送
		evt = NewWAMEvent(ET_WithDocumentDetection, WithDocumentDetection(media.Image.Width, media.Image.Height))
		cache.AddEvent(evt)

		evt = NewWAMEvent(ET_WamEventIphoneImageExport, WithIphoneImageExportOption(media.Image.Width, media.Image.Height,
			media.Image.Size, media.FirstScan, media.LowScan, media.MidScan))
		cache.AddEvent(evt)

		evt = NewWAMEvent(ET_WamEventMediaPicker, WithEventMediaPicker(MediaImage))
		cache.AddEvent(evt)

		evt = NewWAMEvent(ET_WamEventOptimisticUploadIndividual, nil)
		cache.AddEvent(evt)

		evt = NewWAMEvent(ET_WamEventMediaPickerPerf, nil)
		cache.AddEvent(evt)

		evt = NewWAMEvent(ET_WamEventMalformedMediaResponse, nil)
		cache.AddEvent(evt)

		evt = NewWAMEvent(ET_WamEventMediaUpload2, WithMediaUpload2Option(media.Image.Size, MediaImage))
		cache.AddEvent(evt)
	case MediaVideo:
		// 选中视频
		evt = NewWAMEvent(ET_WamEventIphoneVideoCaching, nil)
		cache.AddEvent(evt)

		evt = NewWAMEvent(ET_WamEventVideoTranscoder, WithVideoTranscoderOption(media.Video.Width, media.Video.Height, media.Video.Size))
		cache.AddEvent(evt)
		// 点击发送
		evt = NewWAMEvent(ET_WamEventMediaPicker, WithEventMediaPicker(MediaVideo))
		cache.AddEvent(evt)

		evt = NewWAMEvent(ET_WamEventMediaPickerPerf, nil)
		cache.AddEvent(evt)

		evt = NewWAMEvent(ET_WamEventMp4Repair, nil)
		cache.AddEvent(evt)

		evt = NewWAMEvent(ET_WamEventMalformedMediaResponse, nil)
		cache.AddEvent(evt)

		evt = NewWAMEvent(ET_WamEventMediaUpload2, WithMediaUpload2Option(media.Video.Size, MediaVideo))
		cache.AddEvent(evt)
	case MediaVoice:
		evt = NewWAMEvent(ET_WamEventForwardPicker, nil)
		cache.AddEvent(evt)

		evt = NewWAMEvent(ET_WamEventMediaUpload2, WithMediaUpload2Option(media.Voice.Size, MediaVoice))
		cache.AddEvent(evt)

		// ET_EventMessageSend
	}
}

func (m *manager) LogSendVCards(container containerInterface.IAppIocContainer, contacts []string) {
	cache := container.ResolveChannel0EventCache()
	cache2 := container.ResolveChannel2EventCache()
	_ = cache2

	var evt containerInterface.WaEvent
	// 点击输入框触发一次
	evt = NewWAMEvent(ET_WamEventChatComposerAction, WithChatComposerActionOption(TargetText))
	cache2.AddEvent(evt)

	evt = NewWAMEvent(ET_WamEventContactSend, WithEventContactSend(len(contacts) > 0, int32(len(contacts))))
	cache.AddEvent(evt)
}
