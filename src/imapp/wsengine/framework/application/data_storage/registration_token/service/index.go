package registrationTokenService

import (
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/data_storage/registration_token/database"
	"ws/framework/plugin/database"
	"ws/framework/plugin/database/database_tools"
	"ws/framework/utils"
)

var _ containerInterface.IRegistrationTokenService = &RegistrationToken{}

// RegistrationToken .
type RegistrationToken struct {
	containerInterface.BaseService
}

// Context .
func (r *RegistrationToken) Context() *registrationTokenDB.QueryResult {
	where := registrationTokenDB.RegistrationToken{JID: r.JID.User}
	result := registrationTokenDB.QueryResult{}
	err := databaseTools.Find(database.MasterDB(), &where, &result)
	if err != nil {
		// TODO 兼容旧数据 一个月后或者一个版本后可以删除了
		where.JID = r.JID.String()
		err = databaseTools.Find(database.MasterDB(), &where, &result)
	}

	if err != nil {
		r.Logger.Error(err)
		return nil
	}

	return &result
}

// Create .
func (r *RegistrationToken) Create() (*registrationTokenDB.RegistrationToken, error) {
	context := registrationTokenDB.RegistrationToken{JID: r.JID.User}

	context.RecoveryToken = utils.RandBytes(16)
	context.BackupToken = utils.RandBytes(20)
	context.BackupKey = make([]byte, 0)
	context.BackupKey2 = make([]byte, 0)

	_, err := databaseTools.Create(database.MasterDB(), &context)
	if err != nil {
		r.Logger.Error(err)
	}

	return &context, err
}

// Import .
func (r *RegistrationToken) Import(reverToken *registrationTokenDB.RegistrationToken) error {
	_, err := databaseTools.Create(database.MasterDB(), reverToken)
	if err != nil {
		r.Logger.Error(err)
	}

	return err
}

// RefreshToken .
func (r *RegistrationToken) RefreshToken() error {
	where := registrationTokenDB.RegistrationToken{JID: r.JID.User}
	context := registrationTokenDB.RegistrationToken{}

	context.RecoveryToken = utils.RandBytes(16)
	context.BackupToken = utils.RandBytes(20)
	context.BackupKey = utils.CreateBackupKey()
	context.BackupKey2 = utils.CreateBackupKey()

	_, err := databaseTools.CreateOrSave(database.MasterDB(), &where, &context)
	if err != nil {
		r.Logger.Error(err)
	}

	return err
}

// OnJIDChangeWhenRegisterSuccess .
func (r *RegistrationToken) OnJIDChangeWhenRegisterSuccess(newJID types.JID) {
	context := registrationTokenDB.RegistrationToken{JID: r.JID.User}
	context.UpdateJID(newJID.User)

	_, err := databaseTools.Save(database.MasterDB(), &context)
	if err != nil {
		r.Logger.Error(err)
	}
}

// CleanupAllData .
func (r *RegistrationToken) CleanupAllData() {
	context := registrationTokenDB.RegistrationToken{JID: r.JID.User}

	_, _ = databaseTools.DeleteByPrimaryKey(database.MasterDB(), &context)
}
