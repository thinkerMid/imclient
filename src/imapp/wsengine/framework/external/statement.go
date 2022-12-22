package external

// ProfileUpdate .
type ProfileUpdate struct {
	JIDNumber string
	Content   string
}

// Contact .
type Contact struct {
	JIDNumber   string
	PhoneNumber string
}

// ChatMessage .
type ChatMessage struct {
	GroupNumber string
	JIDNumber   string

	MessageID string

	MediaUrl      string
	Mimetype      string
	CBCKey        []byte
	CBCIv         []byte
	JpegThumbnail []byte // 预览图
	Seconds       uint32 // 时长

	Conversation string
}

// PrivateChatContactMessage .
type PrivateChatContactMessage struct {
	DisplayName string
	Contents    []string
}

// MessageStateChange .
type MessageStateChange struct {
	MessageID string
	JIDNumber string
}

// GroupInfo .
type GroupInfo struct {
	GroupNumber string
	Description GroupDesc
	Title       GroupTitle
	Members     []GroupMember
	CreateTime  int64
}

// GroupDesc 群描述
type GroupDesc struct {
	Text     string
	EditTime int64
	Editor   string
	EditId   string
}

// GroupTitle 群名称
type GroupTitle struct {
	Text     string
	EditTime int64
	Editor   string
}

type GroupMember struct {
	MemberNumber string
	Right        int32 // 0:普通成员 1:管理员 2:超管
}

// OnlineNotify .
type OnlineNotify struct {
	JIDNumber string
	Time      int64
}
