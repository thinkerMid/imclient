package msisdnDB

// QueryResult .
type QueryResult struct {
	MCC string `gorm:"column:mcc"`
	MNC string `gorm:"column:mnc"`
}

func (m *QueryResult) TableName() string {
	return "msisdn"
}
