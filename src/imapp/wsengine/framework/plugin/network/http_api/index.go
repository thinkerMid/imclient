package httpApi

import (
	"bytes"
	"context"
	hertzClient "github.com/cloudwego/hertz/pkg/app/client"
	hertzProtocol "github.com/cloudwego/hertz/pkg/protocol"
	"ws/framework/plugin/json"
	"ws/framework/plugin/logger"
	functionTools "ws/framework/utils/function_tools"
)

var httpLogger = logger.New("HTTP")

// RequestOptionsFn .
type RequestOptionsFn func(*hertzProtocol.Request) error

// Url .
func Url(url string) RequestOptionsFn {
	return func(req *hertzProtocol.Request) error {
		req.SetRequestURI(url)
		return nil
	}
}

// UserAgent .
func UserAgent(agent string) RequestOptionsFn {
	return func(req *hertzProtocol.Request) error {
		req.Header.SetUserAgentBytes([]byte(agent))
		return nil
	}
}

// Method .
func Method(key string) RequestOptionsFn {
	return func(req *hertzProtocol.Request) error {
		req.Header.SetMethod(key)

		return nil
	}
}

// Header .
func Header(key string, value string) RequestOptionsFn {
	return func(req *hertzProtocol.Request) error {
		req.Header.Set(key, value)

		return nil
	}
}

// QueryParams .
func QueryParams(key string, value string) RequestOptionsFn {
	return func(req *hertzProtocol.Request) error {
		args := req.URI().QueryArgs()
		args.Set(key, value)

		return nil
	}
}

// Body .
func Body(reqBody []byte) RequestOptionsFn {
	return func(req *hertzProtocol.Request) error {
		req.SetBody(reqBody)

		return nil
	}
}

// IBody .
func IBody(body interface{}) RequestOptionsFn {
	return func(req *hertzProtocol.Request) error {
		buff, _ := json.Marshal(body)
		req.SetBody(buff)

		return nil
	}
}

// RawBody .
func RawBody(reqBody []byte) RequestOptionsFn {
	return func(req *hertzProtocol.Request) error {
		buf := bytes.NewBuffer(reqBody)
		req.SetBodyStream(buf, buf.Len())

		return nil
	}
}

// DoAndBind .
func DoAndBind(client *hertzClient.Client, bind interface{}, opts ...RequestOptionsFn) (err error) {
	req := hertzProtocol.AcquireRequest()
	resp := hertzProtocol.AcquireResponse()

	defer hertzProtocol.ReleaseRequest(req)
	defer hertzProtocol.ReleaseResponse(resp)

	err = run(client, req, resp, opts...)
	if err != nil {
		return
	}

	body := resp.Body()

	if logger.EnabledDebug() {
		if len(body) > 4096 {
			httpLogger.Debug("http response body is greater than 4096b. not printable")
		} else {
			httpLogger.Debug(functionTools.B2S(body))
		}
	}

	if bind == nil {
		return
	}

	return json.Unmarshal(body, bind)
}

// Do .
func Do(client *hertzClient.Client, opts ...RequestOptionsFn) ([]byte, error) {
	req := hertzProtocol.AcquireRequest()
	resp := hertzProtocol.AcquireResponse()

	defer hertzProtocol.ReleaseRequest(req)
	defer hertzProtocol.ReleaseResponse(resp)

	err := run(client, req, resp, opts...)
	if err != nil {
		return nil, err
	}

	body := resp.Body()

	if logger.EnabledDebug() {
		if len(body) > 4096 {
			httpLogger.Debug("http response body is greater than 4096b. not printable")
		} else {
			httpLogger.Debug(functionTools.B2S(body))
		}
	}

	return body, nil
}

func run(client *hertzClient.Client, req *hertzProtocol.Request, resp *hertzProtocol.Response, opts ...RequestOptionsFn) (err error) {
	for _, fn := range opts {
		if err = fn(req); err != nil {
			return
		}
	}

	if err = client.Do(context.Background(), req, resp); err != nil {
		httpLogger.Warn(err)
		return
	}

	return
}
