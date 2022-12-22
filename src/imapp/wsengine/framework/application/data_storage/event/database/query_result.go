package eventDB

// EventBuffer .
type EventBuffer struct {
	EventLog []byte `gorm:"column:event_log"`
}

// TableName .
func (s *EventBuffer) TableName() string {
	return "event_cache"
}
