package store

import (
	"ws/framework/application/libsignal/protocol"
	"ws/framework/application/libsignal/state/record"
)

// ISessionStore store is an interface for the persistent storage of session
// state information for remote clients.
type ISessionStore interface {
	CreateSession(address *protocol.SignalAddress, record *record.Session)
	FindSession(address *protocol.SignalAddress) (*record.Session, error)
	SaveEncryptSession(remoteAddress *protocol.SignalAddress, record *record.Session)
	SaveDecryptSession(remoteAddress *protocol.SignalAddress, record *record.Session)
	SaveRebuildSession(address *protocol.SignalAddress, record *record.Session)
	SaveSession(address *protocol.SignalAddress, record *record.Session)
}
