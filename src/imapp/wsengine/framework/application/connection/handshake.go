package connection

import (
	"crypto/cipher"
	"fmt"
	"github.com/cloudwego/netpoll"
	"google.golang.org/protobuf/proto"
	"ws/framework/application/constant"
	"ws/framework/application/constant/binary/proto"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/data_storage/device/database"
	functionTools "ws/framework/utils/function_tools"
	"ws/framework/utils/keys"
)

// Handshake .
type Handshake struct {
	containerInterface.BaseService

	// 二次重连后需要使用握手内容
	edgeRouting []byte
}

// Init .
func (m *Handshake) Init() {
	routingInfo := m.AppIocContainer.ResolveRoutingInfoService().Context()
	if routingInfo == nil {
		_ = m.AppIocContainer.ResolveRoutingInfoService().Create(make([]byte, 0))
		return
	}

	m.edgeRouting = routingInfo.Content
}

// Do .
func (m *Handshake) Do(conn netpoll.Connection) (cipher.AEAD, cipher.AEAD, error) {
	device := m.AppIocContainer.ResolveDeviceService().Context()
	if len(device.ServerStaticKey) == 0 {
		return m.noiseFullHandshake(conn)
	}

	return m.noiseResumeHandshake(conn)
}

func (m *Handshake) noiseFullHandshake(conn netpoll.Connection) (cipher.AEAD, cipher.AEAD, error) {
	configuration := m.AppIocContainer.ResolveWhatsappConfiguration()
	keyPair := keys.NewKeyPair()

	data, _ := proto.Marshal(&waProto.HandshakeMessage{
		ClientHello: &waProto.ClientHello{
			Ephemeral: keyPair.Pub[:],
		},
	})

	handshakeResponse, err := m.doHandshake(conn, data)
	if err != nil {
		return nil, nil, err
	}

	handshakeCrypto := HandshakeCrypto{}
	handshakeCrypto.Start(configuration.NoiseFullPattern, configuration.HandshakeHeader)
	handshakeCrypto.Authenticate(keyPair.Pub[:])

	serverEphemeral := handshakeResponse.GetServerHello().GetEphemeral()
	serverStaticCiphertext := handshakeResponse.GetServerHello().GetStatic()
	certificateCiphertext := handshakeResponse.GetServerHello().GetPayload()
	if len(serverEphemeral) != 32 || serverStaticCiphertext == nil || certificateCiphertext == nil {
		return nil, nil, fmt.Errorf("missing parts of Handshake response")
	}

	serverEphemeralArr := functionTools.SliceTo32SizeArray(serverEphemeral)
	handshakeCrypto.Authenticate(serverEphemeral)
	err = handshakeCrypto.MixSharedSecretIntoKey(*keyPair.Priv, serverEphemeralArr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to mix server ephemeral key in: %w", err)
	}

	serverStaticKey, err := handshakeCrypto.Decrypt(serverStaticCiphertext)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decrypt server static ciphertext: %w", err)
	} else if len(serverStaticKey) != 32 {
		return nil, nil, fmt.Errorf("unexpected length of server static plaintext %d (expected 32)", len(serverStaticKey))
	}

	serverStaticKeyArray := functionTools.SliceTo32SizeArray(serverStaticKey)
	err = handshakeCrypto.MixSharedSecretIntoKey(*keyPair.Priv, serverStaticKeyArray)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to mix server static key in: %w", err)
	}

	certDecrypted, err := handshakeCrypto.Decrypt(certificateCiphertext)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decrypt noise certificate ciphertext: %w", err)
	}

	var cert waProto.NoiseCertificate
	err = proto.Unmarshal(certDecrypted, &cert)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal noise certificate: %w", err)
	}

	certDetailsRaw := cert.GetDetails()
	certSignature := cert.GetSignature()
	if certDetailsRaw == nil || certSignature == nil {
		return nil, nil, fmt.Errorf("missing parts of noise certificate")
	}

	var certDetails waProto.NoiseCertificateDetails
	err = proto.Unmarshal(certDetailsRaw, &certDetails)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal noise certificate details: %w", err)
	} else if !functionTools.SliceEqual(certDetails.GetKey(), serverStaticKey) {
		return nil, nil, fmt.Errorf("cert key doesn't match decrypted static")
	}

	deviceService := m.AppIocContainer.ResolveDeviceService()
	device := deviceService.Context()

	clientFinishStatic := handshakeCrypto.Encrypt(device.ClientStaticPubKey)
	clientStaticPriKey := functionTools.SliceTo32SizeArray(device.ClientStaticPriKey)
	err = handshakeCrypto.MixSharedSecretIntoKey(clientStaticPriKey, serverEphemeralArr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to mix noise private key in: %w", err)
	}

	clientFinishPayloadBytes, err := proto.Marshal(deviceService.GetClientPayload())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal client finish payload: %w", err)
	}

	data, _ = proto.Marshal(&waProto.HandshakeMessage{
		ClientFinish: &waProto.ClientFinish{
			Static:  clientFinishStatic,
			Payload: handshakeCrypto.Encrypt(clientFinishPayloadBytes),
		},
	})

	err = m.sendHandShakeFinished(conn, data)
	if err != nil {
		return nil, nil, err
	}

	// save serverStaticKey
	deviceService.ContextExecute(func(device *deviceDB.Device) {
		device.UpdateServerStaticKey(serverStaticKey)
	})

	return handshakeCrypto.Gen()
}

func (m *Handshake) noiseResumeHandshake(conn netpoll.Connection) (cipher.AEAD, cipher.AEAD, error) {
	configuration := m.AppIocContainer.ResolveWhatsappConfiguration()
	deviceService := m.AppIocContainer.ResolveDeviceService()
	device := deviceService.Context()

	keyPair := keys.NewKeyPair()

	handshakeCrypto := HandshakeCrypto{}
	handshakeCrypto.Start(configuration.NoiseResumePattern, configuration.HandshakeHeader)
	handshakeCrypto.Authenticate(device.ServerStaticKey)
	handshakeCrypto.Authenticate(keyPair.Pub[:])

	serverStaticKey := functionTools.SliceTo32SizeArray(device.ServerStaticKey)
	err := handshakeCrypto.MixSharedSecretIntoKey(*keyPair.Priv, serverStaticKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to mix server static key in: %w", err)
	}

	clientHelloStatic := handshakeCrypto.Encrypt(device.ClientStaticPubKey)
	clientStaticPriKey := functionTools.SliceTo32SizeArray(device.ClientStaticPriKey)
	err = handshakeCrypto.MixSharedSecretIntoKey(clientStaticPriKey, serverStaticKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to mix server static key in: %w", err)
	}

	clientFinishPayloadBytes, err := proto.Marshal(deviceService.GetClientPayload())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal client finish payload: %w", err)
	}

	data, _ := proto.Marshal(&waProto.HandshakeMessage{
		ClientHello: &waProto.ClientHello{
			Ephemeral: keyPair.Pub[:],
			Static:    clientHelloStatic,
			Payload:   handshakeCrypto.Encrypt(clientFinishPayloadBytes),
		},
	})

	handshakeResponse, err := m.doHandshake(conn, data)
	if err != nil {
		return nil, nil, err
	}

	serverEphemeral := handshakeResponse.GetServerHello().GetEphemeral()
	certificateCiphertext := handshakeResponse.GetServerHello().GetPayload()
	if len(serverEphemeral) != 32 || certificateCiphertext == nil {
		return nil, nil, fmt.Errorf("missing parts of Handshake response")
	}

	handshakeCrypto.Authenticate(serverEphemeral)

	serverEphemeralArr := functionTools.SliceTo32SizeArray(serverEphemeral)
	err = handshakeCrypto.MixSharedSecretIntoKey(*keyPair.Priv, serverEphemeralArr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to mix server ephemeral key in: %w", err)
	}

	err = handshakeCrypto.MixSharedSecretIntoKey(clientStaticPriKey, serverEphemeralArr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to mix server ephemeral key in: %w", err)
	}

	// 这里存在解密失败的情况 原本是需要走fallback握手 现在是重置key并让上层重连走full握手
	_, err = handshakeCrypto.Decrypt(certificateCiphertext)
	if err != nil {
		// clear serverStaticKey
		deviceService.ContextExecute(func(device *deviceDB.Device) {
			device.UpdateServerStaticKey(make([]byte, 0))
		})

		return nil, nil, constant.RetryHandshakeError
	}

	return handshakeCrypto.Gen()
}

func (m *Handshake) doHandshake(conn netpoll.Connection, data []byte) (handshakeResponse waProto.HandshakeMessage, err error) {
	// send
	err = m.sendHandshakeRequest(conn, data)
	if err != nil {
		return
	}

	// await
	reader := conn.Reader()
	defer reader.Release()

	// read packet header
	header, NextErr := reader.Next(3)
	if NextErr != nil {
		err = NextErr
		return
	}

	packetSize := builtinDecodeLen(header)

	// read packet body
	next, NextErr := reader.Next(packetSize)
	if NextErr != nil || len(next) == 0 {
		err = NextErr
		return
	}

	// unmarshal response
	err = proto.Unmarshal(next, &handshakeResponse)
	if err != nil {
		return
	}

	return
}

func (m *Handshake) sendHandshakeRequest(conn netpoll.Connection, data []byte) error {
	configuration := m.AppIocContainer.ResolveWhatsappConfiguration()
	w := conn.Writer()
	bodyLength := len(data)

	// 默认的握手就发这个
	if len(m.edgeRouting) == 0 {
		_, _ = w.WriteBinary(configuration.HandshakeHeader)
	} else {
		// 指定连接区域代码的
		_, _ = w.WriteBinary(configuration.EdInfo)
		_, _ = w.WriteBinary(configuration.EdLen)
		_, _ = w.WriteBinary(m.edgeRouting)
		_, _ = w.WriteBinary(configuration.HandshakeHeader)
	}

	// packet header
	_ = w.WriteByte(byte(bodyLength >> 16))
	_ = w.WriteByte(byte(bodyLength >> 8))
	_ = w.WriteByte(byte(bodyLength))

	// packet body
	_, _ = w.WriteBinary(data)

	return w.Flush()
}

func (m *Handshake) sendHandShakeFinished(conn netpoll.Connection, data []byte) error {
	w := conn.Writer()
	bodyLength := len(data)

	// packet header
	_ = w.WriteByte(byte(bodyLength >> 16))
	_ = w.WriteByte(byte(bodyLength >> 8))
	_ = w.WriteByte(byte(bodyLength))

	// packet body
	_, _ = w.WriteBinary(data)

	return w.Flush()
}

// SetEdgeRouting .
func (m *Handshake) SetEdgeRouting(v []byte) {
	m.edgeRouting = v
	m.AppIocContainer.ResolveRoutingInfoService().Save(v)
}
