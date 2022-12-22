package prekeyDB

import "ws/framework/plugin/database/database_tools"

// PreKey .
type PreKey struct {
	databaseTools.ChangeExtension

	JID      string `gorm:"column:jid;primaryKey"`
	KeyId    uint32 `gorm:"column:keyId"`
	KeyBuff  []byte `gorm:"column:keyBuff"`
	IsUpload bool   `gorm:"column:isUpload"`
}

// TableName .
func (m *PreKey) TableName() string {
	return "prekey"
}
