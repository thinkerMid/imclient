package common

import (
	"fmt"
	"time"
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
)

// SubscribeStatus .
type SubscribeStatus struct {
	processor.BaseAction
	UserID    string
	TcToToken string
	ToGroup   bool
}

/**
[presence type=<'subscribe'> to=<8618898739725@s.whatsapp.net> ]
[presence from=<8618898739725@s.whatsapp.net> last=<10 bytes> type=<'unavailable'> ]
<presence to="85268067387@s.whatsapp.net" type="subscribe"><tctoken>01010e64ff3fb6a728023c62b89e203088dcd43c8442c8fc3682a619025ca890807a51</tctoken></presence>
*/

// Start .
func (m *SubscribeStatus) Start(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	defer next()

	// 如果有订阅过的标记 不发
	key := fmt.Sprintf("%s_subscribe_%s", context.ResolveJID().User, m.UserID)
	_, ok := context.ResolveMemoryCache().FindInCache(key)
	if ok {
		return
	}

	toServer := types.DefaultUserServer
	if m.ToGroup {
		toServer = types.GroupServer
	}

	node := waBinary.Node{
		Tag: "presence",
		Attrs: waBinary.Attrs{
			"type": "subscribe",
			"to":   fmt.Sprintf("%s@%s", m.UserID, toServer),
		},
	}

	// 收到privacy_token的通知会有trusted_contact的token内容
	if len(m.TcToToken) > 0 {
		node.Content = []waBinary.Node{
			{Tag: "tctoken", Content: m.TcToToken},
		}
	}

	m.SendMessageId, err = context.SendNode(node)

	context.ResolveMemoryCache().CacheTTL(key, struct{}{}, time.Hour)

	return
}

// Receive .
func (m *SubscribeStatus) Receive(_ containerInterface.IMessageContext, _ containerInterface.NextActionFn) error {
	return nil
}

// Error .
func (m *SubscribeStatus) Error(_ containerInterface.IMessageContext, _ error) {}
