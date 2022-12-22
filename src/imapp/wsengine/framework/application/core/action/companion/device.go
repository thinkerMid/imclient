package companion

import (
	"errors"
	"fmt"
	"google.golang.org/protobuf/proto"
	"math/rand"
	"strconv"
	"time"
	"ws/framework/application/constant/binary"
	waProto "ws/framework/application/constant/binary/proto"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
)

var ErrQRCode = errors.New("error code, please try again")

// CompanionDevice .
type CompanionDevice struct {
	processor.BaseAction
	Content      string
	TargetDevice types.JID
}

func MakeCompanionDevice(content string) *CompanionDevice {
	return &CompanionDevice{
		Content: content,
	}
}

// Start .
func (m *CompanionDevice) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	//identityKeyPair := context.ResolveIdentityService().Context()
	//
	//var pairNode waBinary.Node
	//pairNode, err = scanCompanionDeviceLogin(m.Content, *identityKeyPair, *device.AesKey)
	//if err != nil {
	//	return
	//}
	//
	//m.SendMessageId, err = context.SendIQ(message.InfoQuery{
	//	ID:        context.GenerateRequestID(),
	//	Namespace: "md",
	//	Type:      "set",
	//	To:        types.ServerJID,
	//	Content: []waBinary.Node{
	//		pairNode,
	//	},
	//})

	return
}

// Receive .
func (m *CompanionDevice) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	next()

	context.AppendResult(containerInterface.MessageResult{
		IContent: parsePairDeviceResult(*context.Message()),
	})
	return nil
}

// Error .
func (m *CompanionDevice) Error(context containerInterface.IMessageContext, err error) {
	// <iq from="s.whatsapp.net" id="1656597282-8" type="error"><error code="400" text="bad-request"/></iq>
	// <iq from="s.whatsapp.net" id="1656643650-9" type="error"><error code="419" text="resource-limit"/></iq>
	// 现在无法链接至新设备 请稍后再试

}

type QRCodeData struct {
	Ref       string
	PubKey    []byte
	IdKey     []byte
	SecretKey []byte
}

func concatBytes(data ...[]byte) []byte {
	length := 0
	for _, item := range data {
		length += len(item)
	}
	output := make([]byte, length)
	ptr := 0
	for _, item := range data {
		ptr += copy(output[ptr:ptr+len(item)], item)
	}
	return output
}

//func scanCompanionDeviceLogin(content string, idKeyPair keys.KeyPair, aesKey store.AesKey) (node waBinary.Node, err error) {
//	err = ErrQRCode
//	var (
//		qrcode        QRCodeData
//		CodeKeyLength = 4
//	)
//
//	strList := strings.Split(content, ",")
//	if len(strList) != CodeKeyLength {
//		return
//	}
//
//	qrcode.Ref = strList[0]
//	qrcode.PubKey, err = base64.StdEncoding.DecodeString(strList[1])
//	if err != nil {
//		return
//	}
//
//	qrcode.IdKey, err = base64.StdEncoding.DecodeString(strList[2])
//	if err != nil {
//		return
//	}
//
//	qrcode.SecretKey, err = base64.StdEncoding.DecodeString(strList[3])
//	if err != nil {
//		return
//	}
//
//	var (
//		pub, key  [32]byte
//		signature [64]byte
//	)
//
//	// key-index-list
//	var (
//		kil  waProto.ADVKeyIndexList
//		skil waProto.ADVSignedKeyIndexList
//
//		dId     waProto.ADVDeviceIdentity
//		sdId    waProto.ADVSignedDeviceIdentity
//		dIdHmac waProto.ADVSignedDeviceIdentityHMAC
//		current = uint64(time.Now().Unix())
//	)
//	// ADVKeyIndexList
//	scanCompanionLogin(&kil)
//
//	// ADVSignedKeyIndexList
//	skil.Details, _ = proto.Marshal(&kil)
//	copy(key[:], aesKey.PriKey)
//	data := concatBytes([]byte{6, 2}, skil.Details)
//	signature = ecc.CalculateSignature(ecc.NewDjbECPrivateKey(key), data)
//	skil.AccountSignature = signature[:]
//
//	// ADVDeviceIdentity
//	dId.RawId = kil.RawId
//	dId.Timestamp = kil.Timestamp
//	dId.KeyIndex = kil.CurrentIndex
//
//	// ADVSignedDeviceIdentity
//	sdId.Details, _ = proto.Marshal(&dId)
//
//	pub = *idKeyPair.Pub
//	sdId.AccountSignatureKey = pub[:]
//
//	copy(key[:], qrcode.IdKey)
//	message := concatBytes([]byte{6, 0}, sdId.Details, qrcode.PubKey)
//	signature = ecc.CalculateSignature(ecc.NewDjbECPrivateKey(key), message)
//	sdId.AccountSignature = signature[:]
//
//	// ADVSignedDeviceIdentityHMAC
//	dIdHmac.Details, _ = proto.Marshal(&sdId)
//
//	h := hmac.New(sha256.New, qrcode.SecretKey)
//	h.Write(dIdHmac.Details)
//	dIdHmac.Hmac = h.Sum(nil)
//
//	keyIndexListData, _ := proto.Marshal(&skil)
//	deviceIdentityData, _ := proto.Marshal(&dIdHmac)
//
//	node = waBinary.Node{
//		Tag: "pair-device",
//		Content: []waBinary.Node{
//			{
//				Tag:     "ref",
//				Content: qrcode.Ref,
//			},
//			{
//				Tag:     "pub-key",
//				Content: qrcode.PubKey,
//			},
//			{
//				Tag: "key-index-list",
//				Attrs: binary.Attrs{
//					"ts": fmt.Sprintf("%d", current),
//				},
//				Content: keyIndexListData,
//			},
//			{
//				Tag:     "device-identity",
//				Content: deviceIdentityData,
//			},
//		},
//	}
//	return node, nil
//}

func parsePairDeviceResult(node waBinary.Node) types.JID {
	device := node.GetChildByTag("device")

	attrs := device.AttrGetter()
	return attrs.JID("jid")
}

// scanCompanionLogin 扫码登陆WEB/APP
func scanCompanionLogin(kil *waProto.ADVKeyIndexList) error {
	current := time.Now().Unix()

	if kil == nil {
		kil = &waProto.ADVKeyIndexList{
			RawId:        proto.Uint32(randRawID()),
			Timestamp:    proto.Uint64(uint64(current)),
			CurrentIndex: proto.Uint32(0),
			ValidIndexes: []uint32{0},
		}
	}

	cnt := len(kil.GetValidIndexes())
	if cnt >= 5 {
		return errors.New("to many companion device")
	}

	kil.Timestamp = proto.Uint64(uint64(current))
	kil.CurrentIndex = proto.Uint32(kil.GetCurrentIndex() + 1)
	kil.ValidIndexes = append(kil.GetValidIndexes(), uint32(cnt))
	return nil
}

func companionLogout(kil *waProto.ADVKeyIndexList, device uint32) {
	current := time.Now().Unix()
	kil.Timestamp = proto.Uint64(uint64(current))

	for idx, num := range kil.ValidIndexes {
		if num == device {
			kil.ValidIndexes = append(kil.ValidIndexes[:idx], kil.ValidIndexes[idx+1:]...)
			break
		}
	}
}

func randRawID() uint32 {
	str := fmt.Sprintf("%10v", rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(100000000000))
	num, _ := strconv.Atoi(str)
	return uint32(num)
}
