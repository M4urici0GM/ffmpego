package main

import (
	"fmt"
	"strings"

	ffmpego "m4urici0gm/ffmpego/pkg"
)

func main() {
	// Example 3: Using ComplexFilterBuilder to compose and apply a graph
	cmd3 := ffmpego.New().
		Input("input.mp4").
		WithFilters(
			ffmpego.NewComplexFilterBuilder().
				With(ffmpego.WithCrop("0:v", "s1", 800, 600, 100, 50)).
				With(ffmpego.WithRotate("s1", "s2", ffmpego.Rotate90)).
				With(ffmpego.WithScale("s2", "s3", 1280, 720)).
				Build()...,
		).
		Output(
			ffmpego.NewOutputBuilder().
				WithFlag(ffmpego.VideoCodecH264).
				WithFlag(ffmpego.AudioCodecAAC).
				WithFlag(ffmpego.CRFGoodQuality).
				WithFlag(ffmpego.PresetFast).
				WithFlag(ffmpego.WithFile("output_builder_complex.mp4")).
				Build(),
		)

	args3, err := cmd3.Build()
	if err != nil {
		fmt.Printf("Build error: %v\n", err)
		return
	}

	fmt.Printf("Command 3: ffmpeg %s\n", strings.Join(args3, " "))
}
