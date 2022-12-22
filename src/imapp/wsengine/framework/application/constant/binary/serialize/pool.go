package nodeSerialize

import (
	"ws/framework/application/constant/binary"
	byteBufferPool "ws/framework/plugin/byte_buffer_pool"
)

// Marshal .
func Marshal(n waBinary.Node) []byte {
	encoder := newEncoder(byteBufferPool.AcquireBuffer())

	body := encoder.Marshal(n)

	byteBufferPool.ReleaseBuffer(encoder.B)

	return body
}

// Unmarshal .
func Unmarshal(data []byte) (*waBinary.Node, error) {
	r := newDecoder(data)
	n, err := r.readNode()
	if err != nil {
		return nil, err
	}
	return n, nil
}
