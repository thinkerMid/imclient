package user

import (
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	"ws/framework/application/core/result/constant"
)

// QueryStatusPrivacyList .
type QueryStatusPrivacyList struct {
	processor.BaseAction
	// 不等待结果，因为有些场景搭配了这个请求，这个查询比较特殊，异常也算是正常，
	//  为了不影响场景流程下 IgnoreResponse 设置为true时会跳过等待回包
	IgnoreResponse bool
}

// Start .
func (m *QueryStatusPrivacyList) Start(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	m.SendMessageId, err = context.SendIQ(message.InfoQuery{
		ID:        context.GenerateRequestID(),
		Namespace: "status",
		Type:      "get",
		To:        types.ServerJID,
		Content: []waBinary.Node{{
			Tag: "privacy",
		}},
	})

	if m.IgnoreResponse {
		next()
	}

	return
}

// Receive .
func (m *QueryStatusPrivacyList) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	next()

	return nil
}

// Error .
func (m *QueryStatusPrivacyList) Error(context containerInterface.IMessageContext, err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.QueryStatusPrivacyList,
		Error:      err,
	})
}
