package xmpp

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/libsignal/protocol"
)

// ----------------------------------------------------------------------------

// ChatMessageType 消息类型
type ChatMessageType string

const (
	// TextMessageType .
	TextMessageType ChatMessageType = "text"
	// MediaMessageType .
	MediaMessageType ChatMessageType = "media"
)

// MessageChildType  消息的子类型
type MessageChildType string

const (
	// NoneChild .
	NoneChild MessageChildType = ""
	// ImageMedia .
	ImageMedia MessageChildType = "image"
	// AudioMedia .
	AudioMedia MessageChildType = "ptt"
	// VideoMedia .
	VideoMedia MessageChildType = "video"
	// ContactCard .
	ContactCard MessageChildType = "vcard"
	// ContactCardArray .
	ContactCardArray MessageChildType = "contact_array"
)

// ----------------------------------------------------------------------------

// ParseMessageInfo .
func ParseMessageInfo(srcJid types.JID, node *waBinary.Node) (types.MessageInfo, error) {
	var info types.MessageInfo
	var err error
	var ok bool

	info.MessageSource, err = ParseMessageSource(srcJid, node)
	if err != nil {
		return info, err
	}

	info.ID, ok = node.Attrs["id"].(string)
	if !ok {
		return info, fmt.Errorf("didn't find valid `id` attribute in message")
	}

	ts, ok := node.Attrs["t"].(string)
	if !ok {
		return info, fmt.Errorf("didn't find valid `t` (timestamp) attribute in message")
	}

	tsInt, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return info, fmt.Errorf("didn't find valid `t` (timestamp) attribute in message: %w", err)
	}

	info.Timestamp = time.Unix(tsInt, 0)
	info.PushName, _ = node.Attrs["notify"].(string)
	info.Category, _ = node.Attrs["category"].(string)

	return info, nil
}

// ParseMessageSource .
func ParseMessageSource(srcJid types.JID, node *waBinary.Node) (source types.MessageSource, err error) {
	from, ok := node.Attrs["from"].(types.JID)
	if !ok {
		err = fmt.Errorf("didn't find valid `from` attribute in message")
	} else if from.Server == types.GroupServer || from.Server == types.BroadcastServer {
		source.IsGroup = true
		source.Chat = from
		sender, ok := node.Attrs["participant"].(types.JID)
		if !ok {
			err = fmt.Errorf("didn't find valid `participant` attribute in group message")
		} else {
			source.Sender = sender
			if source.Sender.User == srcJid.User {
				source.IsFromMe = true
			}
		}
	} else if from.User == srcJid.User {
		source.IsFromMe = true
		source.Sender = from
		recipient, ok := node.Attrs["recipient"].(types.JID)
		if !ok {
			source.Chat = from.ToNonAD()
		} else {
			source.Chat = recipient
		}
	} else {
		source.Chat = from.ToNonAD()
		source.Sender = from
	}
	return
}

// CreateMessageNode
//
//	<message id="3A2D74932F061888A331" to="1000@s.whatsapp.net" type="text"></message>
func CreateMessageNode(to types.JID, messageType ChatMessageType) waBinary.Node {
	return waBinary.Node{
		Tag: "message",
		Attrs: waBinary.Attrs{
			"id":   GenerateMessageID(),
			"type": string(messageType),
			"to":   to.String(),
		},
	}
}

// EncryptMessageForJID
//
//	<enc type="pkmsg" v="2"><!-- 0 bytes --></enc>
//	<enc type="msg" v="2"><!-- 0 bytes --></enc>
//	<enc type="pkmsg" v="2" mediatype="?"><!-- 0 bytes --></enc>
//	<enc type="msg" v="2" mediatype="?"><!-- 0 bytes --></enc>
func EncryptMessageForJID(spf containerInterface.ISignalProtocolService, to types.JID, mediaType MessageChildType, plaintext []byte) (node waBinary.Node, err error) {
	var ciphertext protocol.CiphertextMessage

	ciphertext, err = spf.EncryptPrivateChatMessage(to, plaintext)
	if err != nil {
		return node, fmt.Errorf("%s cipher encryption failed:  %w", to, err)
	}

	encType := "msg"
	if ciphertext.Type() == protocol.PREKEY_TYPE {
		encType = "pkmsg"
	}

	node.Tag = "enc"
	node.Attrs = waBinary.Attrs{
		"v":    "2",
		"type": encType,
	}

	if mediaType != NoneChild {
		node.Attrs["mediatype"] = string(mediaType)
	}

	node.Content = ciphertext.Serialize()

	return
}

// GenerateMessageID .
// 固定byte的第一位 0x3A
// 3A59A66D1DBA04FEEAE1 3A固定后面随机
func GenerateMessageID() types.MessageID {
	id := make([]byte, 10)
	id[0] = 0x3A

	_, err := rand.Read(id[1:])
	if err != nil {
		// Out of entropy
		panic(err)
	}

	return strings.ToUpper(hex.EncodeToString(id))
}
