package constant

import "errors"

// NewIgnoreError .
func NewIgnoreError(msg string) *IgnoreError {
	return &IgnoreError{Message: msg}
}

// IgnoreError .
type IgnoreError struct {
	Message string `json:"message"`
}

// Error makes it compatible with the `error` interface.
func (e *IgnoreError) Error() string {
	return e.Message
}

// AbortedError 用于action流程中断
var AbortedError = NewIgnoreError("aborted")

// ConnectionClosedError 连接断开
var ConnectionClosedError = NewIgnoreError("connection closed")

// LogoutError 主动关闭连接
var LogoutError = NewIgnoreError("logout")

// ConnectionConnectFailureError 连接失败
var ConnectionConnectFailureError = NewIgnoreError("connection connect failure")

// ConnectionHandshakeFailureError 握手失败
var ConnectionHandshakeFailureError = NewIgnoreError("connection handshake failure")

// RetryHandshakeError 重试握手
var RetryHandshakeError = errors.New("retry handshake error")
