package common

import (
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/types"
	containerInterface "ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	"ws/framework/utils/xmpp"
)

// QueryDisappearingMode 似乎是用于查询双方聊天消息消失时限的设置
type QueryDisappearingMode struct {
	processor.BaseAction
	UserID string
}

// Start .
func (q *QueryDisappearingMode) Start(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	jid := types.NewJID(q.UserID, types.DefaultUserServer)

	//<iq id="1663066499-8" to="84564844255@s.whatsapp.net" type="get" xmlns="usync">
	//<usync context="interactive" index="0" last="true" mode="query" sid="1663066797-2851852557-2">
	//<query><disappearing_mode/></query>
	//<list><user jid="85268067387@s.whatsapp.net"/></list>
	//</usync>
	//</iq>
	q.SendMessageId, err = context.SendIQ(
		xmpp.UsyncIQTemplate(context, "query", "interactive",
			[]waBinary.Node{
				{Tag: "query", Content: []waBinary.Node{
					{Tag: "disappearing_mode"},
				}},
				{Tag: "list", Content: []waBinary.Node{
					{Tag: "user", Attrs: waBinary.Attrs{"jid": jid.String()}},
				}},
			},
		),
	)

	return
}

// Receive .
func (q *QueryDisappearingMode) Receive(_ containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	next()
	return
}

// Error .
func (q *QueryDisappearingMode) Error(_ containerInterface.IMessageContext, _ error) {
}
