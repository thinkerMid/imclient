package mediaContent

import (
	"context"
	"fmt"
	"strings"
	"time"
	"ws/framework/plugin/media_decode/ffprobe"
)

var supportedAudioMimeTypes = "ogg;mp3;wav;"

// Audio .
type Audio struct {
	Duration uint32
	Mimetype string
}

// NewAudioContent .
func NewAudioContent(container []byte) (Audio, error) {
	var i Audio

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	reader, err := ffprobe.CodecAudio(ctx, container)
	if err != nil {
		return i, err
	}

	if !strings.Contains(supportedAudioMimeTypes, reader.Format.FormatName) {
		return i, fmt.Errorf("unsupported mime extension: %v", reader.Format.FormatName)
	}

	i.Duration = uint32(reader.Format.Duration().Seconds())

	if i.Duration == 0 {
		return i, fmt.Errorf("can not parse duration")
	}

	i.Mimetype = "audio/mpeg"

	return i, nil
}
