// Author: HyperparameterSI

package ui

import (
	"testing"
	"time"
)

func TestSpacingScalesWidgets(t *testing.T) {
	type sizes struct{ button, check, tab Rect }
	measure := func(spacing float32) sizes {
		tt := coreNewTester(func(c *context) {
			th := *c.Theme()
			th.Spacing = spacing
			c.SetTheme(&th)
			coreColumn(c).AlignItems(Start).Children(func() {
				coreButton(c, "OK").Label("button")
				on, tab := false, 0
				coreCheckbox(c, &on, "").Label("check")
				coreTabs(c, &tab, "One", "Two")
			})
		}, 400, 300)
		var s sizes
		s.button, _ = tt.Find("button")
		s.check, _ = tt.Find("check")
		s.tab, _ = tt.Find("One")
		return s
	}
	compact, normal, roomy, unset := measure(3), measure(4), measure(5), measure(0)
	if normal.check.W != 16 || normal.check.H != 16 {
		t.Errorf("the check box is %vx%v at the default spacing, not 16x16", normal.check.W, normal.check.H)
	}
	if unset != normal {
		t.Errorf("no spacing lays widgets out as %+v, not as the default %+v", unset, normal)
	}
	for name, got := range map[string][3]Rect{
		"button":    {compact.button, normal.button, roomy.button},
		"check box": {compact.check, normal.check, roomy.check},
	} {
		if !(got[0].W < got[1].W && got[1].W < got[2].W && got[0].H < got[1].H && got[1].H < got[2].H) {
			t.Errorf("the %s does not grow with the spacing: %v", name, got)
		}
	}
	// The text of a tab sits in its padding: 3 units to the left.
	for _, s := range []sizes{compact, normal, roomy} {
		if s.tab.W <= 0 {
			t.Fatal("no tab")
		}
	}
	if !(compact.tab.X < normal.tab.X && normal.tab.X < roomy.tab.X) {
		t.Errorf("the tabs' padding does not follow the spacing: %v, %v, %v", compact.tab.X, normal.tab.X, roomy.tab.X)
	}
}

func TestParseHex(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want Color
		ok   bool
	}{
		{"#2563eb", Color{R: 0x25, G: 0x63, B: 0xeb, A: 255}, true},
		{"2563EB", Color{R: 0x25, G: 0x63, B: 0xeb, A: 255}, true},
		{" #00ff0080 ", Color{R: 0, G: 255, B: 0, A: 0x80}, true},
		{"#fA0", Color{R: 255, G: 0xaa, B: 0, A: 255}, true},
		{"#fa08", Color{R: 255, G: 0xaa, B: 0, A: 0x88}, true},
		{"#12345", Color{}, false},
		{"#1234567", Color{}, false},
		{"#123456789", Color{}, false},
		{"#12345g", Color{}, false},
		{"#+12345", Color{}, false},
		{"#é12", Color{}, false},
		{"", Color{}, false},
	} {
		got, err := parseHex(tc.in)
		if got != tc.want || (err == nil) != tc.ok {
			t.Errorf("parseHex(%q) = %v, %v; want %v, ok %v", tc.in, got, err, tc.want, tc.ok)
		}
	}
	if n := testing.AllocsPerRun(10, func() { Hex("#2563eb") }); n != 0 {
		t.Errorf("Hex allocates %v times", n)
	}
}

func TestInverseColors(t *testing.T) {
	th := *DarkTheme()
	if fill, text := th.inverse(); fill != th.Text || text != th.Background {
		t.Errorf("a theme without inverse colors gives %v on %v, not its text on its background", text, fill)
	}
	th.Inverse = Hex("#3a3a40")
	if fill, text := th.inverse(); fill != Hex("#3a3a40") || text != th.Background {
		t.Errorf("Inverse alone gives %v on %v", text, fill)
	}
	th.InverseText = Hex("#f2f2f7")
	if fill, text := th.inverse(); fill != Hex("#3a3a40") || text != Hex("#f2f2f7") {
		t.Errorf("both inverse colors give %v on %v", text, fill)
	}
}

// inverseFills shows a tooltip and a toast in a window of theme th, and
// returns the colors they are filled with.
func inverseFills(t *testing.T, th *Theme) (tooltip, toast Color) {
	t.Helper()
	tt := coreNewTester(func(c *context) {
		c.SetTheme(th)
		coreColumn(c).Padding(40).Children(func() {
			if coreButton(c, "Save").Tooltip("Save the note").Clicked() {
				c.Toast("Saved")
			}
		})
	}, 400, 300)
	// Left of the text, inside the padding.
	fill := func(s string) Color {
		r, ok := tt.Find(s)
		if !ok {
			t.Fatalf("no %q; texts %q", s, tt.Texts())
		}
		px := tt.Image().RGBAAt(int(r.X)-3, int(r.Y+r.H/2))
		return Color{R: px.R, G: px.G, B: px.B, A: px.A}
	}
	b, _ := tt.Find("Save")
	tt.Move(b.X+b.W/2, b.Y+b.H/2)
	tt.rt.tips.hoverSince = tt.now.Add(-time.Second)
	tt.Frame()
	tooltip = fill("Save the note")
	tt.Click("Save")
	// Past its fading in.
	tt.rt.toasts[0].at = tt.now.Add(-time.Second)
	tt.Frame()
	return tooltip, fill("Saved")
}

func TestInverseFillsTooltipsAndToasts(t *testing.T) {
	th := *DarkTheme()
	if tooltip, toast := inverseFills(t, &th); tooltip != th.Text || toast != th.Text {
		t.Errorf("a dark theme fills a tooltip with %v and a toast with %v, not its text color %v", tooltip, toast, th.Text)
	}
	th.Inverse, th.InverseText = Hex("#3a3a40"), Hex("#f2f2f7")
	if tooltip, toast := inverseFills(t, &th); tooltip != th.Inverse || toast != th.Inverse {
		t.Errorf("a theme fills a tooltip with %v and a toast with %v, not its Inverse %v", tooltip, toast, th.Inverse)
	}
}

// TestSemanticThemeFallbacks protects applications that supply their own
// themes without the newer input and navigation tokens.
func TestSemanticThemeFallbacks(t *testing.T) {
	th := &Theme{Border: Hex("#345678"), Surface: Hex("#abcdef")}
	if th.inputBorderColor() != th.Border || th.sidebarColor() != th.Surface {
		t.Fatal("unset semantic colors must inherit existing custom colors")
	}
	th.InputBorder, th.Sidebar = RGBA(255, 255, 255, .15), Hex("#171717")
	if th.inputBorderColor() != th.InputBorder || th.sidebarColor() != th.Sidebar {
		t.Fatal("explicit semantic colors must take precedence")
	}
	for _, dark := range []bool{false, true} {
		base := LightTheme()
		if dark {
			base = DarkTheme()
		}
		preset := FormaTheme(dark)
		base.FontSize = 14
		if *preset != *base {
			t.Fatal("native presets must share the same semantic palette")
		}
		preset.Radius = 99
		if FormaTheme(dark).Radius == 99 {
			t.Fatal("customizing one theme must not mutate other themes")
		}
	}
}

func TestStandardControlHeightsAgree(t *testing.T) {
	for _, dark := range []bool{false, true} {
		value, selected := "", "One"
		tt := coreNewTester(func(c *context) {
			c.SetTheme(FormaTheme(dark))
			coreColumn(c).Children(func() {
				coreButton(c, "Secondary").Label("secondary")
				corePrimaryButton(c, "Primary").Label("primary")
				coreTextInput(c, &value).Label("input")
				coreSelect(c, &selected, []string{"One", "Two"}).Label("select")
			})
		}, 400, 400)
		want := at(t, tt, "secondary").H
		for _, name := range []string{"primary", "input", "select"} {
			if got := at(t, tt, name).H; got != want {
				t.Errorf("%s height %v differs from secondary %v (dark=%v)", name, got, want, dark)
			}
		}
	}
}
