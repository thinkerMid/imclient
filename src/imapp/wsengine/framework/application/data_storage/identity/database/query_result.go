package identityDB

type QueryResult struct {
	Identity []byte `gorm:"column:identity"`
}

// TableName .
func (m *QueryResult) TableName() string {
	return "identity"
}
