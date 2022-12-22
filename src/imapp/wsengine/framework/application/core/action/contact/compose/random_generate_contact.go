package contactComposeAction

import (
	containerInterface "ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	phoneNumberGenerate "ws/framework/lib/phone_number_generate"
	"ws/framework/utils"
)

// UploadLocalAddressBook 上传随机生成的通讯录号码
type UploadLocalAddressBook struct {
	processor.BaseAction
	Min int64
	Max int64
}

// Start .
func (c *UploadLocalAddressBook) Start(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	ios2 := context.ResolveDeviceService().Context().Country

	generator := phoneNumberGenerate.AcquireGenerator(ios2)
	defer phoneNumberGenerate.ReleaseGenerator(generator)

	count := utils.RandInt64(c.Min, c.Max)

	userIDs, err := generator.GenerateMultipleNumber(count)
	if err != nil {
		return err
	}

	context.AddMessageProcessor(processor.NewOnceProcessor(
		[]containerInterface.IAction{
			&BatchAddContactAndCreateSession{UserIDs: userIDs},
		},
	))

	next()

	return nil
}

// Receive .
func (c *UploadLocalAddressBook) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	return nil
}

// Error .
func (c *UploadLocalAddressBook) Error(_ containerInterface.IMessageContext, _ error) {
}
