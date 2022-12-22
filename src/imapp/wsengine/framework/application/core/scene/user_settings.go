package scene

import (
	"ws/framework/application/container/abstract_interface"
	userSettings "ws/framework/application/core/action/user_settings"
	"ws/framework/application/core/monitor"
	"ws/framework/application/core/processor"
)

// UserSettingsScene .
type UserSettingsScene struct {
	ActionList []containerInterface.IAction
}

// NewUserSettingsScene .
func NewUserSettingsScene() UserSettingsScene {
	return UserSettingsScene{}
}

// Build .
func (u *UserSettingsScene) Build() containerInterface.IProcessor {
	return processor.NewOnceIgnoreErrorProcessor(
		u.ActionList,
		processor.AliasName("userSettings"),
		processor.AttachMonitor(&monitor.UserSettings{}),
	)
}

// ModifyNickName 修改昵称
func (u *UserSettingsScene) ModifyNickName(content string) {
	u.ActionList = append(u.ActionList, &userSettings.QueryAvatarPreview{})
	u.ActionList = append(u.ActionList, &userSettings.ModifyNickName{
		DstName: content,
	})
}

// ModifySignature 修改签名
func (u *UserSettingsScene) ModifySignature(content string) {
	u.ActionList = append(u.ActionList, &userSettings.QueryAvatarPreview{})
	u.ActionList = append(u.ActionList, &userSettings.ModifySignature{Content: content})
}

// ModifyAvatar 修改头像
func (u *UserSettingsScene) ModifyAvatar(content []byte) {
	u.ActionList = append(u.ActionList, &userSettings.QueryAvatarPreview{})
	u.ActionList = append(u.ActionList, &userSettings.ModifyAvatar{Content: content})
	u.ActionList = append(u.ActionList, &userSettings.QueryAvatarUrl{})
}

// QueryQrCode 查二维码
func (u *UserSettingsScene) QueryQrCode() {
	u.ActionList = append(u.ActionList, &userSettings.QueryQrCode{})
}

// QueryNickname 查昵称(本地)
func (u *UserSettingsScene) QueryNickname() {
	u.ActionList = append(u.ActionList, &userSettings.QueryNickname{})
}

// QuerySignature 查签名(本地)
func (u *UserSettingsScene) QuerySignature() {
	u.ActionList = append(u.ActionList, &userSettings.QuerySignature{})
}

// QueryAvatar 查头像(如果本地有)
func (u *UserSettingsScene) QueryAvatar() {
	u.ActionList = append(u.ActionList, &userSettings.QueryAvatarUrl{})
}
