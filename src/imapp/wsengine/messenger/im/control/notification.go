package control

import (
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/action/account"
	"ws/framework/application/core/action/event"
	"ws/framework/application/core/action/user"
	"ws/framework/application/core/monitor"
	"ws/framework/application/core/notification"
	accountNotification "ws/framework/application/core/notification/account"
	contactNotification "ws/framework/application/core/notification/contact"
	groupNotification "ws/framework/application/core/notification/group"
	identityNotification "ws/framework/application/core/notification/identity"
	privateChatNotification "ws/framework/application/core/notification/private_chat"
	"ws/framework/application/core/notification/xml_stream"
	"ws/framework/application/core/processor"
)

// initNotificationProcessor 处理服务器主动发送过来的消息
func initNotificationProcessor(channel containerInterface.IMessageChannel) {
	// 登出
	channel.AddMessageProcessor(processor.NewOnceProcessor(
		[]containerInterface.IAction{
			&user.Presence{PresenceState: types.PresenceUnavailable},
		},
		processor.TriggerTag(message.LogoutEvent),
	))

	// 离线
	channel.AddMessageProcessor(processor.NewOnceProcessor(
		[]containerInterface.IAction{
			&accountState.Offline{},
		},
	))

	// 上传上次未上传日志
	channel.AddMessageProcessor(processor.NewOnceProcessor(
		[]containerInterface.IAction{
			&event.UploadChannel0RecordEvent{},
			&event.UploadChannel2RecordEvent{},
		},
		processor.TriggerTag(message.LoginSuccess),
	))

	// 封号
	channel.AddMessageProcessor(processor.NewNotificationProcessor(
		func() []containerInterface.INotification {
			return []containerInterface.INotification{
				accountNotification.Ban{},
			}
		},
		processor.TriggerTag(message.Failure),
	))

	// 初始化登录数据
	channel.AddMessageProcessor(processor.NewNotificationProcessor(
		func() []containerInterface.INotification {
			return []containerInterface.INotification{
				&accountNotification.LoginData{},
			}
		},
		processor.TriggerTag(message.LoginSuccess),
	))

	// 消息流关闭
	channel.AddMessageProcessor(processor.NewNotificationProcessor(
		func() []containerInterface.INotification {
			return []containerInterface.INotification{
				&xmlStreamNotification.ReceiveStreamEnd{},
			}
		},
		processor.TriggerTag(message.StreamEnd),
	))

	// 消息流错误码
	channel.AddMessageProcessor(processor.NewNotificationProcessor(
		func() []containerInterface.INotification {
			return []containerInterface.INotification{
				&xmlStreamNotification.StreamErrorStatus{},
			}
		},
		processor.TriggerTag(message.StreamError),
	))

	// 解析edge_routing里的握手头部内容
	channel.AddMessageProcessor(processor.NewNotificationProcessor(
		func() []containerInterface.INotification {
			return []containerInterface.INotification{
				xmlStreamNotification.RoutingInfo{},
			}
		},
		processor.TriggerTag(message.IB),
	))

	// 响应服务器下发的ping
	channel.AddMessageProcessor(processor.NewNotificationProcessor(
		func() []containerInterface.INotification {
			return []containerInterface.INotification{
				notification.PingPong{},
			}
		},
		processor.TriggerTag(message.IQ),
	))

	// 语音通话
	channel.AddMessageProcessor(processor.NewNotificationProcessor(
		func() []containerInterface.INotification {
			return []containerInterface.INotification{
				notification.Ack{},
			}
		},
		processor.TriggerTag(message.Call),
	))

	// 通知处理
	channel.AddMessageProcessor(processor.NewNotificationProcessor(
		func() []containerInterface.INotification {
			return []containerInterface.INotification{
				notification.Ack{},
				accountNotification.PreKeyNotEnough{},
				identityNotification.DeleteSession{},
				identityNotification.UpdateDeviceID{},
				&contactNotification.AvatarUpdate{},
				contactNotification.SignatureUpdate{},
				&contactNotification.BeAddedOrDeleted{},
				privateChatNotification.PrivacyToken{},
				groupNotification.CreateGroup{},
				groupNotification.LeftGroup{},
				groupNotification.JoinGroup{},
				groupNotification.UpdateIcon{},
				groupNotification.UpdateAdmin{},
				groupNotification.UpdateDesc{},
				groupNotification.UpdateChatPermission{},
				groupNotification.UpdateEditDescPermission{},
			}
		},
		processor.TriggerTag(message.Notification),
		processor.AttachMonitor(&monitor.NotificationMonitor{}),
	))

	// 消息回执
	channel.AddMessageProcessor(processor.NewNotificationProcessor(
		func() []containerInterface.INotification {
			return []containerInterface.INotification{
				privateChatNotification.MessageReceiptAck{},
				privateChatNotification.HandleTrustedContact{},
			}
		},
		processor.TriggerTag(message.MessageReceipt),
	))

	// 收到对方消息-ack
	// 接收消息-解析内容
	channel.AddMessageProcessor(processor.NewNotificationProcessor(
		func() []containerInterface.INotification {
			return []containerInterface.INotification{
				notification.Ack{},
				contactNotification.TakeNicknameFromMessage{},
				privateChatNotification.ReceiveMessage{},
				groupNotification.ReceiveMessage{},
			}
		},
		processor.TriggerTag(message.ReceiveMessage),
	))

	// 最后上线时间
	channel.AddMessageProcessor(processor.NewNotificationProcessor(
		func() []containerInterface.INotification {
			return []containerInterface.INotification{
				contactNotification.UserOnlineNotify{},
			}
		},
		processor.TriggerTag(message.Presence),
	))
}
