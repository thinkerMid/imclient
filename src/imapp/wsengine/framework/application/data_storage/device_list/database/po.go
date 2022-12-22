package deviceListDB

import (
	"ws/framework/plugin/database/database_tools"
)

// Device .
type Device struct {
	databaseTools.ChangeExtension

	JID                   string `gorm:"column:our_jid;primaryKey"`
	TheirJID              string `gorm:"column:their_jid;primaryKey"`
	ID                    *uint8 `gorm:"column:device_id;primaryKey"`
	RegistrationID        uint32 `gorm:"column:registration_id"`          // 设备注册的ID
	Identity              []byte `gorm:"column:identity"`                 // 设备身份
	Initialization        bool   `gorm:"column:initialization"`           // 初始化标志
	UnacknowledgedState   bool   `gorm:"column:unacknowledged_state"`     // 未确认的会话;pkmsg来源
	PendingPreKeyID       uint32 `gorm:"column:pending_prekey_id"`        // pkmsg密钥ID
	PendingSignedPreKeyID uint32 `gorm:"column:pending_signed_prekey_id"` // pkmsg密钥ID
	PendingBaseKey        []byte `gorm:"column:pending_base_key"`         // pkmsg密钥
	PreviousCounter       uint32 `gorm:"column:previous_counter"`         // 接收用的计数器
	ReceiverPublicKey     []byte `gorm:"column:receiver_public_key"`      // 接收用的密钥
	ReceiverPrivateKey    []byte `gorm:"column:receiver_private_key"`     // 接收用的密钥
	ReceiverChainIndex    uint32 `gorm:"column:receiver_chain_index"`     // 接收用的计数器
	ReceiverChainKey      []byte `gorm:"column:receiver_chain_key"`       // 接收用的密钥
	SenderPublicKey       []byte `gorm:"column:sender_public_key"`        // 发送用的密钥
	SenderPrivateKey      []byte `gorm:"column:sender_private_key"`       // 发送用的密钥
	SenderChainIndex      uint32 `gorm:"column:sender_chain_index"`       // 发送用的计数器
	SenderChainKey        []byte `gorm:"column:sender_chain_key"`         // 发送用的密钥
	SenderBaseKey         []byte `gorm:"column:sender_base_key"`          // 会话密钥
	KdfRootKey            []byte `gorm:"column:kdf_root_key"`             // 会话密钥
	SessionVersion        uint32 `gorm:"column:session_version"`          // 会话版本
}

// TableName .
func (d *Device) TableName() string {
	return "device_list"
}

// UpdateRegistrationID .
func (d *Device) UpdateRegistrationID(v uint32) {
	//if d.RegistrationID == v {
	//	return
	//}

	d.RegistrationID = v
	d.Update("registration_id", v)
}

// UpdateIdentity .
func (d *Device) UpdateIdentity(v []byte) {
	//if functionTools.SliceEqual(d.Identity, v) {
	//	return
	//}

	d.Identity = v
	d.Update("identity", v)
}

// UpdateInitialization .
func (d *Device) UpdateInitialization(v bool) {
	//if d.Initialization == v {
	//	return
	//}

	d.Initialization = v
	d.Update("initialization", v)
}

// UpdateUnacknowledgedState .
func (d *Device) UpdateUnacknowledgedState(v bool) {
	//if d.UnacknowledgedState == v {
	//	return
	//}

	d.UnacknowledgedState = v
	d.Update("unacknowledged_state", v)
}

// UpdatePendingPreKeyID .
func (d *Device) UpdatePendingPreKeyID(v uint32) {
	//if d.PendingPreKeyID == v {
	//	return
	//}

	d.PendingPreKeyID = v
	d.Update("pending_prekey_id", v)
}

// UpdatePendingSignedPreKeyID .
func (d *Device) UpdatePendingSignedPreKeyID(v uint32) {
	//if d.PendingSignedPreKeyID == v {
	//	return
	//}

	d.PendingSignedPreKeyID = v
	d.Update("pending_signed_prekey_id", v)
}

// UpdatePendingBaseKey .
func (d *Device) UpdatePendingBaseKey(v []byte) {
	//if functionTools.SliceEqual(d.PendingBaseKey, v) {
	//	return
	//}

	d.PendingBaseKey = v
	d.Update("pending_base_key", v)
}

// UpdatePreviousCounter .
func (d *Device) UpdatePreviousCounter(v uint32) {
	//if d.PreviousCounter == v {
	//	return
	//}

	d.PreviousCounter = v
	d.Update("previous_counter", v)
}

// UpdateReceiverPublicKey .
func (d *Device) UpdateReceiverPublicKey(v []byte) {
	//if functionTools.SliceEqual(d.ReceiverPublicKey, v) {
	//	return
	//}

	d.ReceiverPublicKey = v
	d.Update("receiver_public_key", v)
}

// UpdateReceiverPrivateKey .
func (d *Device) UpdateReceiverPrivateKey(v []byte) {
	//if functionTools.SliceEqual(d.ReceiverPrivateKey, v) {
	//	return
	//}

	d.ReceiverPrivateKey = v
	d.Update("receiver_private_key", v)
}

// UpdateReceiverChainIndex .
func (d *Device) UpdateReceiverChainIndex(v uint32) {
	//if d.ReceiverChainIndex == v {
	//	return
	//}

	d.ReceiverChainIndex = v
	d.Update("receiver_chain_index", v)
}

// UpdateReceiverChainKey .
func (d *Device) UpdateReceiverChainKey(v []byte) {
	//if functionTools.SliceEqual(d.ReceiverChainKey, v) {
	//	return
	//}

	d.ReceiverChainKey = v
	d.Update("receiver_chain_key", v)
}

// UpdateSenderPublicKey .
func (d *Device) UpdateSenderPublicKey(v []byte) {
	//if functionTools.SliceEqual(d.SenderPublicKey, v) {
	//	return
	//}

	d.SenderPublicKey = v
	d.Update("sender_public_key", v)
}

// UpdateSenderPrivateKey .
func (d *Device) UpdateSenderPrivateKey(v []byte) {
	//if functionTools.SliceEqual(d.SenderPrivateKey, v) {
	//	return
	//}

	d.SenderPrivateKey = v
	d.Update("sender_private_key", v)
}

// UpdateSenderChainIndex .
func (d *Device) UpdateSenderChainIndex(v uint32) {
	//if d.SenderChainIndex == v {
	//	return
	//}

	d.SenderChainIndex = v
	d.Update("sender_chain_index", v)
}

// UpdateSenderChainKey .
func (d *Device) UpdateSenderChainKey(v []byte) {
	//if functionTools.SliceEqual(d.SenderChainKey, v) {
	//	return
	//}

	d.SenderChainKey = v
	d.Update("sender_chain_key", v)
}

// UpdateSenderBaseKey .
func (d *Device) UpdateSenderBaseKey(v []byte) {
	//if functionTools.SliceEqual(d.SenderBaseKey, v) {
	//	return
	//}

	d.SenderBaseKey = v
	d.Update("sender_base_key", v)
}

// UpdateKdfRootKey .
func (d *Device) UpdateKdfRootKey(v []byte) {
	//if functionTools.SliceEqual(d.KdfRootKey, v) {
	//	return
	//}

	d.KdfRootKey = v
	d.Update("kdf_root_key", v)
}

// UpdateSessionVersion .
func (d *Device) UpdateSessionVersion(v uint32) {
	//if d.SessionVersion == v {
	//	return
	//}

	d.SessionVersion = v
	d.Update("session_version", v)
}
