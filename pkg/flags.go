package ffmpego

type FfmpegOptions struct {
	flags []FfmpegFlagParser
}

func NewFfmpegOptions(flags ...FfmpegFlagFn) *FfmpegOptions {
    options := &FfmpegOptions{flags:make([]FfmpegFlagParser, 0)}
    for _, fn := range flags {
        fn(options)
    }

    return options
}

func (ff *FfmpegOptions) Add(option FfmpegFlagParser) {
	ff.flags = append(ff.flags, option)
}

func (ff *FfmpegOptions) BuildAndValidate() ([]string, error) {
	var args []string
	for _, flag := range ff.flags {
		if err := flag.Validate(); err != nil {
			return []string{}, err
		}

		args = append(args, flag.Parse()...)
	}

	return args, nil
}

type FfmpegFlagFn = func(*FfmpegOptions)

type FfmpegFlagParser interface {
	Validate() error
	Parse() []string
}

// GlobalOptions represents global FFmpeg options
type GlobalOptions struct {
	Overwrite bool
	LogLevel  string
	JSON      bool
}

type Input []string

func (in Input) Validate() error {
	return nil
}

func (input Input) Parse() []string {
	var inputs []string
	for _, i := range input {
		inputs = append(inputs, "-i", i)
	}

	return inputs
}

type Overwrite struct{}

func (ow Overwrite) Validate() error {
	return nil
}

func (ow Overwrite) Parse() []string {
	return []string{"-y"}
}

type LogLevel string

func (ll LogLevel) Validate() error {
	return nil
}

func (ll LogLevel) Parse() []string {
	return []string{"-loglevel", string(ll)}
}

type Output string

func (ll Output) Validate() error {
	return nil
}

func (ll Output) Parse() []string {
	return []string{"-progress", string(ll)}
}
