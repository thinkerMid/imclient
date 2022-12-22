package types

import (
	"fmt"
	"time"
)

type MessageSource struct {
	Chat     JID  // The chat where the message was sent.
	Sender   JID  // The user who sent the message.
	IsFromMe bool // Whether the message was sent by the current user instead of someone else.
	IsGroup  bool // Whether the chat is a group chat or broadcast list.
}

type DeviceSentMeta struct {
	DestinationJID string // The destination user. This should match the MessageInfo.Recipient field.
	Phash          string
}

type MessageInfo struct {
	MessageSource
	ID        string
	Type      string
	PushName  string
	Timestamp time.Time
	Category  string

	DeviceSentMeta *DeviceSentMeta // Metadata for direct messages sent from another one of the user's own devices.
}

func (ms *MessageSource) SourceString() string {
	if ms.Sender != ms.Chat {
		return fmt.Sprintf("%s in %s", ms.Sender, ms.Chat)
	} else {
		return ms.Chat.String()
	}
}
