// Author: HyperparameterSI
package forma

import (
	"math"

	"github.com/HyperparameterSI/golet/backend/ui"
)

// Preset is a native palette inspired by Tailwind's color families. These are
// Golet presets, not a claim to reproduce every upstream shadcn preset.
type Preset string

const (
	Neutral Preset = "Neutral"
	Zinc    Preset = "Zinc"
	Stone   Preset = "Stone"
	Blue    Preset = "Blue"
	Violet  Preset = "Violet"
	Emerald Preset = "Emerald"
	Rose    Preset = "Rose"
)

// Presets returns an independent list suitable for a theme picker.
func Presets() []Preset { return []Preset{Neutral, Zinc, Stone, Blue, Violet, Emerald, Rose} }

// ThemeOptions configures a native theme. Start from a preset and override only
// the semantic tokens you need; light and dark configurations are independent.
// Radius is a pointer so zero deliberately requests square corners. Its default
// is 16 DIP, preserving Golet's existing Maia geometry. FontSize defaults to 14.
type ThemeOptions struct {
	Preset   Preset
	Dark     bool
	Radius   *float32
	FontSize float32
	Font     string
	Tokens   map[string]ui.Color
}

// NewTheme resolves a caller-owned theme once, outside the view's frame loop.
// Overrides use shadcn names ("primary", "primary-foreground", "card", etc.).
// No CSS parser, webview, global registry or background work is involved.
func NewTheme(o ThemeOptions) *ui.Theme {
	t := ui.FormaTheme(o.Dark)
	t.Radius, t.Spacing, t.FontSize = 16, 4, 14
	if o.Radius != nil && finite(*o.Radius) && *o.Radius >= 0 {
		t.Radius = *o.Radius
	}
	if finite(o.FontSize) && o.FontSize > 0 {
		t.FontSize = o.FontSize
	}
	t.Font = o.Font
	if o.Preset == Zinc {
		if o.Dark {
			t.Background = ui.Hex("#09090b")
			t.Surface = ui.Hex("#18181b")
			t.SurfaceHover = ui.Hex("#27272a")
			t.TextMuted = ui.Hex("#a1a1aa")
			t.Border = ui.Hex("#3f3f46")
		} else {
			t.Text = ui.Hex("#18181b")
			t.SurfaceHover = ui.Hex("#f4f4f5")
			t.TextMuted = ui.Hex("#71717a")
			t.Border = ui.Hex("#e4e4e7")
		}
	}
	if o.Preset == Stone {
		if o.Dark {
			t.Background = ui.Hex("#0c0a09")
			t.Surface = ui.Hex("#1c1917")
			t.SurfaceHover = ui.Hex("#292524")
			t.TextMuted = ui.Hex("#a8a29e")
			t.Border = ui.Hex("#44403c")
		} else {
			t.Background = ui.Hex("#fafaf9")
			t.Text = ui.Hex("#1c1917")
			t.SurfaceHover = ui.Hex("#f5f5f4")
			t.TextMuted = ui.Hex("#78716c")
			t.Border = ui.Hex("#e7e5e4")
		}
		t.Sidebar = t.Surface
		if !o.Dark {
			t.Sidebar = t.Background
		}
	}
	if o.Preset == Zinc && o.Dark {
		t.Sidebar = t.Surface
	}
	t.InputBorder = t.Border
	light, dark := "", ""
	switch o.Preset {
	case Blue:
		light, dark = "#2563eb", "#60a5fa"
	case Violet:
		light, dark = "#7c3aed", "#a78bfa"
	case Emerald:
		light, dark = "#047857", "#34d399"
	case Rose:
		light, dark = "#be123c", "#fb7185"
	}
	if light != "" {
		t.Accent = ui.Hex(light)
		t.AccentText = ui.Hex("#ffffff")
		if o.Dark {
			t.Accent = ui.Hex(dark)
			t.AccentText = ui.Hex("#09090b")
		}
		t.Focus = t.Accent
	}
	values := map[string]ui.Color{
		"background": t.Background, "foreground": t.Text, "card": t.Surface, "card-foreground": t.Text,
		"popover": t.Surface, "popover-foreground": t.Text, "primary": t.Accent, "primary-foreground": t.AccentText,
		"secondary": t.SurfaceHover, "secondary-foreground": t.Text, "muted": t.SurfaceHover, "muted-foreground": t.TextMuted,
		"accent": t.SurfaceHover, "accent-foreground": t.Text, "destructive": t.Danger, "border": t.Border, "input": t.InputBorder, "ring": t.Focus,
		"sidebar": t.Sidebar, "sidebar-foreground": t.Text, "sidebar-primary": t.Accent, "sidebar-primary-foreground": t.AccentText,
		"sidebar-accent": t.SurfaceHover, "sidebar-accent-foreground": t.Text, "sidebar-border": t.Border, "sidebar-ring": t.Focus,
	}
	charts := []string{"#2563eb", "#7c3aed", "#059669", "#d97706", "#e11d48"}
	if o.Dark {
		charts = []string{"#60a5fa", "#a78bfa", "#34d399", "#fbbf24", "#fb7185"}
	}
	for i, name := range []string{"chart-1", "chart-2", "chart-3", "chart-4", "chart-5"} {
		values[name] = ui.Hex(charts[i])
	}
	for name, color := range o.Tokens {
		values[name] = color
	}
	namespaced := make(map[string]ui.Color, len(values))
	for name, color := range values {
		namespaced["forma."+name] = color
	}
	t.Tokens = ui.NewThemeTokens(namespaced)
	// Map shared roles into core UI too, including text selection and base widgets.
	t.Background, t.Text = values["background"], values["foreground"]
	t.Surface, t.SurfaceHover = values["card"], values["muted"]
	t.TextMuted, t.Border, t.InputBorder = values["muted-foreground"], values["border"], values["input"]
	t.Accent, t.AccentText, t.Focus, t.Danger = values["primary"], values["primary-foreground"], values["ring"], values["destructive"]
	t.Sidebar = values["sidebar"]
	t.AccentHover = t.Accent.Mix(t.Text, .08)
	t.AccentPressed = t.Accent.Mix(t.Text, .16)
	t.SurfacePressed = t.SurfaceHover.Mix(t.Text, .08)
	t.Selection = t.Accent.Alpha(.25)
	return t
}

func finite(v float32) bool { return !math.IsNaN(float64(v)) && !math.IsInf(float64(v), 0) }
func token(c *ui.Context, name string, fallback ui.Color) ui.Color {
	if value, ok := c.Theme().Tokens.Lookup("forma." + name); ok {
		return value
	}
	return fallback
}

// Color resolves a semantic token for custom native compositions. Legacy themes
// fall back to their shared roles; unknown names return transparent.
func Color(c *ui.Context, name string) ui.Color {
	t := c.Theme()
	fallback := ui.Transparent
	switch name {
	case "background":
		fallback = t.Background
	case "card", "popover":
		fallback = t.Surface
	case "foreground", "card-foreground", "popover-foreground", "secondary-foreground", "accent-foreground", "sidebar-foreground", "sidebar-accent-foreground":
		fallback = t.Text
	case "primary", "sidebar-primary":
		fallback = t.Accent
	case "primary-foreground", "sidebar-primary-foreground":
		fallback = t.AccentText
	case "secondary", "muted", "accent", "sidebar-accent":
		fallback = t.SurfaceHover
	case "muted-foreground":
		fallback = t.TextMuted
	case "border", "sidebar-border":
		fallback = t.Border
	case "input":
		fallback = t.InputBorder
		if fallback == ui.Transparent {
			fallback = t.Border
		}
	case "ring", "sidebar-ring":
		fallback = t.Focus
	case "destructive":
		fallback = t.Danger
	case "sidebar":
		fallback = t.Sidebar
		if fallback == ui.Transparent {
			fallback = t.Surface
		}
	}
	return token(c, name, fallback)
}

// ChartColors returns a copy of the five semantic chart colors.
func ChartColors(c *ui.Context) []ui.Color {
	result := make([]ui.Color, 5)
	for i, name := range []string{"chart-1", "chart-2", "chart-3", "chart-4", "chart-5"} {
		result[i] = token(c, name, c.Theme().Accent)
	}
	return result
}
func radius(c *ui.Context, base float32) float32 {
	r := c.Theme().Radius
	if !finite(r) || r < 0 {
		r = 16
	}
	return base * r / 16
}
func fontSize(c *ui.Context, base float32) float32 {
	v := c.Theme().FontSize
	if !finite(v) || v <= 0 {
		v = 14
	}
	return base * v / 14
}
