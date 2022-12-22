package aesKeyDB

type QueryResult struct {
	AesKey []byte `gorm:"column:aes_key"`
	PubKey []byte `gorm:"column:pub_key"`
}

// TableName .
func (s *QueryResult) TableName() string {
	return "aes_key"
}
