package record

import (
	"labs/src/imapp/wsa/types/ecc"
	"labs/src/imapp/wsa/types/keys/chain"
	"labs/src/imapp/wsa/types/keys/kdf"
)

// NewReceiverChainPair will return a new ReceiverChainPair object.
func NewReceiverChainPair(receiverChain *Chain) *ReceiverChainPair {
	return &ReceiverChainPair{
		ReceiverChain: receiverChain,
	}
}

// ReceiverChainPair is a structure for a receiver chain key and index number.
type ReceiverChainPair struct {
	ReceiverChain *Chain
}

// NewChain returns a new Chain structure for SessionState.
func NewChain(senderRatchetKeyPair *ecc.ECKeyPair, chainKey *chain.Key) *Chain {
	return &Chain{
		senderRatchetKeyPair: senderRatchetKeyPair,
		chainKey:             chainKey,
	}
}

// NewChainFromStructure will return a new Chain with the given
// chain structure.
func NewChainFromStructure(structure *ChainStructure) (*Chain, error) {
	// Build the sender ratchet key from bytes.
	senderRatchetKeyPublic := ecc.NewDjbECPublicKey(structure.SenderRatchetKeyPublic)
	senderRatchetKeyPrivate := ecc.NewDjbECPrivateKey(structure.SenderRatchetKeyPrivate)
	senderRatchetKeyPair := ecc.NewECKeyPair(senderRatchetKeyPublic, senderRatchetKeyPrivate)

	// Build our new chain state.
	chainState := NewChain(
		senderRatchetKeyPair,
		chain.NewKeyFromStruct(structure.ChainKey, kdf.DeriveSecrets),
	)

	return chainState, nil
}

// ChainStructure is a serializeable structure for chain states.
type ChainStructure struct {
	SenderRatchetKeyPublic  [32]byte
	SenderRatchetKeyPrivate [32]byte
	ChainKey                *chain.KeyStructure
}

// Chain is a structure used inside the SessionState that keeps
// track of an ongoing ratcheting chain for a session.
type Chain struct {
	senderRatchetKeyPair *ecc.ECKeyPair
	chainKey             *chain.Key
}

// SenderRatchetKey returns the sender's EC keypair.
func (c *Chain) SenderRatchetKey() *ecc.ECKeyPair {
	return c.senderRatchetKeyPair
}

// SetSenderRatchetKey will set the chain state with the given EC
// key pair.
func (c *Chain) SetSenderRatchetKey(key *ecc.ECKeyPair) {
	c.senderRatchetKeyPair = key
}

// ChainKey will return the chain key in the chain state.
func (c *Chain) ChainKey() *chain.Key {
	return c.chainKey
}

// SetChainKey will set the chain state's chain key.
func (c *Chain) SetChainKey(key *chain.Key) {
	c.chainKey = key
}

// structure returns a serializeable structure of the chain state.
func (c *Chain) structure() *ChainStructure {
	// Convert our sender ratchet key private
	var senderRatchetKeyPrivate [32]byte
	if c.senderRatchetKeyPair.PrivateKey() != nil {
		senderRatchetKeyPrivate = c.senderRatchetKeyPair.PrivateKey().Serialize()
	}

	// Build the chain structure.
	return &ChainStructure{
		SenderRatchetKeyPublic:  c.senderRatchetKeyPair.PublicKey().PublicKey(),
		SenderRatchetKeyPrivate: senderRatchetKeyPrivate,
		ChainKey:                chain.NewStructFromKey(c.chainKey),
	}
}
