// Author: HyperparameterSI
package browser

import (
	"github.com/HyperparameterSI/golet/backend/ui"
	"github.com/HyperparameterSI/golet/backend/ui/forma"
	"testing"
)

func TestFormaThemeSnapshot(t *testing.T) {
	for _, preset := range forma.Presets() {
		for _, dark := range []bool{false, true} {
			source := forma.NewTheme(forma.ThemeOptions{Preset: preset, Dark: dark, Tokens: map[string]ui.Color{"border": ui.Hex("#12345678")}})
			result := Theme(source)
			if result.Dark != dark || result.Border != "12345678" || result.FontSize != 14 || result.Radius != 16 {
				t.Fatalf("%s dark=%v: %+v", preset, dark, result)
			}
			source.Border = ui.RGB(0, 0, 0)
			if result.Border != "12345678" {
				t.Fatal("theme was not a snapshot")
			}
		}
	}
	if Theme(nil).FontSize < 10 {
		t.Fatal("invalid default font")
	}
}
