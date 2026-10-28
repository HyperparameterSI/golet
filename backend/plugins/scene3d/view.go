// Author: HyperparameterSI

package scene3d

import (
	"fmt"
	"math"
	"time"

	"github.com/HyperparameterSI/golet/backend/ui"
)

// ViewOptions control rendering and integration with Forma or other native UI.
// Revision must change when Scene or mesh instance/material/light data changes.
// Camera and resolution changes are detected automatically. No timer is started.
// Use a stable UI Key when showing/reordering several viewports.
type ViewOptions struct {
	Scene    *Scene
	Camera   Camera
	Revision uint64
	// Scale is render pixels per DIP (default 1, maximum 2). MaxDimension
	// caps either side (default 1024, maximum 2048), preserving aspect ratio.
	Scale        float32
	MaxDimension int
	// NewRenderer is called once per mounted view. nil uses Software.
	NewRenderer func() Renderer
	OnInput     func(ui.InputEvent) bool
	// Animate updates scene transforms before each visible frame. nil leaves
	// the viewport event-driven. Reduced-motion preferences suppress it.
	// It runs on the UI thread; use the supplied clock, not wall-clock timers.
	Animate func(time.Time)
	// OnFrame runs after a render attempt on the UI thread. It must not
	// reenter rendering. Errors are explicit; failed frames are not displayed.
	// After an error, change Revision or another frame input to retry.
	OnFrame func(Stats, error)
}
type viewKey struct{}
type frameKey struct {
	scene    *Scene
	camera   Camera
	revision uint64
	w, h     int
}
type viewState struct {
	target    Target
	renderer  Renderer
	bitmap    *ui.Bitmap
	key       frameKey
	valid     bool
	failed    bool
	wantFocus bool
}

// View creates a normal UI element containing a 3D image. It participates in
// layout, clipping, overlays and input with Forma in the same native window.
// Call only while building a ui.View, and mutate scenes on that UI thread.
// The renderer is CPU-based; the resulting image uses Golet's usual compositor.
func View(c *ui.Context, o ViewOptions) ui.Element {
	animate := o.Animate != nil && !c.Preferences().ReduceMotion
	e := ui.Box(c).Focusable().Label("3D viewport")
	state := ui.Local(e, viewKey{}, func() viewState {
		r := Renderer(Software{})
		if o.NewRenderer != nil {
			r = o.NewRenderer()
		}
		if r == nil {
			r = Software{}
		}
		return viewState{renderer: r}
	})
	if state.wantFocus {
		state.wantFocus = false
		e.Focus()
	}
	e.HandleInput(func(ev ui.InputEvent) bool {
		taken := false
		if ev.Kind == ui.InputPointerDown && ev.Button == 0 {
			state.wantFocus = true
			taken = true
		}
		if o.OnInput != nil {
			return o.OnInput(ev) || taken
		}
		return taken
	})
	return e.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 || !finite(r.W) || !finite(r.H) {
			return
		}
		scale := o.Scale
		if !finite(scale) || scale <= 0 {
			scale = 1
		}
		scale = min(scale, 2)
		limit := o.MaxDimension
		if limit <= 0 {
			limit = 1024
		}
		limit = min(limit, 2048)
		scale = min(scale, float32(limit)/max(r.W, r.H))
		w, h := max(1, int(math.Round(float64(r.W*scale)))), max(1, int(math.Round(float64(r.H*scale))))
		key := frameKey{o.Scene, o.Camera, o.Revision, w, h}
		if animate && (!state.failed || state.key != key) {
			o.Animate(p.Now())
			state.valid = false
		}
		if !state.valid || state.key != key {
			stats, err := state.renderer.Render(&state.target, Request{o.Scene, o.Camera, w, h})
			if err == nil && (state.target.Image == nil || state.target.Image.Bounds().Dx() != w || state.target.Image.Bounds().Dy() != h) {
				err = fmt.Errorf("scene3d: renderer returned an invalid frame size")
			}
			state.key = key
			state.valid = true
			state.failed = err != nil
			if err != nil {
				state.bitmap = nil
			} else if state.bitmap == nil {
				state.bitmap = ui.NewBitmap(state.target.Image)
			} else {
				state.bitmap.Update(state.target.Image)
			}
			if o.OnFrame != nil {
				o.OnFrame(stats, err)
			}
		}
		if state.bitmap != nil {
			p.Image(state.bitmap, r, ui.FillBox)
			if animate {
				p.AnimationFrame()
			}
		}
	})
}
