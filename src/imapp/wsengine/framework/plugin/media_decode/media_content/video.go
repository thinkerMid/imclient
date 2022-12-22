package mediaContent

import (
	"context"
	"fmt"
	"strings"
	"time"
	"ws/framework/plugin/media_decode/ffmpeg"
	"ws/framework/plugin/media_decode/ffprobe"
)

var supportedVideoCodec = "h264;avc1;"

// Video .
type Video struct {
	ThumbnailJPEG []byte
	Duration      uint32
	Height        uint32
	Width         uint32
	Mimetype      string
}

// NewVideoContent .
func NewVideoContent(container []byte) (Video, error) {
	var i Video

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	reader, err := ffprobe.CodecVideo(ctx, container)
	if err != nil {
		return i, err
	}

	stream := reader.TakeSteam()

	if !strings.Contains(supportedVideoCodec, stream.CodecName) {
		return i, fmt.Errorf("unsupported mime extension: %v", reader.Format.FormatName)
	} else if !strings.Contains(supportedVideoCodec, stream.CodecTagString) {
		return i, fmt.Errorf("unsupported mime extension: %v", reader.Format.FormatName)
	}

	i.Duration = uint32(reader.Format.Duration().Seconds())
	i.Height = stream.Height
	i.Width = stream.Width
	i.Mimetype = "video/mp4"

	if i.Duration == 0 {
		return i, fmt.Errorf("can not parse duration")
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	i.ThumbnailJPEG, _ = ffmpeg.VideoThumbnailJPEG(ctx, container)
	if len(i.ThumbnailJPEG) == 0 {
		return i, fmt.Errorf("can not parse video thumbnail")
	}

	return i, nil
}
