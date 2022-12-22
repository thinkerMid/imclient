package record

import (
	"bytes"
	"ws/framework/application/libsignal/ecc"
	"ws/framework/application/libsignal/kdf"
	"ws/framework/application/libsignal/keys/identity"
	"ws/framework/application/libsignal/keys/root"
)

// NewSession creates a new session record and uses the given session and state
// serializers to convert the object into storeable bytes.
func NewSession() *Session {
	record := Session{
		sessionState: &State{},
	}

	return &record
}

// NewSessionFromStructure will return a new session state with the
// given state structure.
func NewSessionFromStructure(structure *SessionStructure) (*Session, error) {
	// Convert our ecc keys from bytes into object form.
	localIdentityPublic := identity.NewKey(ecc.NewDjbECPublicKey(structure.LocalIdentityPublic))
	remoteIdentityPublic := ecc.NewDjbECPublicKey(structure.RemoteIdentityPublic)
	senderBaseKey := ecc.NewDjbECPublicKey(structure.SenderBaseKey)

	var pendingPreKey *PendingPreKey
	if structure.PendingPreKey != nil {
		pendingPreKey = NewPendingPreKeyFromStruct(structure.PendingPreKey)
	}

	senderChain, err := NewChainFromStructure(structure.SenderChain)
	if err != nil {
		return nil, err
	}

	receiverChain, err := NewChainFromStructure(structure.ReceiverChain)
	if err != nil {
		return nil, err
	}

	// Build our state object.
	state := &State{
		localIdentityPublic:  localIdentityPublic,
		localRegistrationID:  structure.LocalRegistrationID,
		pendingPreKey:        pendingPreKey,
		unacknowledgedState:  structure.UnacknowledgedState,
		previousCounter:      structure.PreviousCounter,
		receiverChain:        receiverChain,
		remoteIdentityPublic: identity.NewKey(remoteIdentityPublic),
		remoteRegistrationID: structure.RemoteRegistrationID,
		rootKey:              root.NewKey(kdf.DeriveSecrets, structure.RootKey),
		senderBaseKey:        senderBaseKey,
		senderChain:          senderChain,
		sessionVersion:       structure.SessionVersion,
	}

	return &Session{sessionState: state}, nil
}

// SessionStructure is the structure of a session state. Fields are public
// to be used for serialization and deserialization.
type SessionStructure struct {
	LocalIdentityPublic  [32]byte
	LocalRegistrationID  uint32
	PendingPreKey        *PendingPreKeyStructure
	UnacknowledgedState  bool
	PreviousCounter      uint32
	ReceiverChain        *ChainStructure
	RemoteIdentityPublic [32]byte
	RemoteRegistrationID uint32
	RootKey              []byte
	SenderBaseKey        [32]byte
	SenderChain          *ChainStructure
	SessionVersion       int
}

// Session encapsulates the state of an ongoing session.
type Session struct {
	sessionState *State
}

// SetState sets the session record's current state to the given
// one.
func (s *Session) SetState(sessionState *State) {
	s.sessionState = sessionState
}

// SessionState returns the session state object of the current
// session record.
func (s *Session) SessionState() *State {
	return s.sessionState
}

// HasSessionState will check this record to see if the sender's
// base key exists in the current and previous states.
func (s *Session) HasSessionState(version int, senderBaseKey []byte) bool {
	// Ensure the session state version is identical to this one.
	if s.sessionState.Version() == version && (bytes.Compare(senderBaseKey, s.sessionState.SenderBaseKey()) == 0) {
		return true
	}

	return false
}

// Pack .
func (s *Session) Pack() *SessionStructure {
	return &SessionStructure{
		PendingPreKey:        s.sessionState.pendingPreKey.structure(),
		UnacknowledgedState:  s.sessionState.unacknowledgedState,
		PreviousCounter:      s.sessionState.previousCounter,
		ReceiverChain:        s.sessionState.receiverChain.structure(),
		RemoteIdentityPublic: s.sessionState.remoteIdentityPublic.PublicKey().PublicKey(),
		RemoteRegistrationID: s.sessionState.remoteRegistrationID,
		RootKey:              s.sessionState.rootKey.Bytes(),
		SenderBaseKey:        s.sessionState.senderBaseKey.PublicKey(),
		SenderChain:          s.sessionState.senderChain.structure(),
		SessionVersion:       s.sessionState.sessionVersion,
	}
}

// PackRebuildLogicData .
func (s *Session) PackRebuildLogicData() *SessionStructure {
	return &SessionStructure{
		UnacknowledgedState: s.sessionState.unacknowledgedState,
		PreviousCounter:     s.sessionState.previousCounter,
		ReceiverChain:       s.sessionState.receiverChain.structure(),
		RootKey:             s.sessionState.rootKey.Bytes(),
		SenderChain:         s.sessionState.senderChain.structure(),
	}
}

// PackDecryptLogicData .
func (s *Session) PackDecryptLogicData() *SessionStructure {
	return &SessionStructure{
		UnacknowledgedState: s.sessionState.unacknowledgedState,
		ReceiverChain:       s.sessionState.receiverChain.structure(),
	}
}

// PackEncryptLogicData .
func (s *Session) PackEncryptLogicData() *SessionStructure {
	return &SessionStructure{
		UnacknowledgedState: s.sessionState.unacknowledgedState,
		SenderChain:         s.sessionState.senderChain.structure(),
	}
}
