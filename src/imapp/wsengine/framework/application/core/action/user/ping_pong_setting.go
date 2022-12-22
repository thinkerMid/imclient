package user

import (
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
)

// XMPPPingPongSetting xmpp协议心跳方式设置
type XMPPPingPongSetting struct {
	processor.BaseAction
	Passive bool // ping方式
}

// Start .
func (s *XMPPPingPongSetting) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	tag := "active" // 服务器主动ping
	if s.Passive {
		tag = "passive" // 客户端主动ping
	}

	s.SendMessageId, err = context.SendIQ(message.InfoQuery{
		ID:        context.GenerateRequestID(),
		Namespace: "passive",
		Type:      message.IqSet,
		To:        types.ServerJID,
		Content:   []waBinary.Node{{Tag: tag}},
	})

	return
}

// Receive .
func (s *XMPPPingPongSetting) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	next()
	return
}

// Error .
func (s *XMPPPingPongSetting) Error(context containerInterface.IMessageContext, err error) {}
