package store

import (
	"ws/framework/application/libsignal/state/record"
)

// PreKey store is an interface describing the local storage
// of PreKeyRecords
type PreKey interface {
	// Load a local PreKeyRecord
	FindPreKey(preKeyID uint32) *record.PreKey

	// Store a local PreKeyRecord
	SavePreKey(preKeyID uint32, preKeyRecord *record.PreKey)

	// Check to see if the store contains a PreKeyRecord
	ContainsPreKey(preKeyID uint32) bool

	// Delete a PreKeyRecord from local storage.
	DeletePreKey(preKeyID uint32)
}
