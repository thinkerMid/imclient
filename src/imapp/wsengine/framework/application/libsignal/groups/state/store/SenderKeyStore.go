package store

import (
	"ws/framework/application/libsignal/groups/state/record"
	"ws/framework/application/libsignal/protocol"
)

type SenderKey interface {
	CreateSenderKey(senderKeyName *protocol.SenderKeyName, keyRecord *record.SenderKey)
	UpdateSenderKey(senderKeyName *protocol.SenderKeyName, keyRecord *record.SenderKey)
	ResetSenderKey(senderKeyName *protocol.SenderKeyName, keyRecord *record.SenderKey)
	ContainsSenderKey(senderKeyName *protocol.SenderKeyName) bool
	FindSenderKey(senderKeyName *protocol.SenderKeyName) (*record.SenderKey, error)
}
