package control

import (
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	accountState "ws/framework/application/core/action/account"
	"ws/framework/application/core/action/common"
	"ws/framework/application/core/action/contact/compose"
	"ws/framework/application/core/action/other"
	"ws/framework/application/core/action/user"
	"ws/framework/application/core/action/user_settings"
	xmlStream "ws/framework/application/core/action/xml_stream"
	"ws/framework/application/core/notification/account"
	"ws/framework/application/core/processor"
	"ws/messenger/im/action/log"
)

// initFirstLoginProcessor 首次登录
func initFirstLoginProcessor(container containerInterface.IAppIocContainer) {
	jid := container.ResolveJID()
	channel := container.ResolveMessageChannel()

	// 注册要发的包
	channel.AddMessageProcessor(processor.NewOnceIgnoreErrorProcessor(
		[]containerInterface.IAction{
			&user.QueryCertificatePEM{},
			&user.Presence{PresenceState: types.PresenceUnavailable},
			&xmlStream.SendStreamEnd{},
			&xmlStream.ConnectStream{},
			&log.ClientLog{CurrentScreen: "profile_photo", PreviousScreen: "no_backup_found", ActionTaken: "skip"},
			&user.UploadPreKeyToServer{Init: true},
			&user.XMPPPingPongSetting{},
			&user.UploadPushConfig{},
			&user.Presence{PresenceState: types.PresenceUnavailable},
			&user.UploadRecoveryToken{},
			&user.QueryServerProperties{},
			&user.QueryABProperties{},
			&user.QueryParticipatingGroups{},
			&user.QueryBroadcastLists{},
			&common.QueryAvatarUrl{UserID: jid.User},
			&common.QueryUserStatus{UserID: jid.User},
			&userSettings.ModifySignature{Content: "Hey there! I am using WhatsApp."},
			&user.QueryTwoFactorAuthentication{},
			&user.QueryStatusPrivacyList{},
			&user.QueryBlockReasonsList{},
			&user.QueryPrivacySetting{},
			&user.DisappearingMode{},
			&user.DisappearingMode{},
			&user.DisappearingMode{},
			&user.QueryAcceptSetting{},
			&user.QueryBlockList{},
			&user.Presence{PresenceState: types.PresenceAvailable},
			&log.ClientLog{CurrentScreen: "home", PreviousScreen: "profile_photo", ActionTaken: "continue"},
			&user.Presence{PresenceState: types.PresenceAvailable},
			&user.UploadBackupKey{UploadBackupKeyIndexOne: true},
			&user.UploadBackupKey{UploadBackupKeyIndexOne: false},
			&user.CleanDirtyType{Type: "groups"},
			&other.RemoveCompanionDevice{},
			&user.CleanDirtyType{Type: "groups"},
			&user.CleanDirtyType{Type: "account_sync"},
			&user.QueryMMSEndPoints{},
			&user.DeleteAllData{},
			&accountState.LoginSuccess{FirstLogin: true},
			&accountState.Online{},
			&contactComposeAction.UploadLocalAddressBook{Min: 1, Max: 100},
		},
		processor.AliasName("firstLogin"),
		processor.TriggerTag(message.LoginSuccess),
		processor.Priority(processor.PriorityForeground),
	))
}

// 正常登录
func initLoginProcessor(container containerInterface.IAppIocContainer) {
	channel := container.ResolveMessageChannel()

	// 正常登录发的包
	channel.AddMessageProcessor(processor.NewOnceIgnoreErrorProcessor(
		[]containerInterface.IAction{
			&user.Presence{PresenceState: types.PresenceAvailable},
			&user.QueryMMSEndPoints{},
			&user.QueryPrivacySetting{},
			&user.DisappearingMode{},
			&user.UploadPushConfig{},
			&accountState.LoginSuccess{},
			&accountState.Online{},
		},
		processor.AliasName("defaultLogin"),
		processor.TriggerTag(message.LoginSuccess),
	))

	// 账号同步通知
	channel.AddMessageProcessor(processor.NewNotificationProcessor(
		func() []containerInterface.INotification {
			return []containerInterface.INotification{
				accountNotification.AccountSync{},
			}
		},
		processor.TriggerTag(message.IB),
	))

	// 已经登录过一次就开启两步验证设置
	if container.ResolveAccountService().Context().LoginCount == 1 {
		channel.AddMessageProcessor(processor.NewTimerProcessor(
			func() []containerInterface.IAction {
				return []containerInterface.IAction{
					&user.ModifyTwoFactorAuthentication{Code: "258000"},
				}
			},
			processor.Interval(1),
			processor.IntervalLoop(false),
		))
	}
}
