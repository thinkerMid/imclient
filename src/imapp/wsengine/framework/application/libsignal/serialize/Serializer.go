// Package serialize provides a serialization structure to serialize and
// deserialize Signal objects into storeable and transportable bytes.
package serialize

import (
	"ws/framework/application/libsignal/protocol"
	"ws/framework/application/libsignal/state/record"
)

// NewSerializer will return a new serializer object that will be used
// to encode/decode Signal objects into bytes.
func NewSerializer() *Serializer {
	return &Serializer{}
}

// Serializer is a structure to serialize Signal objects
// into bytes. This allows you to use any serialization format
// to store or send Signal objects.
type Serializer struct {
	SignalMessage                protocol.SignalMessageSerializer
	PreKeySignalMessage          protocol.PreKeySignalMessageSerializer
	SenderKeyMessage             protocol.SenderKeyMessageSerializer
	SenderKeyDistributionMessage protocol.SenderKeyDistributionMessageSerializer
	SignedPreKeyRecord           record.SignedPreKeySerializer
	PreKeyRecord                 record.PreKeySerializer
}
