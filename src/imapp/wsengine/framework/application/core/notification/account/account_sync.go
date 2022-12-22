package accountNotification

import (
	containerInterface "ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/action/common"
	"ws/framework/application/core/action/user"
	"ws/framework/application/core/processor"
)

// AccountSync .
type AccountSync struct{}

// Receive .
func (a AccountSync) Receive(context containerInterface.IMessageContext) (err error) {
	node, ok := context.Message().GetOptionalChildByTag("dirty")
	if !ok {
		return
	}

	if dirtyType, ok := node.Attrs["type"]; !ok || dirtyType != "account_sync" {
		return
	}

	jid := context.ResolveJID()

	context.AddMessageProcessor(processor.NewOnceIgnoreErrorProcessor(
		[]containerInterface.IAction{
			&user.UploadRecoveryToken{},
			&common.QueryAvatarUrl{UserID: jid.User},
			&user.QueryBlockList{},
			&user.DisappearingMode{},
			&common.QueryUserStatus{UserID: jid.User},
			&user.FreshNotice{},
			&common.QueryUserDeviceListLite{UserIDs: []string{jid.User}},
			&user.QueryPrivacySetting{},
			&common.QueryAvatarPreview{UserID: jid.User},
			&user.CleanDirtyType{Type: "account_sync"},
		},
	))

	return
}
