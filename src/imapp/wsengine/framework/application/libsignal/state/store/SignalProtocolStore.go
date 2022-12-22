package store

import (
	"ws/framework/application/libsignal/groups/state/store"
)

// SignalProtocol store is an interface that implements the
// methods for all stores needed in the Signal Protocol.
type SignalProtocol interface {
	IdentityKey
	PreKey
	ISessionStore
	SignedPreKey
	store.SenderKey
}
