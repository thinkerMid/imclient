package xmpp

import (
	"bytes"
	"crypto/sha256"
	"github.com/chenzhuoyu/base64x"
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
)

// GenerateDeviceHash .
//
//	2:I42C3wPe
func GenerateDeviceHash(jid types.JID) string {
	jid.AD = true
	jidADStr := jid.String()

	buf := bytes.NewBuffer(make([]byte, 0))
	sha := sha256.New()

	buf.WriteString(jidADStr)

	_, _ = sha.Write(buf.Bytes())
	shaValue := sha.Sum(nil)

	buf.Reset()
	buf.Write(shaValue[0:6])

	b64 := base64x.RawStdEncoding.EncodeToString(buf.Bytes())

	return "2:" + b64
}

// BatchGenerateDeviceHash .
func BatchGenerateDeviceHash(jidStrList []string) string {
	buf := bytes.NewBuffer(make([]byte, 0))
	sha := sha256.New()

	for i := range jidStrList {
		buf.WriteString(jidStrList[i])
	}

	_, _ = sha.Write(buf.Bytes())
	shaValue := sha.Sum(nil)

	buf.Reset()
	buf.Write(shaValue[0:6])

	b64 := base64x.RawStdEncoding.EncodeToString(buf.Bytes())

	return "2:" + b64
}

// UsyncIQTemplate .
func UsyncIQTemplate(context containerInterface.IMessageContext, mode, contextStr string, nodes []waBinary.Node) message.InfoQuery {
	return message.InfoQuery{
		ID:        context.GenerateRequestID(),
		Namespace: "usync",
		Type:      message.IqGet,
		To:        context.ResolveJID(),
		Content: []waBinary.Node{{
			Tag: "usync",
			Attrs: waBinary.Attrs{
				"sid":     context.GenerateSID(),
				"mode":    mode,
				"last":    "true",
				"index":   "0",
				"context": contextStr,
			},
			Content: nodes,
		}},
	}
}
