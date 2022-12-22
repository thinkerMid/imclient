package sessionDB

import (
	"ws/framework/plugin/database/database_tools"
)

// Session .
type Session struct {
	databaseTools.ChangeExtension
}

// TableName .
func (m *Session) TableName() string {
	return "session"
}
