package accountDB

// AccountJID .
type AccountJID struct {
	JID string `gorm:"column:jid;primaryKey"`
}

// TableName .
func (a *AccountJID) TableName() string {
	return "account_info"
}
