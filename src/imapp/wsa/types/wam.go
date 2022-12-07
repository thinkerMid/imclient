package types

import (
	"bytes"
	"encoding/binary"
	"math/big"
)

const (
	OPTION_STRING uint8 = 2
	OPTION_NUMBER uint8 = 1
	OPTION_UNKNOW uint8 = 0
	CHANNEL_0           = 0
	CHANNEL_1           = 1
	CHANNEL_2           = 1
)

// IntToBigEndianBytes int转换成大端排列的Byte
func IntToBigEndianBytes(n int) []byte {
	return big.NewInt(int64(n)).Bytes()
}

// IntToLittleEndianBytes int转换成小端排列的Byte
func IntToLittleEndianBytes(n int) []byte {
	b := big.NewInt(int64(n)).Bytes()

	size := len(b)
	maxIdx := size - 1

	for i := 0; i < size/2; i++ {
		b[i], b[maxIdx-i] = b[maxIdx-i], b[i]
	}

	return b
}

// SerializeBuf .
//
//	TODO:可加入bytebuffer池优化
func SerializeBuf(option uint8, attrCode int64, number float64, content string, isHeader int) []byte {
	c := GenUnsignedVarInt(attrCode)
	if (c - 3) >= 2 {
		return make([]byte, 0)
	}
	var retBuff []byte
	var p int

	if option == OPTION_STRING {
		//字符串
		cc := GenUnsignedVarInt(int64(len(content)))
		if (cc - 3) < 3 {
			p = 16*(cc-3) + 128
		} else {
			return make([]byte, 0)
		}
	} else if option == OPTION_NUMBER {
		//number
		if number == 0.0 {
			p = 16
		} else if number == 1.0 {
			p = 32
		} else if number == float64(int(number)) {
			cc := GenSignedVarInt(int64(number))
			p = 16 * cc
		} else {
			p = 112
		}
	} else if option == OPTION_UNKNOW {
		p = 0
	}

	magic := p | isHeader
	if c == 4 {
		magic = magic | 8
	}

	if option != OPTION_UNKNOW {
		magicByte := IntToBigEndianBytes(magic)
		cByte := IntToLittleEndianBytes(int(attrCode))
		retBuff = append(retBuff, magicByte...)
		retBuff = append(retBuff, cByte...)
	}

	if option == OPTION_STRING {
		//字符串
		contentLenByte := IntToLittleEndianBytes(len(content))
		retBuff = append(retBuff, contentLenByte...)
		retBuff = append(retBuff, []byte(content)...)

	} else if option == OPTION_NUMBER || option == OPTION_UNKNOW {
		if number == 0.0 || number == 1.0 || number == -1.0 {
			return retBuff
		} else {
			if number != float64(int(number)) {
				//是小数
				byteBuf := bytes.NewBuffer([]byte{})
				binary.Write(byteBuf, binary.LittleEndian, number)
				tmpByte := byteBuf.Bytes()
				retBuff = append(retBuff, tmpByte...)
			} else {
				//是整数
				tmpByte := IntToLittleEndianBytes(int(number))
				retBuff = append(retBuff, tmpByte...)
				//tmpByte := utils.IntToBytesn(int64(number), width)
				//retBuff = append(retBuff,tmpByte...)
			}
		}
	}

	return retBuff
}

func GenUnsignedVarInt(n int64) int {
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

func GenSignedVarInt(n int64) int {
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
