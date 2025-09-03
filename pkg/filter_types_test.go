package ffmpego

import "testing"

func TestScaleFilter_ValidateParse(t *testing.T) {
	f := ScaleFilter{Input: "in", Output: "out", Width: -2, Height: 720}
	if err := f.Validate(); err != nil {
		t.Fatalf("validate failed: %v", err)
	}
	got := f.Parse()
	want := "[in]scale=-2:720[out]"
	if got != want {
		t.Fatalf("parse mismatch: got %q want %q", got, want)
	}
}

func TestScaleFilter_InvalidDims(t *testing.T) {
	cases := []ScaleFilter{
		{Input: "in", Output: "out", Width: 0, Height: 720},
		{Input: "in", Output: "out", Width: 1280, Height: 0},
	}
	for i, f := range cases {
		if err := f.Validate(); err == nil {
			t.Fatalf("case %d: expected error for invalid dims, got nil", i)
		}
	}
}

func TestCropFilter_ValidateParse(t *testing.T) {
	f := CropFilter{Input: "in", Output: "out", W: 800, H: 600, X: 100, Y: 50}
	if err := f.Validate(); err != nil {
		t.Fatalf("validate failed: %v", err)
	}
	got := f.Parse()
	want := "[in]crop=800:600:100:50[out]"
	if got != want {
		t.Fatalf("parse mismatch: got %q want %q", got, want)
	}
}

func TestCropFilter_InvalidXY(t *testing.T) {
	f := CropFilter{Input: "in", Output: "out", W: 800, H: 600, X: -1, Y: 10}
	if err := f.Validate(); err == nil {
		t.Fatalf("expected error for negative X, got nil")
	}
}

func TestRotateFilter_ValidateParse(t *testing.T) {
	f := RotateFilter{Input: "in", Output: "out", Mode: Rotate90}
	if err := f.Validate(); err != nil {
		t.Fatalf("validate failed: %v", err)
	}
	got := f.Parse()
	want := "[in]transpose=1[out]"
	if got != want {
		t.Fatalf("parse mismatch: got %q want %q", got, want)
	}
}

func TestRotateFilter_InvalidMode(t *testing.T) {
	f := RotateFilter{Input: "in", Output: "out", Mode: 4}
	if err := f.Validate(); err == nil {
		t.Fatalf("expected error for invalid mode, got nil")
	}
}

func TestSplitFilter_ValidateParse(t *testing.T) {
	f := SplitFilter{Input: "in", N: 3, Outputs: []string{"o1", "o2", "o3"}}
	if err := f.Validate(); err != nil {
		t.Fatalf("validate failed: %v", err)
	}
	got := f.Parse()
	want := "[in]split=3[o1][o2][o3]"
	if got != want {
		t.Fatalf("parse mismatch: got %q want %q", got, want)
	}
}

func TestSplitFilter_InvalidOutputsCount(t *testing.T) {
	f := SplitFilter{Input: "in", N: 2, Outputs: []string{"o1"}}
	if err := f.Validate(); err == nil {
		t.Fatalf("expected error for outputs count != n, got nil")
	}
}
