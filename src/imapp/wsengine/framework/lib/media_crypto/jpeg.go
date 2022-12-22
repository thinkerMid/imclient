package mediaCrypto

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
	"ws/framework/plugin/media_decode/vips"
)

var sosError = errors.New("sos parse error")

const (
	sep = "ffda00"
	end = "ffd9"
	max = 8

	msgImageQuality  = 75
	msgImageProgress = true
)

// 返回分隔符在原切片中的所有下标
func bytesIndex(s, sep []byte) []float64 {
	var idx []float64

	arr := bytes.Split(s, sep)
	if len(arr) == 0 {
		fmt.Println("arr empty")
		return idx
	}

	length := 0

	for n, item := range arr {
		if n == 0 {
			length = len(item) * 2
		} else if n == len(arr)-1 {
			break
		} else {
			length += len(sep)*2 + len(item)*2
		}
		idx = append(idx, float64(length))
	}
	return idx
}

func sosParse(buff []byte) (sos []float64, err error) {
	sepTag, _ := hex.DecodeString(sep)
	endTag, _ := hex.DecodeString(end)

	idxList := bytesIndex(buff, sepTag)
	endList := bytesIndex(buff, endTag)
	if len(idxList) == 0 || len(endList) != 1 {
		return nil, sosError
	}

	idxList = append(idxList, endList...)
	if len(idxList) < max+1 {
		return nil, sosError
	}

	return idxList, nil
	//s1 := idxList[1] / 2
	//s2 := idxList[6] - idxList[1] / 2
	//s3 := idxList[7] - idxList[6] / 2
	//s4 := idxList[8] - idxList[7] / 2
	//return []float64{s1, s2, s3, s4}, nil
}

// ParseImage .
func ParseImage(content []byte) (buff []byte, sos []float64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	if buff, err = vips.CompressImageWithProgress(ctx, content, msgImageQuality, msgImageProgress); err != nil {
		return
	}

	sos, err = sosParse(buff)

	return
}
