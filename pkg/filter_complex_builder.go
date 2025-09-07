package ffmpego

// FilterGraphBuilder provides a fluent API to compose a filter_complex graph
// using the same builder pattern as OutputBuilder. It collects FilterFn helpers
// and can either build concrete FilterComplexParser units or apply them to a command.
type FilterGraphBuilder struct {
	fns []FilterFn
}

// NewComplexFilterBuilder creates a new empty ComplexFilterBuilder.
func NewComplexFilterBuilder() *FilterGraphBuilder {
	return &FilterGraphBuilder{
		fns: make([]FilterFn, 0),
	}
}

// With appends any custom FilterFn.
func (b *FilterGraphBuilder) Add(fn FilterFn) *FilterGraphBuilder {
	b.fns = append(b.fns, fn)
	return b
}

// WithFilters appends multiple custom FilterFn at once.
func (b *FilterGraphBuilder) WithFilters(fns ...FilterFn) *FilterGraphBuilder {
	b.fns = append(b.fns, fns...)
	return b
}

// Expr adds a raw filter expression (no labels), e.g. "scale=1280:-2".
func (b *FilterGraphBuilder) Expr(expr string) *FilterGraphBuilder {
	return b.Add(WithFilterExpr(expr))
}

// Chain adds a generic labeled chain: "[input]expr[output]".
func (b *FilterGraphBuilder) Chain(input, expr, output string) *FilterGraphBuilder {
	return b.Add(WithFilterChain(input, expr, output))
}

// Build materializes the collected FilterFn into concrete FilterComplexParser units.
func (b *FilterGraphBuilder) Build() *FilterGraph {
	fg := &FilterGraph{Options: make([]FilterComplexParser, 0)}
	for _, fn := range b.fns {
		fn(fg)
	}

	return fg
}
