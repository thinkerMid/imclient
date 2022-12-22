package abKeyDB

type QueryResult struct {
	Content string `gorm:"column:content"`
}

// TableName .
func (s *QueryResult) TableName() string {
	return "ab_key"
}
