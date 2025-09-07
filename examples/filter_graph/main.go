package main

import (
	"fmt"
	"strings"

	ffmpego "m4urici0gm/ffmpego/pkg"
)

func main() {
	// Example 3: Using ComplexFilterBuilder to compose and apply a graph
	options := ffmpego.NewFfmpegOptions(
		ffmpego.WithInput("in.mp4"),
		ffmpego.WithOverwrite())

	filterGraph := ffmpego.NewComplexFilterBuilder().
		Add(ffmpego.WithCrop("0:v", "s1", 800, 600, 100, 50)).
		Add(ffmpego.WithRotate("s1", "s2", ffmpego.Rotate90)).
		Add(ffmpego.WithScale("s2", "s3", 1280, 720)).
		Build()

	output := ffmpego.NewOutputBuilder().
		WithFlag(ffmpego.VideoCodecH264).
		WithFlag(ffmpego.AudioCodecAAC).
		WithFlag(ffmpego.CRFGoodQuality).
		WithFlag(ffmpego.PresetFast).
		WithFlag(ffmpego.WithFile("output_builder_complex.mp4")).
		Build()

	cmd3 := ffmpego.New("").
		WithOptions(options).
		WithFilterGraph(filterGraph).
		Output(output)

	args3, err := cmd3.Build()
	if err != nil {
		fmt.Printf("Build error: %v\n", err)
		return
	}

	fmt.Printf("Command 3: ffmpeg %s\n", strings.Join(args3, " "))
}
