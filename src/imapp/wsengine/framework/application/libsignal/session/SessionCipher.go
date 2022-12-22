package session

import (
	"fmt"
	"ws/framework/application/libsignal/cipher"
	"ws/framework/application/libsignal/ecc"
	"ws/framework/application/libsignal/keys/chain"
	"ws/framework/application/libsignal/keys/message"
	"ws/framework/application/libsignal/protocol"
	"ws/framework/application/libsignal/serialize"
	"ws/framework/application/libsignal/signalerror"
	"ws/framework/application/libsignal/state/record"
	"ws/framework/application/libsignal/state/store"
)

const maxFutureMessages = 2000

// NewCipher constructs a session cipher for encrypt/decrypt operations on a
// session. In order to use the session cipher, a session must have already
// been created and stored using session.Builder.
func NewCipher(builder *Builder, remoteAddress *protocol.SignalAddress) *Cipher {
	cipher := &Cipher{
		sessionStore:            builder.sessionStore,
		preKeyMessageSerializer: builder.serializer.PreKeySignalMessage,
		signalMessageSerializer: builder.serializer.SignalMessage,
		preKeyStore:             builder.preKeyStore,
		remoteAddress:           remoteAddress,
		builder:                 builder,
	}

	return cipher
}

// Cipher is the main entry point for Signal Protocol encrypt/decrypt operations.
// Once a session has been established with session.Builder, this can be used for
// all encrypt/decrypt operations within that session.
type Cipher struct {
	sessionStore            store.ISessionStore
	preKeyMessageSerializer protocol.PreKeySignalMessageSerializer
	signalMessageSerializer protocol.SignalMessageSerializer
	preKeyStore             store.PreKey
	remoteAddress           *protocol.SignalAddress
	builder                 *Builder
}

// Encrypt will take the given message in bytes and return an object that follows
// the CiphertextMessage interface.
func (d *Cipher) Encrypt(plaintext []byte) (protocol.CiphertextMessage, error) {
	sessionRecord, err := d.sessionStore.FindSession(d.remoteAddress)
	if err != nil {
		return nil, fmt.Errorf("not have session to encrypt")
	}

	sessionState := sessionRecord.SessionState()
	chainKey := sessionState.SenderChainKey()
	messageKeys := chainKey.MessageKeys()
	senderEphemeral := sessionState.SenderRatchetKey()
	previousCounter := sessionState.PreviousCounter()
	sessionVersion := sessionState.Version()

	ciphertextBody, err := encrypt(messageKeys, plaintext)
	//logger.Debug("Got ciphertextBody: ", ciphertextBody)
	if err != nil {
		return nil, err
	}

	var ciphertextMessage protocol.CiphertextMessage
	ciphertextMessage, err = protocol.NewSignalMessage(
		sessionVersion,
		chainKey.Index(),
		previousCounter,
		messageKeys.MacKey(),
		senderEphemeral,
		ciphertextBody,
		sessionState.LocalIdentityKey(),
		sessionState.RemoteIdentityKey(),
		d.signalMessageSerializer,
	)
	if err != nil {
		return nil, err
	}

	// If we haven't established a session with the recipient yet,
	// send our message as a PreKeySignalMessage.
	if sessionState.HasUnacknowledgedPreKeyMessage() {
		items, err := sessionState.UnackPreKeyMessageItems()
		if err != nil {
			return nil, err
		}
		localRegistrationID := sessionState.LocalRegistrationID()

		ciphertextMessage, err = protocol.NewPreKeySignalMessage(
			sessionVersion,
			localRegistrationID,
			items.PreKeyID(),
			items.SignedPreKeyID(),
			items.BaseKey(),
			sessionState.LocalIdentityKey(),
			ciphertextMessage.(*protocol.SignalMessage),
			d.preKeyMessageSerializer,
			d.signalMessageSerializer,
		)
		if err != nil {
			return nil, err
		}
	}

	sessionState.SetSenderChainKey(chainKey.NextKey())
	d.sessionStore.SaveEncryptSession(d.remoteAddress, sessionRecord)
	return ciphertextMessage, nil
}

// Decrypt decrypts the given message using an existing session that
// is stored in the session store.
func (d *Cipher) Decrypt(messageBody []byte) ([]byte, error) {
	ciphertextMessage, err := protocol.NewSignalMessageFromBytes(messageBody, serialize.Proto.SignalMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to parse normal message: %w", err)
	}

	plaintext, _, err := d.DecryptAndGetKey(ciphertextMessage)

	return plaintext, err
}

// DecryptAndGetKey decrypts the given message using an existing session that
// is stored in the session store and returns the message keys used for encryption.
func (d *Cipher) DecryptAndGetKey(ciphertextMessage *protocol.SignalMessage) ([]byte, *message.Keys, error) {
	// Load the session record from our session store and decrypt the message.
	sessionRecord, err := d.sessionStore.FindSession(d.remoteAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("%w %s", signalerror.ErrNoSessionForUser, d.remoteAddress.String())
	}

	// remark the sender chain key index from decrypting the message
	// maybe the session chain have been rebuild
	oldSenderIdx := sessionRecord.SessionState().SenderChainKey().Index()

	plaintext, messageKeys, err := d.DecryptWithRecord(sessionRecord, ciphertextMessage)
	if err != nil {
		return nil, nil, err
	}

	// change the sender chain key index need save all data
	if oldSenderIdx != sessionRecord.SessionState().SenderChainKey().Index() {
		// update all the session data
		d.sessionStore.SaveRebuildSession(d.remoteAddress, sessionRecord)
	} else {
		// Store the session record in our session store.
		d.sessionStore.SaveDecryptSession(d.remoteAddress, sessionRecord)
	}

	return plaintext, messageKeys, nil
}

func (d *Cipher) DecryptMessageReturnKey(messageBody []byte) ([]byte, *message.Keys, error) {
	ciphertextMessage, err := protocol.NewPreKeySignalMessageFromBytes(messageBody, d.builder.serializer.PreKeySignalMessage, d.builder.serializer.SignalMessage)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse prekey message: %w", err)
	}

	sessionRecord, loadErr := d.sessionStore.FindSession(d.remoteAddress)
	if loadErr != nil {
		// create session record for this session.
		sessionRecord = record.NewSession()
	}

	unsignedPreKeyID, err := d.builder.Process(sessionRecord, ciphertextMessage)
	if err != nil {
		return nil, nil, err
	}
	plaintext, keys, err := d.DecryptWithRecord(sessionRecord, ciphertextMessage.WhisperMessage())
	if err != nil {
		return nil, nil, err
	}

	if loadErr != nil {
		// create the session record in our session store.
		d.sessionStore.CreateSession(d.remoteAddress, sessionRecord)
	} else {
		d.sessionStore.SaveSession(d.remoteAddress, sessionRecord)
	}

	if !unsignedPreKeyID.IsEmpty {
		d.preKeyStore.DeletePreKey(unsignedPreKeyID.Value)
	}

	return plaintext, keys, nil
}

// DecryptWithKey will decrypt the given message using the given symmetric key. This
// can be used when decrypting messages at a later time if the message key was saved.
func (d *Cipher) DecryptWithKey(ciphertextMessage *protocol.SignalMessage, key *message.Keys) ([]byte, error) {
	//logger.Debug("Decrypting ciphertext body: ", ciphertextMessage.Body())
	plaintext, err := decrypt(key, ciphertextMessage.Body())
	if err != nil {
		//logger.Error("Unable to get plain text from ciphertext: ", err)
		return nil, err
	}

	return plaintext, nil
}

// DecryptWithRecord decrypts the given message using the given session record.
func (d *Cipher) DecryptWithRecord(sessionRecord *record.Session, ciphertext *protocol.SignalMessage) ([]byte, *message.Keys, error) {
	//logger.Debug("Decrypting ciphertext with record: ", sessionRecord)
	sessionState := sessionRecord.SessionState()

	// Try and decrypt the message with the current session state.
	plaintext, messageKeys, err := d.DecryptWithState(sessionState, ciphertext)

	// If we received an error using the current session state, loop
	// through all previous states.
	if err != nil {
		//logger.Warning(err)
		return nil, nil, signalerror.ErrNoValidSessions
	}

	// If decryption was successful, set the session state and return the plain text.
	sessionRecord.SetState(sessionState)

	return plaintext, messageKeys, nil
}

// DecryptWithState decrypts the given message with the given session state.
func (d *Cipher) DecryptWithState(sessionState *record.State, ciphertextMessage *protocol.SignalMessage) ([]byte, *message.Keys, error) {
	//logger.Debug("Decrypting ciphertext with session state: ", sessionState)
	if !sessionState.HasSenderChain() {
		//logger.Error("Unable to decrypt message with state: ", signalerror.ErrUninitializedSession)
		return nil, nil, signalerror.ErrUninitializedSession
	}

	if ciphertextMessage.MessageVersion() != sessionState.Version() {
		//logger.Error("Unable to decrypt message with state: ", signalerror.ErrWrongMessageVersion)
		return nil, nil, signalerror.ErrWrongMessageVersion
	}

	messageVersion := ciphertextMessage.MessageVersion()
	theirEphemeral := ciphertextMessage.SenderRatchetKey()
	counter := ciphertextMessage.Counter()
	chainKey, chainCreateErr := getOrCreateChainKey(sessionState, theirEphemeral)
	if chainCreateErr != nil {
		//logger.Error("Unable to get or create chain key: ", chainCreateErr)
		return nil, nil, fmt.Errorf("failed to get or create chain key: %w", chainCreateErr)
	}

	messageKeys, keysCreateErr := getOrCreateMessageKeys(sessionState, theirEphemeral, chainKey, counter)
	if keysCreateErr != nil {
		//logger.Error("Unable to get or create message keys: ", keysCreateErr)
		return nil, nil, fmt.Errorf("failed to get or create message keys: %w", keysCreateErr)
	}

	err := ciphertextMessage.VerifyMac(messageVersion, sessionState.RemoteIdentityKey(), sessionState.LocalIdentityKey(), messageKeys.MacKey())
	if err != nil {
		//logger.Error("Unable to verify ciphertext mac: ", err)
		return nil, nil, fmt.Errorf("failed to verify ciphertext MAC: %w", err)
	}

	plaintext, err := d.DecryptWithKey(ciphertextMessage, messageKeys)
	if err != nil {
		return nil, nil, err
	}

	sessionState.ClearUnackPreKeyMessage()

	return plaintext, messageKeys, nil
}

func getOrCreateMessageKeys(sessionState *record.State, theirEphemeral ecc.ECPublicKeyable,
	chainKey *chain.Key, counter uint32) (*message.Keys, error) {

	if chainKey.Index() > counter {
		return nil, fmt.Errorf("%w (index: %d, count: %d)", signalerror.ErrOldCounter, chainKey.Index(), counter)
	}

	if counter-chainKey.Index() > maxFutureMessages {
		return nil, signalerror.ErrTooFarIntoFuture
	}

	for chainKey.Index() < counter {
		chainKey = chainKey.NextKey()
	}

	sessionState.SetReceiverChainKey(theirEphemeral, chainKey.NextKey())
	return chainKey.MessageKeys(), nil
}

// getOrCreateChainKey will either return the existing chain key or
// create a new one with the given session state and ephemeral key.
func getOrCreateChainKey(sessionState *record.State, theirEphemeral ecc.ECPublicKeyable) (*chain.Key, error) {
	// If our session state already has a receiver chain, use their
	// ephemeral key in the existing chain.
	if sessionState.HasReceiverChain(theirEphemeral) {
		return sessionState.ReceiverChainKey(), nil
	}

	// If we don't have a chain key, create one with ephemeral keys.
	rootKey := sessionState.RootKey()
	ourEphemeral := sessionState.SenderRatchetKeyPair()
	receiverChain, rErr := rootKey.CreateChain(theirEphemeral, ourEphemeral)
	if rErr != nil {
		return nil, rErr
	}

	// Generate a new ephemeral key pair.
	ourNewEphemeral, gErr := ecc.GenerateKeyPair()
	if gErr != nil {
		return nil, gErr
	}

	// Create a new chain using our new ephemeral key.
	senderChain, cErr := receiverChain.RootKey.CreateChain(theirEphemeral, ourNewEphemeral)
	if cErr != nil {
		return nil, cErr
	}

	// Set our session state parameters.
	sessionState.SetRootKey(senderChain.RootKey)
	sessionState.AddReceiverChain(theirEphemeral, receiverChain.ChainKey)

	senderChainKeyIndex := sessionState.SenderChainKey().Index()
	if senderChainKeyIndex == 0 {
		sessionState.SetPreviousCounter(senderChainKeyIndex)
	} else {
		sessionState.SetPreviousCounter(senderChainKeyIndex - 1)
	}

	sessionState.SetSenderChain(ourNewEphemeral, senderChain.ChainKey)

	return receiverChain.ChainKey.(*chain.Key), nil
}

// decrypt will use the given message keys and ciphertext and return
// the plaintext bytes.
func decrypt(keys *message.Keys, body []byte) ([]byte, error) {
	//logger.Debug("Using cipherKey: ", keys.CipherKey())
	return cipher.DecryptCbc(keys.Iv(), keys.CipherKey(), body)
}

// encrypt will use the given cipher, message keys, and plaintext bytes
// and return ciphertext bytes.
func encrypt(messageKeys *message.Keys, plaintext []byte) ([]byte, error) {
	//logger.Debug("Using cipherKey: ", messageKeys.CipherKey())
	return cipher.EncryptCbc(messageKeys.Iv(), messageKeys.CipherKey(), plaintext)
}

// Max is a uint32 implementation of math.Max
func max(x, y uint32) uint32 {
	if x > y {
		return x
	}
	return y
}
