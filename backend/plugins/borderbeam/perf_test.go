// Author: HyperparameterSI
package borderbeam

import (
	"bytes"
	"github.com/HyperparameterSI/golet/backend/ui"
	"testing"
	"time"
)

// BenchmarkAnimationFrame measures complete animated frames, including texture
// updates and composition, at card and window sizes without lowering quality.
func BenchmarkAnimationFrame(b *testing.B) {
	for _, size := range []struct {
		name string
		w, h float32
	}{{"card", 480, 160}, {"window", 1200, 800}} {
		for _, mode := range []struct {
			name string
			p    Preset
		}{{"medium", Medium}, {"small", Small}, {"line", Line}, {"inner", PulseInner}, {"outside", PulseOutside}} {
			b.Run(size.name+"/"+mode.name, func(b *testing.B) {
				tt := ui.NewTester(func(c *ui.Context) { Box(c, Options{Preset: mode.p, Palette: Colorful}).Size(size.w, size.h) }, int(size.w), int(size.h))
				tt.Frame()
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					tt.Advance(time.Second / 30)
				}
			})
		}
	}
}

func TestOptimizedRenderMatchesReference(t *testing.T) {
	for _, size := range [][3]float64{{96, 64, 20}, {480, 160, 16}, {13.5, 9.5, 40}} {
		for mode := Medium; mode <= PulseOutside; mode++ {
			for _, dark := range []bool{false, true} {
				for _, seconds := range []float64{0, .4, 1.2, 2.2} {
					o := Options{Preset: mode, Palette: Colorful, Width: 3, Glow: 18}
					f := makeField(size[0], size[1], seconds, 2.3, mode, o, dark, nil)
					var actual, want [3]rasterLayer
					renderLayers(&actual, size[0], size[1], size[2], f, o, 1)
					referenceRenderLayers(&want, size[0], size[1], size[2], f, o, 1)
					for i := range actual {
						if !bytes.Equal(actual[i].pixels.Pix, want[i].pixels.Pix) {
							t.Fatalf("mode=%v dark=%v time=%v layer=%v pixels differ", mode, dark, seconds, i)
						}
						for j := range actual[i].scratch {
							a, z := actual[i].scratch[j], want[i].scratch[j]
							if (a == nil) != (z == nil) || a != nil && !bytes.Equal(a.Pix, z.Pix) {
								t.Fatalf("mode=%v layer=%v blur differs", mode, i)
							}
						}
					}
				}
			}
		}
	}
}

// BenchmarkPausedFrame simulates unrelated window repaints while the user has
// paused an effect (also the reduced-motion path).
func BenchmarkPausedFrame(b *testing.B) {
	tt := ui.NewTester(func(c *ui.Context) { Box(c, Options{Paused: true, Palette: Colorful}).Size(480, 160) }, 480, 160)
	tt.Frame()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tt.Frame()
	}
}

func TestRenderCacheInvalidation(t *testing.T) {
	var state clock
	colors := []ui.Color{ui.RGB(20, 80, 160)}
	o := Options{Palette: Colorful}
	render := func(w float32, elapsed time.Duration, dark bool) bool {
		return state.render(w, 40, 12, elapsed, time.Second, Medium, o, dark, colors, .7)
	}
	if !render(80, 0, false) || render(80, 0, false) {
		t.Fatal("unchanged frame was not reused")
	}
	colors[0] = ui.RGB(200, 40, 20)
	if !render(80, 0, false) {
		t.Fatal("mutated caller palette was not detected")
	}
	if !render(90, 0, false) || !render(90, time.Second/30, false) || !render(90, time.Second/30, true) {
		t.Fatal("size, animation or theme change was ignored")
	}
	for _, change := range []func(){
		func() { o.Brightness = 2 }, func() { o.Saturation = 2 }, func() { o.HueRange = 50 },
		func() { o.Width = 3 }, func() { o.Glow = 40 }, func() { o.StaticColors = true }, func() { o.Palette = Ocean },
	} {
		change()
		if !render(90, time.Second/30, true) || render(90, time.Second/30, true) {
			t.Fatal("option change did not invalidate exactly once")
		}
	}
}
