package models

import (
	"labs/src/imapp/wsa/types/ec25519"
	"labs/src/imapp/wsa/types/ecc"
	"labs/src/imapp/wsa/types/keys/identity"
)

type Identity struct {
	RegistrationId uint32
	KeyPair        *ec25519.KeyPair
}

func GenerateIdentity(device Device) Identity {
	var identityPriKey [32]byte
	copy(identityPriKey[:], device.IdentityKey)
	keyPair := ec25519.NewKeyPairFromPrivateKey(identityPriKey)

	return Identity{
		RegistrationId: device.RegistrationId,
		KeyPair:        keyPair,
	}
}

func GenerateIdentityKeyPair(keyPair *ec25519.KeyPair) *identity.KeyPair {
	return identity.NewKeyPair(
		identity.NewKey(ecc.NewDjbECPublicKey(*keyPair.Pub)),
		ecc.NewDjbECPrivateKey(*keyPair.Pri),
	)
}
