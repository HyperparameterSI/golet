// Author: HyperparameterSI

// Package borderbeam ports the five Border Beam presets to Golet native UI.
// Color fields, masks and Gaussian bloom run in Go without a browser runtime.
// Reference source and MIT attribution are in reference/.
package borderbeam

import (
	"image"
	"math"
	"time"

	"github.com/HyperparameterSI/golet/backend/ui"
)

//go:generate go run ./internal/gen

// Preset selects one of the five reference animations.
type Preset uint8

const (
	Medium Preset = iota
	Small
	Line
	PulseInner
	PulseOutside
)

// Mode selects the animation of the edge light.
type Mode uint8

const (
	// Rotate moves the highlight and palette around the perimeter.
	Rotate Mode = iota
	// Pulse selects PulseInner when Preset is Medium; retained for compatibility.
	Pulse
)

// Palette selects a built-in set of edge lights.
type Palette uint8

const (
	// Theme uses the current accent color.
	Theme Palette = iota
	Colorful
	Ocean
	Sunset
	Mono
)

// Options configure a border. Zero values select the defaults documented below.
// Dimensions are in DIPs. Allow at least 12 DIPs of padding for content.
// PulseOutside also needs surrounding space for its exterior glow.
type Options struct {
	// Preset selects Medium, Small, Line, PulseInner, or PulseOutside.
	Preset Preset
	// Brightness and Saturation override reference defaults when positive.
	Brightness, Saturation float32
	// HueRange overrides the rotating presets' hue excursion in degrees.
	HueRange float32
	// StaticColors disables hue changes, while retaining spatial animation.
	StaticColors bool
	// Strength scales the light from 0 to 1; zero selects the default 0.7. Use Disabled to hide it.
	Strength float32
	// Glow is the inward mask depth in DIPs, default 28, maximum 48.
	Glow float32
	// Palette selects the built-in colors when Colors is empty.
	Palette Palette
	// Radius is the container's corner radius; defaults to 16.
	Radius float32
	// Width is the crisp border width; defaults to 1, limited to 8.
	Width float32
	// Duration is one animation cycle; defaults to the reference preset cycle.
	Duration time.Duration
	// Colors is a cyclic palette, copied when Box is built. Empty uses the theme accent.
	Colors []ui.Color
	// Mode selects Rotate (default) or Pulse.
	Mode Mode
	// Paused holds the current animation phase, resuming from it when cleared.
	Paused bool
	// Disabled removes the effect and stops requesting frames.
	Disabled bool
}

type stateKey struct{}
type clock struct {
	start   time.Time
	held    time.Duration
	stopped bool
	layers  [3]rasterLayer
}

// Box creates a normal layout container with decorative light over its border.
// Children remain interactive and accessible. It follows reduced motion and
// only requests repaint frames while visible and animated. Use a stable Key
// on an enclosing element when showing several reorderable containers.
func Box(c *ui.Context, o Options) ui.Element {
	preset := o.Preset
	if preset > PulseOutside {
		preset = Medium
	}
	if preset == Medium && o.Mode == Pulse {
		preset = PulseInner
	}
	radius := dimension(o.Radius, 16, 4096)
	if o.Radius == 0 && preset == Small {
		radius = 32
	}
	duration := o.Duration
	if duration <= 0 {
		duration = 1960 * time.Millisecond
		if preset == Line {
			duration = 3100 * time.Millisecond
		}
		if preset >= PulseInner {
			duration = 2300 * time.Millisecond
		}
	}
	colors := append([]ui.Color(nil), o.Colors...)
	if len(colors) == 0 && o.Palette == Theme {
		colors = []ui.Color{c.Theme().Accent, c.Theme().Accent.Mix(ui.RGB(255, 255, 255), .55)}
	}
	strength := dimension(o.Strength, .7, 1)
	dark := c.Theme().Dark
	e := ui.Box(c).Radius(radius)
	now := c.Now()
	k := ui.Local(e, stateKey{}, func() clock {
		state := clock{start: now}
		for i := range state.layers {
			state.layers[i].bitmap = ui.NewBitmap(image.NewRGBA(image.Rect(0, 0, 1, 1)))
		}
		return state
	})
	stopped := o.Paused || o.Disabled || c.Preferences().ReduceMotion
	if stopped && !k.stopped {
		k.held = now.Sub(k.start)
	}
	if !stopped && k.stopped {
		k.start = now.Add(-k.held)
	}
	k.stopped = stopped
	if o.Disabled {
		return e
	}
	e.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 {
			return
		}
		elapsed := k.held
		if !k.stopped {
			elapsed = p.Now().Sub(k.start)
			p.After(time.Second / 30)
		}
		field := makeField(float64(r.W), float64(r.H), elapsed.Seconds(), duration.Seconds(), preset, o, dark, colors)
		renderLayers(&k.layers, float64(r.W), float64(r.H), float64(radius), field, o, float64(strength))
	})
	e.Children(func() {
		overlay := ui.Box(c).Absolute().Left(0).Right(0).Top(0).Bottom(0).Radius(radius).PassThrough()
		if preset != PulseOutside {
			overlay.Clip()
		}
		overlay.Children(func() {
			for _, i := range []int{1, 0, 2} {
				padding := float32(0)
				if preset == PulseOutside {
					padding = 60
				}
				ui.Image(c, k.layers[i].bitmap).Absolute().Left(-padding).Right(-padding).Top(-padding).Bottom(-padding).Fit(ui.FillBox).PassThrough()
			}
		})
	})
	return e
}

func dimension(v, fallback, limit float32) float32 {
	if v <= 0 || math.IsNaN(float64(v)) || math.IsInf(float64(v), 0) {
		return fallback
	}
	return min(v, limit)
}
