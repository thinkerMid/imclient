// Package groups is responsible for setting up group SenderKey encrypted sessions.
// Once a session has been established, GroupCipher can be used to encrypt/decrypt
// messages in that session.
//
// The built sessions are unidirectional: they can be used either for sending or
// for receiving, but not both. Sessions are constructed per (groupId + senderId +
// deviceId) tuple. Remote logical users are identified by their senderId, and each
// logical recipientId can have multiple physical devices.
package groups

import (
	"fmt"
	"ws/framework/application/libsignal/groups/state/record"
	"ws/framework/application/libsignal/groups/state/store"
	"ws/framework/application/libsignal/protocol"
	"ws/framework/application/libsignal/serialize"
	"ws/framework/application/libsignal/util/keyhelper"
)

// NewGroupSessionBuilder will return a new group session builder.
func NewGroupSessionBuilder(senderKeyStore store.SenderKey,
	serializer *serialize.Serializer) *SessionBuilder {

	return &SessionBuilder{
		senderKeyStore: senderKeyStore,
		serializer:     serializer,
	}
}

// SessionBuilder is a structure for building group sessions.
type SessionBuilder struct {
	senderKeyStore store.SenderKey
	serializer     *serialize.Serializer
}

// Process will process an incoming group key and new group session for it.
func (b *SessionBuilder) Process(senderKeyName *protocol.SenderKeyName, msgBody []byte) error {
	msg, err := protocol.NewSenderKeyDistributionMessageFromBytes(msgBody, b.serializer.SenderKeyDistributionMessage)
	if err != nil {
		return fmt.Errorf("failed to parse sender key")
	}

	senderKey := record.NewSenderKey()
	senderKey.AddSenderKeyState(msg.ID(), msg.Iteration(), msg.ChainKey(), msg.SignatureKey())

	if b.senderKeyStore.ContainsSenderKey(senderKeyName) {
		b.senderKeyStore.ResetSenderKey(senderKeyName, senderKey)
	} else {
		b.senderKeyStore.CreateSenderKey(senderKeyName, senderKey)
	}

	return nil
}

// Create will create a new group session for the given name.
func (b *SessionBuilder) Create(senderKeyName *protocol.SenderKeyName) error {
	senderKey := record.NewSenderKey()

	// If the record is empty, generate new keys.
	signingKey, err := keyhelper.GenerateSenderSigningKey()
	if err != nil {
		return err
	}

	senderKey.SetSenderKeyState(
		keyhelper.GenerateSenderKeyID(), 0,
		keyhelper.GenerateSenderKey(),
		signingKey,
	)

	b.senderKeyStore.CreateSenderKey(senderKeyName, senderKey)

	return nil
	//// Get the senderkey state.
	//state, err := senderKey.SenderKeyState()
	//if err != nil {
	//	return nil, err
	//}
	//
	//// Create the group message to return.
	//senderKeyDistributionMessage := protocol.NewSenderKeyDistributionMessage(
	//	state.KeyID(),
	//	state.SenderChainKey().Iteration(),
	//	state.SenderChainKey().Seed(),
	//	state.SigningKey().PublicKey(),
	//	b.serializer.SenderKeyDistributionMessage,
	//)
	//
	//return senderKeyDistributionMessage, nil
}
