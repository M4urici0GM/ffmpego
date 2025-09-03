package ffmpego

import (
	"strings"
	"testing"
)

func TestBuild_SimpleOutput_NoFilters(t *testing.T) {
	cmd := New().
		Input("in.mp4").
		Output(*NewOutputDescriptor(
			WithVideoCodec("libx264"),
			WithFile("out.mp4"),
		))

	args, err := cmd.Build()
	if err != nil {
		t.Fatalf("Build() error: %v", err)
	}

	want := []string{"-i", "in.mp4", "-c:v", "libx264", "out.mp4"}
	if len(args) != len(want) {
		t.Fatalf("args len mismatch: got %d want %d: %v", len(args), len(want), args)
	}
	for i := range want {
		if args[i] != want[i] {
			t.Fatalf("arg[%d] mismatch: got %q want %q (args=%v)", i, args[i], want[i], args)
		}
	}
}

func TestBuild_WithFilters_FilterComplexJoin(t *testing.T) {
	cmd := New().
		Input("in.mp4").
		WithFilter(NewCrop("0:v", "a", 800, 600, 100, 50)).
		WithFilter(NewScale("a", "b", 1280, 720)).
		Output(*NewOutputDescriptor(
			WithFile("out.mp4"),
		))

	args, err := cmd.Build()
	if err != nil {
		t.Fatalf("Build() error: %v", err)
	}

	// Expect: -filter_complex "[0:v]crop=800:600:100:50[a];[a]scale=1280:720[b]"
	idx := indexOf(args, "-filter_complex")
	if idx < 0 || idx+1 >= len(args) {
		t.Fatalf("missing -filter_complex in args: %v", args)
	}
	gotGraph := args[idx+1]
	wantGraph := "[0:v]crop=800:600:100:50[a];[a]scale=1280:720[b]"
	if gotGraph != wantGraph {
		t.Fatalf("filter graph mismatch:\n got: %s\nwant: %s", gotGraph, wantGraph)
	}

	if args[len(args)-1] != "out.mp4" {
		t.Fatalf("last arg should be the output file, got %q (args=%v)", args[len(args)-1], args)
	}
}

func TestOutputBuilder_Integration(t *testing.T) {
	ob := NewOutputBuilder().
		WithFlag(VideoCodecH264). // same as WithVideoCodec("libx264")
		File("built.mp4")

	cmd := New().
		Input("in.mp4").
		Output(ob.Build())

	args, err := cmd.Build()
	if err != nil {
		t.Fatalf("Build() error: %v", err)
	}

	if indexOf(args, "-c:v") == -1 || indexOf(args, "libx264") == -1 {
		t.Fatalf("expected video codec flags in args, got %v", args)
	}
	if args[len(args)-1] != "built.mp4" {
		t.Fatalf("expected output file 'built.mp4', got %q", args[len(args)-1])
	}
}

func TestComplexFilterBuilder_Apply(t *testing.T) {
	builder := NewComplexFilterBuilder().
		With(WithCrop("0:v", "s1", 800, 600, 100, 50)).
		With(WithRotate("s1", "s2", Rotate90)).
		With(WithScale("s2", "s3", 1280, 720))

	cmd := New().Input("in.mp4")
	builder.Apply(cmd)
	cmd = cmd.Output(*NewOutputDescriptor(WithFile("out.mp4")))

	args, err := cmd.Build()
	if err != nil {
		t.Fatalf("Build() error: %v", err)
	}

	graphIdx := indexOf(args, "-filter_complex")
	if graphIdx < 0 || graphIdx+1 >= len(args) {
		t.Fatalf("missing -filter_complex in args: %v", args)
	}
	graph := args[graphIdx+1]
	expectParts := []string{
		"[0:v]crop=800:600:100:50[s1]",
		"[s1]transpose=1[s2]",
		"[s2]scale=1280:720[s3]",
	}
	for _, part := range expectParts {
		if !strings.Contains(graph, part) {
			t.Fatalf("graph missing expected part %q: %s", part, graph)
		}
	}
	if !strings.Contains(graph, ";") {
		t.Fatalf("graph should join chains with ';': %s", graph)
	}
}

func TestBuild_InvalidFilter_ReturnsError(t *testing.T) {
	// Invalid: missing input label
	bad := ScaleFilter{Input: "", Output: "x", Width: 1280, Height: 720}

	cmd := New().
		Input("in.mp4").
		WithFilter(bad).
		Output(*NewOutputDescriptor(WithFile("out.mp4")))

	_, err := cmd.Build()
	if err == nil {
		t.Fatalf("expected error for invalid filter, got nil")
	}
}

// helper
func indexOf(slice []string, v string) int {
	for i, s := range slice {
		if s == v {
			return i
		}
	}
	return -1
}
