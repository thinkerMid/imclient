package store

import (
	"ws/framework/application/libsignal/state/record"
)

// SignedPreKey store is an interface that describes how to persistently
// store signed PreKeys.
type SignedPreKey interface {
	// LoadSignedPreKey loads a local SignedPreKeyRecord
	FindSignedPreKey(signedPreKeyID uint32) *record.SignedPreKey
}
