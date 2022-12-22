package xmpp

import (
	"encoding/binary"
	"ws/framework/application/constant/binary"
	"ws/framework/utils/keys"
)

// PreKeyToNode .
func PreKeyToNode(key keys.PreKey) waBinary.Node {
	var keyID [4]byte
	binary.BigEndian.PutUint32(keyID[:], key.KeyID)
	node := waBinary.Node{
		Tag: "key",
		Content: []waBinary.Node{
			{Tag: "id", Content: keyID[1:]},
			{Tag: "value", Content: key.Pub[:]},
		},
	}
	if key.Signature != nil {
		node.Tag = "skey"
		node.Content = append(node.GetChildren(), waBinary.Node{
			Tag:     "signature",
			Content: key.Signature[:],
		})
	}
	return node
}

// PreKeysToNodes .
func PreKeysToNodes(preKeys []keys.PreKey) []waBinary.Node {
	nodes := make([]waBinary.Node, len(preKeys))
	for i, key := range preKeys {
		nodes[i] = PreKeyToNode(key)
	}
	return nodes
}
