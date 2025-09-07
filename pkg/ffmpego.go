package ffmpego

import "time"

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
	binary           string
	flags            *FfmpegOptions
	graph            *FilterGraph
	outputs          []*OutputDescriptor
	timeout          time.Duration
	progressCallback ProgressCallback
}

// New creates a new FFmpeg command
func New(binary string) *Ffmpego {
	if binary == "" {
		binary = "/bin/ffmpeg"
	}

	return &Ffmpego{
		binary:           binary,
		graph:            &FilterGraph{Options: make([]FilterComplexParser, 0)},
		flags:            &FfmpegOptions{},
		outputs:          []*OutputDescriptor{},
		timeout:          0,
		progressCallback: nil,
	}
}

func (c *Ffmpego) WithOptions(flags *FfmpegOptions) *Ffmpego {
    c.flags = flags
    return c
}

// Output adds an output configuration
func (c *Ffmpego) Output(output *OutputDescriptor) *Ffmpego {
	c.outputs = append(c.outputs, output)
	return c
}

// WithFilter adds a single filter chain or subgraph to the -filter_complex graph.
// The provided filter will be concatenated with other filters using ';'.
func (c *Ffmpego) WithFilterGraph(graph *FilterGraph) *Ffmpego {
	c.graph = graph
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

	options, err := c.flags.BuildAndValidate()
	if err != nil {
		return []string{}, err
	}

	args = append(args, options...)

	// Filters - combine all filters into a single filter_complex
	if c.graph != nil && len(c.graph.Options) > 0 {
		complexFilter, err := c.graph.BuildAndValidate()
		if err != nil {
			return []string{}, err
		}

		args = append(args, "-filter_complex", complexFilter)
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
