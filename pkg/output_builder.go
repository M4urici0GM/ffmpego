package ffmpego

// OutputBuilder provides a fluent API to compose output options similarly to examples/default.go.
// It collects OutputFlagFn builders and produces an OutputDescriptor on Build().
type OutputBuilder struct {
	opts []OutputFlagFn
}

// NewOutputBuilder creates a new empty OutputBuilder.
func NewOutputBuilder() *OutputBuilder {
	return &OutputBuilder{
		opts: make([]OutputFlagFn, 0),
	}
}

// WithFlag appends an output option builder (e.g., VideoCodecH264, WithBitrate("2M")).
func (b *OutputBuilder) WithFlag(opt OutputFlagFn) *OutputBuilder {
	b.opts = append(b.opts, opt)
	return b
}

// File sets the output file path (appends WithFile(path)).
func (b *OutputBuilder) File(path string) *OutputBuilder {
	b.opts = append(b.opts, WithFile(path))
	return b
}

// Build materializes an OutputDescriptor with all collected options.
// It returns the descriptor by value for compatibility with Ffmpego.Output(OutputDescriptor).
func (b *OutputBuilder) Build() OutputDescriptor {
	desc := NewOutputDescriptor(b.opts...)
	if desc == nil {
		return OutputDescriptor{}
	}
	return *desc
}
