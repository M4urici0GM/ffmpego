package ffmpego

import (
	"strings"
)

// FilterGraph collects FilterComplexParser units before attaching them to Ffmpego.
type FilterGraph struct {
	Options []FilterComplexParser
}

// Add appends a filter unit to the graph.
func (fg *FilterGraph) Add(filter FilterComplexParser) {
	fg.Options = append(fg.Options, filter)
}

type FilterFn = func(*FilterGraph)

// WithFilterExpr adds a raw filter expression (no labels).
// Example: "scale=1280:-2"
func WithFilterExpr(expr string) FilterFn {
	return func(fg *FilterGraph) {
		fg.Add(FilterComplex(strings.TrimSpace(expr)))
	}
}

// WithFilterChain builds a single labeled chain: "[in]expr[out]".
// Empty input/output strings are omitted.
func WithFilterChain(input string, expr string, output string) FilterFn {
	return func(fg *FilterGraph) {
		var ins, outs []string
		if strings.TrimSpace(input) != "" {
			ins = []string{input}
		}
		if strings.TrimSpace(output) != "" {
			outs = []string{output}
		}
		fg.Add(LabeledFilter{Inputs: ins, Expr: expr, Outputs: outs})
	}
}

// WithScale adds a labeled scale filter chain.
// Renders: "[input]scale=width:height[output]"
func WithScale(input string, output string, width, height int) FilterFn {
	return func(fg *FilterGraph) {
		fg.Add(ScaleFilter{
			Input:  strings.TrimSpace(input),
			Output: strings.TrimSpace(output),
			Width:  width,
			Height: height,
		})
	}
}

// WithCrop adds a labeled crop filter chain.
// Renders: "[input]crop=w:h:x:y[output]"
func WithCrop(input string, output string, w, h, x, y int) FilterFn {
	return func(fg *FilterGraph) {
		fg.Add(CropFilter{
			Input:  strings.TrimSpace(input),
			Output: strings.TrimSpace(output),
			W:      w,
			H:      h,
			X:      x,
			Y:      y,
		})
	}
}

// WithRotate adds a labeled transpose-based rotate chain using standard FFmpeg transpose values.
// Common options: 0=90° clockwise and vertical flip, 1=90° clockwise, 2=90° counterclockwise, 3=90° counterclockwise and vertical flip.
// Renders: "[input]transpose=mode[output]"
func WithRotate(input string, output string, transposeMode TransposeMode) FilterFn {
	return func(fg *FilterGraph) {
		fg.Add(RotateFilter{
			Input:  strings.TrimSpace(input),
			Output: strings.TrimSpace(output),
			Mode:   transposeMode,
		})
	}
}

// WithSplit adds a split chain with N outputs.
// Renders: "[input]split=n[out0][out1]...[out{n-1}]"
// The outputs are named by the provided labels. If labels slice length != n, this will validate at build time.
func WithSplit(input string, n int, outputs ...string) FilterFn {
	return func(fg *FilterGraph) {
		outs := make([]string, len(outputs))
		for i, o := range outputs {
			outs[i] = strings.TrimSpace(o)
		}
		fg.Add(SplitFilter{
			Input:   strings.TrimSpace(input),
			N:       n,
			Outputs: outs,
		})
	}
}
