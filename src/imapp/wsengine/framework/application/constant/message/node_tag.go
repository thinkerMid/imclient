package message

// 内部定义的，与whatsapp无关
const (
	// TickClockEvent 滴答（秒时间）
	TickClockEvent = "TickClockEvent"
	// ReconnectEvent 重连结果
	ReconnectEvent = "ReconnectEvent"
	// LogoutEvent 登出事件
	LogoutEvent = "LogoutEvent"
)

// whatsapp xmpp tag
const (
	// LoginSuccess 登录成功
	LoginSuccess string = "success"

	// MessageReceipt 消息回执
	MessageReceipt string = "receipt"

	// ReceiveMessage 接收消息
	ReceiveMessage string = "message"

	// StreamError 通讯出现异常
	StreamError string = "stream:error"

	// StreamEnd 通讯被正式关闭
	StreamEnd string = "xmlstreamend"

	// Notification 通知
	Notification string = "notification"

	// Failure 账号错误状态
	Failure string = "failure"

	// ChatState 聊天状态
	ChatState string = "chatstate"

	// Presence 最后上线在线通知
	Presence string = "presence"

	// IB 服务器主动通知
	IB string = "ib"

	// Call 语音通话
	Call string = "call"

	// IQ 服务器消息请求
	IQ string = "iq"
)
