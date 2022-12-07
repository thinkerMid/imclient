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

var (
	ErrNoSession      = errors.New("can't encrypt message for device: no signal session established")
	ErrIQTimedOut     = errors.New("info query timed out")
	ErrIQDisconnected = errors.New("socket disconnected before info query returned response")
	ErrNotConnected   = errors.New("socket not connected")
	ErrNotLoggedIn    = errors.New("the store doesn't contain a device UserID")

	ErrAlreadyConnected = errors.New("socket is already connected")

	ErrQRAlreadyConnected = errors.New("GetQRChannel must be called before connecting")
	ErrQRStoreContainsID  = errors.New("GetQRChannel can only be called when there's no user id in the client's Store")

	ErrNoPushName = errors.New("can't send presence without PushName set")
)

var (
	ErrProfilePictureUnauthorized  = errors.New("the user has hidden their profile picture from you")
	ErrGroupInviteLinkUnauthorized = errors.New("you don't have the permission to get the group's invite link")
	ErrNotInGroup                  = errors.New("you're not participating in that group")
	ErrGroupNotFound               = errors.New("that group does not exist")
	ErrInviteLinkInvalid           = errors.New("that group invite link is not valid")
	ErrInviteLinkRevoked           = errors.New("that group invite link has been revoked")
	ErrBusinessMessageLinkNotFound = errors.New("that business message link does not exist or has been revoked")
)

var (
	ErrBroadcastListUnsupported = errors.New("sending to broadcast lists is not yet supported")
	ErrUnknownServer            = errors.New("can't send message to unknown server")
	ErrRecipientADJID           = errors.New("message recipient must be normal (non-AD) UserID")
)

var (
	ErrMediaDownloadFailedWith404 = errors.New("download failed with status code 404")
	ErrMediaDownloadFailedWith410 = errors.New("download failed with status code 410")
	ErrNoURLPresent               = errors.New("no url present")
	ErrFileLengthMismatch         = errors.New("file length does not match")
	ErrTooShortFile               = errors.New("file too short")
	ErrInvalidMediaHMAC           = errors.New("invalid media hmac")
	ErrInvalidMediaEncSHA256      = errors.New("hash of media ciphertext doesn't match")
	ErrInvalidMediaSHA256         = errors.New("hash of media plaintext doesn't match")
	ErrUnknownMediaType           = errors.New("unknown media type")
	ErrNothingDownloadableFound   = errors.New("didn't find any attachments in message")
)

type wrappedIQError struct {
	HumanError error
	IQError    error
}

func (err *wrappedIQError) Error() string {
	return err.HumanError.Error()
}

func (err *wrappedIQError) Is(other error) bool {
	return errors.Is(other, err.HumanError)
}

func (err *wrappedIQError) Unwrap() error {
	return err.IQError
}

func WrapIQError(human, iq error) error {
	return &wrappedIQError{human, iq}
}

type IQError struct {
	Code      int
	Text      string
	ErrorNode *Node
	RawNode   *Node
}

var (
	ErrIQNotAuthorized error = &IQError{Code: 401, Text: "not-authorized"}
	ErrIQForbidden     error = &IQError{Code: 403, Text: "forbidden"}
	ErrIQNotFound      error = &IQError{Code: 404, Text: "item-not-found"}
	ErrIQNotAcceptable error = &IQError{Code: 406, Text: "not-acceptable"}
	ErrIQGone          error = &IQError{Code: 410, Text: "gone"}
)

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

type ElementMissingError struct {
	Tag string
	In  string
}

func (eme *ElementMissingError) Error() string {
	return fmt.Sprintf("missing <%s> element in %s", eme.Tag, eme.In)
}
