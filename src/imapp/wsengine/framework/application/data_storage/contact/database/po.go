package contactDB

import "ws/framework/plugin/database/database_tools"

// Contact .
type Contact struct {
	databaseTools.ChangeExtension

	JID            string `gorm:"column:src_number;primaryKey"`
	DstPhoneNumber string `gorm:"column:dst_phone_number;primaryKey"` // 手机号
	DstJIDUser     string `gorm:"column:dst_jid_number;primaryKey"`   // JID的号码
	TrustedContact bool   `gorm:"column:trusted_contact"`
	AddTime        int64  `gorm:"column:add_time"`        // 添加时间
	HaveAvatar     bool   `gorm:"column:have_avatar"`     // 有没头像
	ChatWith       bool   `gorm:"column:chat_with"`       // 发送过消息
	ReceiveChat    bool   `gorm:"column:receive_chat"`    // 收到过消息
	InAddressBook  bool   `gorm:"column:in_address_book"` // 已存储到通讯录
}

// TableName .
func (c *Contact) TableName() string {
	return "contacts"
}

// UpdateTrustedContact .
func (c *Contact) UpdateTrustedContact(v bool) {
	if v == c.TrustedContact {
		return
	}

	c.TrustedContact = v
	c.Update("trusted_contact", v)
}

// UpdateAddTime .
func (c *Contact) UpdateAddTime(v int64) {
	if v == c.AddTime {
		return
	}

	c.AddTime = v
	c.Update("add_time", v)
}

// UpdateHaveAvatar .
func (c *Contact) UpdateHaveAvatar(v bool) {
	if v == c.HaveAvatar {
		return
	}

	c.HaveAvatar = v
	c.Update("have_avatar", v)
}

// UpdateChatWith .
func (c *Contact) UpdateChatWith(v bool) {
	if v == c.ChatWith {
		return
	}

	c.ChatWith = v
	c.Update("chat_with", v)
}

// UpdateReceiveChat .
func (c *Contact) UpdateReceiveChat(v bool) {
	if v == c.ReceiveChat {
		return
	}

	c.ReceiveChat = v
	c.Update("receive_chat", v)
}

// UpdateInAddressBook .
func (c *Contact) UpdateInAddressBook(v bool) {
	if v == c.InAddressBook {
		return
	}

	c.InAddressBook = v
	c.Update("in_address_book", v)
}

// NewContactByDstPhoneNumber 注意使用场景，因存在目标手机号和whatsapp的JID不一致的情况 按需使用
func NewContactByDstPhoneNumber(src, dst string) Contact {
	return Contact{
		JID:            src,
		DstPhoneNumber: dst,
	}
}

// NewContactByDstJIDNumber 注意使用场景，因存在目标手机号和whatsapp的JID不一致的情况 按需使用
func NewContactByDstJIDNumber(src, dst string) Contact {
	return Contact{
		JID:        src,
		DstJIDUser: dst,
	}
}
