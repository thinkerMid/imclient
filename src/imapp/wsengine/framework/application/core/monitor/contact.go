package monitor

import (
	containerInterface "ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/action/contact"
	. "ws/framework/application/core/wam"
	signalProtocol "ws/framework/application/libsignal/protocol"
)

// ContactMonitor .
type ContactMonitor struct {
	haveSession   bool
	contactAvatar bool
}

// OnStart .
func (p *ContactMonitor) OnStart(ioc containerInterface.IAppIocContainer) {
	LogManager().SwitchAppMenu(ioc, PageSession)
}

// OnActionStartBefore .
func (p *ContactMonitor) OnActionStartBefore(_ interface{}, _ containerInterface.IMessageContext) {
}

// OnActionStartAfter .
func (p *ContactMonitor) OnActionStartAfter(_ interface{}, _ containerInterface.IMessageContext) {
}

func (p *ContactMonitor) OnActionStartFail(action interface{}, context containerInterface.IMessageContext) {
	switch action.(type) {
	case *contact.Check:
		// 添加联系人失败时，记录日志 [号码错误，或者 已存在该联系人]
		LogManager().LogContactAdd(context, false, false)
	}
}

func (p *ContactMonitor) OnActionStartSuccess(action interface{}, context containerInterface.IMessageContext) {
	switch action.(type) {
	case *contact.Delete:
		// 删除联系人时，记录日志
		query := action.(*contact.Delete)

		c := context.ResolveContactService().FindByJID(query.UserID)
		if c != nil {
			p.contactAvatar = c.HaveAvatar

			addr := signalProtocol.NewSignalAddress(c.JID, 0)
			p.haveSession = context.ResolveDeviceListService().ContainsSession(addr)
		}
	}
}

// ActionExecuteSuccess .
func (p *ContactMonitor) ActionExecuteSuccess(action interface{}, context containerInterface.IMessageContext) {
	switch action.(type) {
	case *contact.Add:
		// 添加联系人成功时，查询头像信息, 记录日志
		query := action.(*contact.Add)
		contact := context.ResolveContactService().FindByJID(query.UserID)
		if contact != nil {
			LogManager().LogContactAdd(context, true, contact.HaveAvatar)
		}
	case *contact.Delete:
		// 删除联系人时，记录日志
		LogManager().LogDeleteContact(context, p.haveSession, p.contactAvatar)
	}
}

func (p *ContactMonitor) ActionExecuteFailure(action interface{}, context containerInterface.IMessageContext) {
	switch action.(type) {
	case *contact.Query:
		// 添加联系人失败时，记录日志 [号码错误，或者 已存在该联系人]
		LogManager().LogContactAdd(context, false, false)
	}
}

// OnExit .
func (p *ContactMonitor) OnExit(containerInterface.IAppIocContainer) {

}
