package containerInterface

import waProto "ws/framework/application/constant/binary/proto"

// WhatsappConfiguration .
type WhatsappConfiguration struct {
	Platform waProto.UserAgent_UserAgentPlatform `validate:"required,eq=1|eq=12"`

	// ---------------注册用的---------------

	HttpUrl                string `validate:"required,min=1"`
	UserAgent              string `validate:"required,min=1"`
	AESPassword            string `validate:"required,min=1"`
	AESCurve25519PublicKey []byte `validate:"required,min=1"`
	BuildHash              string `validate:"required,min=1"`
	CommitHash             string `validate:"required,min=1"`

	// ---------------即时通讯---------------

	TCPAddress           string `validate:"required,min=1"`
	HandshakeHeader      []byte `validate:"required,min=1"`
	EdInfo               []byte `validate:"required,min=1"`
	EdLen                []byte `validate:"required,min=1"`
	NoiseFullPattern     string `validate:"required,min=1"`
	NoiseResumePattern   string `validate:"required,min=1"`
	NoiseCallbackPattern string `validate:"required,min=1"`

	// ---------------渠道2日志---------------

	PrivateStatsURL         string `validate:"required,min=1"`
	PrivateStatsAccessToken string `validate:"required,min=1"`
	PrivateStatsBoundary    string `validate:"required,min=1"`

	// ---------------事件日志序列化---------------
	CommonSerializeCode []int64 `validate:"required,min=1"`

	// 版本代码
	VersionCode   []uint32 `validate:"required,min=1"`
	VersionString string   `validate:"required,min=1"`

	NoticeId        string `validate:"required,min=1"`
	BuildSDKVersion string `validate:"required,min=1"`
}
