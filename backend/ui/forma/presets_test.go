// Author: HyperparameterSI
package forma

import (
	"math"
	"testing"

	"github.com/HyperparameterSI/golet/backend/ui"
)

func TestThemeTokensAreIndependent(t *testing.T) {
	overrides := map[string]ui.Color{"primary": ui.Hex("#123456"), "popover": ui.Hex("#abcdef"), "secondary": ui.Transparent}
	a := NewTheme(ThemeOptions{Preset: Blue, Tokens: overrides})
	b := NewTheme(ThemeOptions{Preset: Blue})
	overrides["primary"] = ui.Hex("#ffffff")
	if a.Accent != ui.Hex("#123456") || a.Accent == b.Accent {
		t.Fatal("override leaked between themes")
	}
	_ = *a == *b // Theme must remain comparable for existing render caches.
	tt := ui.NewTester(func(c *ui.Context) {
		c.SetTheme(a)
		bg, _, _, _ := colors(c, Secondary)
		if bg != ui.Transparent {
			t.Fatal("transparent override treated as absent")
		}
		if Color(c, "popover") != ui.Hex("#abcdef") {
			t.Fatal("popover override lost")
		}
		original := c.Theme()
		selected := "one"
		Sidebar(c, &selected, func() { SidebarItem(c, "one", "One") })
		if c.Theme() != original {
			t.Fatal("sidebar leaked its local theme")
		}
	}, 400, 300)
	tt.Frame()
}

func TestPresetAppearanceAndMetrics(t *testing.T) {
	for _, preset := range Presets() {
		for _, dark := range []bool{false, true} {
			t.Run(string(preset)+map[bool]string{false: "/light", true: "/dark"}[dark], func(t *testing.T) {
				r := float32(0)
				theme := NewTheme(ThemeOptions{Preset: preset, Dark: dark, Radius: &r, FontSize: 16})
				tt := ui.NewTester(func(c *ui.Context) {
					c.SetTheme(theme)
					if radius(c, 32) != 0 || fontSize(c, 14) != 16 {
						t.Fatal("custom metrics not applied")
					}
					ui.Column(c).Children(func() {
						Button(c, "Primary")
						Button(c, "Secondary", ButtonOptions{Variant: Secondary})
						Card(c).Children(func() { CardTitle(c, "Card") })
						v := ""
						Input(c, &v)
					})
				}, 400, 300)
				tt.Frame()
				if !tt.HasText("Primary") {
					t.Fatal("missing themed button")
				}
			})
		}
	}
}

func TestThemeInvalidMetricsFallBack(t *testing.T) {
	r := float32(math.NaN())
	theme := NewTheme(ThemeOptions{Radius: &r, FontSize: float32(math.Inf(1))})
	if theme.Radius != 16 || theme.FontSize != 14 {
		t.Fatalf("invalid theme metrics: %v %v", theme.Radius, theme.FontSize)
	}
}
