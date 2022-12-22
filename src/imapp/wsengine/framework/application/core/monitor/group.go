package monitor

import (
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/action/group"
	"ws/framework/application/core/action/group/compose"
	privateChatCommon "ws/framework/application/core/action/private_chat/common"
	. "ws/framework/application/core/wam"
)

// GroupMonitor .
type GroupMonitor struct {
	hasAvatar      bool
	sendStateCount int32
}

// OnStart .
func (p *GroupMonitor) OnStart(ioc containerInterface.IAppIocContainer) {
	LogManager().SwitchAppMenu(ioc, PageSession)
}

// OnActionStartBefore .
func (p *GroupMonitor) OnActionStartBefore(_ interface{}, _ containerInterface.IMessageContext) {}

// OnActionStartAfter .
func (p *GroupMonitor) OnActionStartAfter(_ interface{}, _ containerInterface.IMessageContext) {}

func (p *GroupMonitor) OnActionStartFail(action interface{}, context containerInterface.IMessageContext) {

}

func (p *GroupMonitor) OnActionStartSuccess(action interface{}, context containerInterface.IMessageContext) {
	switch action.(type) {
	case *group.ModifyIcon:
		r := action.(*group.ModifyIcon)

		g := context.ResolveGroupService().Find(r.GroupID)
		p.hasAvatar = g != nil && g.HaveGroupIcon
	}
}

// ActionExecuteSuccess .
func (p *GroupMonitor) ActionExecuteSuccess(action interface{}, context containerInterface.IMessageContext) {
	switch action.(type) {
	case *privateChatCommon.SimulateInputChatState:
		state := action.(*privateChatCommon.SimulateInputChatState)
		p.sendStateCount = state.TotalCount
	case *groupComposeAction.CreateGroup:
		r := action.(*groupComposeAction.CreateGroup)

		LogManager().LogGroupCreate(context, uint32(len(r.Icon)), uint32(len(r.JoinUserIDs)))
	case *group.ModifyIcon:
		r := action.(*group.ModifyIcon)

		LogManager().LogGroupEditAvatar(context, p.hasAvatar, uint32(len(r.Icon)))
	case *group.ModifyGroupAdmin:
		r := action.(*group.ModifyGroupAdmin)

		g := context.ResolveGroupService().Find(r.GroupID)
		avatar := g != nil && g.HaveGroupIcon

		LogManager().LogGroupEditAdmins(context, avatar)
	case *group.ModifyGroupMember:
		LogManager().LogGroupEditMembers(context)
	case *group.SendText:
		LogManager().LogGroupSendText(context, p.sendStateCount, false)
	}

}

// ActionExecuteFailure .
func (p *GroupMonitor) ActionExecuteFailure(action interface{}, context containerInterface.IMessageContext) {
}

// OnExit .
func (p *GroupMonitor) OnExit(ioc containerInterface.IAppIocContainer) {}
