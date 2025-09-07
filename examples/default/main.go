package main

import (
	"context"
	"log"
	ffmpego "m4urici0gm/ffmpego/pkg"
)

func main() {
	ctx := context.Background()

	// Create a simple crop and scale operation
	cmd := ffmpego.New("").
		WithOptions(
			ffmpego.NewFfmpegOptions(
				ffmpego.WithInput("in.mp4"),
				ffmpego.WithOverwrite(),
			),
		).
		WithFilterGraph(
			ffmpego.NewComplexFilterBuilder().
				Add(ffmpego.WithCrop("0:v", "cropped", 800, 600, 100, 50)).
				Add(ffmpego.WithScale("cropped", "scaled", 1280, 720)).
				Build(),
		).
		Output(
			ffmpego.NewOutputBuilder().
				File("output_basic.mp4").
				WithFlag(ffmpego.VideoCodecH264).
				WithFlag(ffmpego.AudioCodecAAC).
				WithFlag(ffmpego.CRFGoodQuality).
				WithFlag(ffmpego.PresetMedium).
				Build(),
		)

	executor := ffmpego.NewRunner(cmd)
	if err := executor.Run(ctx); err != nil {
		log.Fatalf("error when trying to execute ffmpeg. %v", err)
	}
}
