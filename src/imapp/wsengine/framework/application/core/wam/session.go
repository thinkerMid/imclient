package wam

import (
	containerInterface "ws/framework/application/container/abstract_interface"
	. "ws/framework/application/core/wam/events"
)

func (m *manager) LogSessionNew(container containerInterface.IAppIocContainer) {
	cache := container.ResolveChannel0EventCache()
	cache2 := container.ResolveChannel2EventCache()
	_ = cache2

	var evt containerInterface.WaEvent

	evt = NewWAMEvent(ET_WamEventIphoneContactListStartNewChat, nil)
	cache.AddEvent(evt)

	// 反复创建同一会话则不发
	//evt = NewWAMEvent(ET_WamEventPrekeysFetch, nil)
	//cache.AddEvent(evt)

	//evt3 := NewWAMEvent(ET_WamEventPrekeysDepletion)
	//cache.AddLogEvent(evt3.GetChannel(), evt3.Pack())

	// 渠道二
	evt = NewWAMEvent(ET_WamEventChatComposerAction, WithChatComposerActionOption(TargetText))
	cache2.AddEvent(evt)
}
