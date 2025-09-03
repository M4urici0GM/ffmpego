package ffmpego

import (
	"strconv"
)

// NewScale constructs a ScaleFilter: "[input]scale=width:height[output]"
func NewScale(input, output string, width, height int) FilterComplexParser {
	return ScaleFilter{
		Input:  input,
		Output: output,
		Width:  width,
		Height: height,
	}
}

// NewCrop constructs a CropFilter: "[input]crop=w:h:x:y[output]"
func NewCrop(input, output string, w, h, x, y int) FilterComplexParser {
	return CropFilter{
		Input:  input,
		Output: output,
		W:      w,
		H:      h,
		X:      x,
		Y:      y,
	}
}

// NewRotate constructs a RotateFilter using FFmpeg transpose mode: "[input]transpose=mode[output]"
func NewRotate(input, output string, mode TransposeMode) FilterComplexParser {
	return RotateFilter{
		Input:  input,
		Output: output,
		Mode:   mode,
	}
}

// ScaleMode enumerates preset scale sizes. Extend as needed.
type ScaleMode int

const (
	// 1280x720
	ScaleHD ScaleMode = iota
	// 1920x1080
	ScaleFHD
	// 854x480
	ScaleSD480
)

func (m ScaleMode) dims() (int, int) {
	switch m {
	case ScaleHD:
		return 1280, 720
	case ScaleFHD:
		return 1920, 1080
	case ScaleSD480:
		return 854, 480
	default:
		// Fallback to HD
		return 1280, 720
	}
}

// NewScaleMode constructs a ScaleFilter from a ScaleMode preset.
func NewScaleMode(input, output string, mode ScaleMode) FilterComplexParser {
	w, h := mode.dims()
	return ScaleFilter{
		Input:  input,
		Output: output,
		Width:  w,
		Height: h,
	}
}

// NewSplit constructs a SplitFilter and auto-generates output labels "1", "2", ..., "n".
// Renders: "[input]split=n[1][2]...[n]"
func NewSplit(input string, n int) FilterComplexParser {
	outs := make([]string, n)
	for i := 0; i < n; i++ {
		outs[i] = strconv.Itoa(i + 1)
	}
	return SplitFilter{
		Input:   input,
		N:       n,
		Outputs: outs,
	}
}
