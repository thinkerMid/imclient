package handshake

import (
	"crypto/cipher"
	"encoding/hex"
	"github.com/cloudwego/netpoll"
	"github.com/golang/protobuf/proto"
	. "labs/src/imapp/wsa/cache"
	"labs/src/imapp/wsa/config"
	waProto "labs/src/imapp/wsa/types/binary/proto"
	"labs/src/imapp/wsa/types/ec25519"
	"labs/src/imapp/wsa/types/ecc"
	"time"
)

type Handshake struct {
	cipher *Cipher
}

func (hs *Handshake) Do() (r cipher.AEAD, w cipher.AEAD, err error) {
	var conn netpoll.Connection

	dialer := netpoll.NewDialer()
	conn, err = dialer.DialConnection("tcp", config.WSUrl, 10*time.Second)
	if err != nil {
		return
	}

	_ = conn.SetReadTimeout(10 * time.Second)
	_ = conn.SetIdleTimeout(10 * time.Second)

	err = hs.sendClientHello(conn.Writer())
	if err != nil {
		return
	}

	var resp waProto.HandshakeMessage
	resp, err = hs.waitServerHello(conn.Reader())
	if err != nil {
		return
	}

	hs.sendClientFinish(conn.Writer())

	_ = conn.Close()
	_ = resp
	return
}

func (hs *Handshake) waitServerHello(reader netpoll.Reader) (resp waProto.HandshakeMessage, err error) {
	defer func() {
		_ = reader.Release()
	}()

	const (
		HeaderSize = 3
	)

	var buff []byte
	buff, err = reader.Next(HeaderSize)
	if err != nil {
		return
	}

	length := int(buff[0])<<16 + int(buff[1])<<8 + int(buff[2])
	buff, err = reader.Next(length)
	if err != nil || len(buff) == 0 {
		return
	}

	err = proto.Unmarshal(buff, &resp)
	if err != nil {
		return
	}
	return
}

func (hs *Handshake) sendClientHello(writer netpoll.Writer) error {
	hello := hs.packClientHello()
	buffer, _ := proto.Marshal(&hello)
	length := len(buffer)
	arr := []byte{byte(length >> 16), byte(length >> 8), byte(length)}

	header, _ := hex.DecodeString("45440001000004")
	routing, _ := hex.DecodeString("080d0805")
	version, _ := hex.DecodeString("57410502")

	data := []byte("")
	data = append(data, header...)
	data = append(data, routing...)
	data = append(data, version...)
	data = append(data, arr...)
	data = append(data, buffer...)

	n, err := writer.WriteBinary(data)
	if err != nil {
		return err
	}
	_ = n

	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}

func (hs *Handshake) sendClientFinish(writer netpoll.Writer) {
	const (
		NoiseStartPattern = "Noise_XX_25519_AESGCM_SHA256\x00\x00\x00"
		Header            = "57410502"
	)
	hs.cipher.UpdateHash([]byte(NoiseStartPattern + Header))
	hs.cipher.UpdateHash(hs.keyPair.Pub[:])
}

func (hs *Handshake) packClientHello() waProto.HandshakeMessage {
	keyPair := ec25519.NewKeyPair()
	ckPair, _ := ecc.GenerateKeyPair()

	hello := waProto.HandshakeMessage{
		ClientHello: &waProto.ClientHello{
			Ephemeral: keyPair.Pub[:],
		},
	}

	sk := Cache().GetServerStaticKey()
	if sk == nil {
		return hello
	}

	ver, _ := hex.DecodeString(config.WSHeader)
	c := createCipher(config.NoiseFullPattern, ver)

	c.UpdateHash(sk)

	// encryptEphemeralKey
	c.UpdateHash(keyPair.Pub[:])

	// encryptStaticKey
	cPub := ckPair.PublicKey().Serialize()
	cPri := ckPair.PrivateKey().Serialize()

	c.Encrypt(cPub)
	_ = c.SetKey(cPri[:], sk)

	// encryptPayload

	return hello
}
