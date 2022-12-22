package wam

import (
	"encoding/hex"
	"fmt"
	"ws/framework/application/constant/types"
	containerInterface "ws/framework/application/container/abstract_interface"
	. "ws/framework/application/core/wam/events"
)

const (
	ET_MIN int32 = iota
	// 基础包
	ET_WAMCommon
	// 登陆
	ET_WamEventIphoneIcloudUbt
	ET_WamEventAppLaunch
	ET_WamEventLogin
	// 注册
	ET_WamEventPsIdCreate
	ET_WamEventICloudRestore
	ET_WamEventProfilePicDownload
	ET_WamEventStickerCommonQueryToStaticServer
	ET_WamEventMediaDownload2
	ET_WamEventMessageReceive
	ET_WamEventRegistrationComplete
	//ET_WamEventIphoneMediaDownloadTransform
	ET_WamEventStickerPackDownload
	ET_WamEventStatusDaily
	ET_WamEventArchiveStateDaily
	ET_WamEventPttDaily
	ET_WamEventDaily
	ET_WamEventMdAppStateDataDeletion
	ET_WamEventMdAppStateCompanionsRemoval
	ET_WamEventIphoneMediaDownloadTransform
	// 个人信息
	ET_WamEventSettingsClick
	ET_WamEventProfilePicUpload
	ET_WamEventUiUsage

	// 好友
	ET_WamEventIphoneAddContactEvent // 打开或关闭联系人信息界面都会触发该日志，包括新建联系人，查看/编辑联系人信息

	// 会话
	ET_WamEventIphoneContactListStartNewChat // 新建会话时才需要发
	ET_WamEventPrekeysFetch
	ET_WamEventPrekeysDepletion // 版本升级弃用
	// 私聊
	ET_EventAndroidMessageSendPerf
	ET_EventMessageSend
	ET_EventE2eMessageSend
	// 群组
	ET_WamEventGroupCreateInit
	ET_WamEventIphoneGroupCreate
	// channel2
	ET_TestAnonymousDaily

	ET_WamEventLocationPicker
	ET_WamEventChatComposerAction
	ET_WamEventMalformedMediaResponse
	ET_WamEventMediaUpload2
	ET_WamEventCrashLog
	ET_WamEventSaveToCameraDaily
	ET_WamEventContactInfo
	ET_WamEventMediaBrowser
	ET_WamEventIphonePjpegEncoding
	ET_WithDocumentDetection
	ET_WamEventIphoneImageExport
	ET_WamEventMediaPicker
	ET_WamEventOptimisticUploadIndividual
	ET_WamEventMediaPickerPerf
	ET_WamEventIphoneVideoCaching
	ET_WamEventVideoTranscoder
	ET_WamEventMp4Repair
	ET_WamEventForwardPicker
	ET_WamEventBadInteraction
	ET_WamEventAdvPrimaryIdentityMissing
	ET_WamEventNotificationSetting
	ET_WamEventCommunityTabAction
	ET_WamEventIphoneTestRealtime
	ET_WamEventWamTestAnonymous0
	ET_WamEventTestAnonymousDailyId
	ET_WamEventGroupProfilePicture
	ET_WamEventGroupInfo
	ET_WamEventContactSend
	ET_WamEventPrivacyHighlightDaily
	ET_MAX
)

type EventCreator func() containerInterface.WaEvent

type LogEventItem struct {
	EventType int32
	Instance  EventCreator

	Channel uint8
	Code    int64
	Weight  float64
}

func (lei LogEventItem) insert2Map() {
	if _, ok := logEventMap[lei.EventType]; ok {
		return
	} else {
		logEventMap[lei.EventType] = lei
	}
}

var (
	logEventMap map[int32]LogEventItem
)

func init() {
	logEventMap = make(map[int32]LogEventItem)
	LogEventItem{Channel: 0, Code: 1102, Weight: 1, EventType: ET_WamEventIphoneIcloudUbt, Instance: func() containerInterface.WaEvent { return new(WamEventIphoneIcloudUbt) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 1094, Weight: 1, EventType: ET_WamEventAppLaunch, Instance: func() containerInterface.WaEvent { return new(WamEventAppLaunch) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 460, Weight: 10, EventType: ET_WamEventLogin, Instance: func() containerInterface.WaEvent { return new(WamEventLogin) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 2310, Weight: 1, EventType: ET_WamEventPsIdCreate, Instance: func() containerInterface.WaEvent { return new(WamEventPsIdCreate) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 502, Weight: 1, EventType: ET_WamEventICloudRestore, Instance: func() containerInterface.WaEvent { return new(WamEventICloudRestore) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 848, Weight: 1, EventType: ET_WamEventProfilePicDownload, Instance: func() containerInterface.WaEvent { return new(WamEventProfilePicDownload) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 2740, Weight: 20, EventType: ET_WamEventStickerCommonQueryToStaticServer, Instance: func() containerInterface.WaEvent { return new(WamEventStickerCommonQueryToStaticServer) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 1590, Weight: 5, EventType: ET_WamEventMediaDownload2, Instance: func() containerInterface.WaEvent { return new(WamEventMediaDownload2) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 450, Weight: 20, EventType: ET_WamEventMessageReceive, Instance: func() containerInterface.WaEvent { return new(WamEventMessageReceive) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 1342, Weight: 1, EventType: ET_WamEventRegistrationComplete, Instance: func() containerInterface.WaEvent { return new(WamEventRegistrationComplete) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 1844, Weight: 1, EventType: ET_WamEventStickerPackDownload, Instance: func() containerInterface.WaEvent { return new(WamEventStickerPackDownload) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 1676, Weight: 1, EventType: ET_WamEventStatusDaily, Instance: func() containerInterface.WaEvent { return new(WamEventStatusDaily) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 2810, Weight: 1, EventType: ET_WamEventArchiveStateDaily, Instance: func() containerInterface.WaEvent { return new(WamEventArchiveStateDaily) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 2938, Weight: -1, EventType: ET_WamEventPttDaily, Instance: func() containerInterface.WaEvent { return new(WamEventPttDaily) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 1158, Weight: 1, EventType: ET_WamEventDaily, Instance: func() containerInterface.WaEvent { return new(WamEventDaily) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 468, Weight: 1, EventType: ET_WamEventProfilePicUpload, Instance: func() containerInterface.WaEvent { return new(WamEventProfilePicUpload) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 2194, Weight: 1, EventType: ET_WamEventIphoneAddContactEvent, Instance: func() containerInterface.WaEvent { return new(WamEventIphoneAddContactEvent) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 3090, Weight: 1, EventType: ET_WamEventIphoneContactListStartNewChat, Instance: func() containerInterface.WaEvent { return new(WamEventIphoneContactListStartNewChat) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 2540, Weight: 1, EventType: ET_WamEventPrekeysFetch, Instance: func() containerInterface.WaEvent { return new(WamEventPrekeysFetch) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 3014, Weight: 1, EventType: ET_WamEventPrekeysDepletion, Instance: func() containerInterface.WaEvent { return new(WamEventPrekeysDepletion) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 2214, Weight: -1, EventType: ET_WamEventSettingsClick, Instance: func() containerInterface.WaEvent { return new(WamEventSettingsClick) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 474, Weight: 1, EventType: ET_WamEventUiUsage, Instance: func() containerInterface.WaEvent { return new(WamEventUiUsage) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 3286, Weight: 1, EventType: ET_WamEventGroupCreateInit, Instance: func() containerInterface.WaEvent { return new(EventGroupCreateInit) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 3288, Weight: 1, EventType: ET_WamEventIphoneGroupCreate, Instance: func() containerInterface.WaEvent { return new(EventIphoneGroupCreate) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 2510, Weight: 1, EventType: ET_WamEventMdAppStateDataDeletion, Instance: func() containerInterface.WaEvent { return new(WamEventMdAppStateDataDeletion) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 2508, Weight: 1, EventType: ET_WamEventMdAppStateCompanionsRemoval, Instance: func() containerInterface.WaEvent { return new(WamEventMdAppStateCompanionsRemoval) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 1152, Weight: 1, EventType: ET_WamEventIphoneMediaDownloadTransform, Instance: func() containerInterface.WaEvent { return new(WamEventIphoneMediaDownloadTransform) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 1994, Weight: 1, EventType: ET_EventAndroidMessageSendPerf, Instance: func() containerInterface.WaEvent { return new(EventAndroidMessageSendPerf) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 854, Weight: 5, EventType: ET_EventMessageSend, Instance: func() containerInterface.WaEvent { return new(EventMessageSend) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 476, Weight: 20, EventType: ET_EventE2eMessageSend, Instance: func() containerInterface.WaEvent { return new(EventE2eMessageSend) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 482, Weight: 1, EventType: ET_WamEventLocationPicker, Instance: func() containerInterface.WaEvent { return new(WamEventLocationPicker) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 1752, Weight: 1, EventType: ET_WamEventMalformedMediaResponse, Instance: func() containerInterface.WaEvent { return new(WamEventMalformedMediaResponse) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 1588, Weight: 1, EventType: ET_WamEventMediaUpload2, Instance: func() containerInterface.WaEvent { return new(WamEventMediaUpload2) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 494, Weight: 1, EventType: ET_WamEventCrashLog, Instance: func() containerInterface.WaEvent { return new(WamEventCrashLog) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 3124, Weight: 5, EventType: ET_WamEventContactInfo, Instance: func() containerInterface.WaEvent { return new(WamEventContactInfo) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 1532, Weight: 1, EventType: ET_WamEventMediaBrowser, Instance: func() containerInterface.WaEvent { return new(WamEventMediaBrowser) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 1686, Weight: 1, EventType: ET_WamEventIphonePjpegEncoding, Instance: func() containerInterface.WaEvent { return new(WamEventIphonePjpegEncoding) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 1436, Weight: 1, EventType: ET_WithDocumentDetection, Instance: func() containerInterface.WaEvent { return new(WamEventDocumentDetection) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 1490, Weight: 1, EventType: ET_WamEventIphoneImageExport, Instance: func() containerInterface.WaEvent { return new(WamEventIphoneImageExport) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 1038, Weight: 1, EventType: ET_WamEventMediaPicker, Instance: func() containerInterface.WaEvent { return new(WamEventMediaPicker) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 1488, Weight: 1, EventType: ET_WamEventOptimisticUploadIndividual, Instance: func() containerInterface.WaEvent { return new(WamEventOptimisticUploadIndividual) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 1534, Weight: 1, EventType: ET_WamEventMediaPickerPerf, Instance: func() containerInterface.WaEvent { return new(WamEventMediaPickerPerf) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 2466, Weight: 1, EventType: ET_WamEventIphoneVideoCaching, Instance: func() containerInterface.WaEvent { return new(WamEventIphoneVideoCaching) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 1802, Weight: 1, EventType: ET_WamEventVideoTranscoder, Instance: func() containerInterface.WaEvent { return new(WamEventVideoTranscoder) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 1066, Weight: 1, EventType: ET_WamEventMp4Repair, Instance: func() containerInterface.WaEvent { return new(WamEventMp4Repair) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 1034, Weight: 1, EventType: ET_WamEventForwardPicker, Instance: func() containerInterface.WaEvent { return new(WamEventForwardPicker) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 3316, Weight: -1, EventType: ET_WamEventBadInteraction, Instance: func() containerInterface.WaEvent { return new(WamEventBadInteraction) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 3044, Weight: 1, EventType: ET_WamEventAdvPrimaryIdentityMissing, Instance: func() containerInterface.WaEvent { return new(WamEventAdvPrimaryIdentityMissing) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 3684, Weight: 1, EventType: ET_WamEventNotificationSetting, Instance: func() containerInterface.WaEvent { return new(WamEventNotificationSetting) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 3496, Weight: -1, EventType: ET_WamEventCommunityTabAction, Instance: func() containerInterface.WaEvent { return new(WamEventCommunityTabAction) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 1972, Weight: 1, EventType: ET_WamEventIphoneTestRealtime, Instance: func() containerInterface.WaEvent { return new(WamEventIphoneTestRealtime) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 2240, Weight: 1, EventType: ET_WamEventWamTestAnonymous0, Instance: func() containerInterface.WaEvent { return new(WamEventWamTestAnonymous0) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 2958, Weight: 1, EventType: ET_WamEventTestAnonymousDailyId, Instance: func() containerInterface.WaEvent { return new(WamEventTestAnonymousDailyId) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 3652, Weight: 1, EventType: ET_WamEventGroupProfilePicture, Instance: func() containerInterface.WaEvent { return new(WamEventGroupProfilePicture) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 3126, Weight: 5, EventType: ET_WamEventGroupInfo, Instance: func() containerInterface.WaEvent { return new(WamEventGroupInfo) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 2234, Weight: 1, EventType: ET_WamEventContactSend, Instance: func() containerInterface.WaEvent { return new(WamEventContactSend) }}.insert2Map()
	LogEventItem{Channel: 0, Code: 3522, Weight: 1, EventType: ET_WamEventPrivacyHighlightDaily, Instance: func() containerInterface.WaEvent { return new(WamEventPrivacyHighlightDaily) }}.insert2Map()

	// CHANNEL 2
	LogEventItem{Channel: 2, Code: 2948, Weight: 1, EventType: ET_WamEventChatComposerAction, Instance: func() containerInterface.WaEvent { return new(WamEventChatComposerAction) }}.insert2Map()
	LogEventItem{Channel: 2, Code: 3210, Weight: 1, EventType: ET_WamEventSaveToCameraDaily, Instance: func() containerInterface.WaEvent { return new(WamEventSaveToCameraDaily) }}.insert2Map()
	LogEventItem{Channel: 2, Code: 2328, Weight: -1, EventType: ET_TestAnonymousDaily, Instance: func() containerInterface.WaEvent { return new(EventTestAnonymousDaily) }}.insert2Map()

}

// NewWAMEvent 根据日志类型创建日志，复用缓存的日志实例
func NewWAMEvent(eventType int32, option interface{}) containerInterface.WaEvent {
	var event containerInterface.WaEvent

	if item, ok := logEventMap[eventType]; !ok {
		return nil
	} else {
		event = item.Instance()
		event.Init(item.Channel, item.Code, item.Weight)
		event.InitFields(option)
	}
	return event
}

func PrintAllHeader() {
	header := types.SerializeBuf(types.OPTION_NUMBER, 1336, 1, "", 1)
	fmt.Printf("header:%v\n", hex.EncodeToString(header))

	for key, val := range logEventMap {
		header := types.SerializeBuf(types.OPTION_NUMBER, val.Code, val.Weight, "", 1)
		fmt.Printf("idx:%d code:%v weight:%v header:%v\n", key, val.Code, val.Weight, hex.EncodeToString(header))
	}
}
