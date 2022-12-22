package types

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
	signalProtocol "ws/framework/application/libsignal/protocol"
)

const (
	DefaultUserServer = "s.whatsapp.net"
	GroupServer       = "g.us"
	LegacyUserServer  = "c.us"
	BroadcastServer   = "broadcast"
)

var (
	EmptyJID            = JID{}
	GroupServerJID      = NewJID("", GroupServer)
	ServerJID           = NewJID("", DefaultUserServer)
	BroadcastServerJID  = NewJID("", BroadcastServer)
	StatusBroadcastJID  = NewJID("status", BroadcastServer)
	PSAJID              = NewJID("0", LegacyUserServer)
	OfficialBusinessJID = NewJID("16505361212", LegacyUserServer)
)

type MessageID = string

type JID struct {
	User   string
	Agent  uint8
	Device uint8
	Server string
	AD     bool
}

func (jid JID) UserInt() uint64 {
	number, _ := strconv.ParseUint(jid.User, 10, 64)
	return number
}

func (jid JID) ToNonAD() JID {
	if jid.AD {
		return JID{
			User:   jid.User,
			Server: DefaultUserServer,
		}
	} else {
		return jid
	}
}

func (jid JID) SignalAddress() *signalProtocol.SignalAddress {
	user := jid.User
	if jid.Agent != 0 {
		user = fmt.Sprintf("%s_%d", jid.User, jid.Agent)
	}
	return signalProtocol.NewSignalAddress(user, uint32(jid.Device))
}

func NewADJID(user string, agent, device uint8) JID {
	return JID{
		User:   user,
		Agent:  agent,
		Device: device,
		Server: DefaultUserServer,
		AD:     true,
	}
}

func parseADJID(user string) (JID, error) {
	var fullJID JID
	fullJID.AD = true
	fullJID.Server = DefaultUserServer

	dotIndex := strings.IndexRune(user, '.')
	colonIndex := strings.IndexRune(user, ':')
	if dotIndex < 0 || colonIndex < 0 || colonIndex+1 <= dotIndex {
		return fullJID, fmt.Errorf("failed to parse ADJID: missing separators")
	}

	fullJID.User = user[:dotIndex]
	agent, err := strconv.Atoi(user[dotIndex+1 : colonIndex])
	if err != nil {
		return fullJID, fmt.Errorf("failed to parse agent from UserID: %w", err)
	} else if agent < 0 || agent > 255 {
		return fullJID, fmt.Errorf("failed to parse agent from UserID: invalid value (%d)", agent)
	}
	device, err := strconv.Atoi(user[colonIndex+1:])
	if err != nil {
		return fullJID, fmt.Errorf("failed to parse device from UserID: %w", err)
	} else if device < 0 || device > 255 {
		return fullJID, fmt.Errorf("failed to parse device from UserID: invalid value (%d)", device)
	}
	fullJID.Agent = uint8(agent)
	fullJID.Device = uint8(device)
	return fullJID, nil
}

func ParseJID(jid string) (JID, error) {
	parts := strings.Split(jid, "@")
	if len(parts) == 1 {
		return NewJID("", parts[0]), nil
	} else if strings.ContainsRune(parts[0], ':') && strings.ContainsRune(parts[0], '.') && parts[1] == DefaultUserServer {
		return parseADJID(parts[0])
	}
	return NewJID(parts[0], parts[1]), nil
}

func NewJID(user, server string) JID {
	return JID{
		User:   user,
		Server: server,
	}
}

func (jid JID) String() string {
	if jid.AD {
		return fmt.Sprintf("%s.%d:%d@%s", jid.User, jid.Agent, jid.Device, jid.Server)
	} else if len(jid.User) > 0 {
		return fmt.Sprintf("%s@%s", jid.User, jid.Server)
	} else {
		return jid.Server
	}
}

func (jid JID) IsEmpty() bool {
	return len(jid.Server) == 0
}

var _ sql.Scanner = (*JID)(nil)

func (jid *JID) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	var out JID
	var err error
	switch val := src.(type) {
	case string:
		out, err = ParseJID(val)
	case []byte:
		out, err = ParseJID(string(val))
	default:
		err = fmt.Errorf("unsupported type %T for scanning UserID", val)
	}
	if err != nil {
		return err
	}
	*jid = out
	return nil
}

func (jid JID) Value() (driver.Value, error) {
	if len(jid.Server) == 0 {
		return nil, nil
	}
	return jid.String(), nil
}
