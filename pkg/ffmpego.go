package ffmpego

import (
	"strings"
	"time"
)

// OutputFlagParser interface for all output flag operations
type OutputFlagParser interface {
	Validate() error
	Parse() []string
}

// FilterComplexParser represents a unit of a filter_complex graph.
// Parse should return a single chain or subgraph (e.g. "[0:v]scale=1280:-2[outv]")
// without the "-filter_complex" flag. Multiple units are joined with ';'.
type FilterComplexParser interface {
	Validate() error
	Parse() string
}

// Progress represents FFmpeg progress information
type Progress struct {
	Frame     int     `json:"frame,omitempty"`
	FPS       float64 `json:"fps,omitempty"`
	Bitrate   string  `json:"bitrate,omitempty"`
	TotalSize int64   `json:"total_size,omitempty"`
	OutTime   string  `json:"out_time,omitempty"`
	OutTimeMS int64   `json:"out_time_ms,omitempty"`
	Speed     string  `json:"speed,omitempty"`
	Progress  string  `json:"progress,omitempty"`
}

// ProgressCallback is a function type for handling progress updates
type ProgressCallback func(Progress)

// Ffmpego represents an FFmpeg command builder
type Ffmpego struct {
	inputs           []string
	outputs          []OutputDescriptor
	graph            *FilterGraph
	flags            []FfmpegFlagParser
	timeout          time.Duration
	progressCallback ProgressCallback
}

// New creates a new FFmpeg command
func New() *Ffmpego {
	return &Ffmpego{
		inputs:  make([]string, 0),
		outputs: make([]OutputDescriptor, 0),
		graph:   &FilterGraph{Options: make([]FilterComplexParser, 0)},
		timeout: 30 * time.Minute,
	}
}

// Input adds an input file
func (c *Ffmpego) Input(file string) *Ffmpego {
	c.inputs = append(c.inputs, file)
	return c
}

// Output adds an output configuration
func (c *Ffmpego) Output(output OutputDescriptor) *Ffmpego {
	c.outputs = append(c.outputs, output)
	return c
}

// WithFilter adds a single filter chain or subgraph to the -filter_complex graph.
// The provided filter will be concatenated with other filters using ';'.
func (c *Ffmpego) WithFilter(filter FilterComplexParser) *Ffmpego {
	c.addFilters(filter)
	return c
}

// WithFilters adds multiple filters at once to the -filter_complex graph.
func (c *Ffmpego) WithFilters(filters ...FilterComplexParser) *Ffmpego {
	c.addFilters(filters...)
	return c
}

// addFilters initializes the filter graph if needed and appends the provided filters.
func (c *Ffmpego) addFilters(filters ...FilterComplexParser) {
	if c.graph == nil {
		c.graph = &FilterGraph{Options: make([]FilterComplexParser, 0)}
	}
	for _, f := range filters {
		c.graph.Add(f)
	}
}

// Timeout sets the command timeout
func (c *Ffmpego) Timeout(duration time.Duration) *Ffmpego {
	c.timeout = duration
	return c
}

// WithProgressCallback sets a callback function for progress updates
func (c *Ffmpego) WithProgressCallback(callback ProgressCallback) *Ffmpego {
	c.progressCallback = callback
	return c
}

// Build constructs the FFmpeg command arguments
func (c *Ffmpego) Build() ([]string, error) {
	args := make([]string, 0)

	// Input files
	for _, input := range c.inputs {
		args = append(args, "-i", input)
	}

	for _, flag := range c.flags {
		if err := flag.Validate(); err != nil {
			return []string{}, err
		}

		args = append(args, flag.Parse()...)
	}

	// Filters - combine all filters into a single filter_complex
	if c.graph != nil && len(c.graph.Options) > 0 {
		var filterParts []string
		for _, filter := range c.graph.Options {
			if err := filter.Validate(); err != nil {
				return []string{}, err
			}
			filterParts = append(filterParts, filter.Parse())
		}
		if len(filterParts) > 0 {
			args = append(args, "-filter_complex", strings.Join(filterParts, ";"))
		}
	}

	// Output configurations
	for _, output := range c.outputs {
		ooArgs, err := output.Build()
		if err != nil {
			return []string{}, err
		}

		// Output file
		args = append(args, ooArgs...)
	}

	return args, nil
}
