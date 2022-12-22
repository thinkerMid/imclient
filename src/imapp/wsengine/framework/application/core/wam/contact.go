package wam

import (
	containerInterface "ws/framework/application/container/abstract_interface"
	. "ws/framework/application/core/wam/events"
	"ws/framework/utils"
)

func randSessionId() float64 {
	buff := utils.RandBytes(8)
	val := utils.BigEndianBytesToInt(buff)

	return float64(val & 0x1FFFFFFFFFFFFF)
}

func (m *manager) LogContactAdd(container containerInterface.IAppIocContainer, success, contactHasAvatar bool) {
	cache := container.ResolveChannel0EventCache()
	cache2 := container.ResolveChannel2EventCache()
	_ = cache2

	var evt containerInterface.WaEvent
	sessionId := randSessionId()

	// 1.打开添加联系人面板
	evt = NewWAMEvent(ET_WamEventIphoneAddContactEvent, WithIphoneAddContactEvent(ContactOpen, sessionId, OpenContactAdd))
	cache.AddEvent(evt)
	// 2.关闭面板
	act := ContactClose
	if !success {
		act = ContactCancel
	}
	evt = NewWAMEvent(ET_WamEventIphoneAddContactEvent, WithIphoneAddContactEvent(act, sessionId, OpenContactAdd))
	cache.AddEvent(evt)

	if success {
		if !contactHasAvatar {
			evt = NewWAMEvent(ET_WamEventProfilePicDownload, WithProfilePicDownloadOption(PTFull, 404, false, 0))
			cache.AddEvent(evt)
		}
	}
}

func (m *manager) LogDeleteContact(container containerInterface.IAppIocContainer, hasSession, contactHasAvatar bool) {
	cache := container.ResolveChannel0EventCache()
	cache2 := container.ResolveChannel2EventCache()
	_ = cache2

	var evt containerInterface.WaEvent

	// session不存在则创建
	if !hasSession {
		m.LogSessionNew(container)
	}

	// 点击session头像查看信息
	if !contactHasAvatar {
		evt = NewWAMEvent(ET_WamEventProfilePicDownload, WithProfilePicDownloadOption(PTFull, 404, false, 0))
		cache.AddEvent(evt)
	}

	// 点编辑删除联系人
	evt = NewWAMEvent(ET_WamEventIphoneAddContactEvent, WithIphoneAddContactEvent(ContactOpen, randSessionId(), OpenContactView))
	cache.AddEvent(evt)

	// 删除后返回信息面板
	if !contactHasAvatar {
		evt = NewWAMEvent(ET_WamEventProfilePicDownload, WithProfilePicDownloadOption(PTFull, 404, false, 0))
		cache.AddEvent(evt)
	}
}

//func (m *manager) LogContactView(container containerInterface.IAppIocContainer, hasSession, contactHasAvatar bool) {
//	cache := container.ResolveChannel0EventCache()
//	cache2 := container.ResolveChannel2EventCache()
//	_ = cache2
//
//	var evt  containerInterface.WaEvent
//
//	// session不存在则创建
//	if !hasSession {
//		m.LogSessionNew(container)
//	}
//
//	if !contactHasAvatar {
//		evt = NewWAMEvent(ET_WamEventProfilePicDownload, WithProfilePicDownloadOption(PTFull, 404, false, 0))
//		cache.AddEvent(evt)
//	}
//}

func (m *manager) LogNotifyContactAvatar(container containerInterface.IAppIocContainer, avatarSize int32) {
	cache := container.ResolveChannel0EventCache()
	cache2 := container.ResolveChannel2EventCache()
	_ = cache2

	var evt containerInterface.WaEvent
	// 每次记录日志次数不确定。

	// delete
	if avatarSize == 0 {
		evt = NewWAMEvent(ET_WamEventProfilePicDownload, WithProfilePicDownloadOption(PTThumbnail, 404, true, 0))
		cache.AddEvent(evt)

		evt = NewWAMEvent(ET_WamEventProfilePicDownload, WithProfilePicDownloadOption(PTFull, 404, false, 0))
		cache.AddEvent(evt)
		//
		//evt = NewWAMEvent(ET_WamEventProfilePicDownload, WithProfilePicDownloadOption(PTFull, 404, false, 0))
		//cache.AddEvent(evt)
	} else {
		// set
		evt = NewWAMEvent(ET_WamEventProfilePicDownload, WithProfilePicDownloadOption(PTThumbnail, 0, true, avatarSize))
		cache.AddEvent(evt)

		evt = NewWAMEvent(ET_WamEventProfilePicDownload, WithProfilePicDownloadOption(PTFull, 0, false, avatarSize))
		cache.AddEvent(evt)
		//
		//evt = NewWAMEvent(ET_WamEventProfilePicDownload, WithProfilePicDownloadOption(PTFull, 0, false, avatarSize))
		//cache.AddEvent(evt)
	}
}

func (m *manager) LogNotifyContactAddedOrDeleted(container containerInterface.IAppIocContainer) {
	cache := container.ResolveChannel0EventCache()
	cache2 := container.ResolveChannel2EventCache()
	_ = cache2

	var evt containerInterface.WaEvent

	evt = NewWAMEvent(ET_WamEventProfilePicDownload, WithProfilePicDownloadOption(PTThumbnail, 0, false, 0))
	cache.AddEvent(evt)
}
