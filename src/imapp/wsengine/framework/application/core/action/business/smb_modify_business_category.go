package business

import (
	waBinary "ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	containerInterface "ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
)

// SMBModifyBusinessCategory .
type SMBModifyBusinessCategory struct {
	processor.BaseAction
	CategoryID   string
	CategoryName string
}

// Start .
func (m *SMBModifyBusinessCategory) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	m.SendMessageId, err = context.SendIQ(message.InfoQuery{
		ID:        context.GenerateRequestID(),
		Namespace: "w:biz",
		Type:      message.IqSet,
		To:        types.ServerJID,
		Content: []waBinary.Node{{
			Tag: "business_profile",
			Attrs: waBinary.Attrs{
				"mutation_type": "delta",
				"v":             "372",
			},
			Content: []waBinary.Node{{
				Tag: "categories",
				Content: []waBinary.Node{{
					Tag: "category",
					Attrs: waBinary.Attrs{
						"id": m.CategoryID,
					},
					Content: m.CategoryName,
				}},
			}},
		}},
	})

	return
}

// Receive .
func (m *SMBModifyBusinessCategory) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	next()

	return nil
}

// Error .
func (m *SMBModifyBusinessCategory) Error(context containerInterface.IMessageContext, err error) {
}
