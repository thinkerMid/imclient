package sms

const (
	apiUrl  = "https://api.sms-activate.org:443/stubs/handler_api.php"
	apiKey  = "f2AA1704bec645df0cd455e24630c845"
	service = "api.sms-activate.org"
)

type rentPhoneSettingStatus int

const (
	sentSmsCodeToPhoneType   rentPhoneSettingStatus = 1
	resentSmsCodeToPhoneType rentPhoneSettingStatus = 3
	recoverRentPhoneType     rentPhoneSettingStatus = 6
	discardRentPhoneType     rentPhoneSettingStatus = 8
)

// PriceInfo .
type PriceInfo struct {
	Id    int // 国家ID
	Name  string
	Price float32
	Count float32
}

type Country struct {
	Id   int    `json:"id"`
	Name string `json:"eng"`
}

type Phone struct {
	Id     string `json:"id"`
	Number string `json:"number"`
}
