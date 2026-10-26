// Author: HyperparameterSI

package ui

import "runtime"

// Theme holds the colors and metrics widgets use. Change a copy of
// LightTheme or DarkTheme and set it with Context.SetTheme.
type Theme struct {
	Dark bool
	// Tokens holds optional, namespaced semantic colors for component libraries.
	// Core controls do not interpret them. The immutable pointer also preserves
	// Theme's comparability for render caches.
	Tokens *ThemeTokens
	// Optional materials replace styled control and panel backgrounds.
	// Nil keeps the ordinary solid theme. Materials belong to the caller;
	// core UI never imports a renderer plugin. High contrast uses solid fills.
	ControlMaterial Material
	AccentMaterial  Material
	PanelMaterial   Material
	// Background fills the window.
	Background Color
	// Surface is the face of buttons, inputs and other controls;
	// SurfaceHover and SurfacePressed while hovered or pressed.
	Surface        Color
	SurfaceHover   Color
	SurfacePressed Color
	// Border outlines cards, buttons and popovers.
	Border Color
	// InputBorder outlines editable fields; zero falls back to Border.
	InputBorder Color
	// Sidebar fills navigation panels; zero falls back to Surface.
	Sidebar Color
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

// inputBorderColor preserves themes created before editable fields had a
// separate token. Resolve it at use time so custom Border values still work.
func (t *Theme) inputBorderColor() Color {
	if t.InputBorder == Transparent {
		return t.Border
	}
	return t.InputBorder
}

func (t *Theme) sidebarColor() Color {
	if t.Sidebar == Transparent {
		return t.Surface
	}
	return t.Sidebar
}

// LightTheme returns the neutral shadcn palette for a light appearance.
// Colors are precomputed sRGB equivalents of the shadcn OKLCH tokens,
// keeping gamut conversion out of the default startup path.
func LightTheme() *Theme {
	return &Theme{
		Background:     Hex("#ffffff"),
		Surface:        Hex("#ffffff"),
		SurfaceHover:   Hex("#f5f5f5"),
		SurfacePressed: Hex("#e5e5e5"),
		Border:         Hex("#e5e5e5"),
		Sidebar:        Hex("#fafafa"),
		Text:           Hex("#171717"),
		TextMuted:      Hex("#737373"),
		Accent:         Hex("#171717"),
		AccentHover:    Hex("#262626"),
		AccentPressed:  Hex("#404040"),
		AccentText:     Hex("#fafafa"),
		Danger:         Hex("#e7000b"),
		Warning:        Hex("#d97706"),
		Success:        Hex("#16a34a"),
		Selection:      RGBA(23, 23, 23, 0.12),
		Focus:          Hex("#a3a3a3"),
		Scrollbar:      RGBA(0, 0, 0, 0.32),
		Radius:         10,
		Spacing:        4,
		FontSize:       defaultFontSize(),
	}
}

// DarkTheme returns the theme for a dark appearance.
func DarkTheme() *Theme {
	return &Theme{
		Dark:           true,
		Background:     Hex("#0a0a0a"),
		Surface:        Hex("#171717"),
		SurfaceHover:   Hex("#262626"),
		SurfacePressed: Hex("#404040"),
		Border:         RGBA(255, 255, 255, 0.10),
		InputBorder:    RGBA(255, 255, 255, 0.15),
		Sidebar:        Hex("#171717"),
		Text:           Hex("#fafafa"),
		TextMuted:      Hex("#a3a3a3"),
		Accent:         Hex("#e5e5e5"),
		AccentHover:    Hex("#d4d4d4"),
		AccentPressed:  Hex("#a3a3a3"),
		AccentText:     Hex("#171717"),
		Danger:         Hex("#ff6467"),
		Warning:        Hex("#f59e0b"),
		Success:        Hex("#22c55e"),
		Selection:      RGBA(250, 250, 250, 0.18),
		Focus:          Hex("#737373"),
		Scrollbar:      RGBA(255, 255, 255, 0.35),
		Radius:         10,
		Spacing:        4,
		FontSize:       defaultFontSize(),
	}
}
