package abKeyDB

import (
	"ws/framework/plugin/database/database_tools"
)

// ABKey .
type ABKey struct {
	databaseTools.ChangeExtension

	JID     string `gorm:"column:jid;primaryKey"`
	Content string `gorm:"column:content"`
}

// TableName .
func (s *ABKey) TableName() string {
	return "ab_key"
}

// UpdateContent .
func (s *ABKey) UpdateContent(v string) {
	if v == s.Content {
		return
	}

	s.Content = v
	s.Update("content", s.Content)
}
