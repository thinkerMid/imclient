package containerInterface

import (
	"crypto/cipher"
	"github.com/cloudwego/netpoll"
	"ws/framework/application/constant/binary"
)

// ConnectionEventListener .
type ConnectionEventListener interface {
	OnConnectionClose()
	OnResponse(node *waBinary.Node)
}

// IConnection .
type IConnection interface {
	Connect() error
	Write(node waBinary.Node) error
	Close()
	IsClosed() bool
	SetEventListener(ConnectionEventListener)
}

// IHandshakeHandler .
type IHandshakeHandler interface {
	Do(netpoll.Connection) (cipher.AEAD, cipher.AEAD, error)
	SetEdgeRouting([]byte)
}
