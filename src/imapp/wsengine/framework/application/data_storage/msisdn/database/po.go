package msisdnDB

// Msisdn .
type Msisdn struct {
	PhoneNumber string `gorm:"column:phone_number;primaryKey"`
	MCC         string `gorm:"column:mcc"`
	MNC         string `gorm:"column:mnc"`
}

// TableName .
func (m *Msisdn) TableName() string {
	return "msisdn"
}
