package sessionService

import (
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/data_storage/session/database"
	"ws/framework/plugin/database"
)

var _ containerInterface.ISessionService = &Session{}

// Session
//
//	Deprecated
type Session struct {
	containerInterface.BaseService
}

// CleanupAllData .
func (s *Session) CleanupAllData() {
	_, _ = sessionDB.DeleteByJID(database.MasterDB(), s.JID.User)
}
