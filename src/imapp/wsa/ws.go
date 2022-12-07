package wsa

import (
	"context"
	"crypto/cipher"
	"github.com/cloudwego/netpoll"
	"github.com/hyahm/golog"
	"labs/src/imapp"
	"labs/src/imapp/wsa/config"
	. "labs/src/imapp/wsa/handshake"
	. "labs/src/imapp/wsa/register"
	"labs/src/imapp/wsa/types"
	"labs/src/imclient"
	"net"
	"time"
)

type WSA struct {
	imapp.App

	cancel context.CancelFunc

	JId       types.JID
	register  Register
	handshake Handshake
	conn      net.Conn

	encoder cipher.AEAD
	decoder cipher.AEAD
}

func CreateWSApp(eventChan chan imclient.EventId) *WSA {
	return &WSA{}
}

func (wa *WSA) Info() string {
	return "WhatsApp 2.23.76"
}

func (wa *WSA) Kill() {
	if wa.cancel != nil {
		wa.cancel()
	}
}

func (wa *WSA) Launch() {
	if err := wa.register.Run(); err != nil {
		return
	}

	r, w, err := wa.handshake.Do()
	if err != nil {
		return
	}

	_, _ = r, w
	func() {
		defer func() {
			golog.Info("whatsapp run over")
		}()
		golog.Info("whatsapp run start")

		ctx, cancel := context.WithCancel(context.Background())
		wa.cancel = cancel

		for true {
			select {
			case <-ctx.Done():
				return
			}
		}
	}()
	return
}

func (wa *WSA) Connect(edgeRouting []byte) error {
	dialer := netpoll.NewDialer()
	conn, err := dialer.DialConnection("tcp", config.WSUrl, 10*time.Second)
	if err != nil {
		return err
	}

	_ = conn.SetReadTimeout(10 * time.Second)
	_ = conn.SetIdleTimeout(10 * time.Second)

	wa.conn = conn

	//_, _, ok := wa.handshake.handshake(conn, edgeRouting)
	//if ok != nil {
	//	return ok
	//}
	//
	//golog.Infof("user:%v handshake success\n", wa.JId.String())
	//
	//_ = conn.SetOnRequest(wa.onRead)
	//_ = conn.AddCloseCallback(wa.onClosed)
	//
	//wa.decoder = r
	//wa.encoder = w
	//
	//wa.eventChan <- events.Event{
	//	Id: events.EventConnOk,
	//}
	return nil
}

func (wa *WSA) Write(buff []byte) error {
	_, err := wa.conn.Write(buff)
	if err != nil {
		return err
	}
	return nil
}

func (wa *WSA) DisConnect() {
	_ = wa.conn.Close()
}

func (wa *WSA) onClosed(connection netpoll.Connection) error {
	//wa.eventChan <- events.Event{
	//	Id: events.EventConnClosed,
	//}
	return nil
}

func (wa *WSA) onRead(_ context.Context, connection netpoll.Connection) (err error) {
	reader := connection.Reader()
	defer func() {
		_ = reader.Release()
		if err != nil {
			golog.Infof("user:%v read err:%v\n", err)
		}
	}()
	/*
		// 消息头
		header, err := reader.Next(consts.FrameLengthSize)
		if err != nil {
			return err
		}
		// 消息体
		l := (int(header[0]) << 16) + (int(header[1]) << 8) + int(header[2])
		next, err := reader.Next(l)
		if err != nil {
			return err
		}
		if len(next) == 0 {
			err = fmt.Errorf("onRead receive length zero")
			return err
		}
		// 消息解密
		buffer := wa.decode(next)
		defer bytePool.Free(buffer)

		// xmpp解压缩
		decompressed, err := waBinary.Unpack(buffer)
		if err != nil {
			return err
		}
		// xmpp反序列化
		node, err := waBinary.Unmarshal(decompressed)
		if err != nil {
			return err
		}

		wa.eventChan <- events.Event{
			Id:   events.EventMsg,
			Xmpp: *node,
		}
	*/
	return nil
}

// 解密
func (wa *WSA) decode(body []byte) []byte {
	/*dst := bytePool.Alloc(len(body))

	plaintext, _ := wa.decoder.Open(dst[:0], functionTools.GenerateIV(wa.r), body, nil)
	wa.r++

	return plaintext
	*/
	return nil
}

// 加密
func (wa *WSA) encode(body []byte) []byte {
	/*dst := bytePool.Alloc(len(body) + 16)
	defer bytePool.Free(dst)

	ciphertext := wa.encoder.Seal(dst[:0], functionTools.GenerateIV(wa.w), body, nil)
	wa.w++

	encodeBuffer := functionTools.ComposeHeaderAndBody(ciphertext, true)
	return encodeBuffer*/
	return nil
}
