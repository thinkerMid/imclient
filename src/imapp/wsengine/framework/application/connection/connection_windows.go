//go:build windows

package connection

import (
	"context"
	"crypto/cipher"
	"errors"
	"fmt"
	"github.com/cloudwego/netpoll"
	"ws/framework/application/constant"
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/binary/serialize"
	"ws/framework/application/container/abstract_interface"
	networkConstant "ws/framework/plugin/network/constant"
	netpollPlugin "ws/framework/plugin/network/netpoll"
	networkWatch "ws/framework/plugin/network/watch"
)

// ----------------------------------------------------------------------------

// Connection .
type Connection struct {
	containerInterface.BaseService

	conn        netpoll.Connection
	bufferCodec bufferCodec
	listener    containerInterface.ConnectionEventListener

	// 有效连接
	activeConnection bool
}

// SetEventListener .
func (c *Connection) SetEventListener(listener containerInterface.ConnectionEventListener) {
	c.listener = listener

	go func() {
		for {
			reader := c.conn.Reader()

			err := c.onRead(nil, reader)

			_ = reader.Release()

			if err != nil {
				c.Close()
				c.onConnectionCloseCallback(nil)
				break
			}
		}
	}()
}

func (c *Connection) createConnection() (netpoll.Connection, error) {
	configuration := c.AppIocContainer.ResolveWhatsappConfiguration()

	config := c.AppIocContainer.ResolveConnectionConfig()
	config.Address = configuration.TCPAddress

	conn, err := netpollPlugin.NewNIO(config)
	if err != nil {
		return nil, err
	}

	return conn.(netpoll.Connection), nil
}

// Connect .
func (c *Connection) Connect() (err error) {
	c.conn, err = c.createConnection()
	if err != nil {
		return constant.ConnectionConnectFailureError
	}

	// 监控协议握手超时
	watchID := networkWatch.Instance().AddConnection(c.conn, networkConstant.DialTimeout)

	var readKey, writeKey cipher.AEAD

	readKey, writeKey, err = c.AppIocContainer.ResolveHandshakeHandler().Do(c.conn)
	if err != nil {
		// 需要重新握手则再次连接
		if errors.Is(err, constant.RetryHandshakeError) {
			return c.Connect()
		}

		return constant.ConnectionHandshakeFailureError
	}

	c.activeConnection = true
	c.bufferCodec = newBufferCodec(readKey, writeKey)
	networkWatch.Instance().RemoveConnection(watchID)

	return
}

func (c *Connection) onRead(_ context.Context, reader netpoll.Reader) error {
	// next 是带阻塞的 如果储存不够 会必须读完这些数量
	header, err := reader.Next(3)
	if err != nil {
		return err
	}

	packetSize := builtinDecodeLen(header)

	// next 是带阻塞的 如果储存不够 会必须读完这些数量
	next, err := reader.Next(packetSize)
	if err != nil {
		return err
	} else if len(next) == 0 {
		return fmt.Errorf("empty buffer for frame")
	}

	decodeBody, err := c.bufferCodec.decode(next)
	if err != nil {
		return err
	}

	// xml node zlib check
	decompressed, err := unpack(decodeBody)
	if err != nil {
		return err
	}

	node, err := nodeSerialize.Unmarshal(decompressed)
	if err != nil {
		return err
	}

	c.listener.OnResponse(node)

	return nil
}

// Write .
func (c *Connection) Write(node waBinary.Node) error {
	if !c.conn.IsActive() {
		return constant.ConnectionClosedError
	}

	body := nodeSerialize.Marshal(node)
	bodyLength := len(body) + 16 // crypto encoded body length

	writer := c.conn.Writer()

	// packet header
	_ = writer.WriteByte(byte(bodyLength >> 16))
	_ = writer.WriteByte(byte(bodyLength >> 8))
	_ = writer.WriteByte(byte(bodyLength))

	// pre-malloc from buffer
	packetBody, _ := writer.Malloc(bodyLength)
	// packet body
	c.bufferCodec.encode(body, packetBody)

	return writer.Flush()
}

// ----------------------------------------------------------------------------

// Close .
func (c *Connection) Close() {
	if !c.activeConnection {
		return
	}

	c.activeConnection = false

	_ = c.conn.Close()
}

// IsClosed .
func (c *Connection) IsClosed() bool {
	if !c.activeConnection {
		return true
	}

	return !c.conn.IsActive()
}

func (c *Connection) onConnectionCloseCallback(_ netpoll.Connection) error {
	c.activeConnection = false
	c.listener.OnConnectionClose()
	return nil
}
