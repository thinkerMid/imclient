package mediaContent

import (
	"context"
	"fmt"
	"time"
	"ws/framework/plugin/media_decode/ffmpeg"
	"ws/framework/plugin/media_decode/ffprobe"
)

var imageMimetype = map[string]string{
	"mjpeg": "image/jpeg",
	"png":   "image/png",
}

// Image .
type Image struct {
	ThumbnailJPEG []byte
	Mimetype      string
	Width         uint32
	Height        uint32
}

// NewImageContent .
func NewImageContent(container []byte) (Image, error) {
	var i Image

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	reader, err := ffprobe.CodecImage(ctx, container)
	if err != nil {
		return i, err
	}

	stream := reader.TakeSteam()
	i.Width = stream.Width
	i.Height = stream.Height

	if typeName, ok := imageMimetype[stream.CodecName]; !ok {
		return i, fmt.Errorf("unsupported mime extension: %v", stream.CodecName)
	} else {
		i.Mimetype = typeName
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	i.ThumbnailJPEG, _ = ffmpeg.ImageThumbnailJPEG(ctx, container)
	if len(i.ThumbnailJPEG) == 0 {
		return i, fmt.Errorf("can not parse thumbnail")
	}

	return i, nil
}
