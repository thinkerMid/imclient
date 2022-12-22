package eventDB

// EventInfo .
type EventInfo struct {
	JID           string `gorm:"column:jid;primaryKey"`
	SerialNumber  int64  `gorm:"column:serial_number"`
	AutoIncrement int32  `gorm:"column:auto_increment_id"`
	ChannelID     byte   `gorm:"column:channel_id"`
	EventLog      []byte `gorm:"column:event_log"`
}

// TableName .
func (s *EventInfo) TableName() string {
	return "event_cache"
}
