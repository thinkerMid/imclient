package scene

import (
	containerInterface "ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/action/private_chat/common"
	"ws/framework/application/core/processor"
)

// NewMessageReply .
func NewMessageReply(userID string) MessageReply {
	return MessageReply{UserID: userID}
}

// MessageReply .
type MessageReply struct {
	ActionList []containerInterface.IAction
	UserID     string
}

// Build .
func (c *MessageReply) Build() containerInterface.IProcessor {
	return processor.NewOnceProcessor(c.ActionList, processor.AliasName("messageReply"))
}

// MarkRead .
func (c *MessageReply) MarkRead(messageIDs []string) {
	c.ActionList = append(c.ActionList, &privateChatCommon.ReceiveMessageMarkRead{UserID: c.UserID, MessageIDs: messageIDs})
}
