package signalProtocolService

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"google.golang.org/protobuf/proto"
	"ws/framework/application/constant/binary/proto"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/libsignal/groups"
	groupStore "ws/framework/application/libsignal/groups/state/store"
	"ws/framework/application/libsignal/protocol"
	"ws/framework/application/libsignal/serialize"
	"ws/framework/application/libsignal/session"
	"ws/framework/application/libsignal/state/store"
)

var _ containerInterface.ISignalProtocolService = &Factory{}

// Factory .
type Factory struct {
	containerInterface.BaseService
}

// CreateGroupSession 创建群组会话
func (f *Factory) CreateGroupSession(groupID string) {
	builder := groups.NewGroupSessionBuilder(f.AppIocContainer.ResolveSenderKeyService(), serialize.Proto)

	err := builder.Create(protocol.NewSenderKeyName(groupID, f.JID.SignalAddress()))
	if err != nil {
		f.Logger.Error(err)
	}
}

// EncryptGroupMessage
// []byte skmsg
// []byte pkmsg内所需的senderkey
func (f *Factory) EncryptGroupMessage(groupID string, plaintext []byte) ([]byte, []byte, error) {
	senderKeyName := protocol.NewSenderKeyName(groupID, f.JID.SignalAddress())

	cipher := groups.NewGroupCipher(senderKeyName, f.AppIocContainer.ResolveSenderKeyService(), serialize.Proto)
	encryptPlaintext, senderKeyDistribution, err := cipher.Encrypt(padMessage(plaintext))
	if err != nil {
		return nil, nil, err
	}

	return encryptPlaintext.SignedSerialize(), senderKeyDistribution.Serialize(), nil
}

// DecryptGroupSenderKey .
func (f *Factory) DecryptGroupSenderKey(senderKeyName *protocol.SenderKeyName, body []byte) (err error) {
	var msg waProto.Message
	err = proto.Unmarshal(body, &msg)
	if err != nil {
		return
	}

	builder := groups.NewGroupSessionBuilder(f.AppIocContainer.ResolveSenderKeyService(), serialize.Proto)

	sdkBody := msg.GetSenderKeyDistributionMessage().GetAxolotlSenderKeyDistributionMessage()
	if len(sdkBody) == 0 {
		return fmt.Errorf("unknown %s sender key. hex:%s", senderKeyName.Sender().String(), hex.EncodeToString(body))
	}

	err = builder.Process(senderKeyName, sdkBody)
	return
}

// DecryptGroupMessage .
func (f *Factory) DecryptGroupMessage(senderKeyName *protocol.SenderKeyName, body []byte) ([]byte, error) {
	cipher := groups.NewGroupCipher(senderKeyName, f.AppIocContainer.ResolveSenderKeyService(), serialize.Proto)

	plaintext, err := cipher.Decrypt(body)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt group message: %w", err)
	}

	return unpackMessage(plaintext)
}

// EncryptPrivateChatMessage .
func (f *Factory) EncryptPrivateChatMessage(dstJID types.JID, plaintext []byte) (protocol.CiphertextMessage, error) {
	signalAddress := dstJID.SignalAddress()

	builder := session.NewBuilderFromSignal(f.Context(), signalAddress, serialize.Proto)
	cipher := session.NewCipher(builder, signalAddress)

	return cipher.Encrypt(padMessage(plaintext))
}

// DecryptPrivateChatMessage .
func (f *Factory) DecryptPrivateChatMessage(dstJID types.JID, body []byte, pkmsg bool) ([]byte, error) {
	builder := session.NewBuilderFromSignal(f.Context(), dstJID.SignalAddress(), serialize.Proto)
	cipher := session.NewCipher(builder, dstJID.SignalAddress())

	var plaintext []byte
	var err error

	// 可能没有设备会话 pkmsg消息类型来源于 1.可能是最近没发过消息 2.第一次发送消息没建立过设备信息
	if pkmsg {
		plaintext, _, err = cipher.DecryptMessageReturnKey(body)
		// 发送过消息 存在设备会话
	} else {
		plaintext, err = cipher.Decrypt(body)
	}

	if err != nil {
		return nil, err
	}

	return unpackMessage(plaintext)
}

func isValidPadding(plaintext []byte) bool {
	lastByte := plaintext[len(plaintext)-1]
	expectedPadding := bytes.Repeat([]byte{lastByte}, int(lastByte))

	return bytes.HasSuffix(plaintext, expectedPadding)
}

func unpackMessage(plaintext []byte) ([]byte, error) {
	if !isValidPadding(plaintext) {
		return nil, fmt.Errorf("plaintext doesn't have expected padding")
	}

	return plaintext[:len(plaintext)-int(plaintext[len(plaintext)-1])], nil
}

func padMessage(plaintext []byte) []byte {
	var pad [1]byte

	_, err := rand.Read(pad[:])
	if err != nil {
		panic(err)
	}

	pad[0] &= 0xf
	if pad[0] == 0 {
		pad[0] = 0xf
	}

	plaintext = append(plaintext, bytes.Repeat(pad[:], int(pad[0]))...)

	return plaintext
}

// ----------------------------------------------------------------------------

type service struct {
	store.IdentityKey
	store.PreKey
	store.ISessionStore
	store.SignedPreKey
	groupStore.SenderKey
}

// Context .
func (f *Factory) Context() store.SignalProtocol {
	return &service{
		IdentityKey:   f.AppIocContainer.ResolveIdentityService(),
		PreKey:        f.AppIocContainer.ResolvePreKeyService(),
		ISessionStore: f.AppIocContainer.ResolveDeviceListService(),
		SignedPreKey:  f.AppIocContainer.ResolveSignedPreKeyService(),
		SenderKey:     f.AppIocContainer.ResolveSenderKeyService(),
	}
}
