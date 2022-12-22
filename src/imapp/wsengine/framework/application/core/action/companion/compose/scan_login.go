package compose

import (
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/action/companion"
	"ws/framework/application/core/action/user"
	"ws/framework/application/core/processor"
)

type CompanionScanLogin struct {
	processor.BaseAction
	query   containerInterface.IAction
	Content string
	jid     types.JID
}

func (m *CompanionScanLogin) ReceiveID() string {
	if m.query != nil {
		return m.query.ReceiveID()
	}

	return m.SendMessageId
}

func (m *CompanionScanLogin) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	//m.query = companion.MakeCompanionDevice(m.Content)
	m.query = &companion.CompanionRemove{}
	return m.query.Start(context, func() {})
}

func (m *CompanionScanLogin) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	//next()
	err := m.query.Receive(context, func() {})
	if err != nil {
		return err
	}

	switch m.query.(type) {
	case *companion.CompanionRemove:
		m.query = companion.MakeCompanionDevice(m.Content)
		return m.query.Start(context, func() {})
	case *companion.CompanionDevice:
		m.jid = context.VisitResult(0).IContent.(types.JID)

		m.query = &companion.CompanionPresence{JID: m.jid}
		m.query.Start(context, func() {})

		m.query = &companion.CompanionPairSession{JID: m.jid}
		return m.query.Start(context, func() {})
	case *companion.CompanionPresence:
		return nil
	case *companion.CompanionPairSession:
		m.query = &user.QueryMMSEndPoints{}
		return m.query.Start(context, func() {})
	case *user.QueryMMSEndPoints:
		next()
		return nil
	}

	return nil
}

func (m *CompanionScanLogin) Error(context containerInterface.IMessageContext, err error) {

}
