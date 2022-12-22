package eventSerialize

import (
	"bytes"
	byteBufferPool "ws/framework/plugin/byte_buffer_pool"
)

// IEventBuffer .
type IEventBuffer interface {
	Reset()
	Byte() []byte
	ByteBuffer() *bytes.Buffer
	Println(tag ...interface{})
	Header() IEventBuffer
	Body() IEventBuffer
	Footer() IEventBuffer
	Common() IEventBuffer
	SerializeString(attrCode int64, text string) IEventBuffer    // 序列化事件参数
	SerializeNumber(attrCode int64, number float64) IEventBuffer // 序列化事件参数
	Write(buf []byte)
	WriteByte(b byte)
	WriteString(s string)
	WriteLittleEndianInt16(n int16, removeZero bool)
	WriteLittleEndianUint16(n uint16, removeZero bool)
	WriteLittleEndianInt32(n int32, removeZero bool)
	WriteLittleEndianUint32(n uint32, removeZero bool)
	WriteLittleEndianInt64(n int64, removeZero bool)
	WriteLittleEndianUint64(n uint64, removeZero bool)
	WriteLittleEndianFloat64(n float64, removeZero bool)
}

// AcquireEventBuffer returns an empty byte buffer from the pool.
//
// Got byte buffer may be returned to the pool via Put call.
// This reduces the number of memory allocations required for byte buffer
// management.
func AcquireEventBuffer() IEventBuffer {
	return newEventBuffer(byteBufferPool.AcquireBuffer())
}

// ReleaseEventBuffer returns byte buffer to the pool.
//
// ByteBuffer.B mustn't be touched after returning it to the pool.
// Otherwise data races will occur.
func ReleaseEventBuffer(e IEventBuffer) {
	byteBufferPool.ReleaseBuffer(e.ByteBuffer())
}
