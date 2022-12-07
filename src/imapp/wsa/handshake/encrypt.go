package handshake

func (c *Cipher) encryptEphemeralKey(pub []byte) {
	c.UpdateHash(pub)
}

func (c *Cipher) encryptStaticKey(pri, pub, sk []byte) {
	c.Encrypt(pub)
	_ = c.SetKey(pri, sk)
}
