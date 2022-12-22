package deviceListService

import (
	"fmt"
	"strconv"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/data_storage/device_list/database"
	"ws/framework/application/libsignal/keys/chain"
	"ws/framework/application/libsignal/protocol"
	"ws/framework/application/libsignal/state/record"
	"ws/framework/application/libsignal/util/optional"
	"ws/framework/plugin/database"
	"ws/framework/plugin/database/database_tools"
	functionTools "ws/framework/utils/function_tools"
)

var _ containerInterface.IDeviceListService = &DeviceList{}

// DeviceList
//
//	对方设备管理 & 与对方设备的会话管理
type DeviceList struct {
	containerInterface.BaseService
}

func (c *DeviceList) formatCacheKey(address *protocol.SignalAddress) string {
	return c.JID.User + "_session_" + address.Name() + ":" + strconv.Itoa(int(address.DeviceID()))
}

// CreateSession .
func (c *DeviceList) CreateSession(address *protocol.SignalAddress, record *record.Session) {
	// 先创建设备
	c.AddDevice(address.Name(), uint8(address.DeviceID()))
	// 保存会话信息
	c.SaveSession(address, record)
}

// FindSession .
func (c *DeviceList) FindSession(address *protocol.SignalAddress) (*record.Session, error) {
	cacheKey := c.formatCacheKey(address)

	s, ok := c.AppIocContainer.ResolveMemoryCache().FindInCache(cacheKey)
	if ok {
		return s.(*record.Session), nil
	}

	deviceID := uint8(address.DeviceID())
	model := deviceListDB.Device{}
	where := deviceListDB.SessionQuery{
		JID:      c.JID.User,
		TheirJID: address.Name(),
		ID:       &deviceID,
	}

	err := databaseTools.Find(database.MasterDB(), &where, &model)
	if err != nil {
		return nil, err
	}

	if !model.Initialization {
		return nil, fmt.Errorf("%s session not initialization", address.String())
	}

	structure := record.SessionStructure{
		LocalIdentityPublic: *c.AppIocContainer.ResolveIdentityService().Context().Pub,
		LocalRegistrationID: c.AppIocContainer.ResolveDeviceService().Context().RegistrationId,
		PendingPreKey: &record.PendingPreKeyStructure{
			PreKeyID:       optional.NewOptionalUint32(model.PendingPreKeyID),
			SignedPreKeyID: model.PendingSignedPreKeyID,
			BaseKey:        functionTools.SliceTo32SizeArray(model.PendingBaseKey),
		},
		UnacknowledgedState: model.UnacknowledgedState,
		PreviousCounter:     model.PreviousCounter,
		ReceiverChain: &record.ChainStructure{
			SenderRatchetKeyPublic:  functionTools.SliceTo32SizeArray(model.ReceiverPublicKey),
			SenderRatchetKeyPrivate: functionTools.SliceTo32SizeArray(model.ReceiverPrivateKey),
			ChainKey: &chain.KeyStructure{
				Key:   model.ReceiverChainKey,
				Index: model.ReceiverChainIndex,
			},
		},
		RemoteIdentityPublic: functionTools.SliceTo32SizeArray(model.Identity),
		RemoteRegistrationID: model.RegistrationID,
		RootKey:              model.KdfRootKey,
		SenderBaseKey:        functionTools.SliceTo32SizeArray(model.SenderBaseKey),
		SenderChain: &record.ChainStructure{
			SenderRatchetKeyPublic:  functionTools.SliceTo32SizeArray(model.SenderPublicKey),
			SenderRatchetKeyPrivate: functionTools.SliceTo32SizeArray(model.SenderPrivateKey),
			ChainKey: &chain.KeyStructure{
				Key:   model.SenderChainKey,
				Index: model.SenderChainIndex,
			},
		},
		SessionVersion: int(model.SessionVersion),
	}

	session, err := record.NewSessionFromStructure(&structure)
	if err != nil {
		return nil, err
	}

	c.AppIocContainer.ResolveMemoryCache().Cache(c.formatCacheKey(address), session)

	return session, nil
}

// SaveSession .
func (c *DeviceList) SaveSession(address *protocol.SignalAddress, record *record.Session) {
	deviceID := uint8(address.DeviceID())
	structure := record.Pack()

	model := deviceListDB.Device{
		JID:      c.JID.User,
		TheirJID: address.Name(),
		ID:       &deviceID,
	}

	model.UpdateRegistrationID(structure.RemoteRegistrationID)
	model.UpdateIdentity(structure.RemoteIdentityPublic[:])
	model.UpdateInitialization(true)
	model.UpdateUnacknowledgedState(structure.UnacknowledgedState)
	model.UpdatePendingPreKeyID(structure.PendingPreKey.PreKeyID.Value)
	model.UpdatePendingSignedPreKeyID(structure.PendingPreKey.SignedPreKeyID)
	model.UpdatePendingBaseKey(structure.PendingPreKey.BaseKey[:])
	model.UpdatePreviousCounter(structure.PreviousCounter)
	model.UpdateReceiverPublicKey(structure.ReceiverChain.SenderRatchetKeyPublic[:])
	model.UpdateReceiverPrivateKey(structure.ReceiverChain.SenderRatchetKeyPrivate[:])
	model.UpdateReceiverChainIndex(structure.ReceiverChain.ChainKey.Index)
	model.UpdateReceiverChainKey(structure.ReceiverChain.ChainKey.Key)
	model.UpdateSenderPublicKey(structure.SenderChain.SenderRatchetKeyPublic[:])
	model.UpdateSenderPrivateKey(structure.SenderChain.SenderRatchetKeyPrivate[:])
	model.UpdateSenderChainIndex(structure.SenderChain.ChainKey.Index)
	model.UpdateSenderChainKey(structure.SenderChain.ChainKey.Key)
	model.UpdateSenderBaseKey(structure.SenderBaseKey[:])
	model.UpdateKdfRootKey(structure.RootKey)
	model.UpdateSessionVersion(uint32(structure.SessionVersion))

	_, err := databaseTools.Save(database.MasterDB(), &model)
	if err != nil {
		c.Logger.Error(err)
	}

	c.AppIocContainer.ResolveMemoryCache().Cache(c.formatCacheKey(address), record)
}

// SaveRebuildSession 保存会话重建的数据
func (c *DeviceList) SaveRebuildSession(address *protocol.SignalAddress, record *record.Session) {
	deviceID := uint8(address.DeviceID())
	structure := record.PackRebuildLogicData()

	model := deviceListDB.UpdateRebuildSession{
		JID:      c.JID.User,
		TheirJID: address.Name(),
		ID:       &deviceID,
	}

	model.UpdateUnacknowledgedState(structure.UnacknowledgedState)
	model.UpdatePreviousCounter(structure.PreviousCounter)
	model.UpdateReceiverPublicKey(structure.ReceiverChain.SenderRatchetKeyPublic[:])
	model.UpdateReceiverChainIndex(structure.ReceiverChain.ChainKey.Index)
	model.UpdateReceiverChainKey(structure.ReceiverChain.ChainKey.Key)
	model.UpdateSenderPublicKey(structure.SenderChain.SenderRatchetKeyPublic[:])
	model.UpdateSenderPrivateKey(structure.SenderChain.SenderRatchetKeyPrivate[:])
	model.UpdateSenderChainIndex(structure.SenderChain.ChainKey.Index)
	model.UpdateSenderChainKey(structure.SenderChain.ChainKey.Key)
	model.UpdateKdfRootKey(structure.RootKey)

	_, err := databaseTools.Save(database.MasterDB(), &model)
	if err != nil {
		c.Logger.Error(err)
	}
}

// SaveEncryptSession 保存加密操作后的会话数据
func (c *DeviceList) SaveEncryptSession(address *protocol.SignalAddress, record *record.Session) {
	deviceID := uint8(address.DeviceID())
	structure := record.PackEncryptLogicData()

	model := deviceListDB.UpdateEncryptSession{
		JID:      c.JID.User,
		TheirJID: address.Name(),
		ID:       &deviceID,
	}

	model.UpdateSenderChainIndex(structure.SenderChain.ChainKey.Index)
	model.UpdateSenderChainKey(structure.SenderChain.ChainKey.Key)

	_, err := databaseTools.Save(database.MasterDB(), &model)
	if err != nil {
		c.Logger.Error(err)
	}
}

// SaveDecryptSession 保存解密操作后的会话数据
func (c *DeviceList) SaveDecryptSession(address *protocol.SignalAddress, record *record.Session) {
	deviceID := uint8(address.DeviceID())
	structure := record.PackDecryptLogicData()

	model := deviceListDB.UpdateDecryptSession{
		JID:      c.JID.User,
		TheirJID: address.Name(),
		ID:       &deviceID,
	}

	// 接收过第一条消息 这个可能是true要更新
	if structure.ReceiverChain.ChainKey.Index == 1 {
		model.UpdateUnacknowledgedState(structure.UnacknowledgedState)
	}

	model.UpdateReceiverChainIndex(structure.ReceiverChain.ChainKey.Index)
	model.UpdateReceiverChainKey(structure.ReceiverChain.ChainKey.Key)

	_, err := databaseTools.Save(database.MasterDB(), &model)
	if err != nil {
		c.Logger.Error(err)
	}
}

// ContainsSession .
func (c *DeviceList) ContainsSession(address *protocol.SignalAddress) bool {
	deviceID := uint8(address.DeviceID())
	init := true

	where := deviceListDB.SessionQuery{
		JID:            c.JID.User,
		TheirJID:       address.Name(),
		ID:             &deviceID,
		Initialization: &init,
	}

	total, _ := databaseTools.Count(database.MasterDB(), &where)

	return total > 0
}

// DeleteSession 重置会话的数据
func (c *DeviceList) DeleteSession(address *protocol.SignalAddress) {
	deviceID := uint8(address.DeviceID())

	model := deviceListDB.SessionQuery{
		JID:      c.JID.User,
		TheirJID: address.Name(),
		ID:       &deviceID,
	}

	model.UpdateInitialization(false)

	_, err := databaseTools.Save(database.MasterDB(), &model)
	if err != nil {
		c.Logger.Error(err)
	}
}

// FindUnInitSessionDeviceIDList .
func (c *DeviceList) FindUnInitSessionDeviceIDList(dstJID string) (idList []uint8) {
	var init bool

	where := deviceListDB.SessionQuery{JID: c.JID.User, TheirJID: dstJID, Initialization: &init}
	result := make([]deviceListDB.DeviceID, 0)

	err := databaseTools.Find(database.MasterDB(), &where, &result)
	if err != nil || len(result) == 0 {
		return
	}

	for i := range result {
		idList = append(idList, result[i].ID)
	}

	return
}

// AddDevice .
func (c *DeviceList) AddDevice(dstJID string, deviceID uint8) {
	emptySlice := make([]byte, 0)

	model := deviceListDB.Device{
		JID:                c.JID.User,
		TheirJID:           dstJID,
		ID:                 &deviceID,
		Identity:           emptySlice,
		PendingBaseKey:     emptySlice,
		ReceiverPublicKey:  emptySlice,
		ReceiverPrivateKey: emptySlice,
		ReceiverChainKey:   emptySlice,
		SenderPublicKey:    emptySlice,
		SenderPrivateKey:   emptySlice,
		SenderChainKey:     emptySlice,
		SenderBaseKey:      emptySlice,
		KdfRootKey:         emptySlice,
	}

	_, err := databaseTools.Create(database.MasterDB(), &model)
	if err != nil {
		c.Logger.Error(err)
	}
}

// DeleteDevice .
func (c *DeviceList) DeleteDevice(dstJID string, deviceID uint8) {
	primaryKey := deviceListDB.DeviceQuery{JID: c.JID.User, TheirJID: dstJID, ID: &deviceID}

	_, err := databaseTools.DeleteByPrimaryKey(database.MasterDB(), &primaryKey)
	if err != nil {
		c.Logger.Error(err)
	}

	// 移除缓存
	address := protocol.NewSignalAddress(dstJID, uint32(deviceID))
	c.AppIocContainer.ResolveMemoryCache().UnCache(c.formatCacheKey(address))
}

// BatchDeleteDevice .
func (c *DeviceList) BatchDeleteDevice(dst string, deviceIDList []uint8) {
	where := deviceListDB.DeleteDevice{
		JID:      c.JID.User,
		TheirJID: dst,
	}

	_, err := databaseTools.BatchDeleteUint8IDByPrimaryKey(database.MasterDB(), &where, deviceIDList)
	if err != nil {
		c.Logger.Error(err)
	}

	// 移除缓存
	for _, v := range deviceIDList {
		address := protocol.NewSignalAddress(dst, uint32(v))
		c.AppIocContainer.ResolveMemoryCache().UnCache(c.formatCacheKey(address))
	}
}

// BatchCreateDevice .
func (c *DeviceList) BatchCreateDevice(dst string, deviceIDList []uint8) {
	deviceIDs := make([]deviceListDB.Device, 0)
	emptySlice := make([]byte, 0)

	for i := range deviceIDList {
		deviceIDs = append(deviceIDs, deviceListDB.Device{
			JID:                c.JID.User,
			TheirJID:           dst,
			ID:                 &deviceIDList[i],
			Identity:           emptySlice,
			PendingBaseKey:     emptySlice,
			ReceiverPublicKey:  emptySlice,
			ReceiverPrivateKey: emptySlice,
			ReceiverChainKey:   emptySlice,
			SenderPublicKey:    emptySlice,
			SenderPrivateKey:   emptySlice,
			SenderChainKey:     emptySlice,
			SenderBaseKey:      emptySlice,
			KdfRootKey:         emptySlice,
		})
	}

	_, err := databaseTools.BatchCreate(database.MasterDB(), deviceIDs)
	if err != nil {
		c.Logger.Error(err)
	}
}

// FindDeviceIDList .
func (c *DeviceList) FindDeviceIDList(dstJID string) (idList []uint8) {
	where := deviceListDB.DeviceQuery{JID: c.JID.User, TheirJID: dstJID}
	result := make([]deviceListDB.DeviceID, 0)

	err := databaseTools.Find(database.MasterDB(), &where, &result)
	if err != nil || len(result) == 0 {
		return
	}

	for i := range result {
		idList = append(idList, result[i].ID)
	}

	return idList
}

// HaveMultiDevice .
func (c *DeviceList) HaveMultiDevice(dstJID string) bool {
	where := deviceListDB.DeviceQuery{JID: c.JID.User, TheirJID: dstJID}

	total, _ := databaseTools.Count(database.MasterDB(), &where)

	return total > 1
}

// UpdateDeviceList 如果deviceIDList是空的 要走删除 而不是更新
func (c *DeviceList) UpdateDeviceList(dstJID string, deviceIDList []uint8) {
	if len(deviceIDList) == 0 {
		return
	}

	oldDeviceList := c.FindDeviceIDList(dstJID)

	if len(oldDeviceList) > 0 {
		add, remove := functionTools.ScanDifferentUint8(deviceIDList, oldDeviceList)

		if len(add) > 0 {
			c.BatchCreateDevice(dstJID, add)
		}

		if len(remove) > 0 {
			c.BatchDeleteDevice(dstJID, remove)
		}
	} else {
		c.BatchCreateDevice(dstJID, deviceIDList)
	}
}

// CleanupAllData .
func (c *DeviceList) CleanupAllData() {
	_, _ = deviceListDB.DeleteByJID(database.MasterDB(), c.JID.User)
}
