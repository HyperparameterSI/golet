// Author: HyperparameterSI

package glass

import (
	"github.com/HyperparameterSI/golet/backend/ui"
	"testing"
)

// A styled control must sample its backdrop; high contrast must not. This
// also catches accidental loss of the public MaterialBuilder adapter.
func TestThemeControlsUseBackdrop(t *testing.T) {
	for _, dark := range []bool{false, true} {
		for _, contrast := range []bool{false, true} {
			th := Theme(dark)
			backdrop := ui.Hex("#2468bc")
			tt := ui.NewTester(func(c *ui.Context) {
				c.SetTheme(th)
				ui.Box(c).Fill().Background(backdrop).Children(func() {
					ui.Button(c, "Action").Size(160, 44)
				})
			}, 200, 80)
			tt.SetPreferences(ui.Preferences{HighContrast: contrast, ReduceMotion: true})
			first := pixelAt(tt.Image(), 30, 22)
			backdrop = ui.Hex("#c45c28")
			tt.Frame()
			second := pixelAt(tt.Image(), 30, 22)
			if contrast && first != second {
				t.Fatalf("high contrast changed with backdrop: %v -> %v", first, second)
			}
			if !contrast && first == second {
				t.Fatalf("glass did not sample backdrop (dark=%v): %v", dark, first)
			}
		}
	}
}

func TestThemeControlsRemainInteractive(t *testing.T) {
	th := Theme(false)
	calls := 0
	choice := "One"
	tt := ui.NewTester(func(c *ui.Context) {
		c.SetTheme(th)
		ui.Column(c).Children(func() {
			if ui.Button(c, "Action").Clicked() {
				calls++
			}
			if ui.Button(c, "Disabled").Disabled(true).Clicked() {
				calls++
			}
			ui.Select(c, &choice, []string{"One", "Two"}).Label("Choice")
		})
	}, 240, 200)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	if err := tt.Click("Action"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Disabled"); err != nil {
		t.Fatal(err)
	}
	if calls != 1 {
		t.Fatalf("enabled/disabled actions: %d", calls)
	}
	if err := tt.Click("Choice"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Two"); err != nil {
		t.Fatal(err)
	}
	if choice != "Two" {
		t.Fatalf("glass popup selection: %s", choice)
	}
}

// Theme slots must invoke MaterialBuilder, not just PaintMaterial.
func TestThemeButtonPressMaterial(t *testing.T) {
	theme := Theme(false)
	tt := ui.NewTester(func(c *ui.Context) {
		c.SetTheme(theme)
		ui.Box(c).Fill().Background(ui.Hex("#000000")).Children(func() {
			ui.Button(c, "Action").Absolute().Left(50).Top(50).Size(200, 60)
		})
	}, 300, 160)
	tt.SetScale(2)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	outside := func() int { return pixelAt(tt.Image(), 99, 160)[0] }
	before := outside()
	tt.Press(150, 80)
	tt.Frame()
	if outside() <= before+50 {
		t.Fatal("theme material did not grow on press")
	}
	tt.Release(150, 80)
	tt.Frame()
	if outside() != before {
		t.Fatal("theme material did not settle after release")
	}
	tt.Frame()
	if tt.FramePending() {
		t.Fatal("settled glass requested another frame")
	}
}

// The emphasized preset must have a more visible rim against its body,
// in both appearances, without changing the default Regular material.
func TestVividRimContrast(t *testing.T) {
	for _, dark := range []bool{false, true} {
		contrast := func(style Style) int {
			th := ui.ShadcnTheme(dark)
			tt := ui.NewTester(func(c *ui.Context) {
				c.SetTheme(th)
				ui.Box(c).Fill().Background(ui.RGB(64, 64, 64)).Children(func() {
					ui.Box(c).Absolute().Left(20).Top(20).Size(200, 100).Radius(16).Material(Glass{Style: style})
				})
			}, 240, 140)
			im := tt.Image()
			return int(im.RGBAAt(120, 20).R) - int(im.RGBAAt(120, 70).R)
		}
		if regular, vivid := contrast(Regular), contrast(Vivid); vivid <= regular {
			t.Fatalf("rim contrast did not increase (dark=%v): regular=%d vivid=%d", dark, regular, vivid)
		}
	}
}
