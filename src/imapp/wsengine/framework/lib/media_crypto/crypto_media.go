package mediaCrypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"time"
	"ws/framework/utils/cbcutil"
	"ws/framework/utils/hkdfutil"
)

var SOSError = errors.New("SOS parse error")
var HashError = errors.New("buffer hash error")
var dataError = errors.New("media data error")
var zipError = errors.New("media zip error")

type MediaType string

const (
	TMediaImage    MediaType = "WhatsApp Image Keys"
	TMediaVideo    MediaType = "WhatsApp Video Keys"
	TMediaAudio    MediaType = "WhatsApp Audio Keys"
	TMediaDocument MediaType = "WhatsApp Document Keys"
	TMediaHistory  MediaType = "WhatsApp History Keys"
	TMediaAppState MediaType = "WhatsApp App State Keys"

	TMediaLinkThumbnail MediaType = "WhatsApp Link Thumbnail Keys"
)

var MediaTypeToMMSType = map[MediaType]string{
	TMediaImage:         "image",
	TMediaAudio:         "ptt",
	TMediaVideo:         "video",
	TMediaDocument:      "document",
	TMediaHistory:       "md-msg-hist",
	TMediaAppState:      "md-app-state",
	TMediaLinkThumbnail: "thumbnail-link",
}

var MediaTypeToMIME = map[MediaType]string{
	TMediaImage: "image",
	TMediaAudio: "ptt",
	TMediaVideo: "video/mp4",
}

type MediaKey struct {
	Iv  []byte
	Enc []byte
	Mac []byte
}

type Media interface {
	parse()
	sidecar([]byte, []float64) []byte
}

type ImageMedia interface {
	Media

	midQualitySha256(sos []float64) []byte
}

type VoiceMedia interface {
	Media
}

type VideoMedia interface {
	Media
}

// ParseMediaKey .
func ParseMediaKey(mediaKey []byte, appInfo MediaType) MediaKey {
	mediaKeyExpanded := hkdfutil.SHA256(mediaKey, nil, []byte(appInfo), 112)
	return MediaKey{
		Iv:  mediaKeyExpanded[:16],
		Enc: mediaKeyExpanded[16:48],
		Mac: mediaKeyExpanded[48:80],
	}
}

func randMediaKey(appInfo MediaType) ([]byte, MediaKey) {
	var keyData MediaKey

	buff := make([]byte, 32)
	_, err := rand.Read(buff)
	if err != nil {
		return nil, keyData
	}

	mediaKeyExpanded := hkdfutil.SHA256(buff, nil, []byte(appInfo), 112)
	return buff, MediaKey{
		Iv:  mediaKeyExpanded[:16],
		Enc: mediaKeyExpanded[16:48],
		Mac: mediaKeyExpanded[48:80],
	}
}

type File struct {
	MediaType MediaType

	RandKey []byte
	KeyData MediaKey
	//CompressText 		[]byte
	//CipherText 			[]byte

	UploadPath        string
	UploadBuff        []byte
	MIME              string
	FileEncSHA256     []byte
	FileSHA256        []byte
	FileLength        uint64
	MediaKeyTimestamp int64
}

func (f *File) sidecar(ciphertext []byte, chunks []float64) []byte {
	const (
		FixedChunkSize = 65536
	)

	var result []byte
	mediaKey := f.KeyData

	// hmac
	h1 := hmac.New(sha256.New, mediaKey.Mac)
	h1.Write(mediaKey.Iv)
	h1.Write(ciphertext)

	// chunks align
	if len(chunks) > 0 {
		single := func(number uint64) uint64 {
			return uint64(number+0xF) & 0xFFFFFFFFFFFFFFF0
		}

		var (
			oriSum  uint64
			castSum uint64
		)

		for idx, ori := range chunks {
			oriSum = oriSum + uint64(ori)

			target := oriSum - castSum
			cast := single(target)

			castSum = castSum + cast
			chunks[idx] = float64(cast)
		}
	}

	tail := h1.Sum(nil)[:10]
	key := mediaKey.Iv
	current := 0

	// 分段生成hash，并截取最后10位
	for idx := 0; idx < len(ciphertext); {
		size := 0
		h1.Reset()

		if len(chunks) > 0 {
			if current >= len(chunks) {
				break
			}

			size = int(chunks[current])
			current += 1
		} else {
			size = FixedChunkSize
			if idx+size > len(ciphertext) {
				size = len(ciphertext) - idx
			}
		}

		text := ciphertext[idx : idx+size]
		h1.Write(key)
		h1.Write(text)

		if idx+size == len(ciphertext) {
			h1.Write(tail)
			//fmt.Println("appen tail:", hex.EncodeToString(tail))
		}

		hash := h1.Sum(nil)[:10]
		result = append(result, hash...)

		//fmt.Println("from:", idx, "to:", idx+size, "front:", hex.EncodeToString(text[:4]), "end:", hex.EncodeToString(text[len(text)-4:]))
		//fmt.Println("hash:", hex.EncodeToString(hash), "curr key:", hex.EncodeToString(key), "next key:", hex.EncodeToString(text[len(text)-16:]))
		//fmt.Println(" ")

		key = text[len(text)-16:]
		idx += size
	}

	return result
}

func (f *File) parse(ciphertext, compressive []byte) {
	mediaKey := f.KeyData

	h := hmac.New(sha256.New, mediaKey.Mac)
	h.Write(mediaKey.Iv)
	h.Write(ciphertext)

	f.UploadPath = fmt.Sprintf("/mms/%s", MediaTypeToMMSType[f.MediaType])
	f.UploadBuff = append(ciphertext, h.Sum(nil)[:10]...)

	fileEncSHA256 := sha256.Sum256(f.UploadBuff)
	f.FileEncSHA256 = fileEncSHA256[:]

	plaintextSHA256 := sha256.Sum256(compressive)
	f.FileSHA256 = plaintextSHA256[:]
	f.FileLength = uint64(len(compressive))
	f.MediaKeyTimestamp = time.Now().Unix()
	f.MIME = MediaTypeToMIME[f.MediaType]
}

type Image struct {
	File

	FirstScanLength       uint32
	LowQualityScanLength  uint32
	MidQualityScanLength  uint32
	FullQualityScanLength uint32
	ScanLengths           []uint32 // sos在图片buff的下标索引
	MidQualitySha256      []byte
	ScanSideCar           []byte // 图片内容按sos块生成hash，取其前10位组成一个40byte的字节数组
}

func (img *Image) midQualitySha256(compressive []byte, sosList []float64) []byte {
	// buffer sha256
	h := sha256.New()

	len1 := int(sosList[1]) / 2
	h.Write(compressive[:len1])
	sum1 := h.Sum(nil)

	len2 := int(sosList[6]) / 2
	h.Write(compressive[len1:len2])
	sum2 := h.Sum(nil)

	len3 := int(sosList[7]) / 2
	h.Write(compressive[len2:len3])
	sum3 := h.Sum(nil)

	_ = [][]byte{sum1, sum2, sum3}
	return sum3
}

func (img *Image) Print() {
	fmt.Println("MidQualitySha256:", hex.EncodeToString(img.MidQualitySha256))
	fmt.Println("ScanSideCar:", hex.EncodeToString(img.ScanSideCar))
	fmt.Println("ScanLengths:", img.ScanLengths)

}

func (img *Image) GetScanLengths() []uint32 {
	return []uint32{img.FirstScanLength, img.LowQualityScanLength, img.MidQualityScanLength}
}

func ParseMediaImage(content []byte) (*Image, error) {
	if cap(content) == 0 {
		return nil, dataError
	}

	buff, sos, err := ParseImage(content)
	if err != nil {
		return nil, err
	}

	randKey, keyData := randMediaKey(TMediaImage)
	ciphertext, _ := cbcutil.Encrypt(keyData.Enc, keyData.Iv, buff)

	img := &Image{
		File: File{
			MediaType: TMediaImage,
			RandKey:   randKey,
			KeyData:   keyData,
			//CipherText:        ciphertext,
			MediaKeyTimestamp: time.Now().Unix(),
			//CompressText:      buff,
		},
	}

	img.parse(ciphertext, buff)

	// mid quality sha256
	img.MidQualitySha256 = img.midQualitySha256(buff, sos)

	// sidecar
	s1 := sos[1] / 2
	s2 := (sos[6] - sos[1]) / 2
	s3 := (sos[7] - sos[6]) / 2
	s4 := (sos[8] - sos[7]) / 2
	qualitys := []float64{s1, s2, s3, s4}
	img.ScanSideCar = img.sidecar(ciphertext, qualitys)

	// quality
	img.FirstScanLength = uint32(qualitys[0])
	img.LowQualityScanLength = uint32(qualitys[1])
	img.MidQualityScanLength = uint32(qualitys[2])
	img.FullQualityScanLength = uint32(qualitys[3])
	img.ScanLengths = []uint32{uint32(qualitys[0]), uint32(qualitys[1]), uint32(qualitys[2]), uint32(qualitys[3])}
	return img, nil
}

type Voice struct {
	File

	PlaySeconds   int32
	PTT           bool
	StreamSideCar []byte
}

func ParseMediaVoice(content []byte) (*Voice, error) {
	if len(content) == 0 {
		return nil, dataError
	}

	randKey, keyData := randMediaKey(TMediaAudio)
	ciphertext, _ := cbcutil.Encrypt(keyData.Enc, keyData.Iv, content)

	voice := &Voice{
		File: File{
			MediaType:         TMediaAudio,
			RandKey:           randKey,
			KeyData:           keyData,
			MediaKeyTimestamp: time.Now().Unix(),
		},
	}

	voice.parse(ciphertext, content)
	voice.StreamSideCar = voice.sidecar(ciphertext, nil)
	return voice, nil
}

type Video struct {
	File

	PlaySeconds   int32
	Height        int32
	Width         int32
	ThumbnailJPEG []byte
	StreamSideCar []byte
}

func ParseMediaVideo(content []byte) (*Video, error) {
	if len(content) == 0 {
		return nil, dataError
	}

	randKey, keyData := randMediaKey(TMediaVideo)
	ciphertext, _ := cbcutil.Encrypt(keyData.Enc, keyData.Iv, content)

	video := &Video{
		File: File{
			MediaType:         TMediaVideo,
			RandKey:           randKey,
			KeyData:           keyData,
			MediaKeyTimestamp: time.Now().Unix(),
		},
	}

	video.parse(ciphertext, content)
	video.StreamSideCar = video.sidecar(ciphertext, nil)
	return video, nil
}
