package ec25519

import (
	"crypto/rand"
	"fmt"
	"golang.org/x/crypto/curve25519"
)

type KeyPair struct {
	Pub *[32]byte
	Pri *[32]byte
}

func NewKeyPairFromPrivateKey(pri [32]byte) *KeyPair {
	var pub [32]byte
	curve25519.ScalarBaseMult(&pub, &pri)

	var kp KeyPair
	kp.Pri = &pri
	kp.Pub = &pub
	return &kp
}

func NewKeyPair() *KeyPair {
	var pri [32]byte

	_, err := rand.Read(pri[:])
	if err != nil {
		panic(fmt.Errorf("failed to generate curve25519 private key: %w", err))
	}

	pri[0] &= 248
	pri[31] &= 127
	pri[31] |= 64

	return NewKeyPairFromPrivateKey(pri)
}
