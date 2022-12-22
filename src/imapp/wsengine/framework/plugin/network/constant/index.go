package networkConstant

import (
	"crypto/tls"
	"time"
)

// 重试配置
const (
	// DialTimeout .
	DialTimeout uint8 = 10
	// InitialInterval .
	InitialInterval uint8 = 3
	// MaxRetries .
	MaxRetries uint64 = 3
)

// ConnectionType .
type ConnectionType int8

// 连接类型
const (
	// Socket .
	Socket ConnectionType = iota + 1
	// Socks5 .
	Socks5
)

// ConnectionConfig .
type ConnectionConfig struct {
	Address                          string
	Type                             ConnectionType
	ProxyAddress, Username, Password string
	ConnectionTimeout                time.Duration
	Tls                              *tls.Config
}

// 读写超时设置
var (
	// TCPRWTimeout .
	TCPRWTimeout = time.Second * 10
	// TCPKeepAliveTime .
	TCPKeepAliveTime = time.Second * 15
)
