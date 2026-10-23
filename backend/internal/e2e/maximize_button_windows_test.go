// Author: HyperparameterSI

//go:build windows

package e2e

import (
	"testing"

	"github.com/HyperparameterSI/golet/backend"
	win "github.com/HyperparameterSI/golet/backend/internal/windows"
	"github.com/HyperparameterSI/golet/backend/ui"
)

// The maximize button a view draws, in a window without the system's
// buttons, answers Windows' hit test as the system's does, which brings
// the snap layouts, and still clicks as the view's own button.
func TestViewMaximizeButton(t *testing.T) {
	clicks := 0
	w := newWindow(t, golet.WindowOptions{X: 80, Y: 80, Width: 400, Height: 240,
		TitleBarStyle: golet.TitleBarHidden, HideWindowButtons: true,
		Content: ui.View(func(c *ui.Context) {
			ui.Row(c).Fill().Children(func() {
				ui.Spacer(c)
				b := ui.ButtonBase(c.Key("max")).Size(40, 30).Label("Maximize").MaximizeButton()
				if b.Clicked() {
					clicks++
				}
			})
		})})
	at := func(press bool) (hit uintptr) {
		golet.RunOnMain(func() { hit = win.TestMaximizeButton(w.NativeHandle(), press) })
		return hit
	}
	const htMaxButton = 9
	eventually(t, "a window over the view's maximize button", func() bool { return at(false) == htMaxButton })
	if hit := at(true); hit != htMaxButton {
		t.Fatalf("hit test = %d, want HTMAXBUTTON", hit)
	}
	eventually(t, "the press to click the view's button", func() bool {
		var n int
		golet.RunOnMain(func() { n = clicks })
		return n == 1
	})
}
