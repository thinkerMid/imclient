package messageResultType

const (
	// SetAvatar 设置头像
	SetAvatar uint8 = iota
	// GetAvatar 获取头像
	GetAvatar

	// SetNickName 设置昵称
	SetNickName
	// GetNickName 获取昵称
	GetNickName

	// SetSignature 设置状态
	SetSignature
	// GetSignature 获取状态
	GetSignature

	// GetQrCode 获取二维码
	GetQrCode
	// SetQrCode 更改二维码
	SetQrCode

	// QueryStatusPrivacyList 查询隐私状态列表
	QueryStatusPrivacyList

	// 联系人相关 START ----------------------------------------------------------------------

	// CheckContact 联系人检测
	CheckContact
	// QueryContact 查询联系人信息
	QueryContact
	// AddContact 添加联系人信息
	AddContact
	// DeleteContact 删除联系人
	DeleteContact
	// BatchAddContact 批量添加联系人信息
	BatchAddContact
	// GetUserDefaultDevice 获取用户默认设备信息
	GetUserDefaultDevice
	// GetUserMultiDevice 获取用户多个设备信息
	GetUserMultiDevice
	// GetUserDeviceList 获取用户设备列表数量
	GetUserDeviceList
	// ContactNicknameUpdate 联系人昵称更新
	ContactNicknameUpdate
	// ContactAvatarUpdate 联系人头像更新
	ContactAvatarUpdate
	// ContactSignatureUpdate 联系人签名更新
	ContactSignatureUpdate
	// CheckStranger 陌生人检测
	CheckStranger
	// QueryStranger 查询陌生人信息
	QueryStranger

	// 私信相关 START ----------------------------------------------------------------------

	// InputChatState 输入状态
	InputChatState
	// SendPrivateChatText 发送私信文本
	SendPrivateChatText
	// SendPrivateChatImage 发送私信图片
	SendPrivateChatImage
	// SendPrivateChatAudio 发送私信语音
	SendPrivateChatAudio
	// SendPrivateChatVideo 发送私信视频
	SendPrivateChatVideo
	// SendPrivateChatVCard 发送名片
	SendPrivateChatVCard
	// SendPrivateChatTemp 发送临时消息
	SendPrivateChatTemp
	// UserLastOnlineNotify 用户最近上线通知
	UserLastOnlineNotify
	// ReceiveMessage 接收消息
	ReceiveMessage
	// SentMessageMarkDelivery 消息送达
	SentMessageMarkDelivery
	// SentMessageMarkRead 消息已读
	SentMessageMarkRead
	// TrustedContactToken 信任联系人
	TrustedContactToken
	// ReceiveMessageMarkRead 接收到的消息已读
	ReceiveMessageMarkRead
	// VCardCheck 名片校验
	VCardCheck

	// 私信状态相关 END ----------------------------------------------------------------------

	// 账号状态相关 START ----------------------------------------------------------------------

	// Online 上线
	Online
	// Offline 离线
	Offline
	// Ban 封禁
	Ban
	// StreamEnd 消息流关闭
	StreamEnd
	// StreamErrorCode 消息流异常码
	StreamErrorCode

	// 账号状态相关 END ------------------------------------------------------------------------

	// 群组 START ----------------------------------------------------------------------

	// CreateGroup 创建群组
	CreateGroup
	// ExitGroup 退出群组
	ExitGroup
	// GroupMemberChange 群组成员变更
	GroupMemberChange
	// GroupChatPermissionChange 聊天权限变更
	GroupChatPermissionChange
	// GroupEditDescPermissionChange 修改描述权限变更
	GroupEditDescPermissionChange
	// ModifyGroupAdmin 变更管理员
	ModifyGroupAdmin
	// ModifyGroupIcon 修改群头像
	ModifyGroupIcon
	// ModifyGroupName 修改群名称
	ModifyGroupName
	// ModifyGroupDesc 修改群描述
	ModifyGroupDesc
	// QueryGroupInfo 查询群组信息
	QueryGroupInfo
	// QueryGroupIcon 查询群组头像
	QueryGroupIcon
	// JoinGroup 加入群组
	JoinGroup
	// LeftGroup 离开群组
	LeftGroup
	// SendGroupChatText 发送私信文本
	SendGroupChatText
	// SendGroupChatImage 发送私信图片
	SendGroupChatImage
	// SendGroupChatAudio 发送私信语音
	SendGroupChatAudio
	// SendGroupChatVideo 发送私信视频
	SendGroupChatVideo
	// SendGroupChatTemp 发送临时消息
	SendGroupChatTemp
	// ReceiveGroupChatMessage 群组聊天消息
	ReceiveGroupChatMessage
	// 群组 END ------------------------------------------------------------------------
)
