package ffmpego

import (
	"fmt"
)

// Common output flag presets
var (
	// Video codecs
	VideoCodecH264 = WithVideoCodec("libx264")
	VideoCodecH265 = WithVideoCodec("libx265")
	VideoCodecVP9  = WithVideoCodec("libvpx-vp9")
	VideoCodecAV1  = WithVideoCodec("libaom-av1")

	// Audio codecs
	AudioCodecAAC  = WithAudioCodec("aac")
	AudioCodecMP3  = WithAudioCodec("libmp3lame")
	AudioCodecOpus = WithAudioCodec("libopus")

	// Presets
	PresetUltraFast = WithPreset("ultrafast")
	PresetFast      = WithPreset("fast")
	PresetMedium    = WithPreset("medium")
	PresetSlow      = WithPreset("slow")
	PresetVeryslow  = WithPreset("veryslow")

	// Quality levels
	CRFHighQuality   = WithCRF(18)
	CRFGoodQuality   = WithCRF(23)
	CRFMediumQuality = WithCRF(28)
	CRFLowQuality    = WithCRF(35)
)

type OutputDescriptor struct {
	Options []OutputFlagParser
}

type OutputFlagFn = func(*OutputDescriptor)

func (oo *OutputDescriptor) Add(option OutputFlagParser) {
	oo.Options = append(oo.Options, option)
}

func (oo *OutputDescriptor) Build() ([]string, error) {
	var args []string
	for _, flag := range oo.Options {
		if err := flag.Validate(); err != nil {
			return []string{}, err
		}

		args = append(args, flag.Parse()...)
	}

	return args, nil
}

func NewOutputDescriptor(opts ...OutputFlagFn) *OutputDescriptor {
	outputOptions := &OutputDescriptor{Options: make([]OutputFlagParser, 0)}
	for _, fn := range opts {
		fn(outputOptions)
	}

	return outputOptions
}

type VideoCodec string

func (vc VideoCodec) Parse() []string {
	return []string{"-c:v", string(vc)}
}

func (vc VideoCodec) Validate() error {
	return nil
}

type AudioCodec string

// Parse returns the audio codec flag arguments
func (f AudioCodec) Parse() []string {
	return []string{"-c:a", string(f)}
}

func (f AudioCodec) Validate() error {
	return nil
}

// WithVideoCodec creates a new video output codec flag
func WithVideoCodec(codec string) OutputFlagFn {
	return func(options *OutputDescriptor) {
		options.Add(VideoCodec(codec))
	}
}

// WithAudioCodec created a new audio output codec flag.
func WithAudioCodec(codec string) OutputFlagFn {
	return func(options *OutputDescriptor) {
		options.Add(AudioCodec(codec))
	}
}

// WithCRF creates a new CRF flag
func WithCRF(crf int) OutputFlagFn {
	return func(options *OutputDescriptor) {
		options.Add(CRFFlag(crf))
	}
}

// WithBitrate creates a new video bitrate flag
func WithBitrate(bitrate string) OutputFlagFn {
	return func(options *OutputDescriptor) {
		options.Add(BitrateFlag(bitrate))
	}
}

// WithPreset creates a new preset flag
func WithPreset(preset string) OutputFlagFn {
	return func(options *OutputDescriptor) {
		options.Add(PresetFlag(preset))
	}
}

// WithFormat creates a new format flag
func WithFormat(format string) OutputFlagFn {
	return func(options *OutputDescriptor) {
		options.Add(FormatFlag(format))
	}
}

// WithAudioBitrate creates a new audio bitrate flag
func WithAudioBitrate(bitrate string) OutputFlagFn {
	return func(options *OutputDescriptor) {
		options.Add(AudioBitrateFlag(bitrate))
	}
}

// WithSampleRate creates a new sample rate flag
func WithSampleRate(rate int) OutputFlagFn {
	return func(options *OutputDescriptor) {
		options.Add(SampleRateFlag(rate))
	}
}

// WithChannels creates a new channels flag
func WithChannels(channels int) OutputFlagFn {
	return func(options *OutputDescriptor) {
		options.Add(ChannelsFlag(channels))
	}
}

// WithMap creates a new map flag
func WithMap(stream string) OutputFlagFn {
	return func(options *OutputDescriptor) {
		options.Add(MapFlag(stream))
	}
}

// File represents an output file path
type File string

func (f File) Parse() []string {
	return []string{string(f)}
}

func (f File) Validate() error {
	if f == "" {
		return fmt.Errorf("file path cannot be empty")
	}
	return nil
}

// WithFile creates a new file path flag
func WithFile(filepath string) OutputFlagFn {
	return func(options *OutputDescriptor) {
		options.Add(File(filepath))
	}
}

// CRFFlag represents a constant rate factor option
type CRFFlag int

// Parse returns the CRF flag arguments
func (f CRFFlag) Parse() []string {
	return []string{"-crf", fmt.Sprintf("%d", int(f))}
}

// Validate validates the CRF flag
func (f CRFFlag) Validate() error {
	if f < 0 || f > 51 {
		return fmt.Errorf("CRF must be between 0 and 51, got %d", f)
	}
	return nil
}

// BitrateFlag represents a video bitrate option
type BitrateFlag string

// Parse returns the bitrate flag arguments
func (f BitrateFlag) Parse() []string {
	return []string{"-b:v", string(f)}
}

// Validate validates the bitrate flag
func (f BitrateFlag) Validate() error {
	if f == "" {
		return fmt.Errorf("bitrate cannot be empty")
	}
	return nil
}

// PresetFlag represents an encoding preset option
type PresetFlag string

// Parse returns the preset flag arguments
func (f PresetFlag) Parse() []string {
	return []string{"-preset", string(f)}
}

// Validate validates the preset flag
func (f PresetFlag) Validate() error {
	if f == "" {
		return fmt.Errorf("preset cannot be empty")
	}
	return nil
}

// FormatFlag represents an output format option
type FormatFlag string

// Parse returns the format flag arguments
func (f FormatFlag) Parse() []string {
	return []string{"-f", string(f)}
}

// Validate validates the format flag
func (f FormatFlag) Validate() error {
	if f == "" {
		return fmt.Errorf("format cannot be empty")
	}
	return nil
}

// AudioBitrateFlag represents an audio bitrate option
type AudioBitrateFlag string

// Parse returns the audio bitrate flag arguments
func (f AudioBitrateFlag) Parse() []string {
	return []string{"-b:a", string(f)}
}

// Validate validates the audio bitrate flag
func (f AudioBitrateFlag) Validate() error {
	if f == "" {
		return fmt.Errorf("audio bitrate cannot be empty")
	}
	return nil
}

// SampleRateFlag represents a sample rate option
type SampleRateFlag int

// Parse returns the sample rate flag arguments
func (f SampleRateFlag) Parse() []string {
	return []string{"-ar", fmt.Sprintf("%d", int(f))}
}

// Validate validates the sample rate flag
func (f SampleRateFlag) Validate() error {
	if f <= 0 {
		return fmt.Errorf("sample rate must be positive, got %d", f)
	}
	return nil
}

// ChannelsFlag represents an audio channels option
type ChannelsFlag int

// Parse returns the channels flag arguments
func (f ChannelsFlag) Parse() []string {
	return []string{"-ac", fmt.Sprintf("%d", int(f))}
}

// Validate validates the channels flag
func (f ChannelsFlag) Validate() error {
	if f <= 0 || f > 8 {
		return fmt.Errorf("channels must be between 1 and 8, got %d", f)
	}
	return nil
}

// MapFlag represents a stream mapping option
type MapFlag string

// Parse returns the map flag arguments
func (f MapFlag) Parse() []string {
	return []string{"-map", string(f)}
}

// Validate validates the map flag
func (f MapFlag) Validate() error {
	if f == "" {
		return fmt.Errorf("stream mapping cannot be empty")
	}
	return nil
}
