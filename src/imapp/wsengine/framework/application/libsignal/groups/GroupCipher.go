package groups

import (
	"fmt"
	"ws/framework/application/libsignal/cipher"
	"ws/framework/application/libsignal/ecc"
	"ws/framework/application/libsignal/groups/ratchet"
	"ws/framework/application/libsignal/groups/state/record"
	"ws/framework/application/libsignal/groups/state/store"
	"ws/framework/application/libsignal/protocol"
	"ws/framework/application/libsignal/serialize"
	"ws/framework/application/libsignal/signalerror"
)

// NewGroupCipher will return a new group message cipher that can be used for
// encrypt/decrypt operations.
func NewGroupCipher(senderKeyID *protocol.SenderKeyName,
	senderKeyStore store.SenderKey, sessionSerializer *serialize.Serializer) *GroupCipher {

	return &GroupCipher{
		senderKeyID:       senderKeyID,
		senderKeyStore:    senderKeyStore,
		sessionSerializer: sessionSerializer,
	}
}

// GroupCipher is the main entry point for group encrypt/decrypt operations.
// Once a session has been established, this can be used for
// all encrypt/decrypt operations within that session.
type GroupCipher struct {
	senderKeyID       *protocol.SenderKeyName
	senderKeyStore    store.SenderKey
	sessionSerializer *serialize.Serializer
}

// Encrypt will take the given message in bytes and return encrypted bytes.
func (c *GroupCipher) Encrypt(plaintext []byte) (protocol.GroupCiphertextMessage, protocol.CiphertextMessage, error) {
	// Load the sender key based on id from our store.
	keyRecord, err := c.senderKeyStore.FindSenderKey(c.senderKeyID)
	if err != nil || keyRecord.IsEmpty() {
		return nil, nil, fmt.Errorf("%w for %s in %s", signalerror.ErrNoSenderKeyForUser, c.senderKeyID.Sender().String(), c.senderKeyID.GroupID())
	}

	senderKeyState, err := keyRecord.SenderKeyState()
	if err != nil {
		return nil, nil, err
	}

	senderKeyDistributionMessage := protocol.NewSenderKeyDistributionMessage(
		senderKeyState.KeyID(),
		senderKeyState.SenderChainKey().Iteration(),
		senderKeyState.SenderChainKey().Seed(),
		senderKeyState.SigningKey().PublicKey(),
		c.sessionSerializer.SenderKeyDistributionMessage,
	)

	// Get the message key from the senderkey state.
	senderKey, err := senderKeyState.SenderChainKey().SenderMessageKey()
	if err != nil {
		return nil, nil, err
	}

	// Encrypt the plaintext.
	ciphertext, err := cipher.EncryptCbc(senderKey.Iv(), senderKey.CipherKey(), plaintext)
	if err != nil {
		return nil, nil, err
	}

	senderKeyMessage := protocol.NewSenderKeyMessage(
		senderKeyState.KeyID(),
		senderKey.Iteration(),
		ciphertext,
		senderKeyState.SigningKey().PrivateKey(),
		c.sessionSerializer.SenderKeyMessage,
	)

	senderKeyState.SetSenderChainKey(senderKeyState.SenderChainKey().Next())

	c.senderKeyStore.UpdateSenderKey(c.senderKeyID, keyRecord)

	return senderKeyMessage, senderKeyDistributionMessage, nil
}

// Decrypt decrypts the given message using an existing session that
// is stored in the senderKey store.
func (c *GroupCipher) Decrypt(messageBody []byte) ([]byte, error) {
	senderKeyMessage, err := protocol.NewSenderKeyMessageFromBytes(messageBody, c.sessionSerializer.SenderKeyMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to parse skmsg message: %w", err)
	}

	keyRecord, err := c.senderKeyStore.FindSenderKey(c.senderKeyID)
	if err != nil || keyRecord.IsEmpty() {
		return nil, fmt.Errorf("%w for %s in %s", signalerror.ErrNoSenderKeyForUser, c.senderKeyID.Sender().String(), c.senderKeyID.GroupID())
	}

	// Get the senderkey state by id.
	senderKeyState, err := keyRecord.GetSenderKeyStateByID(senderKeyMessage.KeyID())
	if err != nil {
		return nil, err
	}

	// Verify the signature of the senderkey message.
	verified := c.verifySignature(senderKeyState.SigningKey().PublicKey(), senderKeyMessage)
	if !verified {
		return nil, signalerror.ErrSenderKeyStateVerificationFailed
	}

	senderKey, err := c.getSenderKey(senderKeyState, senderKeyMessage.Iteration())
	if err != nil {
		return nil, err
	}

	// Decrypt the message ciphertext.
	plaintext, err := cipher.DecryptCbc(senderKey.Iv(), senderKey.CipherKey(), senderKeyMessage.Ciphertext())
	if err != nil {
		return nil, err
	}

	// Store the sender key by id.
	c.senderKeyStore.UpdateSenderKey(c.senderKeyID, keyRecord)

	return plaintext, nil
}

// verifySignature will verify the signature of the senderkey message with
// the given public key.
func (c *GroupCipher) verifySignature(signingPubKey ecc.ECPublicKeyable,
	senderKeyMessage *protocol.SenderKeyMessage) bool {

	return ecc.VerifySignature(signingPubKey, senderKeyMessage.Serialize(), senderKeyMessage.Signature())
}

func (c *GroupCipher) getSenderKey(senderKeyState *record.SenderKeyState, iteration uint32) (*ratchet.SenderMessageKey, error) {
	senderChainKey := senderKeyState.SenderChainKey()
	if senderChainKey.Iteration() > iteration {
		return nil, fmt.Errorf("%w (current: %d, received: %d)", signalerror.ErrOldCounter, senderChainKey.Iteration(), iteration)
	}

	if iteration-senderChainKey.Iteration() > 2000 {
		return nil, signalerror.ErrTooFarIntoFuture
	}

	for senderChainKey.Iteration() < iteration {
		senderChainKey = senderChainKey.Next()
	}

	senderKeyState.SetSenderChainKey(senderChainKey.Next())
	return senderChainKey.SenderMessageKey()
}
