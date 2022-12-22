package record

import (
	"fmt"
	ecc "ws/framework/application/libsignal/ecc"
	"ws/framework/application/libsignal/groups/ratchet"
	"ws/framework/application/libsignal/signalerror"
)

// SenderKeyStructure is a serializeable structure of SenderKeyState.
type SenderKeyStructure struct {
	KeyID             uint32
	Iteration         uint32
	ChainKey          []byte
	SigningKeyPrivate [32]byte
	SigningKeyPublic  [32]byte
}

// NewSenderKey record returns a new sender key record that can
// be stored in a SenderKeyStore.
func NewSenderKey() *SenderKey {
	return &SenderKey{
		senderKeyStates: &SenderKeyState{},
	}
}

// NewSenderKeyFromStructure will return a new session state with the
// given state structure. This structure is given back from an
// implementation of the sender key state serializer.
func NewSenderKeyFromStructure(structure *SenderKeyStructure) (*SenderKey, error) {
	// Convert our ecc keys from bytes into object form.
	signingKeyPublic := ecc.NewDjbECPublicKey(structure.SigningKeyPublic)
	signingKeyPrivate := ecc.NewDjbECPrivateKey(structure.SigningKeyPrivate)

	// Build our state object.
	state := &SenderKeyState{
		keyID:          structure.KeyID,
		senderChainKey: ratchet.NewSenderChainKey(structure.Iteration, structure.ChainKey),
		signingKeyPair: ecc.NewECKeyPair(signingKeyPublic, signingKeyPrivate),
	}

	return &SenderKey{senderKeyStates: state}, nil
}

// SenderKey record is a structure for storing pre keys inside
// a SenderKeyStore.
type SenderKey struct {
	senderKeyStates *SenderKeyState
}

// SenderKeyState will return the first sender key state in the record's
// list of sender key states.
func (k *SenderKey) SenderKeyState() (*SenderKeyState, error) {
	if k.senderKeyStates != nil {
		return k.senderKeyStates, nil
	}

	return nil, signalerror.ErrNoSenderKeyStatesInRecord
}

// GetSenderKeyStateByID will return the sender key state with the given
// key id.
func (k *SenderKey) GetSenderKeyStateByID(keyID uint32) (*SenderKeyState, error) {
	if k.senderKeyStates.KeyID() == keyID {
		return k.senderKeyStates, nil
	}

	return nil, fmt.Errorf("%w %d", signalerror.ErrNoSenderKeyStateForID, keyID)
}

// IsEmpty will return false if there is more than one state in this
// senderkey record.
func (k *SenderKey) IsEmpty() bool {
	return k.senderKeyStates == nil
}

// AddSenderKeyState will add a new state to this senderkey record with the given
// id, iteration, chainkey, and signature key.
func (k *SenderKey) AddSenderKeyState(id uint32, iteration uint32,
	chainKey []byte, signatureKey ecc.ECPublicKeyable) {

	newState := NewSenderKeyStateFromPublicKey(id, iteration, chainKey, signatureKey)
	k.senderKeyStates = newState
}

// SetSenderKeyState will  replace the current senderkey states with the given
// senderkey state.
func (k *SenderKey) SetSenderKeyState(id uint32, iteration uint32,
	chainKey []byte, signatureKey *ecc.ECKeyPair) {

	newState := NewSenderKeyState(id, iteration, chainKey, signatureKey)
	k.senderKeyStates = newState
}

// Pack .
func (k *SenderKey) Pack() *SenderKeyStructure {
	// Build and return our state structure.
	s := &SenderKeyStructure{
		KeyID:     k.senderKeyStates.keyID,
		Iteration: k.senderKeyStates.senderChainKey.Iteration(),
		ChainKey:  k.senderKeyStates.senderChainKey.Seed(),
	}

	s.SigningKeyPublic = k.senderKeyStates.signingKeyPair.PublicKey().PublicKey()

	if k.senderKeyStates.signingKeyPair.PrivateKey() != nil {
		s.SigningKeyPrivate = k.senderKeyStates.signingKeyPair.PrivateKey().Serialize()
	} else {
		s.SigningKeyPrivate = [32]byte{}
	}

	return s
}

// PackChainKey .
func (k *SenderKey) PackChainKey() *SenderKeyStructure {
	// Build and return our state structure.
	return &SenderKeyStructure{
		Iteration: k.senderKeyStates.senderChainKey.Iteration(),
		ChainKey:  k.senderKeyStates.senderChainKey.Seed(),
	}
}
