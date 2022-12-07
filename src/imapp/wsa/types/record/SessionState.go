package record

import (
	"labs/src/imapp/wsa/types/ecc"
	functionTools "labs/src/imapp/wsa/types/helpers"
	"labs/src/imapp/wsa/types/keys/chain"
	"labs/src/imapp/wsa/types/keys/identity"
	"labs/src/imapp/wsa/types/keys/kdf"
	"labs/src/imapp/wsa/types/keys/optional"
	"labs/src/imapp/wsa/types/keys/root"
	"labs/src/imapp/wsa/types/keys/session"
)

// State is a session state that contains the structure for
// all sessions. Session states are contained inside session records.
// The session state is implemented as a struct rather than protobuffers
// to allow other serialization methods.
type State struct {
	localIdentityPublic  *identity.Key
	localRegistrationID  uint32
	pendingPreKey        *PendingPreKey
	unacknowledgedState  bool
	previousCounter      uint32
	receiverChain        *Chain
	remoteIdentityPublic *identity.Key
	remoteRegistrationID uint32
	rootKey              *root.Key
	senderBaseKey        ecc.ECPublicKeyable
	senderChain          *Chain
	sessionVersion       int
}

// SenderBaseKey returns the sender's base key in bytes.
func (s *State) SenderBaseKey() []byte {
	if s.senderBaseKey == nil {
		return nil
	}
	return s.senderBaseKey.Serialize()
}

// SetSenderBaseKey sets the sender's base key with the given bytes.
func (s *State) SetSenderBaseKey(senderBaseKey []byte) {
	s.senderBaseKey, _ = ecc.DecodePoint(senderBaseKey, 0)
}

// Version returns the session's version.
func (s *State) Version() int {
	return s.sessionVersion
}

// SetVersion sets the session state's version number.
func (s *State) SetVersion(version int) {
	s.sessionVersion = version
}

// RemoteIdentityKey returns the identity key of the remote user.
func (s *State) RemoteIdentityKey() *identity.Key {
	return s.remoteIdentityPublic
}

// SetRemoteIdentityKey sets this session's identity key for the remote
// user.
func (s *State) SetRemoteIdentityKey(identityKey *identity.Key) {
	s.remoteIdentityPublic = identityKey
}

// LocalIdentityKey returns the session's identity key for the local
// user.
func (s *State) LocalIdentityKey() *identity.Key {
	return s.localIdentityPublic
}

// SetLocalIdentityKey sets the session's identity key for the local
// user.
func (s *State) SetLocalIdentityKey(identityKey *identity.Key) {
	s.localIdentityPublic = identityKey
}

// PreviousCounter returns the counter of the previous message.
func (s *State) PreviousCounter() uint32 {
	return s.previousCounter
}

// SetPreviousCounter sets the counter for the previous message.
func (s *State) SetPreviousCounter(previousCounter uint32) {
	s.previousCounter = previousCounter
}

// RootKey returns the root key for the session.
func (s *State) RootKey() session.RootKeyable {
	return s.rootKey
}

// SetRootKey sets the root key for the session.
func (s *State) SetRootKey(rootKey session.RootKeyable) {
	s.rootKey = rootKey.(*root.Key)
}

// SenderRatchetKey returns the public ratchet key of the sender.
func (s *State) SenderRatchetKey() ecc.ECPublicKeyable {
	return s.senderChain.senderRatchetKeyPair.PublicKey()
}

// SenderRatchetKeyPair returns the public/private ratchet key pair
// of the sender.
func (s *State) SenderRatchetKeyPair() *ecc.ECKeyPair {
	return s.senderChain.senderRatchetKeyPair
}

// HasReceiverChain will check to see if the session state has
// the given ephemeral key.
func (s *State) HasReceiverChain(senderEphemeral ecc.ECPublicKeyable) bool {
	if s.receiverChain == nil {
		return false
	}

	chainSenderRatchetKey, err := ecc.DecodePoint(s.receiverChain.senderRatchetKeyPair.PublicKey().Serialize(), 0)
	if err == nil {
		srcPub := chainSenderRatchetKey.PublicKey()
		dstPub := senderEphemeral.PublicKey()

		return functionTools.SliceEqual(srcPub[:], dstPub[:])
	}

	return false
}

// HasSenderChain will check to see if the session state has a
// sender chain.
func (s *State) HasSenderChain() bool {
	return s.senderChain != nil
}

// ReceiverChainKey will use the given ephemeral key to generate a new
// chain key.
func (s *State) ReceiverChainKey() *chain.Key {
	return chain.NewKey(
		kdf.DeriveSecrets,
		s.receiverChain.chainKey.Key(),
		s.receiverChain.chainKey.Index(),
	)
}

// AddReceiverChain will add the given ratchet key and chain key to the session
// state.
func (s *State) AddReceiverChain(senderRatchetKey ecc.ECPublicKeyable, chainKey session.ChainKeyable) {
	// Create a keypair structure with our sender ratchet key.
	senderKey := ecc.NewECKeyPair(senderRatchetKey, nil)

	// Create a Chain state object that will hold our sender key, chain key, and
	// message keys.
	s.receiverChain = NewChain(senderKey, chainKey.(*chain.Key))
}

// SetSenderChain will set the given ratchet key pair and chain key for this session
// state.
func (s *State) SetSenderChain(senderRatchetKeyPair *ecc.ECKeyPair, chainKey session.ChainKeyable) {
	// Create a Chain state object that will hold our sender key, chain key, and
	// message keys.
	s.senderChain = NewChain(senderRatchetKeyPair, chainKey.(*chain.Key))
}

// SenderChainKey will return the chain key of the session state.
func (s *State) SenderChainKey() session.ChainKeyable {
	chainKey := s.senderChain.chainKey
	return chain.NewKey(kdf.DeriveSecrets, chainKey.Key(), chainKey.Index())
}

// SetSenderChainKey will set the chain key in the chain state for this session to
// the given chain key.
func (s *State) SetSenderChainKey(nextChainKey session.ChainKeyable) {
	senderChain := s.senderChain
	senderChain.SetChainKey(nextChainKey.(*chain.Key))
}

// SetReceiverChainKey sets the session's receiver chain key with the given chain key
// associated with the given senderEphemeral key.
func (s *State) SetReceiverChainKey(senderEphemeral ecc.ECPublicKeyable, chainKey session.ChainKeyable) {
	s.receiverChain.SetChainKey(chainKey.(*chain.Key))
}

// SetUnacknowledgedPreKeyMessage will return unacknowledged pre key message with the
// given key ids and base key.
func (s *State) SetUnacknowledgedPreKeyMessage(preKeyID *optional.Uint32, signedPreKeyID uint32, baseKey ecc.ECPublicKeyable) {
	s.pendingPreKey = NewPendingPreKey(
		preKeyID,
		signedPreKeyID,
		baseKey,
	)

	s.unacknowledgedState = true
}

// HasUnacknowledgedPreKeyMessage will return true if this session has an unacknowledged
// pre key message.
func (s *State) HasUnacknowledgedPreKeyMessage() bool {
	return s.unacknowledgedState
}

// UnackPreKeyMessageItems will return the session's unacknowledged pre key messages.
func (s *State) UnackPreKeyMessageItems() (*UnackPreKeyMessageItems, error) {
	preKeyID := s.pendingPreKey.preKeyID
	signedPreKeyID := s.pendingPreKey.signedPreKeyID
	baseKey, err := ecc.DecodePoint(s.pendingPreKey.baseKey.Serialize(), 0)
	if err != nil {
		return nil, err
	}
	return NewUnackPreKeyMessageItems(preKeyID, signedPreKeyID, baseKey), nil
}

// ClearUnackPreKeyMessage will clear the session's pending pre key.
func (s *State) ClearUnackPreKeyMessage() {
	s.unacknowledgedState = false
}

// SetRemoteRegistrationID sets the remote user's registration id.
func (s *State) SetRemoteRegistrationID(registrationID uint32) {
	s.remoteRegistrationID = registrationID
}

// RemoteRegistrationID returns the remote user's registration id.
func (s *State) RemoteRegistrationID() uint32 {
	return s.remoteRegistrationID
}

// SetLocalRegistrationID sets the local user's registration id.
func (s *State) SetLocalRegistrationID(registrationID uint32) {
	s.localRegistrationID = registrationID
}

// LocalRegistrationID returns the local user's registration id.
func (s *State) LocalRegistrationID() uint32 {
	return s.localRegistrationID
}
