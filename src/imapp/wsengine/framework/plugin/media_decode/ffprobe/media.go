package ffprobe

import (
	"bytes"
	"context"
	"ws/framework/plugin/media_decode/ffprobe/cmd"
)

// CodecImage .
func CodecImage(ctx context.Context, dataSource []byte) (cmd.ProbeData, error) {
	buf := bytes.NewBuffer(dataSource)

	return cmd.ProbeReader(ctx, buf, "-show_streams")
}

// CodecVideo .
func CodecVideo(ctx context.Context, dataSource []byte) (cmd.ProbeData, error) {
	buf := bytes.NewBuffer(dataSource)

	return cmd.ProbeReader(ctx, buf, "-show_streams")
}

// CodecAudio .
func CodecAudio(ctx context.Context, dataSource []byte) (cmd.ProbeData, error) {
	buf := bytes.NewBuffer(dataSource)

	return cmd.ProbeReader(ctx, buf, "-show_entries", "packet=pts_time,duration_time")
}
