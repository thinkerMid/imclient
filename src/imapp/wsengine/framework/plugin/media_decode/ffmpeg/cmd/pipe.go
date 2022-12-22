package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
)

// ffmpeg Documentation
// http://ffmpeg.org/ffmpeg-all.html#pipe

var binPath = "ffmpeg"

// SetFFMpegBinPath .
func SetFFMpegBinPath(newBinPath string) {
	binPath = newBinPath
}

// MpegURL .
func MpegURL(ctx context.Context, fileURL string, extraFFMpegOptions ...string) (out []byte, err error) {
	args := append([]string{"-i", fileURL}, extraFFMpegOptions...)
	args = append(args, "pipe:1")

	cmd := exec.CommandContext(ctx, binPath, args...)

	return runMpeg(cmd)
}

// MpegReader .
func MpegReader(ctx context.Context, reader io.Reader, extraFFMpegOptions ...string) (out []byte, err error) {
	args := append([]string{"-i", "pipe:"}, extraFFMpegOptions...)
	args = append(args, "pipe:1")

	cmd := exec.CommandContext(ctx, binPath, args...)
	cmd.Stdin = reader

	return runMpeg(cmd)
}

func runMpeg(cmd *exec.Cmd) (out []byte, err error) {
	var outputBuf bytes.Buffer

	cmd.Stderr = os.Stderr
	cmd.Stdout = &outputBuf

	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("error running %s %w", binPath, err)
	}

	if outputBuf.Len() == 0 {
		return nil, fmt.Errorf("error running %s no output buffer", binPath)
	}

	return outputBuf.Bytes(), nil
}
