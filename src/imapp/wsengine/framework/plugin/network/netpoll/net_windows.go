//go:build windows

package netpoll

import (
	"context"
	"crypto/tls"
	hertzNetwork "github.com/cloudwego/hertz/pkg/network"
	"net"
	"time"
	networkConstant "ws/framework/plugin/network/constant"
	"ws/framework/plugin/network/netpoll/standard"
)

// NewNIO .
func NewNIO(config networkConstant.ConnectionConfig) (conn net.Conn, err error) {
	f := func() error {
		netConn, stdErr := newStd(context.Background(), config)
		if stdErr != nil {
			err = stdErr
			return stdErr
		}

		conn, err = standard.WarpStdConn(netConn)

		return err
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

	c, err := newStd(context.Background(), d.config)
	if err != nil {
		return nil, err
	}

	if d.config.Tls == nil {
		return standard.WarpStdConn(c)
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
