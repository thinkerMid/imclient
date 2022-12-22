package aesKeyDB

import "ws/framework/plugin/database/database_tools"

// AESKey .
type AESKey struct {
	databaseTools.ChangeExtension

	JID    string `gorm:"column:jid;primaryKey"`
	AesKey []byte `gorm:"column:aes_key"`
	PubKey []byte `gorm:"column:pub_key"`
	PriKey []byte `gorm:"column:pri_key"`
}

// TableName .
func (s *AESKey) TableName() string {
	return "aes_key"
}

// UpdateJID .
func (s *AESKey) UpdateJID(name string) {
	if name == s.JID {
		return
	}

	s.Update("jid", name)
}
