package store

import (
	"ws/framework/application/libsignal/keys/identity"
)

// IdentityKey provides an interface to identity information.
type IdentityKey interface {
	// Get the local client's identity key pair.
	GetIdentityKeyPair() *identity.KeyPair

	// Return the local client's registration ID.
	//
	// Clients should maintain a registration ID, a random number between 1 and 16380
	// that's generated once at install time.
	GetLocalRegistrationId() uint32
}
