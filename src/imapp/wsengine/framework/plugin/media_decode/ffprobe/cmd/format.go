package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"time"
	"ws/framework/plugin/json"
)

// ffprobe Documentation
// https://ffmpeg.org/ffprobe.html

var binPath = "ffprobe"

// ProbeData .
type ProbeData struct {
	Packets []packet `json:"packets"`
	Streams []stream `json:"streams"`
	Format  format   `json:"format"`
}

type format struct {
	FormatName      string  `json:"format_name"`
	FormatLongName  string  `json:"format_long_name"`
	DurationSeconds float64 `json:"duration,string"`
	BitRate         int64   `json:"bit_rate,string"`
}

type stream struct {
	CodecName       string  `json:"codec_name"`
	CodecTagString  string  `json:"codec_tag_string"`
	Width           uint32  `json:"width"`
	Height          uint32  `json:"height"`
	DurationSeconds float64 `json:"duration,string"`
}

// packet .
type packet struct {
	PtsTime      string `json:"pts_time"`
	DurationTime string `json:"duration_time"`
}

// TakeSteam .
func (p ProbeData) TakeSteam() stream {
	if len(p.Streams) == 0 {
		return stream{}
	}

	return p.Streams[0]
}

// Duration .
func (f format) Duration() (duration time.Duration) {
	return time.Duration(f.DurationSeconds * float64(time.Second))
}

// Duration .
func (p packet) Duration() time.Duration {
	pts, _ := time.ParseDuration(fmt.Sprintf("%ss", p.PtsTime))
	duration, _ := time.ParseDuration(fmt.Sprintf("%ss", p.DurationTime))

	return pts + duration
}

// Duration .
func (s stream) Duration() time.Duration {
	return time.Duration(s.DurationSeconds * float64(time.Second))
}

// SetFFProbeBinPath .
func SetFFProbeBinPath(newBinPath string) {
	binPath = newBinPath
}

// ProbeURL .
func ProbeURL(ctx context.Context, fileURL string, extraFFProbeOptions ...string) (ProbeData, error) {
	args := append([]string{
		"-loglevel", "fatal",
		"-print_format", "json",
		"-show_format",
	}, extraFFProbeOptions...)

	args = append(args, fileURL)

	cmd := exec.CommandContext(ctx, binPath, args...)

	outputBuf := bytes.Buffer{}
	probe, err := runProbe(cmd, outputBuf)
	if err != nil {
		return ProbeData{}, err
	}

	return probe, nil
}

// ProbeReader .
func ProbeReader(ctx context.Context, reader io.Reader, extraFFProbeOptions ...string) (ProbeData, error) {
	args := append([]string{
		"-loglevel", "fatal",
		"-print_format", "json",
		"-show_format",
	}, extraFFProbeOptions...)

	args = append(args, "-")

	cmd := exec.CommandContext(ctx, binPath, args...)
	cmd.Stdin = reader

	outputBuf := bytes.Buffer{}
	probe, err := runProbe(cmd, outputBuf)
	if err != nil {
		return ProbeData{}, err
	}

	// https://trac.ffmpeg.org/ticket/4358
	// can not parse duration in pipe. if needed that. can search in packet list end
	length := len(probe.Packets)
	if length > 0 {
		probe.Format.DurationSeconds = probe.Packets[length-1].Duration().Seconds()
	}

	return probe, nil
}

func runProbe(cmd *exec.Cmd, outputBuf bytes.Buffer) (data ProbeData, err error) {
	var stdErr bytes.Buffer

	cmd.Stdout = &outputBuf
	cmd.Stderr = &stdErr

	err = cmd.Run()
	if err != nil {
		return data, fmt.Errorf("error running %s [%s] %w", binPath, stdErr.String(), err)
	}

	if stdErr.Len() > 0 {
		return data, fmt.Errorf("ffprobe error: %s", stdErr.String())
	}

	err = json.Unmarshal(outputBuf.Bytes(), &data)
	if err != nil {
		return data, fmt.Errorf("error parsing ffprobe output: %w", err)
	}

	return data, nil
}
