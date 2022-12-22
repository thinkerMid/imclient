package senderKeyDB

import (
	"ws/framework/plugin/database/database_tools"
)

// SenderKey .
type SenderKey struct {
	databaseTools.ChangeExtension

	JID            string `gorm:"column:our_jid;primaryKey"`
	TheirJID       string `gorm:"column:their_jid;primaryKey"`
	GroupID        string `gorm:"column:group_id;primaryKey"`
	DeviceID       *uint8 `gorm:"column:device_id;primaryKey"`
	ChatWith       bool   `gorm:"column:chat_with"`
	KeyID          uint32 `gorm:"column:key_id"`
	Iteration      uint32 `gorm:"column:iteration"`
	ChainKey       []byte `gorm:"column:chain_key"`
	PublicSignKey  []byte `gorm:"column:public_sign_key"`
	PrivateSignKey []byte `gorm:"column:private_sign_key"`
}

// TableName .
func (m *SenderKey) TableName() string {
	return "senderkey"
}

// UpdateKeyID .
func (m *SenderKey) UpdateKeyID(v uint32) {
	//if m.KeyID == v {
	//	return
	//}

	m.KeyID = v
	m.Update("key_id", v)
}

// UpdateIteration .
func (m *SenderKey) UpdateIteration(v uint32) {
	//if m.Iteration == v {
	//	return
	//}

	m.Iteration = v
	m.Update("iteration", v)
}

// UpdateChainKey .
func (m *SenderKey) UpdateChainKey(v []byte) {
	//if functionTools.SliceEqual(m.ChainKey, v) {
	//	return
	//}

	m.ChainKey = v
	m.Update("chain_key", v)
}

// UpdatePublicSignKey .
func (m *SenderKey) UpdatePublicSignKey(v []byte) {
	//if functionTools.SliceEqual(m.PublicSignKey, v) {
	//	return
	//}

	m.PublicSignKey = v
	m.Update("public_sign_key", v)
}

// UpdatePrivateSignKey .
func (m *SenderKey) UpdatePrivateSignKey(v []byte) {
	//if functionTools.SliceEqual(m.PrivateSignKey, v) {
	//	return
	//}

	m.PrivateSignKey = v
	m.Update("private_sign_key", v)
}
