package groupComposeAction

import (
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/action/common"
	"ws/framework/application/core/action/group"
	"ws/framework/application/core/processor"
)

// CreateGroup 创建群组
type CreateGroup struct {
	processor.BaseAction

	GroupName   string
	Icon        []byte
	JoinUserIDs []string

	groupID string
}

// Start .
func (c *CreateGroup) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) error {
	selfJId := context.ResolveJID()

	if len(c.JoinUserIDs) > 0 {
		for _, u := range c.JoinUserIDs {
			if u != selfJId.User {
				// TODO: 本地有查过则不查，不考虑其头像更新
				c.Query = &common.QueryAvatarPreview{UserID: u}
				_ = c.Query.Start(context, func() {})
			}
		}
	}

	c.Query = &group.Create{
		GroupName:   c.GroupName,
		JoinUserIDs: c.JoinUserIDs,
		HaveIcon:    len(c.Icon) > 0,
	}

	return c.Query.Start(context, func() {})
}

// Receive .
func (c *CreateGroup) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	err := c.Query.Receive(context, func() {})
	if err != nil {
		return err
	}

	switch c.Query.(type) {
	// 创建结果
	case *group.Create:
		result := context.VisitResult(0)
		c.groupID = result.Content

		c.Query = &group.QueryMemberDeviceList{UserIDs: c.JoinUserIDs}
		return c.Query.Start(context, func() {})
	// 查询群成员设备列表结果
	case *group.QueryMemberDeviceList:
		c.Query = &common.QueryMultiDevicesIdentityBatch{UserIDs: c.JoinUserIDs}
		return c.Query.Start(context, func() {})
	}

	return nil
}

// Error .
func (c *CreateGroup) Error(context containerInterface.IMessageContext, err error) {
	if c.Query != nil {
		c.Query.Error(context, err)
	}
}
