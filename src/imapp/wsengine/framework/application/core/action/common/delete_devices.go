package common

import (
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	"ws/framework/utils/xmpp"
)

// DeleteDevices
//
//	不是删除自己的设备，是删除联系人过程所需要的一个步骤
type DeleteDevices struct {
	processor.BaseAction
	UserID string // 手机号
}

// Start .
func (m *DeleteDevices) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	jid := types.NewJID(m.UserID, types.DefaultUserServer)

	m.SendMessageId, err = context.SendIQ(
		xmpp.UsyncIQTemplate(context, "delta", "interactive",
			[]waBinary.Node{
				{Tag: "query", Content: []waBinary.Node{
					{Tag: "business", Content: []waBinary.Node{
						{Tag: "verified_name"},
						{Tag: "profile",
							Attrs: waBinary.Attrs{
								"v": "116",
							},
						},
					}},
					{Tag: "contact"},
					{Tag: "devices", Attrs: waBinary.Attrs{"version": "2"}},
					{Tag: "disappearing_mode"},
					{Tag: "sidelist"},
					{Tag: "status"},
				}},
				{Tag: "list"},
				{Tag: "side_list", Content: []waBinary.Node{{
					Tag:   "user",
					Attrs: waBinary.Attrs{"jid": jid.String()},
					Content: []waBinary.Node{{Tag: "devices",
						Attrs: waBinary.Attrs{"device_hash": xmpp.GenerateDeviceHash(jid)},
					}},
				}}},
			},
		),
	)

	return
}

// Receive .
func (m *DeleteDevices) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	defer next()

	return nil
}

func (m *DeleteDevices) Error(context containerInterface.IMessageContext, err error) {
	context.AppendResult(containerInterface.MessageResult{
		Error: err,
	})
}
