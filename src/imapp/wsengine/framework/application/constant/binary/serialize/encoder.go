package nodeSerialize

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/binary/token"
	"ws/framework/application/constant/types"
)

const tagSize = 1

// BinaryEncoder .
type BinaryEncoder struct {
	B *bytes.Buffer
}

func newEncoder(b *bytes.Buffer) BinaryEncoder {
	return BinaryEncoder{B: b}
}

// Marshal .
func (w *BinaryEncoder) Marshal(n waBinary.Node) []byte {
	w.pushByte(0) // zlib type

	w.writeNode(n) // marshal body

	return w.B.Bytes()
}

func (w *BinaryEncoder) pushByte(b byte) {
	w.B.WriteByte(b)
}

func (w *BinaryEncoder) pushBytes(bytes []byte) {
	w.B.Write(bytes)
}

func (w *BinaryEncoder) pushIntN(value, n int, littleEndian bool) {
	for i := 0; i < n; i++ {
		var curShift int
		if littleEndian {
			curShift = i
		} else {
			curShift = n - i - 1
		}
		w.pushByte(byte((value >> uint(curShift*8)) & 0xFF))
	}
}

func (w *BinaryEncoder) pushInt20(value int) {
	w.pushBytes([]byte{byte((value >> 16) & 0x0F), byte((value >> 8) & 0xFF), byte(value & 0xFF)})
}

func (w *BinaryEncoder) pushInt8(value int) {
	w.pushIntN(value, 1, false)
}

func (w *BinaryEncoder) pushInt16(value int) {
	w.pushIntN(value, 2, false)
}

func (w *BinaryEncoder) pushInt32(value int) {
	w.pushIntN(value, 4, false)
}

func (w *BinaryEncoder) pushString(value string) {
	w.B.WriteString(value)
}

func (w *BinaryEncoder) writeByteLength(length int) {
	if length < 256 {
		w.pushByte(token.Binary8)
		w.pushInt8(length)
	} else if length < (1 << 20) {
		w.pushByte(token.Binary20)
		w.pushInt20(length)
	} else if length < math.MaxInt32 {
		w.pushByte(token.Binary32)
		w.pushInt32(length)
	} else {
		panic(fmt.Errorf("length is too large: %d", length))
	}
}

func (w *BinaryEncoder) writeNode(n waBinary.Node) {
	if n.Tag == "0" {
		w.pushByte(token.List8)
		w.pushByte(token.ListEmpty)
		return
	}

	hasContent := 0
	if n.Content != nil {
		hasContent = 1
	}

	w.writeListStart(2*len(n.Attrs) + tagSize + hasContent)
	w.writeString(n.Tag)
	w.writeAttributes(n.Attrs)
	if n.Content != nil {
		w.write(n.Content)
	}
}

func (w *BinaryEncoder) write(data interface{}) {
	switch typedData := data.(type) {
	case nil:
		w.pushByte(token.ListEmpty)
	case types.JID:
		w.writeJID(typedData)
	case string:
		w.writeString(typedData)
	case int:
		w.writeString(strconv.Itoa(typedData))
	case int32:
		w.writeString(strconv.FormatInt(int64(typedData), 10))
	case uint:
		w.writeString(strconv.FormatUint(uint64(typedData), 10))
	case uint32:
		w.writeString(strconv.FormatUint(uint64(typedData), 10))
	case int64:
		w.writeString(strconv.FormatInt(typedData, 10))
	case uint64:
		w.writeString(strconv.FormatUint(typedData, 10))
	case bool:
		w.writeString(strconv.FormatBool(typedData))
	case []byte:
		w.writeBytes(typedData)
	case []waBinary.Node:
		w.writeListStart(len(typedData))
		for _, n := range typedData {
			w.writeNode(n)
		}
	default:
		panic(fmt.Errorf("%w: %T", waBinary.ErrInvalidType, typedData))
	}
}

func (w *BinaryEncoder) writeString(data string) {
	var dictIndex byte
	if tokenIndex, ok := token.IndexOfSingleToken(data); ok {
		w.pushByte(tokenIndex)
	} else if dictIndex, tokenIndex, ok = token.IndexOfDoubleByteToken(data); ok {
		w.pushByte(token.Dictionary0 + dictIndex)
		w.pushByte(tokenIndex)
	} else if validateNibble(data) {
		w.writePackedBytes(data, token.Nibble8)
	} else if validateHex(data) {
		w.writePackedBytes(data, token.Hex8)
	} else {
		w.writeStringRaw(data)
	}
}

func (w *BinaryEncoder) writeBytes(value []byte) {
	w.writeByteLength(len(value))
	w.pushBytes(value)
}

func (w *BinaryEncoder) writeStringRaw(value string) {
	w.writeByteLength(len(value))
	w.pushString(value)
}

func (w *BinaryEncoder) writeJID(jid types.JID) {
	if jid.AD {
		w.pushByte(token.ADJID)
		w.pushByte(jid.Agent)
		w.pushByte(jid.Device)
		w.writeString(jid.User)
	} else {
		w.pushByte(token.JIDPair)
		if len(jid.User) == 0 {
			w.pushByte(token.ListEmpty)
		} else {
			w.write(jid.User)
		}
		w.write(jid.Server)
	}
}

func (w *BinaryEncoder) writeAttributes(attributes waBinary.Attrs) {
	if attributes == nil {
		return
	}

	for key, val := range attributes {
		w.writeString(key)

		if val == "" || val == nil {
			w.writeStringRaw("")
			continue
		}

		w.write(val)
	}
}

func (w *BinaryEncoder) writeListStart(listSize int) {
	if listSize == 0 {
		w.pushByte(byte(token.ListEmpty))
	} else if listSize < 256 {
		w.pushByte(byte(token.List8))
		w.pushInt8(listSize)
	} else {
		w.pushByte(byte(token.List16))
		w.pushInt16(listSize)
	}
}

func (w *BinaryEncoder) writePackedBytes(value string, dataType byte) {
	if len(value) > token.PackedMax {
		panic(fmt.Errorf("too many bytes to pack: %d", len(value)))
	}

	w.pushByte(dataType)

	roundedLength := byte(math.Ceil(float64(len(value)) / 2.0))
	if len(value)%2 != 0 {
		roundedLength |= 128
	}
	w.pushByte(roundedLength)
	var packer func(byte) byte
	if dataType == token.Nibble8 {
		packer = packNibble
	} else if dataType == token.Hex8 {
		packer = packHex
	} else {
		// This should only be called with the correct values
		panic(fmt.Errorf("invalid packed byte data type %v", dataType))
	}
	for i, l := 0, len(value)/2; i < l; i++ {
		w.pushByte(w.packBytePair(packer, value[2*i], value[2*i+1]))
	}
	if len(value)%2 != 0 {
		w.pushByte(w.packBytePair(packer, value[len(value)-1], '\x00'))
	}
}

func (w *BinaryEncoder) packBytePair(packer func(byte) byte, part1, part2 byte) byte {
	return (packer(part1) << 4) | packer(part2)
}

func validateNibble(value string) bool {
	if len(value) > token.PackedMax {
		return false
	}
	for _, char := range value {
		if !(char >= '0' && char <= '9') && char != '-' && char != '.' {
			return false
		}
	}
	return true
}

func packNibble(value byte) byte {
	switch value {
	case '-':
		return 10
	case '.':
		return 11
	case 0:
		return 15
	default:
		if value >= '0' && value <= '9' {
			return value - '0'
		}
		// This should be validated beforehand
		panic(fmt.Errorf("invalid string to pack as nibble: %d / '%s'", value, string(value)))
	}
}

func validateHex(value string) bool {
	if len(value) > token.PackedMax {
		return false
	}
	for _, char := range value {
		if !(char >= '0' && char <= '9') && !(char >= 'A' && char <= 'F') {
			return false
		}
	}
	return true
}

func packHex(value byte) byte {
	switch {
	case value >= '0' && value <= '9':
		return value - '0'
	case value >= 'A' && value <= 'F':
		return 10 + value - 'A'
	case value == 0:
		return 15
	default:
		// This should be validated beforehand
		panic(fmt.Errorf("invalid string to pack as hex: %d / '%s'", value, string(value)))
	}
}
