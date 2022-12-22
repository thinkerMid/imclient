package privateChatNotification

import (
	"ws/framework/application/container/abstract_interface"
	privateChatCommon "ws/framework/application/core/action/private_chat/common"
	"ws/framework/application/core/processor"
	"ws/framework/utils"
	"ws/framework/utils/xmpp"
)

// DelayReceiveMessageMarkRead .
type DelayReceiveMessageMarkRead struct{}

// Receive .
func (r DelayReceiveMessageMarkRead) Receive(context containerInterface.IMessageContext) (err error) {
	info, parseInfoErr := xmpp.ParseMessageInfo(context.ResolveJID(), context.Message())
	// 						群消息不处理
	if parseInfoErr != nil || info.IsGroup {
		return
	}

	interval := utils.RandInt64(1, 3)

	context.ResolveMessageChannel().AddMessageProcessor(processor.NewTimerProcessor(
		func() []containerInterface.IAction {
			return []containerInterface.IAction{
				&privateChatCommon.ReceiveMessageMarkRead{UserID: info.Sender.User, MessageIDs: []string{info.ID}},
			}
		},
		processor.Interval(uint32(interval)),
		processor.IntervalLoop(false),
	))

	return
}
