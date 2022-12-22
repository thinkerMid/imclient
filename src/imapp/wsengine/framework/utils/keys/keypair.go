package keys

import (
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/curve25519"
)

type KeyPair struct {
	Pub  *[32]byte
	Priv *[32]byte
}

func NewKeyPairFromPrivateKey(priv [32]byte) *KeyPair {
	var kp KeyPair
	kp.Priv = &priv
	var pub [32]byte
	curve25519.ScalarBaseMult(&pub, kp.Priv)
	kp.Pub = &pub
	return &kp
}

func NewKeyPair() *KeyPair {
	var priv [32]byte

	_, err := rand.Read(priv[:])
	if err != nil {
		panic(fmt.Errorf("failed to generate curve25519 private key: %w", err))
	}

	priv[0] &= 248
	priv[31] &= 127
	priv[31] |= 64

	return NewKeyPairFromPrivateKey(priv)
}

type PreKey struct {
	KeyPair
	KeyID     uint32
	Signature *[64]byte
}

func NewPreKey(keyID uint32) *PreKey {
	return &PreKey{
		KeyPair: *NewKeyPair(),
		KeyID:   keyID,
	}
}
