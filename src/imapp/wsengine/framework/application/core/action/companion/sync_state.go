package companion

import (
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
)

type SyncMode int

const (
	SyncMode_DelAll SyncMode = iota
	SyncMode_Null
	SyncMode_CriticalLow
	SyncMode_RegularLow
)

type CompanionSyncState struct {
	processor.BaseAction
	SyncMode SyncMode
}

func (m *CompanionSyncState) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	var node *waBinary.Node
	if m.SyncMode == SyncMode_DelAll {
		node = m.generateDeleteAll()
	} else {
		node = m.generateSyncNode(m.SyncMode)
	}

	m.SendMessageId, err = context.SendIQ(
		message.InfoQuery{
			ID:        context.GenerateRequestID(),
			Namespace: "w:sync:app:state",
			Type:      message.IqSet,
			To:        types.ServerJID,
			Content: []waBinary.Node{
				*node,
			},
		},
	)
	return nil
}

func (m *CompanionSyncState) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	next()
	return nil
}

func (m *CompanionSyncState) Error(context containerInterface.IMessageContext, err error) {

}

func (m *CompanionSyncState) generateSyncNode(mode SyncMode) *waBinary.Node {
	switch mode {
	case SyncMode_Null:
		return nil
	}
	return &waBinary.Node{}
}

func (m *CompanionSyncState) generateDeleteAll() *waBinary.Node {
	return &waBinary.Node{
		Tag: "delete_all_data",
	}
}
