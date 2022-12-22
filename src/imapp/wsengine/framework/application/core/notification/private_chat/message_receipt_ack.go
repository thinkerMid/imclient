package privateChatNotification

import (
	"ws/framework/application/constant/binary"
	"ws/framework/application/container/abstract_interface"
	messageResultType "ws/framework/application/core/result/constant"
	"ws/framework/external"
)

// MessageReceiptAck .
type MessageReceiptAck struct{}

// Receive .
func (r MessageReceiptAck) Receive(context containerInterface.IMessageContext) (err error) {
	node := context.Message()
	attrGetter := node.AttrGetter()

	receiptType := attrGetter.String("type")
	messageID := attrGetter.String("id")
	from := attrGetter.JID("from")

	// 消息回复
	attrs := waBinary.Attrs{
		"class": node.Tag,
		"id":    messageID,
		"to":    from.String(),
	}

	if participant, ok := node.Attrs["participant"]; ok {
		attrs["participant"] = participant
	}
	if recipient, ok := node.Attrs["recipient"]; ok {
		attrs["recipient"] = recipient
	}
	if len(receiptType) > 0 {
		attrs["type"] = receiptType
	}

	_, err = context.SendNode(waBinary.Node{
		Tag:   "ack",
		Attrs: attrs,
	})

	r.pushResult(context, messageID, from.User, receiptType)

	// 存在多个消息的情况
	if val, ok := node.GetOptionalChildByTag("list"); ok {
		child := val.GetChildren()
		for i := range child {
			attrGetter = child[i].AttrGetter()
			messageID = attrGetter.String("id")

			r.pushResult(context, messageID, from.User, receiptType)
		}
	}

	return
}

func (r *MessageReceiptAck) pushResult(context containerInterface.IMessageContext, messageID, fromJID, receiptType string) {
	var messageResult uint8

	if receiptType == "read" {
		messageResult = messageResultType.SentMessageMarkRead
	} else {
		messageResult = messageResultType.SentMessageMarkDelivery
	}

	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResult,
		IContent: external.MessageStateChange{
			MessageID: messageID,
			JIDNumber: fromJID,
		},
	})
}
