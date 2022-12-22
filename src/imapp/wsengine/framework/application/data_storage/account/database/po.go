package accountDB

import (
	"time"
	"ws/framework/application/data_storage/account/constant"
	"ws/framework/plugin/database/database_tools"
)

// Account .
type Account struct {
	databaseTools.ChangeExtension

	JID                    string `gorm:"column:jid;primaryKey"`
	Status                 int16  `gorm:"column:status"`
	Signature              string `gorm:"column:signature"`
	HaveAvatar             bool   `gorm:"column:have_avatar"`
	LogoutTime             int64  `gorm:"column:logout_time"`
	SendChannel0EventTime  int64  `gorm:"column:send_event_0_time"`
	SendChannel0EventCount int32  `gorm:"column:send_event_0_count"`
	SendChannel2EventTime  int64  `gorm:"column:send_event_2_time"`
	SendChannel2EventCount int32  `gorm:"column:send_event_2_count"`
	LoginCount             int64  `gorm:"column:login_count"`
	AppPage                int32  `gorm:"column:app_page"`
	BusinessAccount        bool   `gorm:"column:business_account"`
}

// TableName .
func (a *Account) TableName() string {
	return "account_info"
}

// FirstLogin .
func (a *Account) FirstLogin() bool {
	return a.LoginCount == 0
}

// AddLoginCount .
func (a *Account) AddLoginCount() {
	a.LoginCount++
	a.Update("login_count", a.LoginCount)
}

// SetCurrentAppPage .
func (a *Account) SetCurrentAppPage(page int32) {
	if a.AppPage == page {
		return
	}

	a.AppPage = page
	a.Update("app_page", a.AppPage)
}

// AddChannel0EventCount .
func (a *Account) AddChannel0EventCount() {
	a.SendChannel0EventCount++
	a.Update("send_event_0_count", a.SendChannel0EventCount)
}

// AddChannel2EventCount .
func (a *Account) AddChannel2EventCount() {
	a.SendChannel2EventCount++
	a.Update("send_event_2_count", a.SendChannel2EventCount)
}

// UpdateLogoutTime .
func (a *Account) UpdateLogoutTime() {
	a.LogoutTime = time.Now().Unix()
	a.Update("logout_time", a.LogoutTime)
}

// UpdateSendChannel0Time .
func (a *Account) UpdateSendChannel0Time() {
	a.SendChannel0EventTime = time.Now().Unix()
	a.Update("send_event_0_time", a.SendChannel0EventTime)
}

// UpdateSendChannel2Time .
func (a *Account) UpdateSendChannel2Time() {
	a.SendChannel2EventTime = time.Now().Unix()
	a.Update("send_event_2_time", a.SendChannel2EventTime)
}

// UpdateAccountStatus .
func (a *Account) UpdateAccountStatus(value int16) {
	if a.Status == value {
		return
	}

	a.Status = value
	a.Update("status", a.Status)
}

// UpdateSignature .
func (a *Account) UpdateSignature(value string) {
	if a.Signature == value {
		return
	}

	a.Signature = value
	a.Update("signature", a.Signature)
}

// UpdateHaveAvatar .
func (a *Account) UpdateHaveAvatar(value bool) {
	if a.HaveAvatar == value {
		return
	}

	a.HaveAvatar = value
	a.Update("have_avatar", a.HaveAvatar)
}

// UpdateBusinessAccount .
func (a *Account) UpdateBusinessAccount(value bool) {
	if a.BusinessAccount == value {
		return
	}

	a.BusinessAccount = value
	a.Update("business_account", a.BusinessAccount)
}

// IsRegistered 是否注册 如果没注册状态是-1的
func (a *Account) IsRegistered() bool {
	return a.Status > accountServiceConstant.Unregistered
}

// AvailableStatus .
func (a *Account) AvailableStatus() bool {
	return a.Status != 401 && a.Status != 403
}

// UpdateJID .
func (a *Account) UpdateJID(name string) {
	if name == a.JID {
		return
	}

	a.Update("jid", name)
}
