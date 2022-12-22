package common

import (
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	"ws/framework/application/core/result/constant"
)

// QueryUserStatus .
type QueryUserStatus struct {
	processor.BaseAction
	UserID string
}

// Start .
func (m *QueryUserStatus) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	userJID := types.NewJID(m.UserID, types.DefaultUserServer)

	m.SendMessageId, err = context.SendIQ(message.InfoQuery{
		ID:        context.GenerateRequestID(),
		Namespace: "status",
		Type:      "get",
		To:        types.ServerJID,
		Content: []waBinary.Node{{
			Tag: "status",
			Content: []waBinary.Node{{
				Tag: "user",
				Attrs: waBinary.Attrs{
					"jid": userJID.String(),
				},
			}},
		}},
	})

	return
}

// Receive .
func (m *QueryUserStatus) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	next()

	//nodes := context.Message().GetChildrenByTag("status")
	//
	//if len(nodes) == 0 {
	//	context.AppendResult(containerInterface.MessageResult{
	//		ResultType: messageResultType.GetSignature,
	//		Content:    "查询用户签名失败",
	//	})
	//
	//	return nil
	//}
	//
	//nodes = nodes[0].GetChildren()
	//
	//for i := range nodes {
	//	n := nodes[i]
	//	if n.Tag != "user" {
	//		continue
	//	}
	//
	//	_, ok := n.Attrs["jid"].(types.UserID)
	//	if !ok {
	//		continue
	//	}
	//
	//	bStatus, _ := n.Content.([]byte)
	//
	//	context.AppendResult(containerInterface.MessageResult{
	//		ResultType: messageResultType.GetSignature,
	//		Content:    string(bStatus),
	//	})
	//}

	return nil
}

// Error .
func (m *QueryUserStatus) Error(context containerInterface.IMessageContext, err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.GetSignature,
		Error:      err,
	})
}
