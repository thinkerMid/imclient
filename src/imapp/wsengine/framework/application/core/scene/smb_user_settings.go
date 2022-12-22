package scene

import (
	"ws/framework/application/constant/types"
	containerInterface "ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/action/business"
	"ws/framework/application/core/action/user"
	userSettings "ws/framework/application/core/action/user_settings"
	xmlStream "ws/framework/application/core/action/xml_stream"
	"ws/framework/application/core/monitor"
	"ws/framework/application/core/processor"
)

// SMBUserSettingsScene .
type SMBUserSettingsScene struct {
	UserSettingsScene
}

// NewSMBUserSettingsScene .
func NewSMBUserSettingsScene() SMBUserSettingsScene {
	return SMBUserSettingsScene{}
}

// 打开个人信息
func (u *SMBUserSettingsScene) recomposeAction() []containerInterface.IAction {
	defaultQuery := []containerInterface.IAction{
		&userSettings.QueryAvatarUrl{},
		&business.SMBQueryBusinessCollections{},
		&business.SMBQueryProductCatalog{},
		&business.SMBQueryBusinessProfile{MagicV: "372"},
		&business.SMBQueryLinkedAccount{},
	}

	defaultQuery = append(defaultQuery, u.ActionList...)

	return defaultQuery
}

// Build .
func (u *SMBUserSettingsScene) Build() containerInterface.IProcessor {
	return processor.NewOnceProcessor(
		u.recomposeAction(),
		processor.AliasName("userSettings"),
		processor.AttachMonitor(&monitor.UserSettings{}),
		processor.Priority(processor.PriorityForeground),
	)
}

// ModifyNickName 修改昵称
func (u *SMBUserSettingsScene) ModifyNickName(content string) {
	u.ActionList = append(u.ActionList, &business.SMBModifyVerifiedName{ModifyName: content})
	u.ActionList = append(u.ActionList, &user.Presence{PresenceState: types.PresenceUnavailable})
	u.ActionList = append(u.ActionList, &xmlStream.SendStreamEnd{})
	u.ActionList = append(u.ActionList, &xmlStream.ConnectStream{})
	u.ActionList = append(u.ActionList, &user.Presence{PresenceState: types.PresenceAvailable})
	u.ActionList = append(u.ActionList, &business.SMBQueryVerifiedName{})
	//u.ActionList = append(u.ActionList, &business.SMBSyncCollection{}) // TODO 缺少这个包
}

// ModifySignature 修改签名
func (u *SMBUserSettingsScene) ModifySignature(content string) {
	u.ActionList = append(u.ActionList, &userSettings.ModifySignature{Content: content})
}

// ModifyAvatar 修改头像
func (u *SMBUserSettingsScene) ModifyAvatar(content []byte) {
	u.ActionList = append(u.ActionList, &userSettings.ModifyAvatar{Content: content})
	u.ActionList = append(u.ActionList, &userSettings.QueryAvatarUrl{})
}

// QueryQrCode 查二维码
func (u *SMBUserSettingsScene) QueryQrCode() {
	u.ActionList = append(u.ActionList, &userSettings.QueryQrCode{})
}

// QueryNickname 查昵称(本地)
func (u *SMBUserSettingsScene) QueryNickname() {
	u.ActionList = append(u.ActionList, &userSettings.QueryNickname{})
}

// QuerySignature 查签名(本地)
func (u *SMBUserSettingsScene) QuerySignature() {
	u.ActionList = append(u.ActionList, &userSettings.QuerySignature{})
}

// QueryAvatar 查头像
func (u *SMBUserSettingsScene) QueryAvatar() {
	// recomposeAction 里有查询头像url 不需要再单独执行一次了
}
