package eventSerialize

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
)

const (
	common uint8 = iota
	header
	body
	footer
)

// eventBuffer .
type eventBuffer struct {
	b *bytes.Buffer

	empty   []byte
	convert []byte

	stage uint8 // 序列化日志的时候会用到 用于切换写入日志的位置类型 头部 内容 尾部 基础包
}

func newEventBuffer(b *bytes.Buffer) *eventBuffer {
	return &eventBuffer{
		b: b,
		// max 8: float64 is 8b
		empty:   make([]byte, 8),
		convert: make([]byte, 8),
	}
}

// Reset .
func (e *eventBuffer) Reset() {
	e.b.Reset()
}

// Byte .
func (e *eventBuffer) Byte() []byte {
	return e.b.Bytes()
}

// ByteBuffer .
func (e *eventBuffer) ByteBuffer() *bytes.Buffer {
	return e.b
}

// Println .
func (e *eventBuffer) Println(tag ...interface{}) {
	fmt.Println(tag, "buffer hex:", hex.EncodeToString(e.b.Bytes()))
}

// Header .
func (e *eventBuffer) Header() IEventBuffer {
	e.stage = header
	return e
}

// Body .
func (e *eventBuffer) Body() IEventBuffer {
	e.stage = body
	return e
}

// Footer .
func (e *eventBuffer) Footer() IEventBuffer {
	e.stage = footer
	return e
}

// Common .
func (e *eventBuffer) Common() IEventBuffer {
	e.stage = common
	return e
}

// SerializeString 序列化事件参数
func (e *eventBuffer) SerializeString(attrCode int64, text string) IEventBuffer {
	// TODO 字符串0值完全不写?
	textSize := len(text)
	if textSize == 0 {
		return e
	}

	p := e.genStringPValue(int64(textSize))
	if p == 0 { // 0值的p不写
		return e
	}

	if !e.writeMagicValue(p, attrCode) {
		return e
	}

	e.writeString(text)

	return e
}

// SerializeNumber 序列化事件参数
func (e *eventBuffer) SerializeNumber(attrCode int64, number float64) IEventBuffer {
	p := e.genNumberPValue(number)

	if !e.writeMagicValue(p, attrCode) {
		return e
	}

	// TODO 0|1|-1 的number只写magic?
	if number == 0 || number == 1 || number == -1 {
		return e
	}

	e.writeNumber(number)

	return e
}

// p max value of 112
func (e *eventBuffer) genNumberPValue(number float64) (p uint8) {
	if number == 0.0 {
		p = 16
	} else if number == 1.0 {
		p = 32
	} else if number == float64(int64(number)) {
		cc := genSignedVarInt(int64(number))
		p = 16 * cc
	} else {
		p = 112
	}

	return p
}

// p max value of 176
func (e *eventBuffer) genStringPValue(strLen int64) (p uint8) {
	cc := genUnsignedVarInt(strLen)
	if (cc - 3) >= 3 {
		return
	}

	p = 16*(cc-3) + 128

	return
}

func (e *eventBuffer) writeMagicValue(p uint8, attrCode int64) (success bool) {
	c := genUnsignedVarInt(attrCode)
	if (c - 3) >= 2 {
		return
	}

	// magic max value for 190 (190 = 176 | 2 | 8 | 4)
	var magic uint8

	switch e.stage {
	// 基础包的
	case common:
		magic = p | 0
	case header:
		magic = p | 1
	default:
		magic = p | 2
	}

	if c == 4 {
		magic = magic | 8
	}

	if e.stage == footer {
		magic = magic | 4
	}

	e.b.WriteByte(magic)
	e.WriteLittleEndianInt64(attrCode, true)

	return true
}

func (e *eventBuffer) writeString(content string) {
	e.WriteLittleEndianInt64(int64(len(content)), true)
	e.b.WriteString(content)
}

func (e *eventBuffer) writeNumber(number float64) {
	intNum := int64(number)

	if number != float64(intNum) {
		e.WriteLittleEndianFloat64(number, true)
	} else {
		e.WriteLittleEndianInt64(intNum, true)
	}
}

func genUnsignedVarInt(n int64) uint8 {
	if n > 0xFF {
		if (n >> 16) > 0 {

			if (n >> 32) > 0 {
				//gen_uint64
				return 6 // 8
			} else {
				// *(_DWORD *)sub_31E224(a1, 4) = a2;
				return 5 // 4
			}

		} else {
			// *(_WORD *)sub_31E224(a1, 2) = a2;
			return 4 // 2
		}
	} else {
		// *sub_31E224(a1, 1) = a2;
		return 3 //1
	}
}

func genSignedVarInt(n int64) uint8 {
	if n == n&0xFF {
		return 3
	} else if n == n&0xFFFF {
		return 4
	} else if n == n&0xFFFFFFFF {
		return 5
	} else {
		return 6
	}
}

// ----------------------------------------------------------------------------

//region 底层写操作API

// Write .
func (e *eventBuffer) Write(buf []byte) {
	e.b.Write(buf)
}

// WriteByte .
func (e *eventBuffer) WriteByte(b byte) {
	e.b.WriteByte(b)
}

// WriteString .
func (e *eventBuffer) WriteString(s string) {
	e.b.WriteString(s)
}

// WriteLittleEndianInt16 .
func (e *eventBuffer) WriteLittleEndianInt16(n int16, removeZero bool) {
	binary.LittleEndian.PutUint16(e.convert, uint16(n))

	e.flushConvertBuffer(2, removeZero)
}

// WriteLittleEndianUint16 .
func (e *eventBuffer) WriteLittleEndianUint16(n uint16, removeZero bool) {
	binary.LittleEndian.PutUint16(e.convert, n)

	e.flushConvertBuffer(2, removeZero)
}

// WriteLittleEndianInt32 .
func (e *eventBuffer) WriteLittleEndianInt32(n int32, removeZero bool) {
	binary.LittleEndian.PutUint32(e.convert, uint32(n))

	e.flushConvertBuffer(4, removeZero)
}

// WriteLittleEndianUint32 .
func (e *eventBuffer) WriteLittleEndianUint32(n uint32, removeZero bool) {
	binary.LittleEndian.PutUint32(e.convert, n)

	e.flushConvertBuffer(4, removeZero)
}

// WriteLittleEndianInt64 .
func (e *eventBuffer) WriteLittleEndianInt64(n int64, removeZero bool) {
	binary.LittleEndian.PutUint64(e.convert, uint64(n))

	e.flushConvertBuffer(8, removeZero)
}

// WriteLittleEndianUint64 .
func (e *eventBuffer) WriteLittleEndianUint64(n uint64, removeZero bool) {
	binary.LittleEndian.PutUint64(e.convert, n)

	e.flushConvertBuffer(8, removeZero)
}

// WriteLittleEndianFloat64 .
func (e *eventBuffer) WriteLittleEndianFloat64(n float64, removeZero bool) {
	binary.LittleEndian.PutUint64(e.convert, math.Float64bits(n))

	e.flushConvertBuffer(8, removeZero)
}

func (e *eventBuffer) flushConvertBuffer(bitSize int, removeZero bool) {
	if removeZero {
		for i := bitSize - 1; i >= 0; i-- {
			if e.convert[i] == 0 {
				continue
			}

			e.b.Write(e.convert[:i+1])
			break
		}
	} else {
		e.b.Write(e.convert[:bitSize])
	}

	copy(e.convert, e.empty)
}

// endregion
