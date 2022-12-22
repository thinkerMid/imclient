package record

import (
	ecc "ws/framework/application/libsignal/ecc"
	"ws/framework/application/libsignal/groups/ratchet"
)

// NewSenderKeyState returns a new SenderKeyState.
func NewSenderKeyState(keyID uint32, iteration uint32, chainKey []byte, signatureKey *ecc.ECKeyPair) *SenderKeyState {
	return &SenderKeyState{
		keyID:          keyID,
		senderChainKey: ratchet.NewSenderChainKey(iteration, chainKey),
		signingKeyPair: signatureKey,
	}
}

// NewSenderKeyStateFromPublicKey returns a new SenderKeyState with the given publicKey.
func NewSenderKeyStateFromPublicKey(keyID uint32, iteration uint32, chainKey []byte, signatureKey ecc.ECPublicKeyable) *SenderKeyState {
	keyPair := ecc.NewECKeyPair(signatureKey, nil)

	return &SenderKeyState{
		keyID:          keyID,
		senderChainKey: ratchet.NewSenderChainKey(iteration, chainKey),
		signingKeyPair: keyPair,
	}
}

// SenderKeyState is a structure for maintaining a senderkey session state.
type SenderKeyState struct {
	keyID          uint32
	senderChainKey *ratchet.SenderChainKey
	signingKeyPair *ecc.ECKeyPair
}

// SigningKey returns the signing key pair of the sender key state.
func (k *SenderKeyState) SigningKey() *ecc.ECKeyPair {
	return k.signingKeyPair
}

// SenderChainKey returns the sender chain key of the state.
func (k *SenderKeyState) SenderChainKey() *ratchet.SenderChainKey {
	return k.senderChainKey
}

// KeyID returns the state's key id.
func (k *SenderKeyState) KeyID() uint32 {
	return k.keyID
}

// SetSenderChainKey will set the state's sender chain key with the given key.
func (k *SenderKeyState) SetSenderChainKey(senderChainKey *ratchet.SenderChainKey) {
	k.senderChainKey = senderChainKey
}
