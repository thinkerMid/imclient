package identityDB

import "ws/framework/plugin/database/database_tools"

// Identity .
type Identity struct {
	databaseTools.ChangeExtension

	JID      string `gorm:"column:ourJid;primaryKey"`
	TheirJID string `gorm:"column:theirJid;primaryKey"`
	Identity []byte `gorm:"column:identity"`
}

// TableName .
func (m *Identity) TableName() string {
	return "identity"
}
