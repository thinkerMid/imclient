package contact

import (
	"fmt"
	containerInterface "ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	messageResultType "ws/framework/application/core/result/constant"
	contactDB "ws/framework/application/data_storage/contact/database"
)

// CheckStranger 检测对方是不是有联系人记录
type CheckStranger struct {
	processor.BaseAction
	UserID  string // ID
	IsExist bool   // 期望的结果  有记录or无记录  结果相同会放行执行next

	FindByJID bool // 使用JID查找
}

// Start .
func (m *CheckStranger) Start(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
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

	find := contact != nil

	// 找到了 并且 期望是有记录
	if find && m.IsExist {
		next()
		return nil
	}

	// 找到了 并且 期望是没记录
	if find && m.IsExist == false {
		return fmt.Errorf("stranger is in contact")
	}

	// 实际没找到 并且 期望是没记录
	if !find && m.IsExist == false {
		// 放行
		next()
		return nil
	}

	return fmt.Errorf("not found stranger in contact")
}

// Receive .
func (m *CheckStranger) Receive(_ containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	return
}

func (m *CheckStranger) Error(context containerInterface.IMessageContext, err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.CheckStranger,
		Error:      err,
	})
}
