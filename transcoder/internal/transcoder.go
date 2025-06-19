package internal

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"transcoder/transcoder/cmd/config"
)

type Transcoder struct {
	streamId    string
	ffmpegCmd   *exec.Cmd
	ffmpegStdin io.WriteCloser
	isActive    bool
}

func NewTranscoder(streamId string, config config.HLSConfig) (*Transcoder, error) {
	cmd, err := buildCommand(streamId, config)
	if err != nil {
		return nil, err
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create pipe: %v", err)
	}

	if err = cmd.Start(); err != nil {
		stdin.Close()
		return nil, fmt.Errorf("failed to start ffmpeg: %v", err)
	}

	return &Transcoder{
		streamId:    streamId,
		ffmpegCmd:   cmd,
		ffmpegStdin: stdin,
		isActive:    true,
	}, nil
}

func buildDirectories(streamId string, config config.HLSConfig) (string, error) {
	timestamp := time.Now().Format(time.RFC3339)
	streamDir := filepath.Join(config.OutputDir, streamId, timestamp)

	for _, r := range config.Resolutions {
		if err := os.MkdirAll(filepath.Join(streamDir, r.Name), os.ModePerm); err != nil {
			return "", fmt.Errorf("failed to create directory [%v]: %v", r, err)
		}
	}

	return streamDir, nil
}

func (t *Transcoder) Write(data []byte) error {
	_, err := t.ffmpegStdin.Write(data)
	return err
}

func (t *Transcoder) Stop() error {
	t.isActive = false
	t.ffmpegStdin.Close()
	return t.ffmpegCmd.Wait()
}

func buildCommand(streamId string, config config.HLSConfig) (*exec.Cmd, error) {
	streamDir, err := buildDirectories(streamId, config)
	if err != nil {
		return nil, err
	}

	args := []string{
		"-i", "pipe:0",
		"-copyts",
		"-loglevel", "error", "-report",
	}

	for _, r := range config.Resolutions {
		args = append(
			args,
			// map video and audio from input 0
			"-map", "0:v",
			"-map", "0:a",

			// video codec
			"-c:v", "libx264",

			// audio coded and bitrate
			"-c:a", "aac",
			"-b:a", "192k",

			// video bitrate
			"-b:v", r.Bitrate, // 6000k
			// video resolution
			"-s", fmt.Sprintf("%dx%d", r.Width, r.Height), // 1920x1080
			// video framerate
			"-r", fmt.Sprintf("%d", r.Framerate), // 60

			// HLS output
			"-f", "hls",
			"-hls_time", "2",
			"-hls_segment_filename",
			filepath.Join(streamDir, r.Name, "segment_%03d.ts"),
			filepath.Join(streamDir, r.Name, "playlist.m3u8"),
		)
	}

	return exec.Command("ffmpeg", args...), nil
}
