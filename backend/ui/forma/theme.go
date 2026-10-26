// Author: HyperparameterSI

// Package forma provides opt-in, pure Go native components based on shadcn's
// Base UI Maia registry. Reference sources and attribution live in docs/references/shadcn-base-maia/.
// Build components inside ui.View; never retain an Element between builds.
package forma

import "github.com/HyperparameterSI/golet/backend/ui"

// Theme returns the neutral Maia palette. Keep a copy per window and pass it to
// Context.SetTheme; the component layer also respects caller theme overrides.
func Theme(dark bool) *ui.Theme {
	t := ui.FormaTheme(dark)
	t.Radius = 16
	t.Spacing = 4
	t.FontSize = 14
	return t
}

// Variant selects the six styles shared by buttons and badges.
type Variant uint8

const (
	Default Variant = iota
	Outline
	Secondary
	Ghost
	Destructive
	Link
)

// Size selects Maia's 24, 32, 36 or 40 DIP button height.
type Size uint8

const (
	Normal Size = iota
	XS
	Small
	Large
	Icon
	IconXS
	IconSmall
	IconLarge
)

// ButtonOptions are applied before input is handled.
type ButtonOptions struct {
	Variant  Variant
	Size     Size
	Disabled bool
	Invalid  bool
}

func colors(c *ui.Context, v Variant) (bg, fg, hover, border ui.Color) {
	t := c.Theme()
	fg = t.Text
	switch v {
	case Default:
		bg = t.Accent
		fg = t.AccentText
		hover = bg.Alpha(.8)
	case Outline:
		bg = Color(c, "input")
		bg = bg.Alpha(.3)
		hover = t.Border.Alpha(.5)
		border = t.Border
	case Secondary:
		bg = token(c, "secondary", t.SurfaceHover)
		fg = token(c, "secondary-foreground", t.Text)
		hover = bg.Mix(t.Text, .05)
	case Ghost:
		hover = token(c, "accent", t.SurfaceHover)
	case Destructive:
		bg = t.Danger.Alpha(.1)
		hover = t.Danger.Alpha(.2)
		fg = t.Danger
		if t.Dark {
			bg = t.Danger.Alpha(.2)
			hover = t.Danger.Alpha(.3)
		}
	case Link:
		fg = t.Accent
	}
	return
}

func focus(c *ui.Context, e ui.Element, radius float32, invalid bool) {
	ring := c.Theme().Focus.Alpha(.5)
	border := c.Theme().Focus
	if invalid {
		ring = c.Theme().Danger.Alpha(.2)
		border = c.Theme().Danger
	}
	e.FocusRing(false).DrawOver(func(p *ui.Painter, r ui.Rect) {
		if e.FocusVisible() {
			outer := ui.Rect{X: r.X - 3, Y: r.Y - 3, W: r.W + 6, H: r.H + 6}
			p.Stroke(outer, ring, radius+3, 3)
			p.Stroke(r, border, radius, 1)
		} else if invalid {
			p.Stroke(r, border, radius, 1)
		}
	})
}

// Button includes pointer/keyboard activation and all Maia variants and sizes.
func Button(c *ui.Context, label string, options ...ButtonOptions) ui.Element {
	o := ButtonOptions{}
	if len(options) > 0 {
		o = options[0]
	}
	height, padding, font := float32(36), float32(12), float32(14)
	switch o.Size {
	case XS, IconXS:
		height, padding, font = 24, 10, 12
	case Small, IconSmall:
		height = 32
	case Large, IconLarge:
		height, padding = 40, 16
	}
	bg, fg, hover, border := colors(c, o.Variant)
	e := ui.ButtonBase(c).Height(height).Padding(0, padding).Gap(6).Radius(radius(c, 32)).FontSize(fontSize(c, font)).FontWeight(500).TextColor(fg).Disabled(o.Disabled).Label(label)
	if o.Size >= Icon {
		e.Width(height).Padding(0)
	}
	if o.Disabled {
		e.Opacity(.5)
	}
	fill := bg
	if e.Hovered() && !o.Disabled {
		fill = hover
		if o.Variant == Ghost {
			e.TextColor(token(c, "accent-foreground", c.Theme().Text))
		}
	}
	e.Background(fill).Border(1, border)
	if e.Pressed() && !o.Disabled {
		e.Translate(0, 1)
	}

	focus(c, e, radius(c, 32), o.Invalid)
	e.Children(func() {
		text := ui.Text(c, label)
		if o.Variant == Link && e.Hovered() {
			text.Underline()
		}
	})
	return e
}

// Badge displays a compact noninteractive status with the button color variants.
func Badge(c *ui.Context, label string, variant ...Variant) ui.Element {
	v := Default
	if len(variant) > 0 {
		v = variant[0]
	}
	bg, fg, _, border := colors(c, v)
	return ui.Row(c).Height(20).Padding(0, 8).Gap(4).Center().Radius(radius(c, 32)).Background(bg).Border(1, border).TextColor(fg).FontSize(fontSize(c, 12)).FontWeight(500).Children(func() { ui.Text(c, label) })
}
