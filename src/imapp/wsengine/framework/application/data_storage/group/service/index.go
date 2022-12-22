package groupService

import (
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/data_storage/group/database"
	"ws/framework/plugin/database"
	"ws/framework/plugin/database/database_tools"
)

var _ containerInterface.IGroupService = &Group{}

// Group .
type Group struct {
	containerInterface.BaseService
}

// CreateGroup .
func (g *Group) CreateGroup(groupID string, isAdmin bool) {
	group := groupDB.Group{
		JID:     g.JID.User,
		GroupID: groupID,
		IsAdmin: isAdmin,
	}

	_, err := databaseTools.Create(database.MasterDB(), &group)
	if err != nil {
		g.Logger.Error(err)
	}

	g.AppIocContainer.ResolveMemoryCache().Cache(group.GroupID, &group)
}

// DeleteGroup .
func (g *Group) DeleteGroup(groupID string) {
	group := groupDB.Group{JID: g.JID.User, GroupID: groupID}

	_, err := databaseTools.DeleteByPrimaryKey(database.MasterDB(), &group)
	if err != nil {
		g.Logger.Error(err)
	}

	g.AppIocContainer.ResolveMemoryCache().UnCache(groupID)
}

// Find .
func (g *Group) Find(groupID string) *groupDB.Group {
	cacheGroup, ok := g.AppIocContainer.ResolveMemoryCache().FindInCache(groupID)
	if ok {
		return cacheGroup.(*groupDB.Group)
	}

	group := groupDB.Group{JID: g.JID.User, GroupID: groupID}

	err := databaseTools.FindByPrimaryKey(database.MasterDB(), &group)
	if err == nil {
		g.AppIocContainer.ResolveMemoryCache().Cache(groupID, &group)
		return &group
	}

	return nil
}

// ContextExecute .
func (g *Group) ContextExecute(groupID string, f func(info *groupDB.Group)) {
	group := g.Find(groupID)
	if group == nil {
		return
	}

	f(group)

	_, err := databaseTools.Save(database.MasterDB(), group)
	if err != nil {
		g.Logger.Error(err)
	}

	g.AppIocContainer.ResolveMemoryCache().Cache(groupID, group)
}

// CleanupAllData .
func (g *Group) CleanupAllData() {
	_, _ = groupDB.DeleteByJID(database.MasterDB(), g.JID.User)
}
