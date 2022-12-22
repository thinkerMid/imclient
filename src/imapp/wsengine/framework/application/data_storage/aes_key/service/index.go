package aesKeyService

import (
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/data_storage/aes_key/database"
	"ws/framework/plugin/database"
	"ws/framework/plugin/database/database_tools"
	"ws/framework/utils"
)

var _ containerInterface.IAesKeyService = &AesKey{}

// AesKey .
type AesKey struct {
	containerInterface.BaseService
}

// Context .
func (a *AesKey) Context() *aesKeyDB.QueryResult {
	where := aesKeyDB.AESKey{JID: a.JID.User}
	result := aesKeyDB.QueryResult{}
	err := databaseTools.Find(database.MasterDB(), &where, &result)
	if err != nil {
		// TODO 兼容旧数据 一个月后或者一个版本后可以删除了
		where.JID = a.JID.String()
		err = databaseTools.Find(database.MasterDB(), &where, &result)
	}

	if err != nil {
		a.Logger.Error(err)
		return nil
	}

	return &result
}

// Create .
func (a *AesKey) Create() (*aesKeyDB.AESKey, error) {
	context := aesKeyDB.AESKey{JID: a.JID.User}

	configuration := a.AppIocContainer.ResolveWhatsappConfiguration()

	context.PubKey, context.PriKey = utils.GenCurve25519KeyPair()
	context.AesKey = utils.CalCurve25519Signature(context.PriKey, configuration.AESCurve25519PublicKey)

	_, err := databaseTools.Create(database.MasterDB(), &context)
	if err != nil {
		a.Logger.Error(err)
	}

	return &context, err
}

// OnJIDChangeWhenRegisterSuccess .
func (a *AesKey) OnJIDChangeWhenRegisterSuccess(newJID types.JID) {
	context := aesKeyDB.AESKey{JID: a.JID.User}
	context.UpdateJID(newJID.User)

	_, err := databaseTools.Save(database.MasterDB(), &context)
	if err != nil {
		a.Logger.Error(err)
	}
}

// CleanupAllData .
func (a *AesKey) CleanupAllData() {
	context := aesKeyDB.AESKey{JID: a.JID.User}

	_, _ = databaseTools.DeleteByPrimaryKey(database.MasterDB(), &context)
}
