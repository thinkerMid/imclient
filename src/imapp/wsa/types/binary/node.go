package waBinary

import "errors"

type Attrs = map[string]interface{}

type Node struct {
	id      string // Attrs下的消息ID，用于简化调用
	err     error  // 消息结果，用于简化调用
	Tag     string
	Attrs   Attrs
	Content interface{}
}

// GetChildren .
func (n *Node) GetChildren() []Node {
	if n.Content == nil {
		return nil
	}
	children, ok := n.Content.([]Node)
	if !ok {
		return nil
	}
	return children
}

// GetChildrenByTag .
func (n *Node) GetChildrenByTag(tag string) (children []Node) {
	for _, node := range n.GetChildren() {
		if node.Tag == tag {
			children = append(children, node)
		}
	}
	return
}

// GetOptionalChildByTag .
func (n *Node) GetOptionalChildByTag(tags ...string) (val Node, ok bool) {
	val = *n
Outer:
	for _, tag := range tags {
		for _, child := range val.GetChildren() {
			if child.Tag == tag {
				val = child
				continue Outer
			}
		}
		// If no matching children are found, return false
		return
	}
	// All iterations of loop found a matching child, return it
	ok = true
	return
}

// GetChildByTag .
func (n *Node) GetChildByTag(tags ...string) Node {
	node, _ := n.GetOptionalChildByTag(tags...)
	return node
}

var emptyError = errors.New("empty")

// HasError .
func (n *Node) HasError() error {
	if n.err != nil && !errors.Is(n.err, emptyError) {
		return n.err
	}

	resType, _ := n.Attrs["type"].(string)
	if resType == "error" {
		n.err = ParseIQError(n)
	} else {
		n.err = emptyError
		return nil
	}

	return n.err
}

// ID .
func (n *Node) ID() string {
	if len(n.id) != 0 {
		return n.id
	}

	nodeId, ok := n.Attrs["id"]
	if !ok {
		return ""
	}

	n.id = nodeId.(string)

	return n.id
}

// Unmarshal .
func Unmarshal(data []byte) (*Node, error) {
	r := newDecoder(data)
	n, err := r.readNode()
	if err != nil {
		return nil, err
	}
	return n, nil
}
