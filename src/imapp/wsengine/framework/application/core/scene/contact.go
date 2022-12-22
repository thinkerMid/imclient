package scene

import (
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/action/common"
	"ws/framework/application/core/action/contact"
	"ws/framework/application/core/action/contact/compose"
	"ws/framework/application/core/action/user"
	"ws/framework/application/core/monitor"
	"ws/framework/application/core/processor"
)

// NewContact .
func NewContact() Contact {
	return Contact{}
}

// Contact .
type Contact struct {
	ActionList []containerInterface.IAction
}

// Build .
func (c *Contact) Build() containerInterface.IProcessor {
	return processor.NewOnceProcessor(c.ActionList,
		processor.AliasName("contact"),
		processor.AttachMonitor(&monitor.ReplaceUserIDWhenContact{}),
		processor.AttachMonitor(&monitor.ContactMonitor{}))
}

// Delete .
func (c *Contact) Delete(jid string) {
	c.ActionList = append(c.ActionList, &contact.Check{UserID: jid, InAddressBook: true})
	c.ActionList = append(c.ActionList, &contact.Delete{UserID: jid})
	c.ActionList = append(c.ActionList, &user.QueryStatusPrivacyList{IgnoreResponse: true})
	c.ActionList = append(c.ActionList, &common.DeleteDevices{UserID: jid})
}

// Search .
func (c *Contact) Search(jid string) {
	c.ActionList = append(c.ActionList, &contact.CheckStranger{UserID: jid, IsExist: false})
	c.ActionList = append(c.ActionList, &contact.QueryStranger{UserID: jid})
}

// Add .
func (c *Contact) Add(jid string) {
	c.ActionList = append(c.ActionList, &contact.Check{UserID: jid, InAddressBook: false})
	c.ActionList = append(c.ActionList, &contact.Add{UserID: jid})
	c.ActionList = append(c.ActionList, &common.QueryAvatarUrl{UserID: jid})
	c.ActionList = append(c.ActionList, &common.QueryDevicesIdentity{UserID: jid})
	c.ActionList = append(c.ActionList, &common.QueryUserDeviceList{UserID: jid})
}

// BatchSearch 用于假添加，纯检索号码是否注册，数量很大都不进行储存联系人关系和设备列表
func (c *Contact) BatchSearch(jids []string) {
	c.ActionList = append(c.ActionList, &contact.BatchAdd{UserIDs: jids})
}

// BatchAdd .
func (c *Contact) BatchAdd(jids []string) {
	c.ActionList = append(c.ActionList, &contactComposeAction.BatchAddContactAndCreateSession{UserIDs: jids})
}
