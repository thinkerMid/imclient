package connection

import (
	"bytes"
	"compress/zlib"
	"crypto/cipher"
	"fmt"
	"io"
	functionTools "ws/framework/utils/function_tools"
)

type bufferCodec struct {
	readKey  cipher.AEAD
	writeKey cipher.AEAD

	writeCounter uint32
	readCounter  uint32
}

func newBufferCodec(readKey cipher.AEAD, writeKey cipher.AEAD) bufferCodec {
	return bufferCodec{readKey: readKey, writeKey: writeKey}
}

// 解密
func (m *bufferCodec) decode(body []byte) (decodeBody []byte, err error) {
	decodeBody, err = m.readKey.Open(body[:0], functionTools.GenerateIV(m.readCounter), body, nil)
	if err != nil {
		return
	}

	m.readCounter++

	return
}

// 加密
func (m *bufferCodec) encode(src []byte, dst []byte) {
	m.writeKey.Seal(dst[:0], functionTools.GenerateIV(m.writeCounter), src, nil)

	m.writeCounter++
}

func unpack(data []byte) ([]byte, error) {
	if len(data) < 1 {
		return nil, fmt.Errorf("not enough bytes to unpack")
	}

	dataType, data := data[0], data[1:]
	if 2&dataType > 0 {
		if decompressor, err := zlib.NewReader(bytes.NewReader(data)); err != nil {
			return nil, fmt.Errorf("failed to create zlib reader: %w", err)
		} else if data, err = io.ReadAll(decompressor); err != nil {
			return nil, err
		}
	}

	return data, nil
}
