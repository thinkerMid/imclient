package signedPreKeyService

import (
	"bytes"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/data_storage/signed_prekey/database"
	"ws/framework/application/libsignal/ecc"
	"ws/framework/application/libsignal/keys/identity"
	"ws/framework/application/libsignal/serialize"
	"ws/framework/application/libsignal/state/record"
	"ws/framework/application/libsignal/util/keyhelper"
	"ws/framework/plugin/database"
	"ws/framework/plugin/database/database_tools"
	"ws/framework/utils/keys"
)

var _ containerInterface.ISignedPreKeyService = &SignedPreKey{}

var privateSignKey = []byte("PrivateSignKey")
var privateKey = []byte("PrivateKey")

// SignedPreKey .
type SignedPreKey struct {
	containerInterface.BaseService
}

// Context .
func (s *SignedPreKey) Context() *record.SignedPreKey {
	where := signedPreKeyDB.SignedPreKey{JID: s.JID.User}
	result := signedPreKeyDB.QueryResult{}
	err := databaseTools.Find(database.MasterDB(), &where, &result)
	if err != nil {
		// TODO 兼容旧数据 一个月后或者一个版本后可以删除了
		where.JID = s.JID.String()
		err = databaseTools.Find(database.MasterDB(), &where, &result)
	}

	if err != nil {
		s.Logger.Error(err)
		return nil
	}

	if bytes.Contains(result.KeyBuff, privateSignKey) {
		result.KeyBuff = bytes.Replace(result.KeyBuff, privateSignKey, privateKey, 1)
	}

	keyPair, err := record.NewSignedPreKeyFromBytes(result.KeyBuff, serialize.Proto.SignedPreKeyRecord)
	if err != nil {
		s.Logger.Error(err)
		return nil
	}

	return keyPair
}

// Create .
func (s *SignedPreKey) Create() (*signedPreKeyDB.SignedPreKey, error) {
	device := s.AppIocContainer.ResolveDeviceService().Context()

	var identityPriKey [32]byte
	copy(identityPriKey[:], device.IdentityKey)
	identityKey := keys.NewKeyPairFromPrivateKey(identityPriKey)

	identityKeyPair := identity.NewKeyPair(identity.NewKey(ecc.NewDjbECPublicKey(*identityKey.Pub)), ecc.NewDjbECPrivateKey(*identityKey.Priv))

	keyPair, err := keyhelper.GenerateSignedPreKey(identityKeyPair, keyhelper.GenerateSignedPreKeyID(), serialize.Proto.SignedPreKeyRecord)
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	context := signedPreKeyDB.SignedPreKey{JID: s.JID.User}
	context.KeyId = keyPair.ID()
	context.KeyBuff = keyPair.Serialize()

	_, err = databaseTools.Create(database.MasterDB(), &context)
	if err != nil {
		s.Logger.Error(err)
	}

	return &context, err
}

// FindSignedPreKey .
func (s *SignedPreKey) FindSignedPreKey(signedPreKeyID uint32) *record.SignedPreKey {
	where := signedPreKeyDB.SignedPreKey{KeyId: signedPreKeyID, JID: s.JID.User}
	result := signedPreKeyDB.QueryResult{}

	err := databaseTools.Find(database.MasterDB(), &where, &result)
	if err != nil {
		//TODO 兼容旧数据 一个月后或者一个版本后可以删除了
		where.JID = s.JID.String()
		err = databaseTools.Find(database.MasterDB(), &where, &result)
	}

	if err != nil {
		return nil
	}

	if bytes.Contains(result.KeyBuff, privateSignKey) {
		result.KeyBuff = bytes.Replace(result.KeyBuff, privateSignKey, privateKey, 1)
	}

	signedPreKey, err := record.NewSignedPreKeyFromBytes(result.KeyBuff, serialize.Proto.SignedPreKeyRecord)
	if err != nil {
		s.Logger.Error(err)
		return nil
	}

	return signedPreKey
}

// SaveSignedPreKeyBuffer .
func (s *SignedPreKey) SaveSignedPreKeyBuffer(signedPreKeyID uint32, signedPreKeyJsonBuffer []byte) {
	context := signedPreKeyDB.SignedPreKey{JID: s.JID.User}
	context.KeyId = signedPreKeyID
	context.KeyBuff = signedPreKeyJsonBuffer

	_, err := databaseTools.Create(database.MasterDB(), &context)
	if err != nil {
		s.Logger.Error(err)
	}
}

// OnJIDChangeWhenRegisterSuccess .
func (s *SignedPreKey) OnJIDChangeWhenRegisterSuccess(newJID types.JID) {
	context := signedPreKeyDB.SignedPreKey{JID: s.JID.User}
	context.UpdateJID(newJID.User)

	_, err := databaseTools.Save(database.MasterDB(), &context)
	if err != nil {
		s.Logger.Error(err)
	}
}

// CleanupAllData .
func (s *SignedPreKey) CleanupAllData() {
	context := signedPreKeyDB.SignedPreKey{JID: s.JID.User}

	_, _ = databaseTools.DeleteByPrimaryKey(database.MasterDB(), &context)
}
