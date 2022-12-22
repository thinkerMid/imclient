package config

import networkConstant "ws/framework/plugin/network/constant"

// OptionsFn .
type OptionsFn func(opts *Options)

// Options .
type Options struct {
	ConnectionConfig    networkConstant.ConnectionConfig
	AutoMessageMarkRead bool
}

// UseSocks .
func UseSocks(ip, port, username, password string) OptionsFn {
	return func(opts *Options) {
		if len(ip) == 0 {
			return
		}

		opts.ConnectionConfig.Type = networkConstant.Socks5
		opts.ConnectionConfig.ProxyAddress = ip + ":" + port
		opts.ConnectionConfig.Username = username
		opts.ConnectionConfig.Password = password
	}
}

// AutoMessageMarkRead 自动回复设置
func AutoMessageMarkRead(enable bool) OptionsFn {
	return func(opts *Options) {
		opts.AutoMessageMarkRead = enable
	}
}
