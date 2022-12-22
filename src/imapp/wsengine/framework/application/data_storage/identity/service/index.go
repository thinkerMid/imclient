package identityService

import (
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/libsignal/ecc"
	"ws/framework/application/libsignal/keys/identity"
	"ws/framework/utils/keys"
)

var _ containerInterface.IIdentityService = &Identity{}

// Identity .
type Identity struct {
	containerInterface.BaseService

	registrationId uint32

	keyPair *keys.KeyPair
}

// Context .
func (m *Identity) Context() *keys.KeyPair {
	if m.keyPair != nil {
		return m.keyPair
	}

	device := m.AppIocContainer.ResolveDeviceService().Context()

	var identityPriKey [32]byte
	copy(identityPriKey[:], device.IdentityKey)
	m.keyPair = keys.NewKeyPairFromPrivateKey(identityPriKey)

	m.registrationId = device.RegistrationId

	return m.keyPair
}

// GetIdentityKeyPair .
func (m *Identity) GetIdentityKeyPair() *identity.KeyPair {
	return identity.NewKeyPair(
		identity.NewKey(ecc.NewDjbECPublicKey(*m.Context().Pub)),
		ecc.NewDjbECPrivateKey(*m.Context().Priv),
	)
}

// GetLocalRegistrationId .
func (m *Identity) GetLocalRegistrationId() uint32 {
	return m.registrationId
}
