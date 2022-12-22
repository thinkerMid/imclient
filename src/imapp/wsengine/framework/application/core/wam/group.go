package wam

import (
	"math/rand"
	containerInterface "ws/framework/application/container/abstract_interface"
	. "ws/framework/application/core/wam/events"
)

// LogGroupCreate .
func (m *manager) LogGroupCreate(container containerInterface.IAppIocContainer, groupAvatarSize, noKeyMembers uint32) {
	cache := container.ResolveChannel0EventCache()
	cache2 := container.ResolveChannel2EventCache()
	_ = cache2

	var evt containerInterface.WaEvent

	evt = NewWAMEvent(ET_WamEventGroupCreateInit, nil)
	cache.AddEvent(evt)

	evt = NewWAMEvent(ET_WamEventIphoneGroupCreate, nil)
	cache.AddEvent(evt)

	evt = NewWAMEvent(ET_WamEventUiUsage, nil)
	cache.AddEvent(evt)

	evt = NewWAMEvent(ET_WamEventIphoneTestRealtime, nil)
	cache.AddEvent(evt)

	evt = NewWAMEvent(ET_WamEventWamTestAnonymous0, nil)
	cache.AddEvent(evt)

	evt = NewWAMEvent(ET_WamEventTestAnonymousDailyId, nil)
	cache.AddEvent(evt)

	evt = NewWAMEvent(ET_WamEventChatComposerAction, WithChatComposerActionOption(TargetText))
	cache2.AddEvent(evt)

	evt = NewWAMEvent(ET_WamEventPrekeysFetch, nil)
	cache.AddEvent(evt)

	// WamEventPrekeysDepletion

	if groupAvatarSize > 0 {
		evt = NewWAMEvent(ET_WamEventProfilePicUpload, WithProfilePicUploadOption(0, groupAvatarSize))
		cache.AddEvent(evt)
	} else {
		// 选中联系人时，查联系人设备信息
		if noKeyMembers > 0 {
			for i := 0; i != int(noKeyMembers); i++ {
				evt = NewWAMEvent(ET_WamEventProfilePicDownload, WithProfilePicDownloadOption(PTThumbnail, 0, true, 0))
				cache.AddEvent(evt)
			}
		}
	}
}

func (m *manager) LogGroupEditAdmins(container containerInterface.IAppIocContainer, groupHasAvatar bool) {
	cache := container.ResolveChannel0EventCache()
	cache2 := container.ResolveChannel2EventCache()
	_ = cache2

	var evt containerInterface.WaEvent

	// 点击会话进入群信息面板
	if !groupHasAvatar {
		for i := 0; i != 4; i++ {
			evt = NewWAMEvent(ET_WamEventProfilePicDownload, WithProfilePicDownloadOption(PTFull, 404, false, 0))
			cache.AddEvent(evt)
		}
	}

	evt = NewWAMEvent(ET_WamEventIphoneTestRealtime, nil)
	cache.AddEvent(evt)

	evt = NewWAMEvent(ET_WamEventWamTestAnonymous0, nil)
	cache.AddEvent(evt)

	evt = NewWAMEvent(ET_WamEventTestAnonymousDailyId, nil)
	cache.AddEvent(evt)

	// 返回群会话
	evt = NewWAMEvent(ET_WamEventGroupInfo, nil)
	cache.AddEvent(evt)
}

// LogNotifyGroupEditAvatar 设置群头像
func (m *manager) LogNotifyGroupEditAvatar(container containerInterface.IAppIocContainer, preHasAvatar bool, setAvatarSize uint32) {
	cache := container.ResolveChannel0EventCache()
	cache2 := container.ResolveChannel2EventCache()
	_ = cache2

	var evt containerInterface.WaEvent

	if preHasAvatar {
		// 删除群头像
		if setAvatarSize == 0 {
			evt = NewWAMEvent(ET_WamEventGroupProfilePicture, WithGroupProfilePictureOption(1, 0, 1))
			cache.AddEvent(evt)

			evt = NewWAMEvent(ET_WamEventGroupProfilePicture, WithGroupProfilePictureOption(6, 0, 1))
			cache.AddEvent(evt)

			evt = NewWAMEvent(ET_WamEventGroupProfilePicture, WithGroupProfilePictureOption(10, 0, 1))
			cache.AddEvent(evt)

			for i := 0; i != 10; i++ {
				evt = NewWAMEvent(ET_WamEventProfilePicDownload, WithProfilePicDownloadOption(PTFull, 404, false, 0))
				cache.AddEvent(evt)
			}

			evt = NewWAMEvent(ET_WamEventProfilePicDownload, WithProfilePicDownloadOption(PTThumbnail, 404, false, 0))
			cache.AddEvent(evt)
		}

		return
	} else {
		// 修改头像，设置头像
		evt = NewWAMEvent(ET_WamEventProfilePicDownload, WithProfilePicDownloadOption(PTFull, 404, false, 0))
		cache.AddEvent(evt)

		evt = NewWAMEvent(ET_WamEventGroupProfilePicture, WithGroupProfilePictureOption(1, 0, 1))
		cache.AddEvent(evt)

		evt = NewWAMEvent(ET_WamEventGroupProfilePicture, WithGroupProfilePictureOption(5, 0, 1))
		cache.AddEvent(evt)

		evt = NewWAMEvent(ET_WamEventGroupProfilePicture, WithGroupProfilePictureOption(10, 0, 1))
		cache.AddEvent(evt)

		evt = NewWAMEvent(ET_WamEventProfilePicUpload, WithProfilePicUploadOption(0, setAvatarSize))
		cache.AddEvent(evt)
	}

	// 返回群会话
	evt = NewWAMEvent(ET_WamEventGroupInfo, nil)
	cache.AddEvent(evt)
}

// LogGroupEditAvatar 设置群头像
func (m *manager) LogGroupEditAvatar(container containerInterface.IAppIocContainer, preHasAvatar bool, setAvatarSize uint32) {
	cache := container.ResolveChannel0EventCache()
	cache2 := container.ResolveChannel2EventCache()
	_ = cache2

	var evt containerInterface.WaEvent

	m.LogNotifyGroupEditAvatar(container, preHasAvatar, setAvatarSize)

	// 返回群会话
	evt = NewWAMEvent(ET_WamEventGroupInfo, nil)
	cache.AddEvent(evt)
}

// LogNotifyGroupEditMembers 群成员变化时记录日志
func (m *manager) LogNotifyGroupEditMembers(container containerInterface.IAppIocContainer) {
	cache := container.ResolveChannel0EventCache()
	cache2 := container.ResolveChannel2EventCache()
	_ = cache2

	var evt containerInterface.WaEvent

	n := rand.Intn(10) + 8
	for i := 0; i != n; i++ {
		evt = NewWAMEvent(ET_WamEventProfilePicDownload, WithProfilePicDownloadOption(PTFull, 404, false, 0))
		cache.AddEvent(evt)
	}
}

// LogGroupEditMembers 群成员变化时记录日志
func (m *manager) LogGroupEditMembers(container containerInterface.IAppIocContainer) {
	cache := container.ResolveChannel0EventCache()
	cache2 := container.ResolveChannel2EventCache()
	_ = cache2

	var evt containerInterface.WaEvent

	m.LogNotifyGroupEditMembers(container)

	// 返回群会话
	evt = NewWAMEvent(ET_WamEventGroupInfo, nil)
	cache.AddEvent(evt)
}

func (m *manager) LogGroupSendText(container containerInterface.IAppIocContainer, stateCount int32, fstMsg bool) {
	m.LogSendText(container, stateCount, fstMsg)
}

func (m *manager) LogGroupSendMedia(container containerInterface.IAppIocContainer, mediaType MediaType, media Media) {
	m.LogSendMedia(container, mediaType, media)
}
