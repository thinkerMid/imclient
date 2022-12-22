package preKeyService

import (
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/data_storage/prekey/database"
	"ws/framework/application/libsignal/serialize"
	"ws/framework/application/libsignal/state/record"
	"ws/framework/application/libsignal/util/keyhelper"
	"ws/framework/plugin/database"
	"ws/framework/plugin/database/database_tools"
	"ws/framework/utils/keys"
)

var _ containerInterface.IPreKeyService = &PreKey{}

// 2022-4-14梳理的逻辑：
// prekey是用于解密陌生人私信的，客户端上传到服务器后，如果收到陌生人私信给自己，服务器会先用prekey进行加密返回给客户端
// 这时客户端就可以用私信中的prekey的keyID去进行解密，之后客户端删除这个key
// 再加上有一个通知类型的消息 会让客户端上传prekey 那么基本可以确定prekey是这么个用途
var wantedPreKeyCount = 812

// PreKey .
type PreKey struct {
	containerInterface.BaseService
}

// StatementPreKeyCount .
func (p *PreKey) StatementPreKeyCount() int {
	return wantedPreKeyCount
}

// SavePreKey .
func (p *PreKey) SavePreKey(preKeyID uint32, preKeyRecord *record.PreKey) {
	model := prekeyDB.PreKey{
		JID:     p.JID.User,
		KeyId:   preKeyID,
		KeyBuff: preKeyRecord.Serialize(),
	}

	_, err := databaseTools.Create(database.MasterDB(), &model)
	if err != nil {
		p.Logger.Error("store pre keys ", err)
	}
}

// ContainsPreKey .
func (p *PreKey) ContainsPreKey(preKeyID uint32) bool {
	return true
}

// DeletePreKey .
func (p *PreKey) DeletePreKey(_ uint32) {
	// 不需要删除 whatsapp服务器是随机使用preKeyID的 一旦删除 客户端就无法找到对应的preKey进行解析内容
	// 如果需要恢复正常，需要让发私信的那个账号，删除与该客户端的设备信息和会话信息，再次重新走私信流程就会正常
}

// FindPreKey .
func (p *PreKey) FindPreKey(preKeyID uint32) *record.PreKey {
	where := prekeyDB.PreKey{JID: p.JID.User, KeyId: preKeyID}
	result := prekeyDB.QueryResult{}

	err := databaseTools.Find(database.MasterDB(), &where, &result)
	if err != nil {
		//TODO 兼容旧数据 一个月后或者一个版本后可以删除了
		where.JID = p.JID.String()
		err = databaseTools.Find(database.MasterDB(), &where, &result)
	}

	if err != nil {
		return nil
	}

	preKey, err := record.NewPreKeyFromBytes(result.KeyBuff, serialize.Proto.PreKeyRecord)

	if err != nil {
		p.Logger.Error("load pre keys ", err)
		return nil
	}

	return preKey
}

// InitPreKeys .
func (p *PreKey) InitPreKeys() ([]keys.PreKey, error) {
	where := prekeyDB.PreKey{JID: p.JID.User}
	remainKey, err := prekeyDB.FindAll(database.MasterDB(), &where, wantedPreKeyCount)
	if err != nil {
		return nil, err
	}

	newKeys := make([]keys.PreKey, wantedPreKeyCount)

	// 有问题的key
	removePreKeyId := make([]uint32, 0)
	defer p.batchRemovePreKey(removePreKeyId)

	// 填充
	var fillIdx int
	for i := range remainKey {
		model := remainKey[i]

		preKey, err := record.NewPreKeyFromBytes(model.KeyBuff, serialize.Proto.PreKeyRecord)
		if err != nil {
			removePreKeyId = append(removePreKeyId, model.KeyId)
			continue
		}

		pri := preKey.KeyPair().PrivateKey().Serialize()

		newKeys[fillIdx].KeyID = model.KeyId
		newKeys[fillIdx].KeyPair = *keys.NewKeyPairFromPrivateKey(pri)

		fillIdx++
	}

	remain := wantedPreKeyCount - fillIdx

	// 足够了
	if remain == 0 {
		return newKeys, nil
	}

	var lastPreKeyID uint32
	if len(remainKey) == 0 {
		lastPreKeyID = 1
	} else {
		lastPreKeyID = remainKey[len(remainKey)-1].KeyId
		lastPreKeyID++
		remain += int(lastPreKeyID)
	}

	// 自增 （GeneratePreKeys 这儿的生成是小于等于循环的 约定的下标是从1开始的 所以start始终是要+1的）
	preKeys, err := keyhelper.GeneratePreKeys(int(lastPreKeyID), remain, serialize.Proto.PreKeyRecord)
	if err != nil {
		panic("unable to generate pre keys!")
	}

	err = p.savePreKeys(preKeys)
	if err != nil {
		return nil, err
	}

	for i := range preKeys {
		pri := preKeys[i].KeyPair().PrivateKey().Serialize()

		newKeys[fillIdx].KeyID = preKeys[i].ID().Value
		newKeys[fillIdx].KeyPair = *keys.NewKeyPairFromPrivateKey(pri)

		fillIdx++
	}

	return newKeys, nil
}

// GeneratePreKeys .
func (p *PreKey) GeneratePreKeys(generateCount int) ([]keys.PreKey, error) {
	lastKey := prekeyDB.PreKey{JID: p.JID.User}
	err := prekeyDB.FindLast(database.MasterDB(), &lastKey)
	if err != nil {
		return nil, err
	}

	newKeys := make([]keys.PreKey, generateCount)
	keyIdx := lastKey.KeyId + 1

	preKeys, err := keyhelper.GeneratePreKeys(int(keyIdx), int(keyIdx)+generateCount-1, serialize.Proto.PreKeyRecord)
	if err != nil {
		panic("unable to generate pre keys!")
	}

	err = p.savePreKeys(preKeys)
	if err != nil {
		return nil, err
	}

	for i := range preKeys {
		pri := preKeys[i].KeyPair().PrivateKey().Serialize()

		newKeys[i].KeyID = preKeys[i].ID().Value
		newKeys[i].KeyPair = *keys.NewKeyPairFromPrivateKey(pri)
	}

	return newKeys, nil
}

func (p *PreKey) batchRemovePreKey(preKeyIDs []uint32) {
	if len(preKeyIDs) == 0 {
		return
	}

	where := prekeyDB.PreKey{JID: p.JID.User}

	err := prekeyDB.DeleteByKeyIDs(database.MasterDB(), &where, preKeyIDs)
	if err != nil {
		p.Logger.Error("batch remove pre keys ", err)
	}
}

// CleanupAllData .
func (p *PreKey) CleanupAllData() {
	context := prekeyDB.PreKey{JID: p.JID.User}

	_, _ = databaseTools.DeleteByPrimaryKey(database.MasterDB(), &context)
}

// Import .
func (p *PreKey) savePreKeys(preKeys []*record.PreKey) error {
	batch := make([]prekeyDB.PreKey, len(preKeys))

	for i := range preKeys {
		preKeyPair := record.NewPreKey(preKeys[i].ID().Value, preKeys[i].KeyPair(), serialize.Proto.PreKeyRecord)

		batch[i] = prekeyDB.PreKey{
			JID:     p.JID.User,
			KeyId:   preKeys[i].ID().Value,
			KeyBuff: preKeyPair.Serialize(),
		}
	}

	_, err := databaseTools.BatchCreate(database.MasterDB(), batch)

	return err
}

// Import .
func (p *PreKey) Import(preKeys []prekeyDB.PreKey) error {
	_, err := databaseTools.BatchCreate(database.MasterDB(), preKeys)

	return err
}
