package abKeyService

import (
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/data_storage/ab_key/database"
	"ws/framework/plugin/database"
	"ws/framework/plugin/database/database_tools"
)

var _ containerInterface.IABKeyService = &ABKey{}

// ABKey .
type ABKey struct {
	containerInterface.BaseService
}

// Context .
func (a *ABKey) Context() *abKeyDB.QueryResult {
	where := abKeyDB.ABKey{JID: a.JID.User}
	result := abKeyDB.QueryResult{}
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
func (a *ABKey) Create(content string) error {
	context := abKeyDB.ABKey{JID: a.JID.User, Content: content}

	_, err := databaseTools.Create(database.MasterDB(), &context)
	if err != nil {
		a.Logger.Error(err)
	}

	return err
}

// CleanupAllData .
func (a *ABKey) CleanupAllData() {
	context := abKeyDB.ABKey{JID: a.JID.User}

	_, _ = databaseTools.DeleteByPrimaryKey(database.MasterDB(), &context)
}
