// Author: HyperparameterSI

package ui

import (
	"runtime"
	"time"
)

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
	// DisableAutoDimming keeps disabled elements at full opacity. Applications
	// opting in must provide their own disabled colors; input, focus and
	// accessibility still use the element's disabled state.
	DisableAutoDimming bool
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
	// InactiveAccent is what the accent fades to in a window that is not
	// drawn as active (Context.WindowActive), as AppKit fades an inactive
	// window's: the fill of selected rows, of checked check boxes, radio
	// buttons and switches, of sliders and progress bars, the line of the
	// selected tab and the icons of a sidebar, with AccentText on it as on
	// the accent; primary buttons show as the others there. Left zero, it
	// is a gray AccentText stands out on; set it to Accent to keep the
	// accent in every window.
	InactiveAccent Color
	Danger         Color
	// Warning and Success color what calls for attention and what went
	// well, as a Meter's levels.
	Warning Color
	Success Color
	// Selection highlights selected text.
	Selection Color
	// Focus is the ring around the control with the keyboard focus,
	// FocusRingWidth DIPs wide (2 by default; 0 is 2) and FocusRingOffset
	// DIPs outside its edge (1 in LightTheme and DarkTheme): 0 puts the
	// ring on the edge, and a negative offset inside it.
	// Element.FocusRingOffset and Element.FocusRingColor change them for
	// one element.
	Focus           Color
	FocusRingWidth  float32
	FocusRingOffset float32
	// Inverse fills tooltips and toasts, the theme turned over so that
	// they stand out from what they are over: dark in a light theme, light
	// in a dark one. InverseText is text on it. Left zero, they are Text
	// and Background.
	Inverse     Color
	InverseText Color
	// Scrollbar colors scroll bar thumbs, ScrollbarWidth DIPs wide (6 by
	// default; 0 is 6): the defaults of ScrollbarStyle, which overrides
	// them where it sets a field. It styles their tracks, hover and drag
	// colors, the width of their bars and thumbs, their corners and
	// shortest thumb; a theme that sets neither draws as it always has.
	Scrollbar      Color
	ScrollbarWidth float32
	ScrollbarStyle ScrollbarStyle
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
	// TooltipDelay is how long the pointer rests on an element before its
	// tooltip shows (Tooltip, TooltipBase): 0 is 0.6 seconds, and a
	// negative delay shows it at once. Within two thirds of it after
	// another tip goes, as the pointer moves along a toolbar, the next one
	// shows at once: 0.4 seconds by default.
	TooltipDelay time.Duration

	// followed is what WithPreferences changed, nil where it has not.
	followed *followedTheme
}

// followedTheme is a theme before and after following the desktop's
// preferences.
type followedTheme struct {
	_        [0]func() // used by pointer, so there is no value equality to link
	from, to Theme
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

// inactiveAccent returns what the accent fades to in a window that is not
// drawn as active (InactiveAccent).
func (t *Theme) inactiveAccent() Color {
	if t.InactiveAccent != Transparent {
		return t.InactiveAccent
	}
	// As far from AccentText as text on it needs, 4.5 times its contrast.
	far := Color{A: 255}
	if t.AccentText.gray().R <= 127 {
		far = Color{R: 255, G: 255, B: 255, A: 255}
	}
	return t.AccentText.Mix(far, 0.55)
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

// focusRingWidth returns how wide the focus ring is (FocusRingWidth).
func (t *Theme) focusRingWidth() float32 {
	if t.FocusRingWidth <= 0 {
		return 2
	}
	return t.FocusRingWidth
}

func (t *Theme) scrollbarWidth() float32 {
	if t.ScrollbarWidth <= 0 {
		return 6
	}
	return t.ScrollbarWidth
}

// tooltipDelay returns how long the pointer rests on an element before its
// tip shows (TooltipDelay), and how soon after a tip goes the next one
// shows at once.
func (t *Theme) tooltipDelay() (delay, warm time.Duration) {
	switch d := t.TooltipDelay; {
	case d == 0:
		return tooltipDelay, tooltipWarm
	case d < 0:
		return 0, 0
	default:
		return d, d * 2 / 3
	}
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
		Background:      Hex("#ffffff"),
		Surface:         Hex("#ffffff"),
		SurfaceHover:    Hex("#f5f5f5"),
		SurfacePressed:  Hex("#e5e5e5"),
		Border:          Hex("#e5e5e5"),
		Sidebar:         Hex("#fafafa"),
		Text:            Hex("#171717"),
		TextMuted:       Hex("#737373"),
		Accent:          Hex("#171717"),
		AccentHover:     Hex("#262626"),
		AccentPressed:   Hex("#404040"),
		AccentText:      Hex("#fafafa"),
		Danger:          Hex("#e7000b"),
		Warning:         Hex("#d97706"),
		Success:         Hex("#16a34a"),
		Selection:       RGBA(23, 23, 23, 0.12),
		Focus:           Hex("#a3a3a3"),
		Scrollbar:       RGBA(0, 0, 0, 0.32),
		Radius:          10,
		Spacing:         4,
		FontSize:        defaultFontSize(),
		FocusRingWidth:  2,
		FocusRingOffset: 1,
	}
}

// DarkTheme returns the theme for a dark appearance.
func DarkTheme() *Theme {
	return &Theme{
		Dark:            true,
		Background:      Hex("#0a0a0a"),
		Surface:         Hex("#171717"),
		SurfaceHover:    Hex("#262626"),
		SurfacePressed:  Hex("#404040"),
		Border:          RGBA(255, 255, 255, 0.10),
		InputBorder:     RGBA(255, 255, 255, 0.15),
		Sidebar:         Hex("#171717"),
		Text:            Hex("#fafafa"),
		TextMuted:       Hex("#a3a3a3"),
		Accent:          Hex("#e5e5e5"),
		AccentHover:     Hex("#d4d4d4"),
		AccentPressed:   Hex("#a3a3a3"),
		AccentText:      Hex("#171717"),
		Danger:          Hex("#ff6467"),
		Warning:         Hex("#f59e0b"),
		Success:         Hex("#22c55e"),
		Selection:       RGBA(250, 250, 250, 0.18),
		Focus:           Hex("#737373"),
		Scrollbar:       RGBA(255, 255, 255, 0.35),
		Radius:          10,
		Spacing:         4,
		FontSize:        defaultFontSize(),
		FocusRingWidth:  2,
		FocusRingOffset: 1,
	}
}
