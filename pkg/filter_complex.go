package ffmpego

import (
	"fmt"
	"strings"
)

// FilterComplex represents a raw -filter_complex graph string.
// Example: "[0:v]scale=1280:-2[outv]"
type FilterComplex string

// Validate ensures the graph is not empty.
func (f FilterComplex) Validate() error {
	if strings.TrimSpace(string(f)) == "" {
		return fmt.Errorf("filter_complex graph cannot be empty")
	}
	return nil
}

// Parse returns the raw filter graph.
func (f FilterComplex) Parse() string {
	return string(f)
}

// LabeledFilter builds a single filter chain with optional input/output labels.
// It renders as "[in1][in2]expr[out1][out2]".
type LabeledFilter struct {
	Inputs  []string
	Expr    string
	Outputs []string
}

func (lf LabeledFilter) Validate() error {
	if strings.TrimSpace(lf.Expr) == "" {
		return fmt.Errorf("filter expression cannot be empty")
	}
	// Ensure labels are non-empty when present
	for _, in := range lf.Inputs {
		if strings.TrimSpace(in) == "" {
			return fmt.Errorf("input label cannot be empty")
		}
	}
	for _, out := range lf.Outputs {
		if strings.TrimSpace(out) == "" {
			return fmt.Errorf("output label cannot be empty")
		}
	}
	return nil
}

func (lf LabeledFilter) Parse() string {
	var b strings.Builder
	for _, in := range lf.Inputs {
		b.WriteByte('[')
		b.WriteString(in)
		b.WriteByte(']')
	}
	b.WriteString(lf.Expr)
	for _, out := range lf.Outputs {
		b.WriteByte('[')
		b.WriteString(out)
		b.WriteByte(']')
	}
	return b.String()
}
