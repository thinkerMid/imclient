package log

import (
	registerAccount "ws/business/anonymous/register_request"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
)

// ClientLog .
type ClientLog struct {
	CurrentScreen, PreviousScreen, ActionTaken string
	processor.BaseAction
}

// Start .
func (c *ClientLog) Start(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	defer next()

	err := registerAccount.SendClientLog(context, c.CurrentScreen, c.PreviousScreen, c.ActionTaken)
	if err != nil {
		return err
	}

	return nil
}

// Receive .
func (c *ClientLog) Receive(_ containerInterface.IMessageContext, _ containerInterface.NextActionFn) error {
	return nil
}

// Error .
func (c *ClientLog) Error(_ containerInterface.IMessageContext, _ error) {}
