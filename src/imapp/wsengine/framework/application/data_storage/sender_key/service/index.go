package senderKeyService

import (
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/data_storage/sender_key/database"
	"ws/framework/application/libsignal/groups/state/record"
	signalProtocol "ws/framework/application/libsignal/protocol"
	"ws/framework/plugin/database"
	"ws/framework/plugin/database/database_tools"
	functionTools "ws/framework/utils/function_tools"
)

var _ containerInterface.ISenderKeyService = &SenderKey{}

// SenderKey
//
//	TODO: 缓存群成员密钥的管理成本过高，暂时不优化
type SenderKey struct {
	containerInterface.BaseService
}

// DeleteDevice 把该设备从所有群内删除
func (s *SenderKey) DeleteDevice(deviceID types.JID) {
	primaryKey := senderKeyDB.DeleteSenderDevice{TheirJID: deviceID.User}
	where := senderKeyDB.DeleteSenderDevice{JID: s.JID.User, DeviceID: &deviceID.Device}

	_, err := databaseTools.Delete(database.MasterDB(), &primaryKey, &where)
	if err != nil {
		s.Logger.Error(err)
	}
}

// SearchSenderInGroupAndCreate 查找目标所在的群组并批量新增对应的设备
func (s *SenderKey) SearchSenderInGroupAndCreate(deviceID types.JID) {
	var defaultDeviceID uint8

	where := senderKeyDB.QuerySenderDeviceInGroup{JID: s.JID.User, TheirJID: deviceID.User, DeviceID: &defaultDeviceID}
	result := make([]senderKeyDB.SenderGroup, 0)

	err := databaseTools.Find(database.MasterDB(), &where, &result)
	if err != nil {
		s.Logger.Error(err)
		return
	}

	senders := make([]senderKeyDB.SenderKey, len(result))
	buffer := make([]byte, 0)

	for i := range result {
		senders[i].JID = s.JID.User
		senders[i].TheirJID = deviceID.User
		senders[i].DeviceID = &deviceID.Device
		senders[i].GroupID = result[i].GroupID
		senders[i].ChainKey = buffer
		senders[i].PublicSignKey = buffer
		senders[i].PrivateSignKey = buffer
	}

	_, err = databaseTools.BatchCreate(database.MasterDB(), senders)
	if err != nil {
		s.Logger.Error(err)
	}
}

// BatchDeleteSenderByGroupIDAndJID 删除群内的成员
func (s *SenderKey) BatchDeleteSenderByGroupIDAndJID(groupID string, senderJID []string) {
	mode := senderKeyDB.DeleteSenderDevice{JID: s.JID.User, GroupID: groupID}

	_, err := databaseTools.BatchDeleteStringByPrimaryKey(database.MasterDB(), &mode, senderJID)
	if err != nil {
		s.Logger.Error(err)
	}
}

// BatchCreateDevice 添加设备
func (s *SenderKey) BatchCreateDevice(groupID string, deviceIDs []types.JID) {
	senders := make([]senderKeyDB.SenderKey, len(deviceIDs))
	buffer := make([]byte, 0)

	for i := range deviceIDs {
		senders[i].JID = s.JID.User
		senders[i].TheirJID = deviceIDs[i].User
		senders[i].DeviceID = &deviceIDs[i].Device
		senders[i].GroupID = groupID
		senders[i].ChainKey = buffer
		senders[i].PublicSignKey = buffer
		senders[i].PrivateSignKey = buffer
	}

	_, err := databaseTools.BatchCreate(database.MasterDB(), senders)
	if err != nil {
		s.Logger.Error(err)
	}
}

// FindDevice 查找设备
func (s *SenderKey) FindDevice(senderKeyName *signalProtocol.SenderKeyName) *senderKeyDB.SenderDevice {
	where := s.createQueryWhereBySenderKeyName(senderKeyName)
	result := senderKeyDB.SenderDevice{}

	err := databaseTools.Find(database.MasterDB(), &where, &result)
	if err != nil {
		return nil
	}

	return &result
}

// FindUnSentMessageDeviceByGroupID 查找群内未发送过消息设备
func (s *SenderKey) FindUnSentMessageDeviceByGroupID(groupID string) ([]senderKeyDB.SenderDevice, error) {
	var chatWith bool

	where := senderKeyDB.QuerySenderDeviceInGroup{JID: s.JID.User, GroupID: groupID, ChatWith: &chatWith}
	result := make([]senderKeyDB.SenderDevice, 0)

	err := databaseTools.Find(database.MasterDB(), &where, &result)

	return result, err
}

// SaveSentMessageByGroupID 设置该群发送过消息
func (s *SenderKey) SaveSentMessageByGroupID(groupID string) {
	data := senderKeyDB.QuerySenderDeviceInGroup{JID: s.JID.User, GroupID: groupID}
	data.UpdateChatWith(true)

	_, err := databaseTools.Save(database.MasterDB(), &data)
	if err != nil {
		s.Logger.Error(err)
	}
}

// DeleteAllDeviceByGroupID 批量删除群组内的设备
func (s *SenderKey) DeleteAllDeviceByGroupID(groupID string) {
	where := senderKeyDB.QuerySenderDeviceInGroup{JID: s.JID.User, GroupID: groupID}

	_, err := databaseTools.DeleteByPrimaryKey(database.MasterDB(), where)
	if err != nil {
		s.Logger.Error(err)
	}
}

// CreateSenderKey 创建群密钥
func (s *SenderKey) CreateSenderKey(senderKeyName *signalProtocol.SenderKeyName, keyRecord *record.SenderKey) {
	model := s.createQueryWhereBySenderKeyName(senderKeyName)

	// 转换
	senderKeyStructure := keyRecord.Pack()
	model.KeyID = senderKeyStructure.KeyID
	model.Iteration = senderKeyStructure.Iteration
	model.ChainKey = senderKeyStructure.ChainKey
	model.PublicSignKey = senderKeyStructure.SigningKeyPublic[:]
	model.PrivateSignKey = senderKeyStructure.SigningKeyPrivate[:]

	// 如果是自己的key 默认设置已发送消息 这个用于区分群员的否发送过消息属性 自身不需要
	if model.TheirJID == s.JID.User {
		model.ChatWith = true
	}

	_, err := databaseTools.Create(database.MasterDB(), &model)
	if err != nil {
		s.Logger.Error(err)
	}
}

// ResetSenderKey 重置群密钥
func (s *SenderKey) ResetSenderKey(senderKeyName *signalProtocol.SenderKeyName, keyRecord *record.SenderKey) {
	model := s.createQueryWhereBySenderKeyName(senderKeyName)

	senderKeyStructure := keyRecord.Pack()

	model.UpdateKeyID(senderKeyStructure.KeyID)
	model.UpdateIteration(senderKeyStructure.Iteration)
	model.UpdateChainKey(senderKeyStructure.ChainKey)
	model.UpdatePublicSignKey(senderKeyStructure.SigningKeyPublic[:])
	model.UpdatePrivateSignKey(senderKeyStructure.SigningKeyPrivate[:])

	_, err := databaseTools.Save(database.MasterDB(), &model)
	if err != nil {
		s.Logger.Error(err)
	}
}

// UpdateSenderKey 更新群密钥
func (s *SenderKey) UpdateSenderKey(senderKeyName *signalProtocol.SenderKeyName, keyRecord *record.SenderKey) {
	model := s.createQueryWhereBySenderKeyName(senderKeyName)

	senderKeyStructure := keyRecord.PackChainKey()

	model.UpdateIteration(senderKeyStructure.Iteration)
	model.UpdateChainKey(senderKeyStructure.ChainKey)

	_, err := databaseTools.Save(database.MasterDB(), &model)
	if err != nil {
		s.Logger.Error(err)
	}
}

// ContainsSenderKey 是否有群密钥
func (s *SenderKey) ContainsSenderKey(senderKeyName *signalProtocol.SenderKeyName) bool {
	where := s.createQueryWhereBySenderKeyName(senderKeyName)

	total, _ := databaseTools.Count(database.MasterDB(), &where)

	return total > 0
}

// FindSenderKey 查找密钥
func (s *SenderKey) FindSenderKey(senderKeyName *signalProtocol.SenderKeyName) (*record.SenderKey, error) {
	where := s.createQueryWhereBySenderKeyName(senderKeyName)
	result := senderKeyDB.QueryCryptoKey{}

	err := databaseTools.Find(database.MasterDB(), &where, &result)
	if err != nil {
		return nil, err
	}

	senderKeyStructure := record.SenderKeyStructure{
		KeyID:             result.KeyID,
		Iteration:         result.Iteration,
		ChainKey:          result.ChainKey,
		SigningKeyPrivate: functionTools.SliceTo32SizeArray(result.PrivateKey),
		SigningKeyPublic:  functionTools.SliceTo32SizeArray(result.PublicSignKey),
	}

	key, err := record.NewSenderKeyFromStructure(&senderKeyStructure)
	if err != nil {
		return nil, err
	}

	return key, nil
}

// CleanupAllData .
func (s *SenderKey) CleanupAllData() {
	_, _ = senderKeyDB.DeleteByJID(database.MasterDB(), s.JID.User)
}

func (s *SenderKey) createQueryWhereBySenderKeyName(senderKeyName *signalProtocol.SenderKeyName) senderKeyDB.SenderKey {
	deviceID := uint8(senderKeyName.Sender().DeviceID())

	return senderKeyDB.SenderKey{
		JID:      s.JID.User,
		TheirJID: senderKeyName.Sender().Name(),
		GroupID:  senderKeyName.GroupID(),
		DeviceID: &deviceID,
	}
}
