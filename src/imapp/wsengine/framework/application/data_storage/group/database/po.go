package groupDB

import "ws/framework/plugin/database/database_tools"

// Group .
type Group struct {
	databaseTools.ChangeExtension

	JID                string `gorm:"column:jid;primaryKey"`
	GroupID            string `gorm:"column:group_id;primaryKey"`
	IsAdmin            bool   `gorm:"column:is_admin"`
	HaveGroupIcon      bool   `gorm:"column:have_group_icon"`
	EditDescKey        string `gorm:"column:last_edit_desc_key"`
	ChatPermission     bool   `gorm:"column:chat_permission"`      // true 所有人 false 管理员
	EditDescPermission bool   `gorm:"column:edit_desc_permission"` // true 所有人 false 管理员
}

// TableName .
func (s *Group) TableName() string {
	return "group"
}

// UpdateHaveGroupIcon .
func (s *Group) UpdateHaveGroupIcon(v bool) {
	if s.HaveGroupIcon == v {
		return
	}

	s.HaveGroupIcon = v
	s.Update("have_group_icon", v)
}

// UpdateEditDescKey .
func (s *Group) UpdateEditDescKey(v string) {
	if s.EditDescKey == v {
		return
	}

	s.EditDescKey = v
	s.Update("last_edit_desc_key", v)
}

// UpdateChatPermission .
func (s *Group) UpdateChatPermission(v bool) {
	if s.ChatPermission == v {
		return
	}

	s.ChatPermission = v
	s.Update("chat_permission", v)
}

// UpdateEditDescPermission .
func (s *Group) UpdateEditDescPermission(v bool) {
	if s.EditDescPermission == v {
		return
	}

	s.EditDescPermission = v
	s.Update("edit_desc_permission", v)
}

// UpdateAdmin .
func (s *Group) UpdateAdmin(v bool) {
	if s.IsAdmin == v {
		return
	}

	s.IsAdmin = v
	s.Update("is_admin", v)
}
