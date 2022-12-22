package common

import (
	"encoding/binary"
	"fmt"
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	"ws/framework/application/core/result/constant"
	"ws/framework/application/libsignal/ecc"
	"ws/framework/application/libsignal/keys/identity"
	"ws/framework/application/libsignal/keys/prekey"
	"ws/framework/application/libsignal/serialize"
	"ws/framework/application/libsignal/session"
	"ws/framework/application/libsignal/util/optional"
	"ws/framework/utils/keys"
)

// QueryMultiDevicesIdentity .
type QueryMultiDevicesIdentity struct {
	processor.BaseAction
	UserID string
}

// Start .
func (c *QueryMultiDevicesIdentity) Start(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	idList := context.ResolveDeviceListService().FindUnInitSessionDeviceIDList(c.UserID)

	deviceNode := make([]waBinary.Node, 0)

	/**
	85252448199@s.whatsapp.net
	85252448199.0:2@s.whatsapp.net
	85252448199.0:1@s.whatsapp.net
	*/
	for _, id := range idList {
		jid := types.NewJID(c.UserID, types.DefaultUserServer)
		jid.Device = id
		jid.AD = id > 0

		deviceNode = append(deviceNode, waBinary.Node{
			Tag:   "user",
			Attrs: waBinary.Attrs{"jid": jid.String()},
		})
	}

	if len(deviceNode) == 0 {
		next()
		return
	}

	c.SendMessageId, err = context.SendIQ(
		message.InfoQuery{
			ID:        context.GenerateRequestID(),
			Namespace: "encrypt",
			Type:      message.IqGet,
			To:        types.ServerJID,
			Content: []waBinary.Node{{
				Tag:     "key",
				Content: deviceNode,
			}},
		},
	)

	return
}

// Receive .
func (c *QueryMultiDevicesIdentity) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	defer next()

	nodes := context.Message().GetChildrenByTag("list")

	if len(nodes) > 0 {
		nodes = nodes[0].GetChildren()
	}

	spf := context.ResolveSignalProtocolFactory().Context()
	logger := context.ResolveLogger()

	for i := range nodes {
		child := nodes[i]
		if child.Tag != "user" {
			continue
		}

		jid := child.AttrGetter().JID("jid")

		bundle, err := nodeToPreKeyBundle(uint32(jid.Device), child)
		if err != nil {
			logger.Warnf("%s %s", jid.SignalAddress().String(), err)
			continue
		}

		builder := session.NewBuilderFromSignal(spf, jid.SignalAddress(), serialize.Proto)

		err = builder.ProcessBundle(bundle)
		if err != nil {
			logger.Warnf("%s %s", jid.SignalAddress().String(), err)
			continue
		}
	}

	return
}

func (c *QueryMultiDevicesIdentity) Error(context containerInterface.IMessageContext, err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.GetUserMultiDevice,
		Error:      err,
	})
}

func nodeToPreKeyBundle(deviceID uint32, node waBinary.Node) (*prekey.Bundle, error) {
	errorNode, ok := node.GetOptionalChildByTag("error")
	if ok && errorNode.Tag == "error" {
		return nil, fmt.Errorf("got error getting prekeys: %s", errorNode.XMLString())
	}

	registrationBytes, ok := node.GetChildByTag("registration").Content.([]byte)
	if !ok || len(registrationBytes) != 4 {
		return nil, fmt.Errorf("invalid registration id in prekey response")
	}
	registrationID := binary.BigEndian.Uint32(registrationBytes)

	keysNode, ok := node.GetOptionalChildByTag("keys")
	if !ok {
		keysNode = node
	}

	identityKeyRaw, ok := keysNode.GetChildByTag("identity").Content.([]byte)
	if !ok || len(identityKeyRaw) != 32 {
		return nil, fmt.Errorf("invalid identity key in prekey response")
	}

	var identityKeyPub [32]byte
	copy(identityKeyPub[:], identityKeyRaw)

	preKey, err := nodeToPreKey(keysNode.GetChildByTag("key"))
	if err != nil {
		return nil, fmt.Errorf("invalid prekey in prekey response: %w", err)
	}
	signedPreKey, err := nodeToPreKey(keysNode.GetChildByTag("skey"))
	if err != nil {
		return nil, fmt.Errorf("invalid signed prekey in prekey response: %w", err)
	}

	return prekey.NewBundle(registrationID, deviceID,
		optional.NewOptionalUint32(preKey.KeyID), signedPreKey.KeyID,
		ecc.NewDjbECPublicKey(*preKey.Pub), ecc.NewDjbECPublicKey(*signedPreKey.Pub), *signedPreKey.Signature,
		identity.NewKey(ecc.NewDjbECPublicKey(identityKeyPub))), nil
}

func nodeToPreKey(node waBinary.Node) (*keys.PreKey, error) {
	key := keys.NewPreKey(0)
	key.Signature = new([64]byte)
	if id := node.GetChildByTag("id"); id.Tag != "id" {
		return nil, fmt.Errorf("prekey node doesn't contain id tag")
	} else if idBytes, ok := id.Content.([]byte); !ok {
		return nil, fmt.Errorf("prekey id has unexpected content (%T)", id.Content)
	} else if len(idBytes) != 3 {
		return nil, fmt.Errorf("prekey id has unexpected number of bytes (%d, expected 3)", len(idBytes))
	} else {
		key.KeyID = binary.BigEndian.Uint32(append([]byte{0}, idBytes...))
	}
	if pubkey := node.GetChildByTag("value"); pubkey.Tag != "value" {
		return nil, fmt.Errorf("prekey node doesn't contain value tag")
	} else if pubkeyBytes, ok := pubkey.Content.([]byte); !ok {
		return nil, fmt.Errorf("prekey value has unexpected content (%T)", pubkey.Content)
	} else if len(pubkeyBytes) != 32 {
		return nil, fmt.Errorf("prekey value has unexpected number of bytes (%d, expected 32)", len(pubkeyBytes))
	} else {
		copy(key.KeyPair.Pub[:], pubkeyBytes)
	}
	if node.Tag == "skey" {
		if sig := node.GetChildByTag("signature"); sig.Tag != "signature" {
			return nil, fmt.Errorf("prekey node doesn't contain signature tag")
		} else if sigBytes, ok := sig.Content.([]byte); !ok {
			return nil, fmt.Errorf("prekey signature has unexpected content (%T)", sig.Content)
		} else if len(sigBytes) != 64 {
			return nil, fmt.Errorf("prekey signature has unexpected number of bytes (%d, expected 64)", len(sigBytes))
		} else {
			copy(key.Signature[:], sigBytes)
		}
	}
	return key, nil
}
