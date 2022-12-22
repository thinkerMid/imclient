package groupComposeAction

import (
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/action/group"
	"ws/framework/application/core/processor"
)

// CheckAndQueryIcon 查询群组Icon
type CheckAndQueryIcon struct {
	processor.BaseAction
	GroupID string
}

// Start .
func (c *CheckAndQueryIcon) Start(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	defer next()

	g := context.ResolveGroupService().Find(c.GroupID)
	if g == nil || g.HaveGroupIcon {
		return nil
	}

	c.Query = &group.QueryIcon{GroupID: c.GroupID}
	err := c.Query.Start(context, func() {})
	if err != nil {
		return err
	}

	return nil
}

// Receive .
func (c *CheckAndQueryIcon) Receive(_ containerInterface.IMessageContext, _ containerInterface.NextActionFn) error {
	return nil
}

// Error .
func (c *CheckAndQueryIcon) Error(_ containerInterface.IMessageContext, _ error) {}
