// Author: HyperparameterSI
package chromium

import (
	"math"
	"testing"
)

func TestThemeValidation(t *testing.T) {
	base := Theme{Background: "ffffffff", Surface: "ffffffff", Muted: "eeeeeeff", Foreground: "000000ff", MutedForeground: "777777ff", Border: "00000020", Primary: "2563ebff", Radius: 16, FontSize: 14}
	if _, err := base.encode(); err != nil {
		t.Fatal(err)
	}
	for _, edit := range []func(*Theme){
		func(t *Theme) { t.Border = "oops" }, func(t *Theme) { t.Primary = "zzzzzzzz" },
		func(t *Theme) { t.Radius = -1 }, func(t *Theme) { t.FontSize = 30 }, func(t *Theme) { t.Radius = float32(math.NaN()) },
	} {
		v := base
		edit(&v)
		if _, err := v.encode(); err == nil {
			t.Fatal("invalid theme accepted")
		}
	}
}
