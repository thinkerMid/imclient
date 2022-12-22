package control

import (
	"ws/business/im/action/log"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	accountState "ws/framework/application/core/action/account"
	"ws/framework/application/core/action/business"
	"ws/framework/application/core/action/common"
	"ws/framework/application/core/action/other"
	"ws/framework/application/core/action/user"
	userSettings "ws/framework/application/core/action/user_settings"
	xmlStream "ws/framework/application/core/action/xml_stream"
	accountNotification "ws/framework/application/core/notification/account"
	"ws/framework/application/core/processor"
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
			&xmlStream.SendStreamEnd{}, // <xmlstreamend/>
			&xmlStream.ConnectStream{},
			&log.ClientLog{CurrentScreen: "profile_photo", PreviousScreen: "no_backup_found", ActionTaken: "skip"},
			&log.SMBClientLog{Step: 13, Sequence: 5}, // http smb log
			&user.Presence{PresenceState: types.PresenceUnavailable},
			&user.SMBUploadPreKeyToServer{Init: true},                                     // -3
			&user.XMPPPingPongSetting{},                                                   // -4
			&user.SMBUploadPushConfig{},                                                   // -5
			&user.QueryServerProperties{},                                                 // -6
			&user.QueryParticipatingGroups{},                                              // -7
			&user.QueryBroadcastLists{},                                                   // -8
			&common.QueryAvatarUrl{UserID: jid.User},                                      // -9
			&common.QueryUserStatus{UserID: jid.User},                                     // -10
			&userSettings.ModifySignature{Content: "Hello. I'm using WhatsApp Business."}, // -11
			&business.SMBQueryBusinessProfile{MagicV: "116"},                              // -12
			&business.SMBQueryLinkedAccount{},                                             // -13
			&user.QueryTwoFactorAuthentication{},                                          // -14
			&user.QueryStatusPrivacyList{},                                                // -15
			&user.QueryBlockReasonsList{},                                                 // -16
			&user.QueryPrivacySetting{},                                                   // -17
			&user.DisappearingMode{},                                                      // -18
			&user.QueryABProperties{},                                                     // -19
			&other.RemoveCompanionDevice{},                                                // -20
			&user.CleanDirtyType{Type: "groups"},                                          // -21
			&user.CleanDirtyType{Type: "groups"},                                          // -22
			&user.QueryBlockList{},                                                        // -23
			// 24 不见了
			&user.UploadRecoveryToken{},              // -25
			&user.DeleteAllData{},                    // -26
			&user.QueryMMSEndPoints{},                // -27
			&business.SMBQueryBusinessCategory{},     // -28
			&log.SMBClientLog{Step: 14, Sequence: 6}, // http smb log
			&business.SMBModifyBusinessCategory{CategoryID: "133436743388217", CategoryName: "艺术与娱乐"}, // -29 TODO 这个两设定值可能需要从查询分类里获取再设置
			&business.SMBQueryBusinessProfile{MagicV: "372"},                                          // -30
			&business.SMBModifyVerifiedName{},                                                         // -31
			&xmlStream.SendStreamEnd{},                                                                // <xmlstreamend/>
			&xmlStream.ConnectStream{},
			&log.SMBClientLog{Step: 19, Sequence: 7}, // http smb log
			&user.Presence{PresenceState: types.PresenceUnavailable},
			&common.QueryAvatarPreview{UserID: jid.User}, // -32
			// 33 不见了
			&business.SMBQueryProductCatalog{ConsumerVisibleOnly: true},  //- 34
			&user.UploadPushConfig{},                                     // -35
			&user.DisappearingMode{},                                     // -36
			&user.QueryBlockList{},                                       // -37
			&user.QueryPrivacySetting{},                                  // -38
			&user.DisappearingMode{},                                     // -39
			&common.QueryAvatarUrl{UserID: jid.User},                     // -40
			&common.QueryUserStatus{UserID: jid.User},                    // -41
			&business.SMBQueryVerifiedName{},                             // -42
			&common.QueryUserDeviceListLite{UserIDs: []string{jid.User}}, // -43
			&common.QueryAvatarPreview{UserID: jid.User},                 // -44
			&user.CleanDirtyType{Type: "account_sync"},                   // -45
			&log.ClientLog{CurrentScreen: "home", PreviousScreen: "profile_photo", ActionTaken: "continue"},
			//&contactComposeAction.UploadLocalAddressBook{Min: 1, Max: 100}, // -46 TODO 这个操作封号率很高
			&business.SMBQueryLinkedAccount{}, // -47
			&user.Presence{PresenceState: types.PresenceAvailable},
			&accountState.LoginSuccess{FirstLogin: true},
			&accountState.Online{},
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
			&business.SMBQueryProductCatalog{ConsumerVisibleOnly: true},
			&user.UploadPushConfig{},
			&user.UploadRecoveryToken{}, // TODO 这个包不是一直发的 有个规律 目前未知
			&user.QueryMMSEndPoints{},
			&user.QueryPrivacySetting{},
			&user.DisappearingMode{},
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
