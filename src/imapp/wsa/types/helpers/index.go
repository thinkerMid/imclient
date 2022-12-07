package helpers

import (
	"encoding/binary"
	"reflect"
	"unsafe"
)

// MakePacketLengthByBody .
func MakePacketLengthByBody(header []byte, body []byte) {
	bodyLength := len(body)

	header[0] = byte(bodyLength >> 16)
	header[1] = byte(bodyLength >> 8)
	header[2] = byte(bodyLength)
}

// GenerateIV .
func GenerateIV(count uint32) []byte {
	iv := make([]byte, 12)
	binary.BigEndian.PutUint32(iv[8:], count)
	return iv
}

// CopyString .
func CopyString(s string) string {
	return string(S2B(s))
}

// S2B converts string to a byte slice without memory allocation.
//
// Note it may break if string and/or slice header will change
// in the future go versions.
//
//	*Unsafe*
func S2B(s string) (b []byte) {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	return
}

// B2S converts byte slice to a string without memory allocation.
// See https://groups.google.com/forum/#!msg/Golang-Nuts/ENgbUzYvCuU/90yGx7GUAgAJ .
//
// Note it may break if string and/or slice header will change
// in the future go versions.
//
//	*Unsafe*
func B2S(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// SliceTo32SizeArray converts byte slice to a byte array without memory allocation.
// See https://github.com/golang/go/issues/46505 .
//
// Note it may break if string and/or slice header will change
// in the future go versions.
//
//	*Unsafe*
func SliceTo32SizeArray(slice []byte) [32]byte {
	if len(slice) != 32 || cap(slice) != 32 {
		var array [32]byte
		copy(array[:], slice)
		return array
	}

	return *(*[32]byte)(slice[:32])
}

// ReflectValueTypeName .
func ReflectValueTypeName(obj interface{}) string {
	if reflect.TypeOf(obj).Kind() == reflect.Ptr {
		return reflect.ValueOf(obj).Elem().Type().Name()
	}

	return reflect.TypeOf(obj).Name()
}

// SliceEqual .
func SliceEqual(a, b []byte) bool {
	return B2S(a) == B2S(b)
}

// ArrayEqual .
func ArrayEqual(a []byte, b [32]byte) bool {
	if len(a) != 32 {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}

// ScanDifferentUint8 .
func ScanDifferentUint8(newSlice, oldSlice []uint8) (add []uint8, remove []uint8) {
	newMapping := make(map[uint8]struct{})
	for _, v := range newSlice {
		newMapping[v] = struct{}{}
	}

	oldMapping := make(map[uint8]struct{})
	for _, v := range oldSlice {
		oldMapping[v] = struct{}{}
	}

	for _, v := range newSlice {
		_, haveInOld := oldMapping[v]
		if !haveInOld {
			add = append(add, v)
			break
		}
	}

	for _, v := range oldSlice {
		_, haveInNew := newMapping[v]
		if !haveInNew {
			remove = append(remove, v)
			break
		}
	}

	return
}
