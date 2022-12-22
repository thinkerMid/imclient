package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"github.com/chenzhuoyu/base64x"
	"github.com/google/uuid"
	"golang.org/x/crypto/curve25519"
	"io"
	"math/big"
	mathRand "math/rand"
	"strings"
	"time"
	"ws/framework/plugin/curve25519_voi/curve"
	"ws/framework/plugin/curve25519_voi/curve/scalar"
	"ws/framework/plugin/curve25519_voi/pkg/elligator"
	"ws/framework/plugin/curve25519_voi/pkg/field"
)

// 获取时间戳
func GetCurTime() int64 {
	timeUnix := time.Now().Unix()
	return timeUnix
}

// 生成uuid4
func GenUUID4() string {
	u4 := uuid.New()
	return strings.ToUpper(u4.String())
}

func ParseUUID4(content string) []byte {
	b, _ := uuid.Parse(content)
	return b[:]
}

// 随机生成指定大小字节
func RandBytes(n int) []byte {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return b
}

// RandInt64 随机范围的数字
//
//	结果包含 min max
func RandInt64(min int64, max int64) int64 {
	var n int64 = -1

	src := mathRand.NewSource(time.Now().UnixNano())
	r := mathRand.New(src)

	max = max + 1

	for n < min || n > max {
		n = r.Int63n(max)
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

func StreamToByte(stream io.Reader) []byte {
	var buf []byte
	io.ReadFull(stream, buf)
	return buf
}

// base64编码
func Base64Encode(content []byte) string {
	return base64x.RawURLEncoding.EncodeToString(content)
}

// base64解码
func UrlBase64Decode(content string) []byte {
	sDec, _ := base64x.RawURLEncoding.DecodeString(content)
	return sDec
}

func GenURLParams(content []byte) string {
	newContent := strings.ReplaceAll(string(content), "{", "")
	newContent = strings.ReplaceAll(newContent, "}", "")
	newContent = strings.ReplaceAll(newContent, "\"", "")
	newContent = strings.ReplaceAll(newContent, ":", "=")
	newContent = strings.ReplaceAll(newContent, ",", "&")
	return newContent
}

func URLEncode(content string) string {
	ret := escape(content, 1)
	return ret
}

type encoding int

const (
	encodePath encoding = 1 + iota
	encodePathSegment
	encodeHost
	encodeZone
	encodeUserPassword
	encodeQueryComponent
	encodeFragment
)

const upperhex = "0123456789ABCDEF"

func escape(s string, mode encoding) string {
	spaceCount, hexCount := 0, 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if shouldEscape(c, mode) {
			if c == ' ' && mode == encodeQueryComponent {
				spaceCount++
			} else {
				hexCount++
			}
		}
	}

	if spaceCount == 0 && hexCount == 0 {
		return s
	}

	var buf [64]byte
	var t []byte

	required := len(s) + 2*hexCount
	if required <= len(buf) {
		t = buf[:required]
	} else {
		t = make([]byte, required)
	}

	if hexCount == 0 {
		copy(t, s)
		for i := 0; i < len(s); i++ {
			if s[i] == ' ' {
				t[i] = '+'
			}
		}
		return string(t)
	}

	j := 0
	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case c == ' ' && mode == encodeQueryComponent:
			t[j] = '+'
			j++
		case shouldEscape(c, mode):
			t[j] = '%'
			t[j+1] = upperhex[c>>4]
			t[j+2] = upperhex[c&15]
			j += 3
		default:
			t[j] = s[i]
			j++
		}
	}
	return string(t)
}

func shouldEscape(c byte, mode encoding) bool {
	// Everything else must be escaped.
	return true
}

func RandPushToken() string {
	buff := RandBytes(32)
	return hex.EncodeToString(buff)
}

// CreateBackupKey .
func CreateBackupKey() []byte {
	backupKey := RandBytes(16)

	h := sha256.New()
	h.Write(backupKey)

	buff := h.Sum(nil)
	return buff
}

func MD5Hex(content string) string {
	h := md5.New()
	h.Write([]byte(content))
	return hex.EncodeToString(h.Sum(nil))
}

func GenCurve25519KeyPair() ([]byte, []byte) {
	var priKey, pubKey [32]byte

	//用随机数填满私钥
	_, err := rand.Reader.Read(priKey[:])
	if err != nil {
		return nil, nil
	}

	curve25519.ScalarBaseMult(&pubKey, &priKey)

	return pubKey[:], priKey[:]
}

func CalCurve25519Signature(priKey []byte, message []byte) []byte {
	for {
		// when provided a low-order point, ScalarMult will set dst to all
		// zeroes, irrespective of the scalar.
		signature, err := curve25519.X25519(priKey, message)
		if err == nil {
			return signature
		}
	}
}

func AesGcmEncrypt(aesKey []byte, content []byte) []byte {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil
	}

	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	nonce := make([]byte, 12)

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil
	}

	return aesgcm.Seal(nil, nonce, content, nil)
}

// GenerateEd25519Credential .
func GenerateEd25519Credential() (scalarBytes, randomBytes, result []byte, err error) {
	// random point
	randomBytes = RandBytes(32)

	temp := sha512.Sum512(randomBytes)
	signBit := (temp[31] & 0x80) >> 7
	temp[31] &= 0x7F

	r, err := new(field.Element).SetBytes(temp[:field.ElementSize])
	if err != nil {
		return
	}

	u, _ := elligator.MontgomeryFlavor(r)
	err = u.ToBytes(temp[:field.ElementSize])
	if err != nil {
		return
	}

	montgomeryPoint, _ := new(curve.MontgomeryPoint).SetBytes(temp[:curve.CompressedPointSize])

	y, err := new(curve.EdwardsPoint).SetMontgomery(montgomeryPoint, signBit)
	if err != nil {
		return
	}

	// random scalar
	s, err := scalar.New().SetRandom(nil)
	if err != nil {
		return
	}

	scalarBytes, err = s.MarshalBinary()
	if err != nil {
		return
	}

	copy(temp[:], scalarBytes)
	temp[0] &= 248
	temp[31] &= 63
	temp[31] |= 64
	_, _ = s.SetBytesModOrder(temp[:scalar.ScalarSize])

	p1 := new(curve.EdwardsPoint).Mul(curve.ED25519_BASEPOINT_POINT, s)
	p2 := new(curve.EdwardsPoint).MulByCofactor(y)
	p3 := p1.Add(p1, p2)
	result, err = p3.MarshalBinary()
	return
}

// GenerateEd25519Signature .
func GenerateEd25519Signature(scalarBytes, acsPublicKey, signedCredential []byte) (result []byte, err error) {
	clampingScalar := make([]byte, len(scalarBytes))
	copy(clampingScalar, scalarBytes)
	clampingScalar[0] &= 248
	clampingScalar[31] &= 63
	clampingScalar[31] |= 64

	s, err := scalar.NewFromBytesModOrder(clampingScalar[:])
	if err != nil {
		return
	}

	p1, p2 := new(curve.EdwardsPoint), new(curve.EdwardsPoint)

	err = p1.UnmarshalBinary(acsPublicKey)
	if err != nil {
		return
	}

	p1.Neg(p1)
	p1.Mul(p1, s)

	err = p2.UnmarshalBinary(signedCredential)
	if err != nil {
		return
	}

	result, err = p1.Add(p1, p2).MarshalBinary()
	return
}
