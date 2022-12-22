//go:build !windows

package netpoll

import (
	"context"
	"crypto/tls"
	hertzNetwork "github.com/cloudwego/hertz/pkg/network"
	hertzNetpoll "github.com/cloudwego/hertz/pkg/network/netpoll"
	"github.com/cloudwego/netpoll"
	"net"
	"time"
	networkConstant "ws/framework/plugin/network/constant"
	"ws/framework/plugin/network/netpoll/standard"
	"ws/framework/plugin/network/socks"
)

// newNIO .
func newNIO(parentCtx context.Context, config networkConstant.ConnectionConfig) (net.Conn, error) {
	var err error
	var conn netpoll.Connection

	if config.ConnectionTimeout == 0 {
		config.ConnectionTimeout = time.Duration(networkConstant.DialTimeout) * time.Second
	}

	dialer := netpoll.NewDialer()

	if config.Type == networkConstant.Socket {
		return dialer.DialConnection("tcp", config.Address, config.ConnectionTimeout)
	}

	// socks5
	conn, err = dialer.DialConnection("tcp", config.ProxyAddress, config.ConnectionTimeout)
	if err != nil {
		return nil, err
	}

	d := socks.NewDialer("", "")

	if config.Username != "" && config.Password != "" {
		up := socks.UsernamePassword{
			Username: config.Username,
			Password: config.Password,
		}
		d.AuthMethods = []socks.AuthMethod{
			socks.AuthMethodNotRequired,
			socks.AuthMethodUsernamePassword,
		}
		d.Authenticate = up.Authenticate
	}

	err = conn.SetReadTimeout(config.ConnectionTimeout)
	if err != nil {
		_ = conn.Close()
		return nil, err
	}

	err = conn.SetIdleTimeout(config.ConnectionTimeout)
	if err != nil {
		_ = conn.Close()
		return nil, err
	}

	_, err = d.DialWithConn(parentCtx, conn, "tcp", config.Address)
	if err != nil {
		_ = conn.Close()
		return nil, err
	}

	return conn, nil
}

// NewNIO .
func NewNIO(config networkConstant.ConnectionConfig) (conn net.Conn, err error) {
	f := func() error {
		conn, err = newNIO(context.Background(), config)
		if err != nil {
			return err
		}

		return nil
	}

	err = retry(f)
	return
}

type dialer struct {
	config networkConstant.ConnectionConfig
}

// DialConnection .
func (d dialer) DialConnection(_, address string, _ time.Duration, _ *tls.Config) (conn hertzNetwork.Conn, err error) {
	d.config.Address = address

	if d.config.Tls == nil {
		c, err := newNIO(context.Background(), d.config)
		if err != nil {
			return nil, err
		}

		return &hertzNetpoll.Conn{Conn: c.(hertzNetwork.Conn)}, nil
	}

	c, err := newStd(context.Background(), d.config)
	if err != nil {
		return nil, err
	}

	return standard.WarpTls(c, d.config.Tls)
}

// DialTimeout .
func (d dialer) DialTimeout(network, address string, timeout time.Duration, tlsConfig *tls.Config) (conn net.Conn, err error) {
	return d.DialConnection(network, address, timeout, tlsConfig)
}

// AddTLS .
func (d dialer) AddTLS(conn hertzNetwork.Conn, _ *tls.Config) (hertzNetwork.Conn, error) {
	return conn, nil
}
