package contact

import (
	"fmt"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	"ws/framework/application/core/result/constant"
	contactDB "ws/framework/application/data_storage/contact/database"
)

// Check 检测手机号是否存了通讯录
type Check struct {
	processor.BaseAction
	UserID        string // ID
	InAddressBook bool   // 期望的结果  存了or没存  结果相同会放行执行next

	FindByJID bool // 使用JID查找
}

// Start .
func (m *Check) Start(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	var contact *contactDB.Contact

	// 使用手机号
	if !m.FindByJID {
		contact = context.ResolveContactService().FindByPhoneNumber(m.UserID)
		if contact != nil {
			// 判断是不是存在手机号和JID号码不一致的情况
			if len(contact.DstJIDUser) > 0 && m.UserID != contact.DstJIDUser {
				m.UserID = contact.DstJIDUser
			}
		}
	} else {
		contact = context.ResolveContactService().FindByJID(m.UserID)
	}

	// AddressBook 是已存储的属性
	find := contact != nil && contact.InAddressBook

	// 找到了 并且 期望是联系人
	if find && m.InAddressBook {
		next()
		return nil
	}

	// 找到了 并且 期望是非联系人
	if find && m.InAddressBook == false {
		return fmt.Errorf("dstUser is in contact")
	}

	// 实际没找到 并且 期望是没找到
	if !find && m.InAddressBook == false {
		// 放行
		next()
		return nil
	}

	return fmt.Errorf("dstUser is not in contact")
}

// Receive .
func (m *Check) Receive(_ containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	return
}

func (m *Check) Error(context containerInterface.IMessageContext, err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.CheckContact,
		Error:      err,
	})
}
