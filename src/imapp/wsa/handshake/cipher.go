package handshake

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/hkdf"
)

type Cipher struct {
	hash    []byte
	salt    []byte
	key     cipher.AEAD
	counter uint32
}

func newGCM(key []byte) (cipher.AEAD, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}
	return gcm, nil
}

func createCipher(pattern string, header []byte) *Cipher {
	cipher := Cipher{}
	cipher.init(pattern, header)
	return &cipher
}

func (c *Cipher) init(pattern string, header []byte) {
	data := []byte(pattern)
	if len(data) == 32 {
		c.hash = data
	} else {
		c.UpdateHash(data)
	}
	c.salt = c.hash
	var err error
	c.key, err = newGCM(c.hash)
	if err != nil {
		panic(err)
	}
	c.UpdateHash(header)
}

func (c *Cipher) UpdateHash(buffer []byte) {
	hash := sha256.Sum256(append(c.hash, buffer...))
	c.hash = hash[:]
}

func (c *Cipher) Encrypt(plaintext []byte) []byte {
	iv := make([]byte, 12)
	binary.BigEndian.PutUint32(iv[8:], c.counter)
	c.counter += 1

	ciphertext := c.key.Seal(nil, iv, plaintext, c.hash)
	c.UpdateHash(ciphertext)
	return ciphertext
}

func (c *Cipher) Decrypt(ciphertext []byte) (plaintext []byte, err error) {
	iv := make([]byte, 12)
	binary.BigEndian.PutUint32(iv[8:], c.counter)
	c.counter += 1

	plaintext, err = c.key.Open(nil, iv, ciphertext, c.hash)
	if err == nil {
		c.UpdateHash(ciphertext)
	}
	return
}

func (c *Cipher) Gen() (r cipher.AEAD, w cipher.AEAD, err error) {
	var (
		read, write []byte
	)

	write, read, err = c.extractAndExpand(c.salt, nil)
	if err != nil {
		return
	}

	w, err = newGCM(write)
	if err != nil {
		return
	}

	r, err = newGCM(read)
	if err != nil {
		return
	}
	return
}

func (c *Cipher) SetKey(priv, pub []byte) error {
	secret, err := curve25519.X25519(priv, pub)
	if err != nil {
		return errors.New("cipher x25519 failed")
	}
	return c.mixIntoKey(secret)
}

func (c *Cipher) mixIntoKey(data []byte) error {
	c.counter = 0
	write, read, err := c.extractAndExpand(c.salt, data)
	if err != nil {
		return err
	}
	c.salt = write
	c.key, err = newGCM(read)
	if err != nil {
		return err
	}
	return nil
}

func (c *Cipher) extractAndExpand(salt, data []byte) (write []byte, read []byte, err error) {
	const (
		KeyLen = 32
	)
	reader := hkdf.New(sha256.New, data, salt, nil)
	write = make([]byte, KeyLen)
	read = make([]byte, KeyLen)

	var (
		rn, wn int
	)
	rn, err = reader.Read(read)
	if err != nil {
		return
	}

	wn, err = reader.Read(write)
	if err != nil {
		return
	}

	if rn != KeyLen || wn != KeyLen {
		err = errors.New("cipher new read and write key error")
		return
	}
	return
}
