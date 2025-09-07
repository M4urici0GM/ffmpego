package ffmpego

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type CommandRunner interface {
	CommandContext(context.Context, string, ...string) *exec.Cmd
}

type NativeCommandHandler struct {
}

func (nc *NativeCommandHandler) CommandContext(ctx context.Context, name string, args ...string) *exec.Cmd {
	return exec.CommandContext(ctx, name, args...)
}

type FfmpegoRunner struct {
	ffmpego       *Ffmpego
	logger        *log.Logger
	commandRunner CommandRunner
}

func NewRunner(ffmpego *Ffmpego, logger ...*log.Logger) *FfmpegoRunner {
	return &FfmpegoRunner{
		ffmpego:       ffmpego,
		logger:        getLogger(logger...),
		commandRunner: &NativeCommandHandler{},
	}
}

// runWithProgress executes FFmpeg with progress monitoring
func (runner *FfmpegoRunner) runWithProgress(cmd *exec.Cmd) error {
	// Create a pipe for stderr to capture progress
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start ffmpeg: %w", err)
	}

	go parseProgress(stderr, runner.ffmpego.progressCallback)

	// Wait for command completion
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("ffmpeg failed: %w", err)
	}

	return nil
}

// Run executes the FFmpeg command
func (c *FfmpegoRunner) Run(ctx context.Context) error {
	args, err := c.ffmpego.Build()
	if err != nil {
		return err
	}

	// Set timeout if context doesn't have one
	if *c.ffmpego.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, *c.ffmpego.timeout)
		defer cancel()
	}

	cmd := c.commandRunner.CommandContext(ctx, "ffmpeg", args...)
	if c.ffmpego.progressCallback != nil {
		return c.runWithProgress(cmd)
	}

	// Standard execution without progress monitoring
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg failed: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// parseProgress parses FFmpeg progress output and calls the callback
func parseProgress(stderr io.ReadCloser, progressCallback ProgressCallback) {
	defer stderr.Close()

	buffer := make([]byte, 1024)
	for {
		n, err := stderr.Read(buffer)
		if err != nil {
			break
		}

		lines := strings.Split(string(buffer[:n]), "\n")
		progress := Progress{}

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}

			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			switch key {
			case "frame":
				if frame, err := strconv.Atoi(value); err == nil {
					progress.Frame = frame
				}
			case "fps":
				if fps, err := strconv.ParseFloat(value, 64); err == nil {
					progress.FPS = fps
				}
			case "bitrate":
				progress.Bitrate = value
			case "total_size":
				if size, err := strconv.ParseInt(value, 10, 64); err == nil {
					progress.TotalSize = size
				}
			case "out_time":
				progress.OutTime = value
			case "out_time_ms":
				if time, err := strconv.ParseInt(value, 10, 64); err == nil {
					progress.OutTimeMS = time
				}
			case "speed":
				progress.Speed = value
			case "progress":
				progress.Progress = value
			}
		}

		if progressCallback != nil && (progress.Frame > 0 || progress.Progress != "") {
			progressCallback(progress)
		}
	}
}

func getLogger(logger ...*log.Logger) *log.Logger {
	if len(logger) != 0 && logger[0] != nil {
		return logger[0]
	}

	return log.New(os.Stdout, "[packlit] ", log.LstdFlags)
}