package ffmpego

type FfmpegOptions struct {
	flags []FfmpegFlagParser
}

func (ff *FfmpegOptions) Add(option FfmpegFlagParser) {
    ff.flags = append(ff.flags, option)
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
