# ffmpego

A clean, fluent Go library for assembling ffmpeg command arguments with validated output flags and filter_complex graphs.

Key files:
- [pkg/ffmpego.go](pkg/ffmpego.go)
- [pkg/output_flags.go](pkg/output_flags.go)
- [pkg/filter_types.go](pkg/filter_types.go)
- [pkg/filter_builders.go](pkg/filter_builders.go)
- [pkg/filter_complex_builder.go](pkg/filter_complex_builder.go)
- [examples/filter_graph.go](examples/filter_graph.go)

Features

- Fluent API for composing ffmpeg arguments
- Strongly-typed, validated filters and flags
- Single -filter_complex assembled from many graph units
- Output builder for convenient output configuration
- Example program and unit tests

Requirements

- Go 1.21+ (module: m4urici0gm/ffmpego)
- ffmpeg installed and available in PATH

Install

Use as a module import:
- module name: m4urici0gm/ffmpego
- import path for code examples: ffmpego "m4urici0gm/ffmpego/pkg"

Quick start

Basic crop + scale with H.264 and AAC:

```go
package main

import (
	"fmt"
	"strings"

	ffmpego "m4urici0gm/ffmpego/pkg"
)

func main() {
	cmd := ffmpego.New().
		Input("input.mp4").
		WithFilter(ffmpego.NewCrop("0:v", "cropped", 800, 600, 100, 50)).
		WithFilter(ffmpego.NewScale("cropped", "scaled", 1280, 720)).
		Output(
			ffmpego.NewOutputBuilder().
				File("output.mp4").
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

Filter graphs

There are two primary ways to build graphs:

1) Add explicit, validated filter units

- Construct validated filter units then add them with WithFilter/WithFilters:

```go
cmd := ffmpego.New().
	Input("input.mp4").
	WithFilter(ffmpego.NewCrop("0:v", "a", 800, 600, 100, 50)).
	WithFilter(ffmpego.NewRotate("a", "b", ffmpego.Rotate90)).
	WithFilter(ffmpego.NewScale("b", "c", 1280, 720))
```

2) Use the helper builder functions with ComplexFilterBuilder

- Compose using With helper functions; the builder materializes to validated units:

```go
builder := ffmpego.NewComplexFilterBuilder().
	With(ffmpego.WithCrop("0:v", "s1", 800, 600, 100, 50)).
	With(ffmpego.WithRotate("s1", "s2", ffmpego.Rotate90)).
	With(ffmpego.WithScale("s2", "s3", 1280, 720))

cmd := ffmpego.New().Input("input.mp4")
builder.Apply(cmd) // attaches all built units into a single -filter_complex
```

Available filter units

- ScaleFilter → renders: [input]scale=width:height[output]
  - NewScale(input, output string, width, height int)
  - Validation: width/height > 0 or -1/-2 for aspect preservation
- CropFilter → renders: [input]crop=w:h:x:y[output]
  - NewCrop(input, output string, w, h, x, y int)
  - Validation: w,h > 0; x,y ≥ 0
- RotateFilter (transpose) → renders: [input]transpose=mode[output]
  - NewRotate(input, output string, mode TransposeMode)
  - Validation: mode in {0,1,2,3}
  - Provided alias: Rotate90 (clockwise)
- SplitFilter → renders: [input]split=n[out0]...[out{n-1}]
  - NewSplit(input string, n int) auto-generates labels "1".."n"
  - Validation: n ≥ 2; outputs count must match n

Helper functions for graph composition

- WithFilterFn/WithFilterFns on the command accept FilterFn builders
- With helper functions that add validated struct units:
  - WithScale(input, output string, width, height int)
  - WithCrop(input, output string, w, h, x, y int)
  - WithRotate(input, output string, mode TransposeMode)
  - WithSplit(input string, n int, outputs ...string)
  - WithFilterExpr(expr string) // low-level fallback, no labels
  - WithFilterChain(input, expr, output string)

Output configuration

- Create outputs with NewOutputBuilder or NewOutputDescriptor:

```go
out := ffmpego.NewOutputBuilder().
	File("out.mp4").
	WithFlag(ffmpego.VideoCodecH264).
	WithFlag(ffmpego.AudioCodecAAC).
	WithFlag(ffmpego.CRFGoodQuality).
	WithFlag(ffmpego.PresetFast).
	Build()

cmd := ffmpego.New().
	Input("in.mp4").
	Output(out)
```

Common flag presets (all validated before being parsed)

- Codecs:
  - VideoCodecH264, VideoCodecH265, VideoCodecVP9, VideoCodecAV1
  - AudioCodecAAC, AudioCodecMP3, AudioCodecOpus
- Presets:
  - PresetUltraFast, PresetFast, PresetMedium, PresetSlow, PresetVeryslow
- Quality (CRF):
  - CRFHighQuality (18), CRFGoodQuality (23), CRFMediumQuality (28), CRFLowQuality (35)
- Other builders:
  - WithBitrate("2500k"), WithAudioBitrate("128k"), WithFormat("mp4"), WithSampleRate(44100), WithChannels(2), WithMap("[vout]")

Examples

Run the example program:

```bash
go run ./examples/filter_graph.go
```

It prints composed ffmpeg command arguments, including a single -filter_complex argument assembled from multiple units.

Validation model

- Every flag or filter unit implements a Validate() error
- The command Build():
  - Validates all flags and filter units
  - Joins filter units with ';' into a single -filter_complex argument
  - Flattens output flags and file targets

Tests

Run all tests:

```bash
go test ./...
```

Look into:
- [pkg/filter_types_test.go](pkg/filter_types_test.go)
- [pkg/ffmpego_test.go](pkg/ffmpego_test.go)

Notes

- Examples may assume ffmpeg is installed on the system PATH
- This library focuses on building arguments; invoking ffmpeg is up to the caller (e.g., using os/exec) after calling Build()


## Extending the library

This section explains how to add new filter-complex units and new output flags so they integrate cleanly with validation, the builder APIs, and the single -filter_complex assembly.

References:
- Filter unit interface [go.declaration()](pkg/ffmpego.go:14)
- Output flag interface [go.declaration()](pkg/ffmpego.go:8)
- Filter graph helpers [go.declaration()](pkg/filter_builders.go:7)
- Output builder [go.declaration()](pkg/output_builder.go:6)

### Adding a new filter (FilterComplexParser)

Every filter graph unit implements:
- Validate() error — check user inputs up-front
- Parse() string — return one chain/subgraph, e.g. "[in]expr[outs]"

Example: add a simple BoxBlur filter using FFmpeg's boxblur

1) Create a struct type and implement Validate/Parse:

```go
// Example file: pkg/filter_types.go

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

2) (Optional) Provide a constructor for direct use with WithFilter/WithFilters:

```go
// Example file: pkg/filter_constructors.go

func NewBoxBlur(input, output string, lr, lp, cr, cp int) FilterComplexParser {
    return BoxBlurFilter{
        Input: input, Output: output,
        LumaRadius: lr, LumaPower: lp, ChromaRadius: cr, ChromaPower: cp,
    }
}
```

3) (Optional) Provide a helper for the filter graph builder pattern:

```go
// Example file: pkg/filter_builders.go

func WithBoxBlur(input, output string, lr, lp, cr, cp int) FilterFn {
    return func(fg *FilterGraph) {
        fg.Add(BoxBlurFilter{
            Input: input, Output: output,
            LumaRadius: lr, LumaPower: lp, ChromaRadius: cr, ChromaPower: cp,
        })
    }
}
```

Usage:
- Constructors with validation:
  cmd.WithFilter(NewBoxBlur("a", "b", 10, 1, 10, 1))
- ComplexFilterBuilder with helpers:
  NewComplexFilterBuilder().With(WithBoxBlur("a", "b", 10, 1, 10, 1)).Build()

Notes and best practices:
- Keep Validate conservative and informative; fail early during Build().
- Always render a single chain/subgraph in Parse(). The command joins units with ';'.
- Prefer typed fields (int, string) with constraints mirrored in Validate.

See existing filters for reference: [go.declaration()](pkg/filter_types.go:25) ScaleFilter, [go.declaration()](pkg/filter_types.go:56) CropFilter, [go.declaration()](pkg/filter_types.go:86) RotateFilter, [go.declaration()](pkg/filter_types.go:110) SplitFilter.

### Adding a new output flag (OutputFlagParser)

Output flags produce flag-value pairs in the final args and also implement:
- Validate() error
- Parse() []string

For single-value flags, use a type alias when possible to keep things simple.

Example: add pixel format flag (-pix_fmt)

1) Define a type alias and implement Parse/Validate:

```go
// Example file: pkg/output_flags.go

type PixelFormatFlag string

func (f PixelFormatFlag) Parse() []string {
    return []string{"-pix_fmt", string(f)}
}
func (f PixelFormatFlag) Validate() error {
    if f == "" {
        return fmt.Errorf("pixel format cannot be empty")
    }
    return nil
}
```

2) Provide a With... builder function to attach it to outputs:

```go
// Example file: pkg/output_flags.go

func WithPixelFormat(pixfmt string) OutputFlagFn {
    return func(o *OutputDescriptor) {
        o.Add(PixelFormatFlag(pixfmt))
    }
}
```

3) (Optional) Add convenience presets if they are common:

```go
// Example file: pkg/output_flags.go

var (
    // Common pixel formats
    PixFmtYUV420P = WithPixelFormat("yuv420p")
    PixFmtYUV444P = WithPixelFormat("yuv444p")
)
```

Usage with OutputBuilder:
```go
out := ffmpego.NewOutputBuilder().
    File("out.mp4").
    WithFlag(ffmpego.VideoCodecH264).
    WithFlag(ffmpego.CRFGoodQuality).
    WithFlag(ffmpego.PixFmtYUV420P). // or WithFlag(ffmpego.WithPixelFormat("yuv420p"))
    Build()
```

### Where to plug things in

- New filter struct + constructor:
  - Struct and Validate/Parse in filter types: [pkg/filter_types.go](pkg/filter_types.go)
  - Constructor in: [pkg/filter_constructors.go](pkg/filter_constructors.go)
  - Helper builder in: [pkg/filter_builders.go](pkg/filter_builders.go)
- New output flag:
  - Type alias + Parse/Validate + WithXXX in: [pkg/output_flags.go](pkg/output_flags.go)

### Testing your additions

- Add unit tests alongside existing suites:
  - Filter units: [pkg/filter_types_test.go](pkg/filter_types_test.go)
  - Command integration paths: [pkg/ffmpego_test.go](pkg/ffmpego_test.go)
- Typical tests should cover:
  - Validate() passes with correct inputs and fails with bad inputs (clear error messages)
  - Parse() exact rendering (compare strings/slices)
  - End-to-end Build() produces expected ffmpeg args and single -filter_complex joining with ';'

Run:
```bash
go test ./...
```
