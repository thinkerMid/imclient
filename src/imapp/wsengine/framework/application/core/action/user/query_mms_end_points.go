package user

import (
	"strconv"
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	mmsConstant "ws/framework/application/data_storage/mms/constant"
)

// QueryMMSEndPoints .
type QueryMMSEndPoints struct {
	processor.BaseAction
	intervalUpdate bool // 定时更新
	lastID         string
}

// Start .
func (m *QueryMMSEndPoints) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	iq := message.InfoQuery{
		ID:        context.GenerateRequestID(),
		Namespace: "w:m",
		Type:      "set",
		To:        types.ServerJID,
	}

	if m.intervalUpdate {
		iq.Content = []waBinary.Node{{
			Tag: "media_conn",
			Attrs: waBinary.Attrs{
				"last_id": m.lastID,
			},
		}}
	} else {
		iq.Content = []waBinary.Node{{
			Tag: "media_conn",
		}}
	}

	m.SendMessageId, err = context.SendIQ(iq)

	return
}

// Receive .
func (m *QueryMMSEndPoints) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	defer next()

	nodes := context.Message().GetChildren()
	mediaConnNode := nodes[0]

	attr := mediaConnNode.AttrGetter()

	mmsEndpoints := mmsConstant.MMSEndpoints{}
	mmsEndpoints.ID = attr.String("id")
	mmsEndpoints.Auth = attr.String("auth")

	nodes = mediaConnNode.GetChildren()
	if len(nodes) > 0 {
		var buckets []mmsConstant.Bucket

		for i := range nodes {
			media := nodes[i]

			if media.Tag != "host" {
				continue
			}

			getter := media.AttrGetter()

			buckets = append(buckets, mmsConstant.Bucket{
				Host:         getter.String("hostname"),
				FallbackHost: getter.String("fallback_hostname"),
			})
		}

		mmsEndpoints.Buckets = buckets
	}

	context.ResolveMultimediaMessagingService().UpdateMMSEndpoints(mmsEndpoints)

	if !m.intervalUpdate {
		authTTL := attr.String("auth_ttl")
		authTTLNumber, _ := strconv.ParseUint(authTTL, 10, 64)

		if authTTLNumber == 0 {
			authTTLNumber = 21600
		}

		// MMS定时更新
		context.AddMessageProcessor(processor.NewTimerProcessor(
			func() []containerInterface.IAction {
				return []containerInterface.IAction{&QueryMMSEndPoints{intervalUpdate: true, lastID: mmsEndpoints.ID}}
			},
			processor.Interval(uint32(authTTLNumber)),
			processor.IntervalLoop(true),
		))
	}

	return nil
}

// Error .
func (m *QueryMMSEndPoints) Error(_ containerInterface.IMessageContext, _ error) {}
