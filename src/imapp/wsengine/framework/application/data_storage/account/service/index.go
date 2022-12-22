package accountService

import (
	"time"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/data_storage/account/constant"
	accountDB "ws/framework/application/data_storage/account/database"
	"ws/framework/plugin/database"
	"ws/framework/plugin/database/database_tools"
)

var _ containerInterface.IAccountService = &Account{}

// Account .
type Account struct {
	containerInterface.BaseService
	// 比较经常用 做个缓存
	context *accountDB.Account
}

// Context .
func (s *Account) Context() *accountDB.Account {
	if s.context != nil {
		return s.context
	}

	acc := accountDB.Account{JID: s.JID.User}

	if err := databaseTools.FindByPrimaryKey(database.MasterDB(), &acc); err != nil {
		s.Logger.Error(err)
		return nil
	}

	s.context = &acc

	return s.context
}

// Create .
func (s *Account) Create() (*accountDB.Account, error) {
	// -1 标记未未注册 and 设置一个登出时间
	acc := accountDB.Account{JID: s.JID.User, Status: accountServiceConstant.Unregistered, LogoutTime: time.Now().Unix()}

	_, err := databaseTools.Create(database.MasterDB(), &acc)

	if err != nil {
		s.Logger.Error("create context ", err)
	}

	s.context = &acc

	return s.context, err
}

// Import .
func (s *Account) Import(acc *accountDB.Account) error {
	_, err := databaseTools.Create(database.MasterDB(), &acc)
	if err != nil {
		s.Logger.Error(err)
		return err
	}

	s.context = acc

	return nil
}

// NeedUploadRecordChannel0Event .
func (s *Account) NeedUploadRecordChannel0Event() bool {
	logoutTime := time.Unix(s.Context().SendChannel0EventTime, 0)
	now := time.Now()

	// 毫秒
	ms := now.Sub(logoutTime)

	if ms < time.Minute*5 {
		return false
	}

	return true
}

// NeedUploadRecordChannel2Event .
func (s *Account) NeedUploadRecordChannel2Event() bool {
	logoutTime := time.Unix(s.Context().SendChannel2EventTime, 0)
	now := time.Now()

	// 毫秒
	ms := now.Sub(logoutTime)

	if ms < time.Minute*10 {
		return false
	}

	return true
}

// ContextExecute .
func (s *Account) ContextExecute(f func(*accountDB.Account)) {
	f(s.Context())

	_, err := databaseTools.Save(database.MasterDB(), s.Context())
	if err != nil {
		s.Logger.Errorf("update context error: %v", err)
	}
}

// OnJIDChangeWhenRegisterSuccess .
func (s *Account) OnJIDChangeWhenRegisterSuccess(newJID types.JID) {
	acc := accountDB.Account{JID: s.JID.User}
	acc.UpdateJID(newJID.User)

	_, err := databaseTools.Save(database.MasterDB(), &acc)
	if err != nil {
		s.Logger.Error(err)
	}
}

// CleanupAllData .
func (s *Account) CleanupAllData() {
	acc := accountDB.Account{JID: s.JID.User}

	_, _ = databaseTools.DeleteByPrimaryKey(database.MasterDB(), &acc)
}
