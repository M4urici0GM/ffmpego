package ffmpego

// ComplexFilterBuilder provides a fluent API to compose a filter_complex graph
// using the same builder pattern as OutputBuilder. It collects FilterFn helpers
// and can either build concrete FilterComplexParser units or apply them to a command.
type ComplexFilterBuilder struct {
	fns []FilterFn
}

// NewComplexFilterBuilder creates a new empty ComplexFilterBuilder.
func NewComplexFilterBuilder() *ComplexFilterBuilder {
	return &ComplexFilterBuilder{
		fns: make([]FilterFn, 0),
	}
}

// With appends any custom FilterFn.
func (b *ComplexFilterBuilder) With(fn FilterFn) *ComplexFilterBuilder {
	b.fns = append(b.fns, fn)
	return b
}

// WithFilters appends multiple custom FilterFn at once.
func (b *ComplexFilterBuilder) WithFilters(fns ...FilterFn) *ComplexFilterBuilder {
	b.fns = append(b.fns, fns...)
	return b
}

// Expr adds a raw filter expression (no labels), e.g. "scale=1280:-2".
func (b *ComplexFilterBuilder) Expr(expr string) *ComplexFilterBuilder {
	return b.With(WithFilterExpr(expr))
}

// Chain adds a generic labeled chain: "[input]expr[output]".
func (b *ComplexFilterBuilder) Chain(input, expr, output string) *ComplexFilterBuilder {
	return b.With(WithFilterChain(input, expr, output))
}

// Build materializes the collected FilterFn into concrete FilterComplexParser units.
func (b *ComplexFilterBuilder) Build() []FilterComplexParser {
	fg := &FilterGraph{Options: make([]FilterComplexParser, 0)}
	for _, fn := range b.fns {
		fn(fg)
	}
	return fg.Options
}

// Apply applies the collected FilterFn directly to the provided command by materializing
// them into FilterComplexParser units and appending via the internal addFilters path.
func (b *ComplexFilterBuilder) Apply(cmd *Ffmpego) *Ffmpego {
	fg := &FilterGraph{Options: make([]FilterComplexParser, 0)}
	for _, fn := range b.fns {
		fn(fg)
	}
	cmd.addFilters(fg.Options...)
	return cmd
}
