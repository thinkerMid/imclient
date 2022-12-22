package contactComposeAction

import (
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/action/common"
	"ws/framework/application/core/action/contact"
	"ws/framework/application/core/processor"
	"ws/framework/application/core/result/constant"
)

// BatchAddContactAndCreateSession 批量添加联系人并获取设备
type BatchAddContactAndCreateSession struct {
	processor.BaseAction
	UserIDs []string
}

// Start .
func (c *BatchAddContactAndCreateSession) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) error {
	c.Query = &contact.BatchAdd{UserIDs: c.UserIDs, SaveResult: true}

	return c.Query.Start(context, func() {})
}

// Receive .
func (c *BatchAddContactAndCreateSession) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	err := c.Query.Receive(context, func() {})
	if err != nil {
		return err
	}

	switch c.Query.(type) {
	// 批量添加结果
	case *contact.BatchAdd:
		result := context.VisitResult(0)

		batchResult := result.IContent.(contact.BatchAddContact)
		if len(batchResult.HaveKeyIndexNumber) == 0 {
			next()
			return nil
		}

		c.UserIDs = batchResult.HaveKeyIndexNumber
		c.Query = &common.QueryMultiDevicesIdentityBatch{UserIDs: c.UserIDs}

		return c.Query.Start(context, next)
	// 批量查询设备会话结果
	case *common.QueryMultiDevicesIdentityBatch:
		c.Query = &common.QueryUserDeviceListBatch{UserIDs: c.UserIDs}

		return c.Query.Start(context, next)
	case *common.QueryUserDeviceListBatch:
		next()
	}

	return nil
}

// Error .
func (c *BatchAddContactAndCreateSession) Error(context containerInterface.IMessageContext, err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.BatchAddContact,
		Error:      err,
	})
}
