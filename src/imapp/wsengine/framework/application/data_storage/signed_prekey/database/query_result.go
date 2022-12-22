package signedPreKeyDB

type QueryResult struct {
	KeyBuff []byte `gorm:"column:keyBuff"`
}

// TableName .
func (m *QueryResult) TableName() string {
	return "signedprekey"
}
