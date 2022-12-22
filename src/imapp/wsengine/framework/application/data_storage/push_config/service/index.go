package pushConfigService

import (
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/data_storage/push_config/database"
	"ws/framework/plugin/database"
	"ws/framework/plugin/database/database_tools"
)

var _ containerInterface.IPushConfigService = &PushConfig{}

// PushConfig .
type PushConfig struct {
	containerInterface.BaseService
}

// Context .
func (a *PushConfig) Context() *pushConfigDB.PushConfig {
	result := pushConfigDB.PushConfig{JID: a.JID.User}

	err := databaseTools.FindByPrimaryKey(database.MasterDB(), &result)
	if err != nil {
		a.Logger.Error(err)
		return nil
	}

	return &result
}

// Create .
func (a *PushConfig) Create() (*pushConfigDB.PushConfig, error) {
	context := pushConfigDB.NewPushConfig(a.JID.User)

	_, err := databaseTools.Create(database.MasterDB(), context)
	if err != nil {
		a.Logger.Error(err)
	}

	return context, err
}

// Import .
func (a *PushConfig) Import(pushConfig *pushConfigDB.PushConfig) error {
	_, err := databaseTools.Create(database.MasterDB(), pushConfig)
	if err != nil {
		a.Logger.Error(err)
	}

	return err
}

// ContextExecute .
func (a *PushConfig) ContextExecute(f func(*pushConfigDB.PushConfig)) {
	p := a.Context()
	if p == nil {
		return
	}

	f(p)

	_, err := databaseTools.Save(database.MasterDB(), p)
	if err != nil {
		a.Logger.Error(err)
	}
}

// CleanupAllData .
func (a *PushConfig) CleanupAllData() {
	context := pushConfigDB.PushConfig{JID: a.JID.User}

	_, _ = databaseTools.DeleteByPrimaryKey(database.MasterDB(), &context)
}
