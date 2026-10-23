// Author: HyperparameterSI

// Package demotheme coordinates the appearance of native example windows.
// All access is on Golet's main UI thread.
package demotheme

import (
	golet "github.com/HyperparameterSI/golet/backend"
	"github.com/HyperparameterSI/golet/backend/plugins/glass"
	"github.com/HyperparameterSI/golet/backend/ui"
)

var liquid bool
var appearance int
var windows = map[*golet.Window]struct{}{}
var presets [2][2]*ui.Theme

// Selection returns the material choice and appearance (system, light or dark).
func Selection() (bool, int) { return liquid, appearance }

// Set updates open windows as well as the initial appearance of new windows.
func Set(useGlass bool, mode int) {
	liquid, appearance = useGlass, mode
	for w := range windows {
		w.Invalidate()
	}
}

// Register subscribes a demo window until it closes. Repeated calls are safe.
func Register(w *golet.Window) {
	if w == nil {
		return
	}
	if _, exists := windows[w]; exists {
		return
	}
	windows[w] = struct{}{}
	w.OnClosed(func() { delete(windows, w) })
}

// Apply installs the shared preset at the start of a native demo's view.
func Apply(c *ui.Context) {
	dark := c.Theme().Dark
	if appearance == 1 {
		dark = false
	}
	if appearance == 2 {
		dark = true
	}
	style, shade := 0, 0
	if liquid {
		style = 1
	}
	if dark {
		shade = 1
	}
	if presets[style][shade] == nil {
		if liquid {
			presets[style][shade] = glass.Theme(dark)
		} else {
			presets[style][shade] = ui.ShadcnTheme(dark)
		}
	}
	c.SetTheme(presets[style][shade])
}
