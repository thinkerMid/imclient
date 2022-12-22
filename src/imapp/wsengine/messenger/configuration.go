package messenger

import (
	"github.com/go-playground/locales/zh"
	translator "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"strings"
	waProto "ws/framework/application/constant/binary/proto"
	containerInterface "ws/framework/application/container/abstract_interface"
)

// whatsapp 参数配置
var configuration = &containerInterface.WhatsappConfiguration{
	Platform: waProto.UserAgent_IOS,

	// ---------------注册用的---------------
	HttpUrl:   "https://v.whatsapp.net",
	UserAgent: "WhatsApp/%s iOS/%v Device/%v",
	// [WARegistrationURLBuilder verificationCodeRequestURLWithMethod:mcc:mnc:jailbroken:context:]
	// aesDecodeWithPassphrase
	AESPassword: "0a1mLfGUIBVrMKF1RdvLI5lkRBvof6vn0fD2QRSM",
	AESCurve25519PublicKey: []byte{
		0x8e, 0x8c, 0x0f, 0x74, 0xc3, 0xeb, 0xc5, 0xd7,
		0xa6, 0x86, 0x5c, 0x6c, 0x3c, 0x84, 0x38, 0x56,
		0xb0, 0x61, 0x21, 0xcc, 0xe8, 0xea, 0x77, 0x4d,
		0x22, 0xfb, 0x6f, 0x12, 0x25, 0x12, 0x30, 0x2d,
	},
	// WAFoundation WABuildHash
	BuildHash: "4174c0243f5277a5d7720ce842cc4ae6",
	// WABuildCommitHash
	CommitHash: "caf9f5f9da9",

	// ---------------即时通讯---------------
	TCPAddress:           "g.whatsapp.net:5222",
	HandshakeHeader:      []byte{'W', 'A', 5, 2},
	EdInfo:               []byte{'E', 'D', 0, 1},
	EdLen:                []byte{0, 0, 4},
	NoiseFullPattern:     "Noise_XX_25519_AESGCM_SHA256\x00\x00\x00\x00",
	NoiseResumePattern:   "Noise_IK_25519_AESGCM_SHA256\x00\x00\x00\x00",
	NoiseCallbackPattern: "Noise_XXfallback_25519_AESGCM_SHA256",

	// ---------------渠道2日志---------------
	// [WAFieldStats privateStatsUploadRequestWithDataBuffer:credentialToken:] 下面这三个参数都在这里
	PrivateStatsURL:         "https://dit.whatsapp.net/deidentified_telemetry",
	PrivateStatsAccessToken: "245118376424571|3e7d275052f1522bf3200afcf53841a7",
	PrivateStatsBoundary:    "pfBvf7MwXQtqX5egJfP0Lg==",

	// ---------------事件日志序列化---------------
	CommonSerializeCode: []int64{
		0x3, 0x5, 0xB, 0xD, 0xF, 0x11, 0x15, 0x17, 0x2F, 0x69,
		0x183, 0x28F, 0x2B1, 0x679, 0x67B, 0x827, 0x85D, 0xA39,
		0xAEB, 0xD69, 0x1179, 0x1199, 0x13A5, 0x1775, 0x186B, 0x1AB1, //这个0x186B 是动态插进去的
		0x1C59, 0x1CA7, 0x1E91, 0x2479,
	},

	// ---------------版本代码---------------
	// WAFoundation WABuildVersion
	VersionCode:   []uint32{2, 0x16, 0x18, 0x51},
	VersionString: "2.22.24.81",

	// ---------------未知---------------
	NoticeId: "20210210", // 大版本更新可能会变的，需要留意一下
	// WAFoundation WABuildSDKVersion
	BuildSDKVersion: "16.0",
}

// ConfigurationVerify
func init() {
	uni := translator.New(zh.New())
	trans, _ := uni.GetTranslator("zh")

	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		return fld.Name
	})

	if err := zhTranslations.RegisterDefaultTranslations(validate, trans); err != nil {
		panic(err.Error())
	}

	if err := validate.Struct(configuration); err != nil {
		stringBuffer := strings.Builder{}

		for _, err := range err.(validator.ValidationErrors) {
			stringBuffer.WriteString(err.Translate(trans))
			stringBuffer.WriteString(";")
		}

		panic(stringBuffer.String())
	}
}
