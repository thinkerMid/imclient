package webhook

import (
	"crypto/tls"
	"fmt"
	hertzConst "github.com/cloudwego/hertz/pkg/protocol/consts"
	"time"
	networkConstant "ws/framework/plugin/network/constant"
	httpApi "ws/framework/plugin/network/http_api"
	"ws/framework/plugin/network/netpoll"
	functionTools "ws/framework/utils/function_tools"
)

// IWebhookEndpoint .
type IWebhookEndpoint interface {
	Push(content string)
}

// ----------------------------------------------------------------------------

var defaultEndpoint IWebhookEndpoint

func init() {
	defaultEndpoint = &defaultWebhookEndpoint{}
}

// SetWebhookEndpoint .
func SetWebhookEndpoint(endpoint IWebhookEndpoint) {
	defaultEndpoint = endpoint
}

// ----------------------------------------------------------------------------

type defaultWebhookEndpoint struct{}

// Push .
func (e *defaultWebhookEndpoint) Push(content string) {
	client := netpoll.HTTP(networkConstant.ConnectionConfig{
		Tls: &tls.Config{InsecureSkipVerify: true},
	})

	_, _ = httpApi.Do(
		client,
		httpApi.Url("https://inbots.zoom.us/incoming/hook/XB1FfnX-XTjuuWUL00MTp5cX"),
		httpApi.Method(hertzConst.MethodPost),
		httpApi.Header("Authorization", "9858Z6AUm16ST-XjdzmRQcFE"),
		httpApi.Body(functionTools.S2B(content)),
	)

	client.CloseIdleConnections()
}

// ----------------------------------------------------------------------------

// PanicPushTemplate .
type PanicPushTemplate struct {
	JID                 string
	Message             string
	ProcessorID         uint32
	ProcessorAlisaName  string
	ProcessorType       string
	CurrentActionName   string
	CurrentActionStatus string
	ActionQueue         string
	PanicError          string
	StackDumpID         string
	Time                time.Time
}

// ----------------------------------------------------------------------------

// Push 基础推送
func Push(content string) {
	defaultEndpoint.Push(content)
}

// PanicPush 崩溃推送
func PanicPush(t PanicPushTemplate) {
	defaultEndpoint.Push(fmt.Sprintf(`Panic: %s
Context:
  [JID]: %s
  [Message]: %s
  [ProcessorID]: %v
  [ProcessorAlisaName]: %s
  [ProcessorType]: %v
  [CurrentActionName]: %s
  [CurrentActionStatus]: %s
  [ActionQueue]: %s
  [StackDumpID]: %s
  [Time]: %s`,
		t.PanicError,
		t.JID,
		t.Message,
		t.ProcessorID,
		t.ProcessorAlisaName,
		t.ProcessorType,
		t.CurrentActionName,
		t.CurrentActionStatus,
		t.ActionQueue,
		t.StackDumpID,
		t.Time.UTC().Format("2006-01-02T15:04:05Z")),
	)
}
