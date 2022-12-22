package registrationTokenDB

import "ws/framework/plugin/database/database_tools"

// RegistrationToken .
type RegistrationToken struct {
	databaseTools.ChangeExtension

	JID           string `gorm:"column:jid;primaryKey"`
	RecoveryToken []byte `gorm:"column:recoverToken"`
	BackupToken   []byte `gorm:"column:backupToken"`
	BackupKey     []byte `gorm:"column:backupKey"`    // 两个key要存着的
	BackupKey2    []byte `gorm:"column:backup_key_2"` // 两个key要存着的
	IsUpload      bool   `gorm:"column:isUpload"`
}

// TableName .
func (m *RegistrationToken) TableName() string {
	return "revertoken"
}

// UpdateJID .
func (m *RegistrationToken) UpdateJID(name string) {
	if name == m.JID {
		return
	}

	m.Update("jid", name)
}
