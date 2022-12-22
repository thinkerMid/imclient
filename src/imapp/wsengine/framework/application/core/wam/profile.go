package wam

import (
	containerInterface "ws/framework/application/container/abstract_interface"
	. "ws/framework/application/core/wam/events"
)

func (m *manager) LogProfileView(container containerInterface.IAppIocContainer) {
	cache := container.ResolveChannel0EventCache()
	cache2 := container.ResolveChannel2EventCache()
	_ = cache2

	account := container.ResolveAccountService().Context()

	var evt containerInterface.WaEvent

	if !account.HaveAvatar {
		evt = NewWAMEvent(ET_WamEventProfilePicDownload, WithProfilePicDownloadOption(PTMine, 404, true, 0))
		cache.AddEvent(evt)

		evt = NewWAMEvent(ET_WamEventProfilePicDownload, WithProfilePicDownloadOption(PTThumbnail, 404, false, 0))
		cache.AddEvent(evt)
	}
}

// LogProfileSetAvatar 修改头像
// @hasAvatar: 当前是否有头像
// @picSize: 要设置的头像大小
func (m *manager) LogProfileSetAvatar(container containerInterface.IAppIocContainer, avatarSize uint32) {
	cache := container.ResolveChannel0EventCache()
	cache2 := container.ResolveChannel2EventCache()
	_ = cache2

	var evt containerInterface.WaEvent

	if avatarSize > 0 {
		// 设置/修改头像
		evt = NewWAMEvent(ET_WamEventProfilePicUpload, WithProfilePicUploadOption(0, avatarSize))
		cache.AddEvent(evt)
	} else if avatarSize == 0 {
		// 重置头像
		evt = NewWAMEvent(ET_WamEventProfilePicDownload, WithProfilePicDownloadOption(PTThumbnail, 404, false, 0))
		cache.AddEvent(evt)

		evt = NewWAMEvent(ET_WamEventProfilePicDownload, WithProfilePicDownloadOption(PTFull, 404, false, 0))
		cache.AddEvent(evt)
	}
}
