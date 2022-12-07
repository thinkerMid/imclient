package utils

import (
	"crypto/tls"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type RequestConfig struct {
	*http.Request
	retry int
}

type RequestOptions func(*RequestConfig)

const (
	DefaultTimeout  = 10 * time.Second // 默认超时时间
	DefaultRetryMax = 3                // 默认最大重试次数
	DefaultIdleConn = 100              // 默认最大空闲链接
)

var (
	transport *http.Transport
	tlsConfig *tls.Config
	client    *http.Client
)

func init() {
	tlsConfig = &tls.Config{InsecureSkipVerify: true}
	transport = &http.Transport{
		TLSClientConfig: tlsConfig,
		MaxIdleConns:    DefaultIdleConn,
	}

	client = &http.Client{
		Transport: transport,
		Timeout:   DefaultTimeout,
	}
}

func defaultConfig(cfg *RequestConfig) *RequestConfig {
	tlsConfig = &tls.Config{InsecureSkipVerify: true}
	transport = &http.Transport{
		TLSClientConfig: tlsConfig,
		MaxIdleConns:    DefaultIdleConn,
	}

	client.Timeout = DefaultTimeout
	client.Transport = transport
	return cfg
}

func WithHeaders(headers map[string]string) RequestOptions {
	return func(config *RequestConfig) {
		for k, v := range headers {
			config.Request.Header.Set(k, v)
		}
	}
}

func WithProxy(uri, user, password string) RequestOptions {
	return func(config *RequestConfig) {
		proxy, err := url.Parse(uri)
		if err != nil {
			return
		}

		if len(user) != 0 && len(password) != 0 {
			proxy.User = url.UserPassword(user, password)
		}

		ts := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           http.ProxyURL(proxy),
		}

		client.Transport = ts
	}
}

func WithRetry(maxRetry int, timeouts []int) RequestOptions {
	return func(config *RequestConfig) {

	}
}

func WithTLSConfig(cfg *tls.Config) RequestOptions {
	return func(config *RequestConfig) {
		ts := &http.Transport{
			TLSClientConfig: cfg,
		}
		client.Transport = ts
	}
}

func HttpRequest(url, method string, body io.Reader, options ...RequestOptions) (buff []byte, err error) {
	var (
		req  *http.Request
		resp *http.Response
	)

	req, err = http.NewRequest(method, url, body)
	if err != nil {
		return
	}

	cfg := RequestConfig{
		Request: req,
	}

	defaultConfig(&cfg)

	for _, opt := range options {
		opt(&cfg)
	}

	resp, err = client.Do(req)
	if err != nil {
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	buff, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return
}
