package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
)

var binPath = "vips"

// SetVipsBinPath .
func SetVipsBinPath(newBinPath string) {
	binPath = newBinPath
}

// VipsReader .
func VipsReader(ctx context.Context, reader io.Reader, extraVipsOptions ...string) (out []byte, err error) {
	cmd := exec.CommandContext(ctx, binPath, extraVipsOptions...)
	cmd.Stdin = reader

	return runVips(cmd)
}

func runVips(cmd *exec.Cmd) (out []byte, err error) {
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
