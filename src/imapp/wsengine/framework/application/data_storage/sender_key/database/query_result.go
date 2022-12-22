package senderKeyDB

import (
	databaseTools "ws/framework/plugin/database/database_tools"
)

// QueryCryptoKey 对应的设备密钥
type QueryCryptoKey struct {
	KeyID         uint32 `gorm:"column:key_id"`
	Iteration     uint32 `gorm:"column:iteration"`
	ChainKey      []byte `gorm:"column:chain_key"`
	PublicSignKey []byte `gorm:"column:public_sign_key"`
	PrivateKey    []byte `gorm:"column:private_sign_key"`
}

// TableName .
func (m *QueryCryptoKey) TableName() string {
	return "senderkey"
}

// ----------------------------------------------------------------------------

// SenderDevice 群设备
type SenderDevice struct {
	TheirJID string `gorm:"column:their_jid;primaryKey"`
	DeviceID uint16 `gorm:"column:device_id;"`
}

// TableName .
func (m *SenderDevice) TableName() string {
	return "senderkey"
}

// ----------------------------------------------------------------------------

// QuerySenderDeviceInGroup 查询群内的所有设备
type QuerySenderDeviceInGroup struct {
	databaseTools.ChangeExtension

	JID      string `gorm:"column:our_jid;primaryKey"`
	GroupID  string `gorm:"column:group_id;primaryKey"`
	TheirJID string `gorm:"column:their_jid"`
	DeviceID *uint8 `gorm:"column:device_id"`
	ChatWith *bool  `gorm:"column:chat_with"`
}

// TableName .
func (m *QuerySenderDeviceInGroup) TableName() string {
	return "senderkey"
}

// UpdateChatWith .
func (m *QuerySenderDeviceInGroup) UpdateChatWith(v bool) {
	if &v == m.ChatWith {
		return
	}

	m.ChatWith = &v
	m.Update("chat_with", v)
}

// ----------------------------------------------------------------------------

// DeleteSenderDevice 删除群内的JID所有设备
type DeleteSenderDevice struct {
	JID      string `gorm:"column:our_jid"`
	TheirJID string `gorm:"column:their_jid;primaryKey"`
	GroupID  string `gorm:"column:group_id"`
	DeviceID *uint8 `gorm:"column:device_id"`
}

// TableName .
func (m *DeleteSenderDevice) TableName() string {
	return "senderkey"
}

// ----------------------------------------------------------------------------

// SenderGroup 查群组
type SenderGroup struct {
	GroupID string `gorm:"column:group_id"`
}

// TableName .
func (m *SenderGroup) TableName() string {
	return "senderkey"
}
