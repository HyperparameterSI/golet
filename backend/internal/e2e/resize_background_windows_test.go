// Author: HyperparameterSI

//go:build windows

package e2e

import (
	"sync/atomic"
	"testing"

	"github.com/HyperparameterSI/golet/backend"
	win "github.com/HyperparameterSI/golet/backend/internal/windows"
	"github.com/HyperparameterSI/golet/backend/ui"
)

// TestNativeUIWindowPixelsUnderSurface checks that the own pixels of a
// window showing native UI on the GPU hold the content's background under
// its surface, from its first frame and after each change of size, by its
// bottom-right corner and by its top-left one. As the window shrinks, DWM
// shows them for a refresh around the surface already fitted to the new
// size, before the window takes it: nothing painted them, and a thin black
// L showed there.
func TestNativeUIWindowPixelsUnderSurface(t *testing.T) {
	var frames atomic.Int32
	w := newWindow(t, golet.WindowOptions{X: 80, Y: 80, Width: 480, Height: 340, TitleBarStyle: golet.TitleBarHidden, HideWindowButtons: true,
		Content: ui.View(func(c *ui.Context) {
			frames.Add(1)
			c.Root().Background(ui.RGB(250, 250, 240))
		})})
	eventually(t, "a frame", func() bool { return frames.Load() > 0 })
	var last string
	under := func() (ok bool) {
		golet.RunOnMain(func() { last, ok = win.TestWindowPixelsUnderSurface(w.NativeHandle()) })
		return ok
	}
	t.Cleanup(func() {
		if t.Failed() && last != "" {
			t.Logf("under the surface, %s", last)
		}
	})
	if !under() {
		t.Skip("the window has no pixels of its own under a swap chain: its frames are drawn in memory or composed")
	}
	held := func(when string) {
		t.Helper()
		eventually(t, "the background under the surface "+when, func() bool { return under() && last == "" })
	}
	held("at the first frame")
	for _, c := range []struct {
		when string
		b    golet.Rectangle
	}{
		{"once it shrank", golet.Rectangle{X: 80, Y: 80, Width: 430, Height: 300}},
		{"once it grew", golet.Rectangle{X: 80, Y: 80, Width: 500, Height: 360}},
		{"once it shrank by its top-left corner", golet.Rectangle{X: 130, Y: 120, Width: 450, Height: 320}},
		{"once it grew by its top-left corner", golet.Rectangle{X: 60, Y: 60, Width: 520, Height: 380}},
	} {
		w.SetBounds(c.b)
		held(c.when)
	}
}
