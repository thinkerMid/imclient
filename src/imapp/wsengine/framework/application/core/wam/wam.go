package wam

import (
	"sync"
	containerInterface "ws/framework/application/container/abstract_interface"
	accountDB "ws/framework/application/data_storage/account/database"
)

/*
* 负责所有账号日志记录
 */

type manager struct {
}

type PageNumber int32

const (
	PageSession PageNumber = iota // 会话页面
	PageProfile                   // 设置信息页面
	PageStatus                    // 状态页面
)

var (
	once sync.Once
	mgr  *manager
)

func LogManager() *manager {
	once.Do(func() {
		mgr = &manager{}
	})
	return mgr
}

func (m *manager) SwitchAppMenu(container containerInterface.IAppIocContainer, to PageNumber) {
	account := container.ResolveAccountService().Context()
	if account.AppPage == int32(to) {
		return
	}

	var evt containerInterface.WaEvent

	cache := container.ResolveChannel0EventCache()
	cache2 := container.ResolveChannel2EventCache()
	_ = cache2

	curr := PageNumber(account.AppPage)
	switch curr {
	case PageSession:
		if to == PageProfile {
			evt = NewWAMEvent(ET_WamEventSettingsClick, nil)
			cache.AddEvent(evt)
		}
	case PageProfile:
	case PageStatus:
	}

	container.ResolveAccountService().ContextExecute(func(act *accountDB.Account) {
		act.SetCurrentAppPage(int32(to))
	})
}
