package waBinary

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidType    = errors.New("unsupported payload type")
	ErrInvalidJIDType = errors.New("invalid UserID type")
	ErrInvalidNode    = errors.New("invalid node")
	ErrInvalidToken   = errors.New("invalid token with tag")
	ErrNonStringKey   = errors.New("non-string key")
)

type IQError struct {
	Code      int
	Text      string
	ErrorNode *Node
	RawNode   *Node
}

func ParseIQError(node *Node) error {
	var err IQError
	err.RawNode = node
	val, ok := node.GetOptionalChildByTag("error")
	if ok {
		err.ErrorNode = &val
		ag := val.AttrGetter()
		err.Code = ag.OptionalInt("code")
		err.Text = ag.OptionalString("text")
	}
	return &err
}

func (iqe *IQError) Error() string {
	if iqe.Code == 0 {
		if iqe.ErrorNode != nil {
			return fmt.Sprintf("info query returned unknown error: %s", iqe.ErrorNode.XMLString())
		} else if iqe.RawNode != nil {
			return fmt.Sprintf("info query returned unexpected response: %s", iqe.RawNode.XMLString())
		} else {
			return "unknown info query error"
		}
	}
	return fmt.Sprintf("info query returned status %d: %s", iqe.Code, iqe.Text)
}

func (iqe *IQError) Is(other error) bool {
	otherIQE, ok := other.(*IQError)
	if !ok {
		return false
	} else if iqe.Code != 0 && otherIQE.Code != 0 {
		return otherIQE.Code == iqe.Code && otherIQE.Text == iqe.Text
	} else if iqe.ErrorNode != nil && otherIQE.ErrorNode != nil {
		return iqe.ErrorNode.XMLString() == otherIQE.ErrorNode.XMLString()
	} else {
		return false
	}
}
