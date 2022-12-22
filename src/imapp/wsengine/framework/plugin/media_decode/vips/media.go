package vips

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"ws/framework/plugin/media_decode/vips/cmd"
)

var (
	MediaCompressError = errors.New("image compress quality error")
)

// CompressImageWithProgress 转换并压缩成jpeg格式
// @quality : 压缩质量
// @progress : 是否渐进式显示
func CompressImageWithProgress(ctx context.Context, dataSource []byte, quality int32, progress bool) ([]byte, error) {
	if quality < 1 || quality > 100 {
		return nil, MediaCompressError
	}

	buff := bytes.NewBuffer(dataSource)
	q := fmt.Sprintf("%d", quality)
	if progress {
		return cmd.VipsReader(ctx, buff, "jpegsave", "stdin", "/dev/stdout", "--interlace", "--Q", q)
	} else {
		return cmd.VipsReader(ctx, buff, "jpegsave", "stdin", "/dev/stdout", "--Q", q)
	}
}
