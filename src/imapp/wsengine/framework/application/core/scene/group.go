package scene

import (
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/action/common"
	"ws/framework/application/core/action/group"
	groupComposeAction "ws/framework/application/core/action/group/compose"
	"ws/framework/application/core/monitor"
	"ws/framework/application/core/processor"
)

// NewCreateGroupScene .
func NewCreateGroupScene() GroupScene {
	return GroupScene{}
}

// NewGroupScene .
func NewGroupScene(groupNumber string) GroupScene {
	return GroupScene{GroupNumber: groupNumber}
}

// GroupScene .
type GroupScene struct {
	GroupNumber string
	ActionList  [][]containerInterface.IAction
}

func (g *GroupScene) applyQueryGroupIcon() *groupComposeAction.CheckAndQueryIcon {
	return &groupComposeAction.CheckAndQueryIcon{
		GroupID: g.GroupNumber,
	}
}

// Create 创建群组
func (g *GroupScene) Create(name string, icon []byte, joinUserIDs []string) {
	g.ActionList = append(g.ActionList, []containerInterface.IAction{
		&groupComposeAction.CreateGroup{
			GroupName:   name,
			Icon:        icon,
			JoinUserIDs: joinUserIDs,
		},
	})
}

// QueryIcon 查询图标
func (g *GroupScene) QueryIcon() {
	g.ActionList = append(g.ActionList, []containerInterface.IAction{
		&group.QueryInfo{GroupID: g.GroupNumber},
		&group.QueryIcon{GroupID: g.GroupNumber},
	})
}

// ModifyIcon 修改图标
func (g *GroupScene) ModifyIcon(icon []byte) {
	g.ActionList = append(g.ActionList, []containerInterface.IAction{
		&group.QueryInfo{GroupID: g.GroupNumber},
		&group.ModifyIcon{GroupID: g.GroupNumber, Icon: icon},
	})
}

// ModifyDescription 修改描述
func (g *GroupScene) ModifyDescription(description string) {
	g.ActionList = append(g.ActionList, []containerInterface.IAction{
		&group.QueryInfo{GroupID: g.GroupNumber},
		&group.ModifyDescription{
			GroupID:     g.GroupNumber,
			DescContent: description,
		},
		g.applyQueryGroupIcon(),
	})
}

// ModifyName 修改名称
func (g *GroupScene) ModifyName(name string) {
	g.ActionList = append(g.ActionList, []containerInterface.IAction{
		&group.QueryInfo{GroupID: g.GroupNumber},
		&group.ModifyName{
			GroupID: g.GroupNumber,
			Name:    name,
		},
		g.applyQueryGroupIcon(),
	})
}

// QueryInfo 查询群组信息
func (g *GroupScene) QueryInfo() {
	g.ActionList = append(g.ActionList, []containerInterface.IAction{
		&group.QueryInfo{GroupID: g.GroupNumber},
	})
}

// AddMembers 新增群成员
func (g *GroupScene) AddMembers(userIDs []string) {
	g.ActionList = append(g.ActionList, []containerInterface.IAction{
		&group.QueryInfo{GroupID: g.GroupNumber},
		&group.ModifyGroupMember{
			GroupID:    g.GroupNumber,
			UserIDs:    userIDs,
			AddOperate: true,
		},
		g.applyQueryGroupIcon(),
	})
}

// RemoveMembers 移除群成员
func (g *GroupScene) RemoveMembers(userIDs []string) {
	g.ActionList = append(g.ActionList, []containerInterface.IAction{
		&group.QueryInfo{GroupID: g.GroupNumber},
		&group.ModifyGroupMember{
			GroupID:    g.GroupNumber,
			UserIDs:    userIDs,
			AddOperate: false,
		},
		g.applyQueryGroupIcon(),
	})
}

// SetAdmins 设置管理员
func (g *GroupScene) SetAdmins(userIDs []string) {
	g.ActionList = append(g.ActionList, []containerInterface.IAction{
		&group.QueryInfo{GroupID: g.GroupNumber},
		&group.ModifyGroupAdmin{
			AddOperate: true,
			GroupID:    g.GroupNumber,
			UserIDs:    userIDs,
		},
		g.applyQueryGroupIcon(),
	})
}

// RemoveAdmins 移除管理员
func (g *GroupScene) RemoveAdmins(userIDs []string) {
	g.ActionList = append(g.ActionList, []containerInterface.IAction{
		&group.QueryInfo{GroupID: g.GroupNumber},
		&group.ModifyGroupAdmin{
			AddOperate: false,
			GroupID:    g.GroupNumber,
			UserIDs:    userIDs,
		},
		g.applyQueryGroupIcon(),
	})
}

// EnabledChatPermission 开启所有人聊天（所有人）
func (g *GroupScene) EnabledChatPermission() {
	g.ActionList = append(g.ActionList, []containerInterface.IAction{
		&group.QueryInfo{GroupID: g.GroupNumber},
		&group.ChatPermission{
			GroupID: g.GroupNumber,
			Enabled: true,
		},
		g.applyQueryGroupIcon(),
	})
}

// DisableChatPermission 关闭所有人聊天（仅管理员可用）
func (g *GroupScene) DisableChatPermission() {
	g.ActionList = append(g.ActionList, []containerInterface.IAction{
		&group.QueryInfo{GroupID: g.GroupNumber},
		&group.ChatPermission{
			GroupID: g.GroupNumber,
			Enabled: false,
		},
		g.applyQueryGroupIcon(),
	})
}

// EnabledEditDescPermission 开启编辑描述权限（所有人）
func (g *GroupScene) EnabledEditDescPermission() {
	g.ActionList = append(g.ActionList, []containerInterface.IAction{
		&group.QueryInfo{GroupID: g.GroupNumber},
		&group.EditDescPermission{
			GroupID: g.GroupNumber,
			Enabled: true,
		},
		g.applyQueryGroupIcon(),
	})
}

// DisableEditDescPermission 关闭编辑描述权限（仅管理员可用）
func (g *GroupScene) DisableEditDescPermission() {
	g.ActionList = append(g.ActionList, []containerInterface.IAction{
		&group.QueryInfo{GroupID: g.GroupNumber},
		&group.EditDescPermission{
			GroupID: g.GroupNumber,
			Enabled: false,
		},
		g.applyQueryGroupIcon(),
	})
}

func (g *GroupScene) MakeTextMessage(messageText string) {
	g.ActionList = append(g.ActionList, []containerInterface.IAction{
		&common.SubscribeStatus{UserID: g.GroupNumber, ToGroup: true},
		&common.InputChatState{UserID: g.GroupNumber, Input: true, ToGroup: true},
		&group.QueryMemberMultiDevicesIdentity{GroupID: g.GroupNumber},
		&group.SendText{GroupID: g.GroupNumber, MessageText: messageText},
	})
}

// Exit 退出群组
func (g *GroupScene) Exit() {
	g.ActionList = append(g.ActionList, []containerInterface.IAction{
		&group.Exit{GroupID: g.GroupNumber},
	})
}

// Build .
func (g *GroupScene) Build() containerInterface.IProcessor {
	return processor.NewOnceComposeProcessor(
		g.ActionList,
		processor.AliasName("group"),
		processor.AttachMonitor(&monitor.GroupMonitor{}),
	)
}
