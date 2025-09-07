package main

import (
	"context"
	"log"

	ffmpego "m4urici0gm/ffmpego/pkg"
)

func main() {
	// Create multiple outputs using split

	options := ffmpego.NewFfmpegOptions(
		ffmpego.WithInput("in.mp4"),
		ffmpego.WithOverwrite())

	filterGraph := ffmpego.NewComplexFilterBuilder().
		Add(ffmpego.WithSplit("0:v", 3)).
		Add(ffmpego.WithScale("1", "480p", 854, 480)).
		Add(ffmpego.WithScale("2", "720p", 1280, 720)).
		Add(ffmpego.WithScale("3", "1080p", 1920, 1080)).
		Build()

	out1 := ffmpego.NewOutputBuilder().
		File("output_480p.mp4").
		WithFlag(ffmpego.VideoCodecH264).
		WithFlag(ffmpego.AudioCodecAAC).
		WithFlag(ffmpego.CRFGoodQuality).
		Build()

	out2 := ffmpego.NewOutputBuilder().
		File("output_720p.mp4").
		WithFlag(ffmpego.VideoCodecH264).
		WithFlag(ffmpego.AudioCodecAAC).
		WithFlag(ffmpego.CRFGoodQuality).
		Build()

	out3 := ffmpego.NewOutputBuilder().
		File("output_1080p.mp4").
		WithFlag(ffmpego.VideoCodecH264).
		WithFlag(ffmpego.AudioCodecAAC).
		WithFlag(ffmpego.CRFGoodQuality).
		Build()

	cmd := ffmpego.New("").
		WithOptions(options).
		WithFilterGraph(filterGraph).
		Output(out1).
		Output(out2).
		Output(out3)

	ctx := context.Background()
	executor := ffmpego.NewRunner(cmd)
	if err := executor.Run(ctx); err != nil {
		log.Fatalf("error when trying to run ffmpeg. %v", err)
	}
}