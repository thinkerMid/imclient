package wam

import (
	containerInterface "ws/framework/application/container/abstract_interface"
	. "ws/framework/application/core/wam/events"
)

/*
* 账号首次注册启动，和 正常启动登录 相关的日志
 */

func (m *manager) LogLogin(container containerInterface.IAppIocContainer) {
	cache := container.ResolveChannel0EventCache()
	cache2 := container.ResolveChannel2EventCache()
	_ = cache2

	var evt containerInterface.WaEvent

	evt = NewWAMEvent(ET_WamEventIphoneIcloudUbt, nil)
	cache.AddEvent(evt)

	//evt = NewWAMEvent(ET_WamEventCrashLog, nil)
	//cache.AddEvent(evt)

	evt = NewWAMEvent(ET_WamEventAppLaunch, nil)
	cache.AddEvent(evt)

	//evt = NewWAMEvent(ET_WamEventBadInteraction, wams.WithBadInteractionOption(wams.InteractionLaunch))
	//cache.AddEvent(evt)
}

// LogRegisterLaunch 用户注册后第一次登录
func (m *manager) LogRegisterLaunch(container containerInterface.IAppIocContainer) {
	device := container.ResolveDeviceService().Context()
	cache := container.ResolveChannel0EventCache()
	cache2 := container.ResolveChannel2EventCache()
	_ = cache2

	uuid := device.FBUuid

	var evt containerInterface.WaEvent

	evt = NewWAMEvent(ET_WamEventIphoneIcloudUbt, nil)
	cache.AddEvent(evt)

	evt = NewWAMEvent(ET_WamEventAppLaunch, nil)
	cache.AddEvent(evt)

	//evt = NewWAMEvent(ET_WamEventBadInteraction, wams.WithBadInteractionOption(wams.InteractionLaunch))
	//cache.AddEvent(evt)

	evt = NewWAMEvent(ET_WamEventPsIdCreate, nil)
	cache.AddEvent(evt)

	evt = NewWAMEvent(ET_WamEventICloudRestore, nil)
	cache.AddEvent(evt)

	//evt = NewWAMEvent(ET_WamEventMessageReceive, WithMessageReceiveOption(MediaImage))
	//cache.AddEvent(evt)

	// 3个WamEventProfilePicDownload
	evt = NewWAMEvent(ET_WamEventProfilePicDownload, WithProfilePicDownloadOption(PTMine, 1, true, 0))
	cache.AddEvent(evt)

	evt = NewWAMEvent(ET_WamEventProfilePicDownload, WithProfilePicDownloadOption(PTMine, 404, true, 0))
	cache.AddEvent(evt)

	evt = NewWAMEvent(ET_WamEventProfilePicDownload, WithProfilePicDownloadOption(PTThumbnail, 404, false, 0))
	cache.AddEvent(evt)

	//idxList := []int32{2, 0, 2, 1}
	//for idx := range idxList {
	//	evt = NewWAMEvent(ET_WamEventStickerCommonQueryToStaticServer, WithStickerCommonQueryToStaticServerOption(int(idxList[idx]), lang, country))
	//	cache.AddEvent(evt)
	//}

	// mode:2
	evt = NewWAMEvent(ET_WamEventMediaDownload2, MediaDownload2Option{Mode: 2})
	cache.AddEvent(evt)

	evt = NewWAMEvent(ET_WamEventAdvPrimaryIdentityMissing, nil)
	cache.AddEvent(evt)

	evt = NewWAMEvent(ET_WamEventMdAppStateDataDeletion, nil)
	cache.AddEvent(evt)

	evt = NewWAMEvent(ET_WamEventMdAppStateCompanionsRemoval, nil)
	cache.AddEvent(evt)

	evt = NewWAMEvent(ET_WamEventStickerCommonQueryToStaticServer, WithStickerCommonQueryToStaticServerOption(1, "", ""))
	cache.AddEvent(evt)

	// mode:1
	for idx := 0; idx != 6; idx++ {
		evt = NewWAMEvent(ET_WamEventMediaDownload2, MediaDownload2Option{Mode: 1})
		cache.AddEvent(evt)
	}

	evt = NewWAMEvent(ET_WamEventStickerPackDownload, nil)
	cache.AddEvent(evt)

	evt = NewWAMEvent(ET_WamEventRegistrationComplete, RegistrationCompleteOption{Identify: uuid})
	cache.AddEvent(evt)

	evt = NewWAMEvent(ET_WamEventIphoneMediaDownloadTransform, nil)
	cache.AddEvent(evt)

	// 以下日志在登陆完成后陆续出现
	evt = NewWAMEvent(ET_WamEventStatusDaily, nil)
	cache.AddEvent(evt)

	evt = NewWAMEvent(ET_WamEventPrivacyHighlightDaily, nil)
	cache.AddEvent(evt)

	evt = NewWAMEvent(ET_WamEventNotificationSetting, nil)
	cache.AddEvent(evt)

	evt = NewWAMEvent(ET_WamEventArchiveStateDaily, nil)
	cache.AddEvent(evt)

	evt = NewWAMEvent(ET_WamEventPttDaily, nil)
	cache.AddEvent(evt)

	evt = NewWAMEvent(ET_WamEventCommunityTabAction, nil)
	cache.AddEvent(evt)

	evt = NewWAMEvent(ET_WamEventDaily, WithDailyOption(device.Language, device.Area))
	cache.AddEvent(evt)

	evt = NewWAMEvent(ET_WamEventLogin, nil)
	cache.AddEvent(evt)

	// 渠道二
	evt = NewWAMEvent(ET_WamEventSaveToCameraDaily, nil)
	cache2.AddEvent(evt)

	evt = NewWAMEvent(ET_TestAnonymousDaily, nil)
	cache2.AddEvent(evt)
}
