package containerInterface

import (
	waProto "ws/framework/application/constant/binary/proto"
	"ws/framework/config"
)

// IIMClient .
type IIMClient interface {
	Connect(optsFn ...config.OptionsFn) (err error)
	EnterScene(scene IScene, resultProcessor LocalResultProcessor)
	AddGlobalResultProcessor(p GlobalResultProcessor)
	EnableAutoMessageMarkRead()
	DisableMessageMarkRead()
	Logout()
	CleanupDataStorage()
	Version() string
	Platform() waProto.UserAgent_UserAgentPlatform
}

// IAnonymousClient .
type IAnonymousClient interface {
	JID() string
	GetSmsCode(optsFn ...config.OptionsFn) (err error)
	SendReceiveSmsCode(smsCode string, optsFn ...config.OptionsFn) (err error)
	CleanupDataStorage()
	Version() string
	Platform() waProto.UserAgent_UserAgentPlatform
}
