package message

import (
	"ws/framework/application/constant/types"
)

const (
	// IqSet .
	IqSet = "set"
	// IqGet .
	IqGet = "get"
)

// InfoQuery .
type InfoQuery struct {
	Namespace string
	Type   string
	To     types.JID
	Target types.JID
	ID     string
	Content   interface{}
}
