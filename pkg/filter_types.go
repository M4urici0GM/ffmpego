package ffmpego

import (
	"fmt"
	"strings"
)

// TransposeMode enumerates FFmpeg transpose options for rotation.
type TransposeMode int

const (
	// 0 = 90° clockwise and vertical flip
	TransposeClockFlip TransposeMode = 0
	// 1 = 90° clockwise
	TransposeClockwise TransposeMode = 1
	// 2 = 90° counterclockwise
	TransposeCounterClockwise TransposeMode = 2
	// 3 = 90° counterclockwise and vertical flip
	TransposeCounterClockFlip TransposeMode = 3

	// Convenience alias commonly used for 90° clockwise
	Rotate90 TransposeMode = TransposeClockwise
)

// ScaleFilter renders: "[input]scale=width:height[output]"
type ScaleFilter struct {
	Input  string
	Output string
	Width  int
	Height int
}

func (f ScaleFilter) Validate() error {
	if strings.TrimSpace(f.Input) == "" {
		return fmt.Errorf("scale: input label cannot be empty")
	}
	if strings.TrimSpace(f.Output) == "" {
		return fmt.Errorf("scale: output label cannot be empty")
	}
	// FFmpeg scale accepts positive ints or -1/-2 for aspect preservation.
	validDim := func(v int) bool { return v > 0 || v == -1 || v == -2 }
	if !validDim(f.Width) {
		return fmt.Errorf("scale: width must be >0 or -1 or -2, got %d", f.Width)
	}
	if !validDim(f.Height) {
		return fmt.Errorf("scale: height must be >0 or -1 or -2, got %d", f.Height)
	}
	return nil
}

func (f ScaleFilter) Parse() string {
	return fmt.Sprintf("[%s]scale=%d:%d[%s]", f.Input, f.Width, f.Height, f.Output)
}

// CropFilter renders: "[input]crop=w:h:x:y[output]"
type CropFilter struct {
	Input  string
	Output string
	W      int
	H      int
	X      int
	Y      int
}

func (f CropFilter) Validate() error {
	if strings.TrimSpace(f.Input) == "" {
		return fmt.Errorf("crop: input label cannot be empty")
	}
	if strings.TrimSpace(f.Output) == "" {
		return fmt.Errorf("crop: output label cannot be empty")
	}
	if f.W <= 0 || f.H <= 0 {
		return fmt.Errorf("crop: width and height must be positive, got %dx%d", f.W, f.H)
	}
	if f.X < 0 || f.Y < 0 {
		return fmt.Errorf("crop: x and y must be non-negative, got %d,%d", f.X, f.Y)
	}
	return nil
}

func (f CropFilter) Parse() string {
	return fmt.Sprintf("[%s]crop=%d:%d:%d:%d[%s]", f.Input, f.W, f.H, f.X, f.Y, f.Output)
}

// RotateFilter renders: "[input]transpose=mode[output]"
type RotateFilter struct {
	Input  string
	Output string
	Mode   TransposeMode
}

func (f RotateFilter) Validate() error {
	if strings.TrimSpace(f.Input) == "" {
		return fmt.Errorf("rotate: input label cannot be empty")
	}
	if strings.TrimSpace(f.Output) == "" {
		return fmt.Errorf("rotate: output label cannot be empty")
	}
	if f.Mode < 0 || f.Mode > 3 {
		return fmt.Errorf("rotate: invalid transpose mode %d (expected 0..3)", f.Mode)
	}
	return nil
}

func (f RotateFilter) Parse() string {
	return fmt.Sprintf("[%s]transpose=%d[%s]", f.Input, int(f.Mode), f.Output)
}

// SplitFilter renders: "[input]split=n[out0][out1]...[out{n-1}]"
type SplitFilter struct {
	Input   string
	N       int
	Outputs []string
}

func (f SplitFilter) Validate() error {
	if strings.TrimSpace(f.Input) == "" {
		return fmt.Errorf("split: input label cannot be empty")
	}
	if f.N < 2 {
		return fmt.Errorf("split: n must be >= 2, got %d", f.N)
	}
	if len(f.Outputs) != f.N {
		return fmt.Errorf("split: number of outputs (%d) must match n (%d)", len(f.Outputs), f.N)
	}
	for i, o := range f.Outputs {
		if strings.TrimSpace(o) == "" {
			return fmt.Errorf("split: output label at index %d cannot be empty", i)
		}
	}
	return nil
}

func (f SplitFilter) Parse() string {
	return fmt.Sprintf("[%s]split=%d[%s]", f.Input, f.N, strings.Join(f.Outputs, "]["))
}
