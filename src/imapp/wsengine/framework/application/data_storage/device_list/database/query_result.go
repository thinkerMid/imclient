package deviceListDB

import databaseTools "ws/framework/plugin/database/database_tools"

// DeviceQuery 设备查询
type DeviceQuery struct {
	JID      string `gorm:"column:our_jid;primaryKey"`
	TheirJID string `gorm:"column:their_jid;primaryKey"`
	ID       *uint8 `gorm:"column:device_id;primaryKey"`
}

// TableName .
func (d *DeviceQuery) TableName() string {
	return "device_list"
}

// ----------------------------------------------------------------------------

// DeleteDevice 删设备
type DeleteDevice struct {
	JID      string `gorm:"column:our_jid;"`
	TheirJID string `gorm:"column:their_jid;"`
	ID       uint8  `gorm:"column:device_id;primaryKey"`
}

// TableName .
func (d *DeleteDevice) TableName() string {
	return "device_list"
}

// ----------------------------------------------------------------------------

// DeviceID 查设备ID
type DeviceID struct {
	ID uint8 `gorm:"column:device_id;"`
}

// TableName .
func (d *DeviceID) TableName() string {
	return "device_list"
}

// ----------------------------------------------------------------------------

// SessionQuery 查会话
type SessionQuery struct {
	databaseTools.ChangeExtension

	JID            string `gorm:"column:our_jid;primaryKey"`
	TheirJID       string `gorm:"column:their_jid;primaryKey"`
	ID             *uint8 `gorm:"column:device_id;primaryKey"`
	Initialization *bool  `gorm:"column:initialization;"`
}

// TableName .
func (d *SessionQuery) TableName() string {
	return "device_list"
}

// UpdateInitialization .
func (d *SessionQuery) UpdateInitialization(v bool) {
	//if d.Initialization == v {
	//	return
	//}

	d.Initialization = &v
	d.Update("initialization", v)
}

// ----------------------------------------------------------------------------

// UpdateEncryptSession 保存加密操作后会话数据
type UpdateEncryptSession struct {
	databaseTools.ChangeExtension

	JID              string `gorm:"column:our_jid;primaryKey"`
	TheirJID         string `gorm:"column:their_jid;primaryKey"`
	ID               *uint8 `gorm:"column:device_id;primaryKey"`
	SenderChainIndex uint32 `gorm:"column:sender_chain_index"` // 发送用的计数器
	SenderChainKey   []byte `gorm:"column:sender_chain_key"`   // 发送用的密钥
}

// TableName .
func (d *UpdateEncryptSession) TableName() string {
	return "device_list"
}

// UpdateSenderChainIndex .
func (d *UpdateEncryptSession) UpdateSenderChainIndex(v uint32) {
	//if d.SenderChainIndex == v {
	//	return
	//}

	d.SenderChainIndex = v
	d.Update("sender_chain_index", v)
}

// UpdateSenderChainKey .
func (d *UpdateEncryptSession) UpdateSenderChainKey(v []byte) {
	//if functionTools.SliceEqual(d.SenderChainKey, v) {
	//	return
	//}

	d.SenderChainKey = v
	d.Update("sender_chain_key", v)
}

// ----------------------------------------------------------------------------

// UpdateDecryptSession 保存解密操作后会话数据
type UpdateDecryptSession struct {
	databaseTools.ChangeExtension

	JID                 string `gorm:"column:our_jid;primaryKey"`
	TheirJID            string `gorm:"column:their_jid;primaryKey"`
	ID                  *uint8 `gorm:"column:device_id;primaryKey"`
	UnacknowledgedState *bool  `gorm:"column:unacknowledged_state"` // 未确认的会话;pkmsg来源
	ReceiverChainIndex  uint32 `gorm:"column:receiver_chain_index"` // 接收用的计数器
	ReceiverChainKey    []byte `gorm:"column:receiver_chain_key"`   // 接收用的密钥
}

// TableName .
func (d *UpdateDecryptSession) TableName() string {
	return "device_list"
}

// UpdateUnacknowledgedState .
func (d *UpdateDecryptSession) UpdateUnacknowledgedState(v bool) {
	//if d.UnacknowledgedState == v {
	//	return
	//}

	d.UnacknowledgedState = &v
	d.Update("unacknowledged_state", v)
}

// UpdateReceiverChainIndex .
func (d *UpdateDecryptSession) UpdateReceiverChainIndex(v uint32) {
	//if d.ReceiverChainIndex == v {
	//	return
	//}

	d.ReceiverChainIndex = v
	d.Update("receiver_chain_index", v)
}

// UpdateReceiverChainKey .
func (d *UpdateDecryptSession) UpdateReceiverChainKey(v []byte) {
	//if functionTools.SliceEqual(d.ReceiverChainKey, v) {
	//	return
	//}

	d.ReceiverChainKey = v
	d.Update("receiver_chain_key", v)
}

// ----------------------------------------------------------------------------

// UpdateRebuildSession .
type UpdateRebuildSession struct {
	databaseTools.ChangeExtension

	JID                 string  `gorm:"column:our_jid;primaryKey"`
	TheirJID            string  `gorm:"column:their_jid;primaryKey"`
	ID                  *uint8  `gorm:"column:device_id;primaryKey"`
	UnacknowledgedState *bool   `gorm:"column:unacknowledged_state"` // 未确认的会话;pkmsg来源
	PreviousCounter     *uint32 `gorm:"column:previous_counter"`     // 接收用的计数器
	ReceiverPublicKey   []byte  `gorm:"column:receiver_public_key"`  // 接收用的密钥
	ReceiverChainIndex  uint32  `gorm:"column:receiver_chain_index"` // 接收用的计数器
	ReceiverChainKey    []byte  `gorm:"column:receiver_chain_key"`   // 接收用的密钥
	SenderPublicKey     []byte  `gorm:"column:sender_public_key"`    // 发送用的密钥
	SenderPrivateKey    []byte  `gorm:"column:sender_private_key"`   // 发送用的密钥
	SenderChainIndex    uint32  `gorm:"column:sender_chain_index"`   // 发送用的计数器
	SenderChainKey      []byte  `gorm:"column:sender_chain_key"`     // 发送用的密钥
	KdfRootKey          []byte  `gorm:"column:kdf_root_key"`         // 会话密钥
}

// TableName .
func (d *UpdateRebuildSession) TableName() string {
	return "device_list"
}

// UpdateUnacknowledgedState .
func (d *UpdateRebuildSession) UpdateUnacknowledgedState(v bool) {
	//if d.UnacknowledgedState == v {
	//	return
	//}

	d.UnacknowledgedState = &v
	d.Update("unacknowledged_state", v)
}

// UpdatePreviousCounter .
func (d *UpdateRebuildSession) UpdatePreviousCounter(v uint32) {
	//if d.PreviousCounter == v {
	//	return
	//}

	d.PreviousCounter = &v
	d.Update("previous_counter", v)
}

// UpdateReceiverPublicKey .
func (d *UpdateRebuildSession) UpdateReceiverPublicKey(v []byte) {
	//if functionTools.SliceEqual(d.ReceiverPublicKey, v) {
	//	return
	//}

	d.ReceiverPublicKey = v
	d.Update("receiver_public_key", v)
}

// UpdateReceiverChainIndex .
func (d *UpdateRebuildSession) UpdateReceiverChainIndex(v uint32) {
	//if d.ReceiverChainIndex == v {
	//	return
	//}

	d.ReceiverChainIndex = v
	d.Update("receiver_chain_index", v)
}

// UpdateReceiverChainKey .
func (d *UpdateRebuildSession) UpdateReceiverChainKey(v []byte) {
	//if functionTools.SliceEqual(d.ReceiverChainKey, v) {
	//	return
	//}

	d.ReceiverChainKey = v
	d.Update("receiver_chain_key", v)
}

// UpdateSenderPublicKey .
func (d *UpdateRebuildSession) UpdateSenderPublicKey(v []byte) {
	//if functionTools.SliceEqual(d.SenderPublicKey, v) {
	//	return
	//}

	d.SenderPublicKey = v
	d.Update("sender_public_key", v)
}

// UpdateSenderPrivateKey .
func (d *UpdateRebuildSession) UpdateSenderPrivateKey(v []byte) {
	//if functionTools.SliceEqual(d.SenderPrivateKey, v) {
	//	return
	//}

	d.SenderPrivateKey = v
	d.Update("sender_private_key", v)
}

// UpdateSenderChainIndex .
func (d *UpdateRebuildSession) UpdateSenderChainIndex(v uint32) {
	//if d.SenderChainIndex == v {
	//	return
	//}

	d.SenderChainIndex = v
	d.Update("sender_chain_index", v)
}

// UpdateSenderChainKey .
func (d *UpdateRebuildSession) UpdateSenderChainKey(v []byte) {
	//if functionTools.SliceEqual(d.SenderChainKey, v) {
	//	return
	//}

	d.SenderChainKey = v
	d.Update("sender_chain_key", v)
}

// UpdateKdfRootKey .
func (d *UpdateRebuildSession) UpdateKdfRootKey(v []byte) {
	//if functionTools.SliceEqual(d.KdfRootKey, v) {
	//	return
	//}

	d.KdfRootKey = v
	d.Update("kdf_root_key", v)
}
