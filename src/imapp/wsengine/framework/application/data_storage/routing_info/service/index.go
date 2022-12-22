package routingInfoService

import (
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/data_storage/routing_info/database"
	"ws/framework/plugin/database"
	"ws/framework/plugin/database/database_tools"
)

var _ containerInterface.IRoutingInfoService = &RoutingInfo{}

// RoutingInfo .
type RoutingInfo struct {
	containerInterface.BaseService
	context *routingInfoDB.RoutingInfo
}

// Context .
func (r *RoutingInfo) Context() *routingInfoDB.RoutingInfo {
	if r.context != nil {
		return r.context
	}

	where := routingInfoDB.RoutingInfo{JID: r.JID.User}
	result := routingInfoDB.RoutingInfo{}

	err := databaseTools.Find(database.MasterDB(), &where, &result)
	if err != nil {
		r.Logger.Error(err)
		return nil
	}

	r.context = &result

	return r.context
}

// Save .
func (r *RoutingInfo) Save(content []byte) {
	r.Context().UpdateContent(content)

	_, err := databaseTools.Save(database.MasterDB(), r.Context())
	if err != nil {
		r.Logger.Error(err)
	}
}

// Create .
func (r *RoutingInfo) Create(content []byte) error {
	context := routingInfoDB.RoutingInfo{JID: r.JID.User, Content: content}

	_, err := databaseTools.Create(database.MasterDB(), &context)
	if err != nil {
		r.Logger.Error(err)
	}

	return err
}

// CleanupAllData .
func (r *RoutingInfo) CleanupAllData() {
	context := routingInfoDB.RoutingInfo{JID: r.JID.User}

	_, _ = databaseTools.DeleteByPrimaryKey(database.MasterDB(), &context)
}
