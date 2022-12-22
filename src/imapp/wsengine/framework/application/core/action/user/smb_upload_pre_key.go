package user

import (
	"encoding/binary"
	waBinary "ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	"ws/framework/application/libsignal/ecc"
	"ws/framework/utils/keys"
	"ws/framework/utils/xmpp"
)

// SMBUploadPreKeyToServer 上传密钥到服务器
type SMBUploadPreKeyToServer struct {
	processor.BaseAction
	Init          bool // 是否初始化的方式上传prekey
	GenerateCount int  // 生成的数量
}

// Start .
func (s *SMBUploadPreKeyToServer) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	var preKeys []keys.PreKey

	// 默认生成一些key
	if s.Init {
		preKeys, err = context.ResolvePreKeyService().InitPreKeys()
		if err != nil {
			return
		}
	} else {
		// 指定生成
		preKeys, err = context.ResolvePreKeyService().GeneratePreKeys(s.GenerateCount)
		if err != nil {
			return
		}
	}

	device := context.ResolveDeviceService().Context()
	identityKeyPair := context.ResolveIdentityService().Context()
	signedPreKeyService := context.ResolveSignedPreKeyService()
	signedPreKeyPair := signedPreKeyService.Context()

	var registrationIDBytes [4]byte
	binary.BigEndian.PutUint32(registrationIDBytes[:], device.RegistrationId)

	keypair := keys.NewKeyPairFromPrivateKey(signedPreKeyPair.KeyPair().PrivateKey().Serialize())

	keyId := signedPreKeyPair.ID()
	keySign := signedPreKeyPair.Signature()
	signedPreKey := keys.PreKey{KeyPair: *keypair, KeyID: keyId, Signature: &keySign}

	s.SendMessageId, err = context.SendIQ(message.InfoQuery{
		ID:        context.GenerateRequestID(),
		Namespace: "encrypt",
		Type:      "set",
		To:        types.ServerJID,
		Content: []waBinary.Node{
			{Tag: "registration", Content: registrationIDBytes[:]},
			{Tag: "type", Content: []byte{ecc.DjbType}},
			{Tag: "identity", Content: identityKeyPair.Pub[:]},
			xmpp.PreKeyToNode(signedPreKey),
			{Tag: "list", Content: xmpp.PreKeysToNodes(preKeys)},
			{Tag: "verified_name", Content: context.ResolveBusinessService().GenerateBusinessVerifiedName(false)},
		},
	})

	return
}

// Receive .
func (s *SMBUploadPreKeyToServer) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	next()

	return
}

// Error .
func (s *SMBUploadPreKeyToServer) Error(_ containerInterface.IMessageContext, _ error) {}
