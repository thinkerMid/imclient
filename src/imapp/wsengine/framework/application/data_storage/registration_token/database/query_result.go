package registrationTokenDB

type QueryResult struct {
	RecoveryToken []byte `gorm:"column:recoverToken"`
	BackupToken   []byte `gorm:"column:backupToken"`
	BackupKey     []byte `gorm:"column:backupKey"`    // 两个key要存着的
	BackupKey2    []byte `gorm:"column:backup_key_2"` // 两个key要存着的
}

// TableName .
func (m *QueryResult) TableName() string {
	return "revertoken"
}
