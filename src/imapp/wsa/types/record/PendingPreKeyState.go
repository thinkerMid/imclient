package record

import (
	"labs/src/imapp/wsa/types/ecc"
	"labs/src/imapp/wsa/types/keys/optional"
)

// NewPendingPreKey will return a new pending pre key object.
func NewPendingPreKey(preKeyID *optional.Uint32, signedPreKeyID uint32,
	baseKey ecc.ECPublicKeyable) *PendingPreKey {

	return &PendingPreKey{
		preKeyID:       preKeyID,
		signedPreKeyID: signedPreKeyID,
		baseKey:        baseKey,
	}
}

// NewPendingPreKeyFromStruct will return a new pending prekey object from the
// given structure.
func NewPendingPreKeyFromStruct(preKey *PendingPreKeyStructure) *PendingPreKey {
	pendingPreKey := NewPendingPreKey(
		preKey.PreKeyID,
		preKey.SignedPreKeyID,
		ecc.NewDjbECPublicKey(preKey.BaseKey),
	)

	return pendingPreKey
}

// PendingPreKeyStructure is a serializeable structure for pending
// prekeys.
type PendingPreKeyStructure struct {
	PreKeyID       *optional.Uint32
	SignedPreKeyID uint32
	BaseKey        [32]byte
}

// PendingPreKey is a structure for pending pre keys
// for a session state.
type PendingPreKey struct {
	preKeyID       *optional.Uint32
	signedPreKeyID uint32
	baseKey        ecc.ECPublicKeyable
}

// structure will return a serializeable structure of the pending prekey.
func (p *PendingPreKey) structure() *PendingPreKeyStructure {
	if p != nil {
		return &PendingPreKeyStructure{
			PreKeyID:       p.preKeyID,
			SignedPreKeyID: p.signedPreKeyID,
			BaseKey:        p.baseKey.PublicKey(),
		}
	}

	return &PendingPreKeyStructure{
		PreKeyID:       optional.NewEmptyUint32(),
		SignedPreKeyID: 0,
		BaseKey:        [32]byte{},
	}
}
