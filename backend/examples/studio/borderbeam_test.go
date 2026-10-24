// Author: HyperparameterSI
package main

import (
	"bytes"
	"image"
	"testing"
	"time"

	"github.com/HyperparameterSI/golet/backend/ui"
)

// Check the real Studio layout, including every card in its grid, using only
// requested repaint frames. Do not force a rebuild to make animation work.
func TestBorderBeamStudioAnimates(t *testing.T) {
	tt := ui.NewTester(newBorderBeamView(), 960, 1000)
	tt.Frame()
	names := []string{"Prism · 多彩棱镜", "Ocean · 海洋蓝", "Sunset · 落日余晖", "Aurora · 极光"}
	frame := func(name string) []byte {
		img := tt.Image()
		rect, ok := tt.Find(name + " card")
		if !ok {
			t.Fatalf("missing card %s", name)
		}
		area := image.Rect(int(rect.X), int(rect.Y), int(rect.X+rect.W), int(rect.Y+rect.H)).Intersect(img.Bounds())
		if area.Empty() {
			t.Fatalf("card %s not visible", name)
		}
		var pixels []byte
		for y := area.Min.Y; y < area.Max.Y; y++ {
			i := img.PixOffset(area.Min.X, y)
			pixels = append(pixels, img.Pix[i:i+area.Dx()*4]...)
		}
		return pixels
	}
	snapshots := func() map[string][]byte {
		m := map[string][]byte{}
		for _, name := range names {
			m[name] = frame(name)
		}
		return m
	}
	before := snapshots()
	// Timed paint callbacks need not request the very next display frame.
	// Advancing without Frame verifies the actual scheduled repaint.
	tt.Advance(575 * time.Millisecond)
	for _, name := range names {
		if bytes.Equal(before[name], frame(name)) {
			t.Errorf("%s did not animate", name)
		}
	}
	tt.Click("暂停动画")
	held := snapshots()
	tt.Advance(575 * time.Millisecond)
	for _, name := range names {
		if !bytes.Equal(held[name], frame(name)) {
			t.Errorf("%s moved while paused", name)
		}
	}
	tt.Click("继续播放")
	tt.Advance(575 * time.Millisecond)
	for _, name := range names {
		if bytes.Equal(held[name], frame(name)) {
			t.Errorf("%s did not resume", name)
		}
	}
}
