package routingInfoDB

import (
	"ws/framework/plugin/database/database_tools"
	functionTools "ws/framework/utils/function_tools"
)

// RoutingInfo .
type RoutingInfo struct {
	databaseTools.ChangeExtension

	JID     string `gorm:"column:jid;primaryKey"`
	Content []byte `gorm:"column:content"`
}

// TableName .
func (s *RoutingInfo) TableName() string {
	return "routing_info"
}

// UpdateContent .
func (s *RoutingInfo) UpdateContent(b []byte) {
	if functionTools.SliceEqual(s.Content, b) {
		return
	}

	s.Content = b
	s.Update("content", s.Content)
}
