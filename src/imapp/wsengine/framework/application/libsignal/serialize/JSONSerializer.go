package serialize

import (
	"ws/framework/application/libsignal/protocol"
	"ws/framework/application/libsignal/state/record"
	"ws/framework/plugin/json"
)

// newJsonSerializer will return a serializer for all Signal objects that will
// be responsible for converting objects to and from json bytes.
func newJsonSerializer() *Serializer {
	serializer := NewSerializer()

	serializer.SignalMessage = &jsonSignalMessageSerializer{}
	serializer.PreKeySignalMessage = &jsonPreKeySignalMessageSerializer{}
	serializer.SignedPreKeyRecord = &jsonSignedPreKeyRecordSerializer{}
	serializer.PreKeyRecord = &jsonPreKeyRecordSerializer{}
	serializer.SenderKeyMessage = &jsonSenderKeyMessageSerializer{}
	serializer.SenderKeyDistributionMessage = &jsonSenderKeyDistributionMessageSerializer{}

	return serializer
}

// jsonSignalMessageSerializer is a structure for serializing signal messages into
// and from json.
type jsonSignalMessageSerializer struct{}

// Serialize will take a signal message structure and convert it to json bytes.
func (j *jsonSignalMessageSerializer) Serialize(signalMessage *protocol.SignalMessageStructure) []byte {
	serialized, _ := json.Marshal(*signalMessage)

	return serialized
}

// Deserialize will take in json bytes and return a signal message structure.
func (j *jsonSignalMessageSerializer) Deserialize(serialized []byte) (*protocol.SignalMessageStructure, error) {
	var signalMessage protocol.SignalMessageStructure
	err := json.Unmarshal(serialized, &signalMessage)
	if err != nil {
		return nil, err
	}

	return &signalMessage, nil
}

// jsonPreKeySignalMessageSerializer is a structure for serializing prekey signal messages
// into and from json.
type jsonPreKeySignalMessageSerializer struct{}

// Serialize will take a prekey signal message structure and convert it to json bytes.
func (j *jsonPreKeySignalMessageSerializer) Serialize(signalMessage *protocol.PreKeySignalMessageStructure) []byte {
	serialized, _ := json.Marshal(signalMessage)

	return serialized
}

// Deserialize will take in json bytes and return a prekey signal message structure.
func (j *jsonPreKeySignalMessageSerializer) Deserialize(serialized []byte) (*protocol.PreKeySignalMessageStructure, error) {
	var preKeySignalMessage protocol.PreKeySignalMessageStructure
	err := json.Unmarshal(serialized, &preKeySignalMessage)
	if err != nil {
		return nil, err
	}

	return &preKeySignalMessage, nil
}

// jsonSignedPreKeyRecordSerializer is a structure for serializing signed prekey records
// into and from json.
type jsonSignedPreKeyRecordSerializer struct{}

// Serialize will take a signed prekey record structure and convert it to json bytes.
func (j *jsonSignedPreKeyRecordSerializer) Serialize(signedPreKey *record.SignedPreKeyStructure) []byte {
	serialized, _ := json.Marshal(signedPreKey)

	return serialized
}

// Deserialize will take in json bytes and return a signed prekey record structure.
func (j *jsonSignedPreKeyRecordSerializer) Deserialize(serialized []byte) (*record.SignedPreKeyStructure, error) {
	var signedPreKeyStructure record.SignedPreKeyStructure
	err := json.Unmarshal(serialized, &signedPreKeyStructure)
	if err != nil {
		return nil, err
	}

	return &signedPreKeyStructure, nil
}

// jsonPreKeyRecordSerializer is a structure for serializing prekey records
// into and from json.
type jsonPreKeyRecordSerializer struct{}

// Serialize will take a prekey record structure and convert it to json bytes.
func (j *jsonPreKeyRecordSerializer) Serialize(preKey *record.PreKeyStructure) []byte {
	serialized, _ := json.Marshal(preKey)

	return serialized
}

// Deserialize will take in json bytes and return a prekey record structure.
func (j *jsonPreKeyRecordSerializer) Deserialize(serialized []byte) (*record.PreKeyStructure, error) {
	var preKeyStructure record.PreKeyStructure
	err := json.Unmarshal(serialized, &preKeyStructure)
	if err != nil {
		return nil, err
	}

	return &preKeyStructure, nil
}

// jsonSenderKeyDistributionMessageSerializer is a structure for serializing senderkey
// distribution records to and from json.
type jsonSenderKeyDistributionMessageSerializer struct{}

// Serialize will take a senderkey distribution message and convert it to json bytes.
func (j *jsonSenderKeyDistributionMessageSerializer) Serialize(message *protocol.SenderKeyDistributionMessageStructure) []byte {
	serialized, _ := json.Marshal(message)

	return serialized
}

// Deserialize will take in json bytes and return a message structure, which can be
// used to create a new SenderKey Distribution object.
func (j *jsonSenderKeyDistributionMessageSerializer) Deserialize(serialized []byte) (*protocol.SenderKeyDistributionMessageStructure, error) {
	var msgStructure protocol.SenderKeyDistributionMessageStructure
	err := json.Unmarshal(serialized, &msgStructure)
	if err != nil {
		return nil, err
	}

	return &msgStructure, nil
}

// jsonSenderKeyMessageSerializer is a structure for serializing senderkey
// messages to and from json.
type jsonSenderKeyMessageSerializer struct{}

// Serialize will take a senderkey message and convert it to json bytes.
func (j *jsonSenderKeyMessageSerializer) Serialize(message *protocol.SenderKeyMessageStructure) []byte {
	serialized, _ := json.Marshal(message)

	return serialized
}

// Deserialize will take in json bytes and return a message structure, which can be
// used to create a new SenderKey message object.
func (j *jsonSenderKeyMessageSerializer) Deserialize(serialized []byte) (*protocol.SenderKeyMessageStructure, error) {
	var msgStructure protocol.SenderKeyMessageStructure
	err := json.Unmarshal(serialized, &msgStructure)
	if err != nil {
		return nil, err
	}

	return &msgStructure, nil
}
