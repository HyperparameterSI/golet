// Author: HyperparameterSI
package borderbeam

import (
	"bytes"
	"testing"
	"time"

	"github.com/HyperparameterSI/golet/backend/ui"
)

func TestAnimationLifecycle(t *testing.T) {
	for _, preset := range []Preset{Medium, Small, Line, PulseInner, PulseOutside} {
		o := Options{Preset: preset, Palette: Colorful, Duration: time.Second}
		tt := ui.NewTester(func(c *ui.Context) { Box(c, o).Size(180, 90) }, 200, 110)
		frame := func() []byte { tt.Frame(); return append([]byte(nil), tt.Image().Pix...) }
		initial := frame()
		tt.Advance(250 * time.Millisecond)
		if bytes.Equal(initial, frame()) {
			t.Fatal("animation did not move")
		}
		o.Paused = true
		held := frame()
		if tt.FramePending() {
			t.Fatal("paused effect requests frames")
		}
		tt.Advance(time.Second)
		if !bytes.Equal(held, frame()) {
			t.Fatal("paused phase changed")
		}
		o.Paused = false
		if !bytes.Equal(held, frame()) {
			t.Fatal("resume jumps")
		}
		tt.Advance(250 * time.Millisecond)
		if bytes.Equal(held, frame()) {
			t.Fatal("resume did not move")
		}
		tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
		held = frame()
		tt.Advance(time.Second)
		if !bytes.Equal(held, frame()) || tt.FramePending() {
			t.Fatal("reduced motion must remain still")
		}
		o.Disabled = true
		if bytes.Equal(held, frame()) || tt.FramePending() {
			t.Fatal("disabled border must disappear and stop")
		}
	}
}

// Straight-edge length must not change the analytic corner geometry.
func TestCornerIndependentOfWidth(t *testing.T) {
	for _, radius := range []float64{8, 16, 32} {
		for y := 0.; y < 40; y++ {
			for x := 0.; x < 40; x++ {
				if a, b := roundedDistance(x, y, 200, 100, radius), roundedDistance(x, y, 3000, 100, radius); a != b {
					t.Fatalf("corner mismatch: %v / %v", a, b)
				}
			}
		}
	}
}

func TestContentRemainsInteractive(t *testing.T) {
	clicks := 0
	tt := ui.NewTester(func(c *ui.Context) {
		Box(c, Options{Paused: true}).Padding(16).Children(func() { ui.Button(c, "Action").OnClick(func() { clicks++ }) })
	}, 200, 100)
	tt.Click("Action")
	if clicks != 1 {
		t.Fatal("border intercepted button click")
	}
}

func TestOffscreenDoesNotAnimate(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Scroll(c).Fill().Children(func() {
			ui.Box(c).Height(300).Shrink(0)
			Box(c, Options{}).Size(180, 90).Shrink(0)
		})
	}, 200, 110)
	tt.Frame()
	if tt.FramePending() {
		t.Fatal("offscreen decoration requests animation frames")
	}
}

func TestSmallSizes(t *testing.T) {
	for _, size := range []float32{1, 8, 16, 32} {
		tt := ui.NewTester(func(c *ui.Context) { Box(c, Options{Paused: true, Glow: 48, Radius: 40, Width: 8}).Size(size, size) }, 64, 64)
		tt.Frame()
	}
}

func BenchmarkCardFrame(b *testing.B) {
	tt := ui.NewTester(func(c *ui.Context) { Box(c, Options{Palette: Colorful}).Size(480, 160) }, 480, 160)
	tt.Frame()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tt.Advance(time.Second / 30)
	}
}

// Every published palette must have complete reference data for every preset.
func TestReferencePalettes(t *testing.T) {
	for _, dark := range []bool{false, true} {
		for preset := Medium; preset <= PulseOutside; preset++ {
			for palette := Colorful; palette <= Mono; palette++ {
				f := makeField(350, 140, 1.1, 2.3, preset, Options{Palette: palette}, dark, nil)
				for layer := 0; layer < 3; layer++ {
					if preset < Line && layer == 2 {
						continue
					} // Medium/Small bloom is a conic highlight.
					if len(f.layers[layer]) == 0 {
						t.Fatalf("missing preset %d palette %d layer %d", preset, palette, layer)
					}
					for _, spot := range f.layers[layer] {
						if spot.shape[2] <= 0 || spot.shape[3] <= 0 || spot.color[3] <= 0 {
							t.Fatalf("invalid reference spot: %+v", spot)
						}
					}
				}
			}
		}
	}
}

func TestOutsideGlowExtendsBeyondCard(t *testing.T) {
	var layers [3]rasterLayer
	o := Options{Preset: PulseOutside, Palette: Colorful}
	f := makeField(180, 90, 1.1, 2.3, PulseOutside, o, true, nil)
	renderLayers(&layers, 180, 90, 16, f, o, 1)
	// The textures include 60 DIPs of outside space, sampled at 2x.
	img := layers[1].pixels
	found := false
	for y := 0; y < img.Rect.Dy(); y++ {
		for x := 0; x < img.Rect.Dx(); x++ {
			if x >= 120 && x < 480 && y >= 120 && y < 300 {
				continue
			}
			if img.Pix[img.PixOffset(x, y)+3] != 0 {
				found = true
			}
		}
	}
	if !found {
		t.Fatal("outside preset has no exterior light")
	}
}
