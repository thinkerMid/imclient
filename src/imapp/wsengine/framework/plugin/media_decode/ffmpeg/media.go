package ffmpeg

import (
	"bytes"
	"context"
	"ws/framework/plugin/media_decode/ffmpeg/cmd"
)

// ConvertOpusAudio .
// ffmpeg aac ogg mp3 wav => opus
func ConvertOpusAudio(ctx context.Context, dataSource []byte) ([]byte, error) {
	buf := bytes.NewBuffer(dataSource)

	// https://zhuanlan.zhihu.com/p/191518845
	return cmd.MpegReader(ctx, buf, "-b:a", "16k", "-ar", "24k", "-mapping_family", "1", "-compression_level", "6", "-f", "opus")
}

// VideoThumbnailJPEG .
func VideoThumbnailJPEG(ctx context.Context, dataSource []byte) ([]byte, error) {
	buf := bytes.NewBuffer(dataSource)

	return cmd.MpegReader(ctx, buf, "-f", "image2", "-frames", "1", "-vf", "scale=320:568", "-q", "75")
}

// ImageThumbnailJPEG .
func ImageThumbnailJPEG(ctx context.Context, dataSource []byte) ([]byte, error) {
	buf := bytes.NewBuffer(dataSource)

	return cmd.MpegReader(ctx, buf, "-f", "image2", "-vf", "scale=320:568", "-q", "75")
}

// ImageToJPEG .
func ImageToJPEG(ctx context.Context, dataSource []byte) ([]byte, error) {
	buf := bytes.NewBuffer(dataSource)

	return cmd.MpegReader(ctx, buf, "-f", "image2")
}
