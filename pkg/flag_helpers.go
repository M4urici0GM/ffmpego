package ffmpego

var (
	PipeProgress = WithProgress("pipe:1")
)

func WithInput(input ...string) FfmpegFlagFn {
	return func(options *FfmpegOptions) {
		options.Add(Input(input))
	}
}

// Adds new '-y' flag to ffmpeg command.
func WithOverwrite() FfmpegFlagFn {
	return func(options *FfmpegOptions) {
		options.Add(Overwrite{})
	}
}

func WithLogLevel(level string) FfmpegFlagFn {
	return func(options *FfmpegOptions) {
		options.Add(LogLevel(level))
	}
}

func WithProgress(progress string) FfmpegFlagFn {
	return func(options *FfmpegOptions) {
		options.Add(Output(progress))
	}
}
