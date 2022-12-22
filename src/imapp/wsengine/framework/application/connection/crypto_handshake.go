package connection

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"fmt"
	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/hkdf"
	"io"
	"sync/atomic"
	"ws/framework/utils/function_tools"
)

type HandshakeCrypto struct {
	hash    []byte
	salt    []byte
	key     cipher.AEAD
	counter uint32
}

func newCipher(key []byte) (cipher.AEAD, error) {
	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(aesCipher)
	if err != nil {
		return nil, err
	}
	return aesGCM, nil
}

func sha256Slice(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

func (nh *HandshakeCrypto) Start(pattern string, header []byte) {
	data := []byte(pattern)
	if len(data) == 32 {
		nh.hash = data
	} else {
		nh.hash = sha256Slice(data)
	}
	nh.salt = nh.hash
	var err error
	nh.key, err = newCipher(nh.hash)
	if err != nil {
		panic(err)
	}
	nh.Authenticate(header)
}

func (nh *HandshakeCrypto) Authenticate(data []byte) {
	nh.hash = sha256Slice(append(nh.hash, data...))
}

func (nh *HandshakeCrypto) postIncrementCounter() uint32 {
	count := atomic.AddUint32(&nh.counter, 1)
	return count - 1
}

func (nh *HandshakeCrypto) Encrypt(plaintext []byte) []byte {
	ciphertext := nh.key.Seal(nil, functionTools.GenerateIV(nh.postIncrementCounter()), plaintext, nh.hash)
	nh.Authenticate(ciphertext)
	return ciphertext
}

func (nh *HandshakeCrypto) Decrypt(ciphertext []byte) (plaintext []byte, err error) {
	plaintext, err = nh.key.Open(nil, functionTools.GenerateIV(nh.postIncrementCounter()), ciphertext, nh.hash)
	if err == nil {
		nh.Authenticate(ciphertext)
	}
	return
}

func (nh *HandshakeCrypto) Gen() (cipher.AEAD, cipher.AEAD, error) {
	if write, read, err := nh.extractAndExpand(nh.salt, nil); err != nil {
		return nil, nil, fmt.Errorf("failed to extract final keys: %w", err)
	} else if writeKey, err := newCipher(write); err != nil {
		return nil, nil, fmt.Errorf("failed to create final write cipher: %w", err)
	} else if readKey, err := newCipher(read); err != nil {
		return nil, nil, fmt.Errorf("failed to create final read cipher: %w", err)
	} else {
		return readKey, writeKey, nil
	}
}

func (nh *HandshakeCrypto) MixSharedSecretIntoKey(priv, pub [32]byte) error {
	secret, err := curve25519.X25519(priv[:], pub[:])
	if err != nil {
		return fmt.Errorf("failed to do x25519 scalar multiplication: %w", err)
	}
	return nh.MixIntoKey(secret)
}

func (nh *HandshakeCrypto) MixIntoKey(data []byte) error {
	nh.counter = 0
	write, read, err := nh.extractAndExpand(nh.salt, data)
	if err != nil {
		return fmt.Errorf("failed to extract keys for mixing: %w", err)
	}
	nh.salt = write
	nh.key, err = newCipher(read)
	if err != nil {
		return fmt.Errorf("failed to create new cipher while mixing keys: %w", err)
	}
	return nil
}

func (nh *HandshakeCrypto) extractAndExpand(salt, data []byte) (write []byte, read []byte, err error) {
	h := hkdf.New(sha256.New, data, salt, nil)
	write = make([]byte, 32)
	read = make([]byte, 32)

	if _, err = io.ReadFull(h, write); err != nil {
		err = fmt.Errorf("failed to read write key: %w", err)
	} else if _, err = io.ReadFull(h, read); err != nil {
		err = fmt.Errorf("failed to read read key: %w", err)
	}
	return
}
