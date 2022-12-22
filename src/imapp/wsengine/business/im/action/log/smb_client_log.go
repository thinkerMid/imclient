package log

import (
	registerRequest "ws/business/anonymous/register_request"
	containerInterface "ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
)

// SMBClientLog .
type SMBClientLog struct {
	Step, Sequence int
	processor.BaseAction
}

// Start .
func (c *SMBClientLog) Start(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	defer next()

	err := registerRequest.SendSMBClientLog(context, c.Step, c.Sequence)
	if err != nil {
		return err
	}

	return nil
}

// Receive .
func (c *SMBClientLog) Receive(_ containerInterface.IMessageContext, _ containerInterface.NextActionFn) error {
	return nil
}

// Error .
func (c *SMBClientLog) Error(_ containerInterface.IMessageContext, _ error) {}
