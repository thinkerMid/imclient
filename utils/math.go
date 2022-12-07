package utils

import (
	"crypto/rand"
	"math/big"
)

// RandBytes 随机生成指定大小字节
func RandBytes(n int) []byte {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return b
}

// RandInt64 随机范围的数字 [min, max]
func RandInt64(min int64, max int64) int64 {
	var val *big.Int
	var err error
	var n int64 = -1

	maxInt := big.NewInt(max + 1)

	for n < min || n > max {
		val, err = rand.Int(rand.Reader, maxInt)
		if err != nil {
			n = 0
			continue
		}

		n = val.Int64()
		if n < min {
			n += min
		}
	}

	return n
}

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

// BigEndianBytesToInt 大端排列的Byte转换成int
func BigEndianBytesToInt(buf []byte) int {
	b := big.NewInt(0).SetBytes(buf)

	return int(b.Int64())
}
