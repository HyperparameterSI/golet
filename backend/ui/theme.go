// Author: HyperparameterSI

package ui

import "runtime"

// Theme holds the colors and metrics widgets use. Change a copy of
// LightTheme or DarkTheme and set it with Context.SetTheme.
type Theme struct {
	Dark bool
	// Background fills the window.
	Background Color
	// Surface is the face of buttons, inputs and other controls;
	// SurfaceHover and SurfacePressed while hovered or pressed.
	Surface        Color
	SurfaceHover   Color
	SurfacePressed Color
	// Border outlines controls.
	Border Color
	// Text and TextMuted color text; TextMuted is for secondary text and
	// placeholders.
	Text      Color
	TextMuted Color
	// Accent colors primary buttons, checked controls and focus rings;
	// AccentText is text on it.
	Accent        Color
	AccentHover   Color
	AccentPressed Color
	AccentText    Color
	Danger        Color
	// Warning and Success color what calls for attention and what went
	// well, as a Meter's levels.
	Warning Color
	Success Color
	// Selection highlights selected text.
	Selection Color
	// Focus is the ring around the control with the keyboard focus.
	Focus Color
	// Inverse fills tooltips and toasts, the theme turned over so that
	// they stand out from what they are over: dark in a light theme, light
	// in a dark one. InverseText is text on it. Left zero, they are Text
	// and Background.
	Inverse     Color
	InverseText Color
	// Scrollbar colors scroll bar thumbs, ScrollbarWidth DIPs wide (6 by
	// default; 0 is 6).
	Scrollbar      Color
	ScrollbarWidth float32
	// Radius rounds the corners of controls.
	Radius float32
	// Spacing is the unit of the room widgets leave inside and between
	// their parts: their paddings and gaps, and the sizes of check boxes,
	// switches, sliders and the rows of tables and trees, are multiples
	// of it. 4 by default: 3 makes widgets compact, 5 roomy; 0 is 4.
	Spacing float32
	// FontSize is the size of text, Font its family ("" is the system's).
	FontSize float32
	Font     string
}

// Space returns n units of the theme's Spacing, to size elements of your
// own in step with the widgets: ui.Column(c).Gap(t.Space(2)).
func (t *Theme) Space(n float32) float32 {
	s := t.Spacing
	if s <= 0 {
		s = 4
	}
	return n * s
}

// inverse returns the fill of tooltips and toasts, and the color of text
// on it.
func (t *Theme) inverse() (fill, text Color) {
	fill, text = t.Inverse, t.InverseText
	if fill == Transparent {
		fill = t.Text
	}
	if text == Transparent {
		text = t.Background
	}
	return fill, text
}

func (t *Theme) scrollbarWidth() float32 {
	if t.ScrollbarWidth <= 0 {
		return 6
	}
	return t.ScrollbarWidth
}

// Rem returns n times the theme's FontSize, as CSS's rem, to size elements
// with the text.
func (t *Theme) Rem(n float32) float32 {
	if t.FontSize <= 0 {
		return n * defaultFontSize()
	}
	return n * t.FontSize
}

func defaultFontSize() float32 {
	if runtime.GOOS == "darwin" {
		return 13
	}
	return 14
}

// LightTheme returns the theme for a light appearance.
func LightTheme() *Theme {
	return &Theme{
		Background:     Hex("#ffffff"),
		Surface:        Hex("#ffffff"),
		SurfaceHover:   Hex("#f4f4f5"),
		SurfacePressed: Hex("#e4e4e7"),
		Border:         Hex("#e4e4e7"),
		Text:           Hex("#18181b"),
		TextMuted:      Hex("#71717a"),
		Accent:         Hex("#18181b"),
		AccentHover:    Hex("#27272a"),
		AccentPressed:  Hex("#3f3f46"),
		AccentText:     Hex("#ffffff"),
		Danger:         Hex("#dc2626"),
		Warning:        Hex("#d97706"),
		Success:        Hex("#16a34a"),
		Selection:      RGBA(24, 24, 27, 0.12),
		Focus:          RGBA(113, 113, 122, 0.45),
		Scrollbar:      RGBA(0, 0, 0, 0.32),
		Radius:         8,
		Spacing:        4,
		FontSize:       defaultFontSize(),
	}
}

// DarkTheme returns the theme for a dark appearance.
func DarkTheme() *Theme {
	return &Theme{
		Dark:           true,
		Background:     Hex("#09090b"),
		Surface:        Hex("#18181b"),
		SurfaceHover:   Hex("#27272a"),
		SurfacePressed: Hex("#3f3f46"),
		Border:         Hex("#303036"),
		Text:           Hex("#f4f4f5"),
		TextMuted:      Hex("#a1a1aa"),
		Accent:         Hex("#fafafa"),
		AccentHover:    Hex("#e4e4e7"),
		AccentPressed:  Hex("#d4d4d8"),
		AccentText:     Hex("#18181b"),
		Danger:         Hex("#ef4444"),
		Warning:        Hex("#f59e0b"),
		Success:        Hex("#22c55e"),
		Selection:      RGBA(244, 244, 245, 0.18),
		Focus:          RGBA(161, 161, 170, 0.5),
		Scrollbar:      RGBA(255, 255, 255, 0.35),
		Radius:         8,
		Spacing:        4,
		FontSize:       defaultFontSize(),
	}
}
