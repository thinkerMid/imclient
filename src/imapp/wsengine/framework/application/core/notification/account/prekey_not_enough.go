package accountNotification

import (
	"ws/framework/application/constant"
	containerInterface "ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/action/user"
	"ws/framework/application/core/processor"
)

// PreKeyNotEnough .
type PreKeyNotEnough struct{}

// Receive .
func (m PreKeyNotEnough) Receive(context containerInterface.IMessageContext) (err error) {
	node := context.Message()

	ag := node.AttrGetter()
	if ag.String("type") != "encrypt" {
		return
	}

	countNode, ok := node.GetOptionalChildByTag("count")
	if !ok {
		return
	}

	remain := countNode.AttrGetter().Int("value")

	generateCount := context.ResolvePreKeyService().StatementPreKeyCount() - remain

	context.AddMessageProcessor(processor.NewOnceProcessor(
		[]containerInterface.IAction{
			&user.UploadPreKeyToServer{GenerateCount: generateCount},
		},
	))

	return constant.AbortedError
}
