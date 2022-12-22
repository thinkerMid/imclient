package netpoll

import (
	"context"
	"github.com/cenkalti/backoff/v4"
	hertzClient "github.com/cloudwego/hertz/pkg/app/client"
	"golang.org/x/net/proxy"
	"net"
	"time"
	"ws/framework/plugin/network/constant"
	"ws/framework/plugin/network/socks"
)

// 重试
func retry(fn func() error) error {
	backoffConfig := backoff.NewExponentialBackOff()
	backoffConfig.InitialInterval = time.Duration(networkConstant.InitialInterval) * time.Second
	bf := backoff.WithMaxRetries(backoffConfig, networkConstant.MaxRetries)

	return backoff.Retry(fn, bf)
}

// HTTP HTTP客户端
func HTTP(config networkConstant.ConnectionConfig) (c *hertzClient.Client) {
	if len(config.ProxyAddress) > 0 {
		config.Type = networkConstant.Socks5
	} else {
		config.Type = networkConstant.Socket
	}

	c, _ = hertzClient.NewClient(
		hertzClient.WithDialer(dialer{config}),
		hertzClient.WithClientReadTimeout(networkConstant.TCPRWTimeout),
		hertzClient.WithWriteTimeout(networkConstant.TCPRWTimeout),
		hertzClient.WithDialTimeout(time.Duration(networkConstant.DialTimeout)*time.Second),
		hertzClient.WithKeepAlive(true),
		hertzClient.WithMaxIdleConnDuration(networkConstant.TCPKeepAliveTime),
	)

	return c
}

func newStd(parentCtx context.Context, config networkConstant.ConnectionConfig) (conn net.Conn, err error) {
	var contextDialer proxy.ContextDialer

	if config.Type == networkConstant.Socket {
		contextDialer = proxy.Direct
	} else {
		dialer := socks.NewDialer("tcp", config.ProxyAddress)

		if config.Username != "" && config.Password != "" {
			up := socks.UsernamePassword{
				Username: config.Username,
				Password: config.Password,
			}
			dialer.AuthMethods = []socks.AuthMethod{
				socks.AuthMethodNotRequired,
				socks.AuthMethodUsernamePassword,
			}
			dialer.Authenticate = up.Authenticate
		}

		contextDialer = dialer
	}

	if config.ConnectionTimeout == 0 {
		config.ConnectionTimeout = time.Duration(networkConstant.DialTimeout) * time.Second
	}

	return contextDialer.DialContext(parentCtx, "tcp", config.Address)
}

// NewStd .
func NewStd(parentCtx context.Context, config networkConstant.ConnectionConfig) (conn net.Conn, err error) {
	f := func() error {
		conn, err = newStd(parentCtx, config)
		if err != nil {
			return err
		}

		return nil
	}

	err = retry(f)
	return
}
