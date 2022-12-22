package signedPreKeyDB

import "ws/framework/plugin/database/database_tools"

// SignedPreKey .
type SignedPreKey struct {
	databaseTools.ChangeExtension

	JID     string `gorm:"column:jid;primaryKey"`
	KeyId   uint32 `gorm:"column:keyId"`
	KeyBuff []byte `gorm:"column:keyBuff"`
}

// TableName .
func (m *SignedPreKey) TableName() string {
	return "signedprekey"
}

// UpdateJID .
func (m *SignedPreKey) UpdateJID(name string) {
	if name == m.JID {
		return
	}

	m.Update("jid", name)
}
