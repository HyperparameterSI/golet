// Author: HyperparameterSI
package gallerydemo

import (
	"testing"

	"github.com/HyperparameterSI/golet/backend/examples/internal/demotheme"
	"github.com/HyperparameterSI/golet/backend/ui"
)

func TestSharedThemeAcrossGalleryWindows(t *testing.T) {
	oldGlass, oldMode := demotheme.Selection()
	t.Cleanup(func() { demotheme.Set(oldGlass, oldMode) })
	demotheme.Set(false, 0)
	first, second := &gallery{sharedTheme: true}, &gallery{sharedTheme: true}
	var one, two *ui.Theme
	a := ui.NewTester(func(c *ui.Context) { first.applyTheme(c); one = c.Theme(); ui.Text(c, "One") }, 200, 100)
	b := ui.NewTester(func(c *ui.Context) { second.applyTheme(c); two = c.Theme(); ui.Text(c, "Two") }, 200, 100)
	first.liquidGlass, first.appearance = true, 2
	first.publishTheme()
	a.Frame()
	b.Frame()
	if one != two || !two.Dark || two.ControlMaterial == nil || !second.liquidGlass {
		t.Fatal("open galleries did not share glass appearance")
	}
	var next *ui.Theme
	ui.NewTester(func(c *ui.Context) { demotheme.Apply(c); next = c.Theme(); ui.Text(c, "New window") }, 200, 100)
	if next != two {
		t.Fatal("new demo did not inherit appearance")
	}
	second.liquidGlass, second.appearance = false, 1
	second.publishTheme()
	a.Frame()
	b.Frame()
	if one != two || two.Dark || two.ControlMaterial != nil || first.liquidGlass {
		t.Fatal("return to neutral did not synchronize")
	}
}
