# ffmpego

A clean, fluent Go library for assembling ffmpeg command arguments with validated output flags and filter_complex graphs. Includes a lightweight runner to execute ffmpeg and optionally parse progress output.

Key files:
- [pkg/ffmpego.go](pkg/ffmpego.go)
- [pkg/flags.go](pkg/flags.go)
- [pkg/flag_helpers.go](pkg/flag_helpers.go)
- [pkg/filter_types.go](pkg/filter_types.go)
- [pkg/filter_builders.go](pkg/filter_builders.go)
- [pkg/filter_complex_builder.go](pkg/filter_complex_builder.go)
- [pkg/output_flags.go](pkg/output_flags.go)
- [pkg/output_builder.go](pkg/output_builder.go)
- [pkg/executor.go](pkg/executor.go)
- Examples:
  - [examples/default/](examples/default/)
  - [examples/filter_graph/](examples/filter_graph/)
  - [examples/multiple_outputs/](examples/multiple_outputs/)

Features

- Fluent API for composing ffmpeg arguments
- Strongly-typed, validated filters and flags
- Single -filter_complex assembled from many graph units
- Output builder for convenient output configuration
- Optional runner to invoke ffmpeg and parse progress
- Example programs and unit tests

Requirements

- Go 1.24+
- ffmpeg installed and available in PATH

Install

Use as a module:
- Module name: m4urici0gm/ffmpego
- Import path for code: ffmpego "m4urici0gm/ffmpego/pkg"

Go get:
```bash
go get m4urici0gm/ffmpego
```

Quick start

Build arguments only:
```go
package main

import (
	"fmt"
	"strings"

	ffmpego "m4urici0gm/ffmpego/pkg"
)

func main() {
	cmd := ffmpego.New(""). // pass "" to use the default binary
		WithOptions(
			ffmpego.NewFfmpegOptions(
				ffmpego.WithInput("in.mp4"),
				ffmpego.WithOverwrite(),
			),
		).
		WithFilterGraph(
			ffmpego.NewComplexFilterBuilder().
				Add(ffmpego.WithCrop("0:v", "cropped", 800, 600, 100, 50)).
				Add(ffmpego.WithScale("cropped", "scaled", 1280, 720)).
				Build(),
		).
		Output(
			ffmpego.NewOutputBuilder().
				File("out.mp4").
				WithFlag(ffmpego.VideoCodecH264).
				WithFlag(ffmpego.AudioCodecAAC).
				WithFlag(ffmpego.CRFGoodQuality).
				WithFlag(ffmpego.PresetMedium).
				Build(),
		)

	args, err := cmd.Build()
	if err != nil {
		panic(err)
	}
	fmt.Println("ffmpeg " + strings.Join(args, " "))
}
```

Execute with the built-in runner:
```go
package main

import (
	"context"
	"log"
	ffmpego "m4urici0gm/ffmpego/pkg"
)

func main() {
	cmd := ffmpego.New("").
		WithOptions(
			ffmpego.NewFfmpegOptions(
				ffmpego.WithInput("in.mp4"),
				ffmpego.WithOverwrite(),
				// Optional: write progress key=value lines to stdout (pipe:1)
				ffmpego.PipeProgress,
			),
		).
		WithFilterGraph(
			ffmpego.NewComplexFilterBuilder().
				Add(ffmpego.WithScale("0:v", "scaled", 1280, 720)).
				Build(),
		).
		Output(
			ffmpego.NewOutputBuilder().
				File("out.mp4").
				WithFlag(ffmpego.VideoCodecH264).
				WithFlag(ffmpego.CRFGoodQuality).
				Build(),
		).
		// Optional: receive parsed progress updates (best when progress is enabled)
		WithProgressCallback(func(p ffmpego.Progress) {
			log.Printf("frame=%d fps=%.2f speed=%s progress=%s", p.Frame, p.FPS, p.Speed, p.Progress)
		})

	ctx := context.Background()
	if err := ffmpego.NewRunner(cmd).Run(ctx); err != nil {
		log.Fatalf("ffmpeg failed: %v", err)
	}
}
```

Filter graphs

Compose with the FilterGraphBuilder:

- Builder: [go.declaration()](pkg/filter_complex_builder.go:6)
- New builder: [go.declaration()](pkg/filter_complex_builder.go:10)
- Filter graph container: [go.declaration()](pkg/filter_builders.go:8)

Example:
```go
filterGraph := ffmpego.NewComplexFilterBuilder().
	Add(ffmpego.WithCrop("0:v", "s1", 800, 600, 100, 50)).
	Add(ffmpego.WithRotate("s1", "s2", ffmpego.Rotate90)).
	Add(ffmpego.WithScale("s2", "s3", 1280, 720)).
	Build()
```

Available labeled helper filters (validated before build):

- WithScale: [go.declaration()](pkg/filter_builders.go:55)
  - Renders: "[input]scale=width:height[output]"
  - Validation: width/height > 0 or -1/-2 for aspect preservation
- WithCrop: [go.declaration()](pkg/filter_builders.go:68)
  - Renders: "[input]crop=w:h:x:y[output]"
  - Validation: w,h > 0; x,y ≥ 0
- WithRotate: [go.declaration()](pkg/filter_builders.go:84)
  - Renders: "[input]transpose=mode[output]"
  - Validation: mode in {0,1,2,3}; alias Rotate90
- WithSplit: [go.declaration()](pkg/filter_builders.go:97)
  - Renders: "[input]split=n[out0]...[out{n-1}]"
  - Validation: n ≥ 2; outputs count must match n

Low-level helpers:

- WithFilterExpr: [go.declaration()](pkg/filter_builders.go:32)
- WithFilterChain: [go.declaration()](pkg/filter_builders.go:40)

Output configuration

Use OutputBuilder to compose output flags:

- OutputBuilder: [go.declaration()](pkg/output_builder.go:5)
- NewOutputBuilder: [go.declaration()](pkg/output_builder.go:10)
- WithFlag: [go.declaration()](pkg/output_builder.go:17)
- File: [go.declaration()](pkg/output_builder.go:23)
- OutputDescriptor: [go.declaration()](pkg/output_flags.go:34)

Example:
```go
out := ffmpego.NewOutputBuilder().
	File("out.mp4").
	WithFlag(ffmpego.VideoCodecH264).
	WithFlag(ffmpego.AudioCodecAAC).
	WithFlag(ffmpego.CRFGoodQuality).
	WithFlag(ffmpego.PresetFast).
	Build()

cmd := ffmpego.New("").
	WithOptions(ffmpego.NewFfmpegOptions(ffmpego.WithInput("in.mp4"))).
	Output(out)
```

Common flag presets (all validated)

- Codecs:
  - Video: VideoCodecH264, VideoCodecH265, VideoCodecVP9, VideoCodecAV1 [go.declaration()](pkg/output_flags.go:9)
  - Audio: AudioCodecAAC, AudioCodecMP3, AudioCodecOpus [go.declaration()](pkg/output_flags.go:15)
- Presets: PresetUltraFast, PresetFast, PresetMedium, PresetSlow, PresetVeryslow [go.declaration()](pkg/output_flags.go:20)
- Quality (CRF): CRFHighQuality(18), CRFGoodQuality(23), CRFMediumQuality(28), CRFLowQuality(35) [go.declaration()](pkg/output_flags.go:27)
- Other builders:
  - WithBitrate [go.declaration()](pkg/output_flags.go:109)
  - WithAudioBitrate [go.declaration()](pkg/output_flags.go:130)
  - WithFormat [go.declaration()](pkg/output_flags.go:123)
  - WithSampleRate [go.declaration()](pkg/output_flags.go:137)
  - WithChannels [go.declaration()](pkg/output_flags.go:144)
  - WithMap [go.declaration()](pkg/output_flags.go:151)
  - WithFile (output path) [go.declaration()](pkg/output_flags.go:172)

Global ffmpeg options

Build via FfmpegOptions:

- FfmpegOptions: [go.declaration()](pkg/flags.go:3)
- NewFfmpegOptions: [go.declaration()](pkg/flags.go:7)

Common helpers:
- WithInput: [go.declaration()](pkg/flag_helpers.go:7)
- WithOverwrite: [go.declaration()](pkg/flag_helpers.go:14)
- WithLogLevel: [go.declaration()](pkg/flag_helpers.go:20)
- WithProgress and preset PipeProgress: [go.declaration()](pkg/flag_helpers.go:26)

Example:
```go
opts := ffmpego.NewFfmpegOptions(
	ffmpego.WithInput("in.mp4"),
	ffmpego.WithOverwrite(),
	ffmpego.WithLogLevel("info"),
	// Writes progress key=value lines to stdout; pair with a progress callback.
	ffmpego.PipeProgress,
)
```

Runner (optional execution)

- Runner type: [go.declaration()](pkg/executor.go:25)
- NewRunner: [go.declaration()](pkg/executor.go:31)
- Run: [go.declaration()](pkg/executor.go:63)
- ProgressCallback on builder: [go.declaration()](pkg/ffmpego.go:78)

Notes:
- The runner invokes the "ffmpeg" binary on PATH with the built args.
- To use a custom binary, execute the args yourself (e.g., with os/exec).

Interfaces

- OutputFlagParser [go.declaration()](pkg/ffmpego.go:5)
- FilterComplexParser [go.declaration()](pkg/ffmpego.go:14)

Examples

Run the examples:
```bash
# Basic crop/scale example
go run ./examples/default

# Filter graph builder example
go run ./examples/filter_graph

# Multiple outputs using split
go run ./examples/multiple_outputs
```

Validation model

- Every flag or filter unit implements Validate() error
- Command Build():
  - Validates all flags and filter units
  - Joins filter units with ';' into a single -filter_complex argument
  - Flattens output flags and file targets

Tests

Run:
```bash
go test ./...
```

Notable suites:
- [pkg/ffmpego_test.go](pkg/ffmpego_test.go)

Extending the library

Add new filter-complex units that implement:
- Validate() error
- Parse() string (returns one chain/subgraph, e.g. "[in]expr[outs]")

References:
- Filter unit interface [go.declaration()](pkg/ffmpego.go:14)
- Filter graph helpers [go.declaration()](pkg/filter_builders.go:7)
- Output flag interface [go.declaration()](pkg/ffmpego.go:5)
- Output builder [go.declaration()](pkg/output_builder.go:5)

Example: add BoxBlur filter using FFmpeg's boxblur

```go
// In pkg/filter_types.go
type BoxBlurFilter struct {
	Input        string
	Output       string
	LumaRadius   int // lr
	LumaPower    int // lp
	ChromaRadius int // cr
	ChromaPower  int // cp
}

func (f BoxBlurFilter) Validate() error {
	if strings.TrimSpace(f.Input) == "" {
		return fmt.Errorf("boxblur: input label cannot be empty")
	}
	if strings.TrimSpace(f.Output) == "" {
		return fmt.Errorf("boxblur: output label cannot be empty")
	}
	if f.LumaRadius < 0 || f.LumaPower < 0 || f.ChromaRadius < 0 || f.ChromaPower < 0 {
		return fmt.Errorf("boxblur: all parameters must be non-negative")
	}
	return nil
}

func (f BoxBlurFilter) Parse() string {
	// boxblur=lr:lp:cr:cp
	return fmt.Sprintf("[%s]boxblur=%d:%d:%d:%d[%s]",
		f.Input, f.LumaRadius, f.LumaPower, f.ChromaRadius, f.ChromaPower, f.Output)
}
```

Add a builder helper:
```go
// In pkg/filter_builders.go
func WithBoxBlur(input, output string, lr, lp, cr, cp int) FilterFn {
	return func(fg *FilterGraph) {
		fg.Add(BoxBlurFilter{
			Input:        strings.TrimSpace(input),
			Output:       strings.TrimSpace(output),
			LumaRadius:   lr,
			LumaPower:    lp,
			ChromaRadius: cr,
			ChromaPower:  cp,
		})
	}
}
```

Usage:
```go
fg := ffmpego.NewComplexFilterBuilder().
	Add(ffmpego.WithBoxBlur("0:v", "b1", 10, 1, 10, 1)).
	Build()
```

Best practices:
- Keep Validate conservative and informative; fail early during Build().
- Always render a single chain/subgraph in Parse(). The command joins units with ';'.
- Prefer typed fields (int, string) with constraints mirrored in Validate.
