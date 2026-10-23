// Author: HyperparameterSI

package ui

import (
	"math"
	"testing"
)

// The toolbar's overflow button is chrome: when every control fits, the
// button is hidden and must not widen the toolbar. Regression test for
// https://github.com/egoist/mygo/issues/491
func TestToolbarEndsAtLastControl(t *testing.T) {
	tt := NewTester(func(c *Context) {
		Row(c).Width(400).Children(func() {
			Text(c, "Changes")
			Spacer(c)
			Toolbar(c, func() {
				Button(c, "One")
				Button(c, "Two").Label("Last")
			}).Padding(0).Label("Tools")
		})
	}, 400, 100)
	bar, _ := tt.Find("Tools")
	last, _ := tt.Find("Last")
	// Fractional theme spacing can round the two sums differently in
	// float32; allow subpixel arithmetic error, not extra layout space.
	if end := last.X + last.W; math.Abs(float64(end-(bar.X+bar.W))) > 1.0/1024 {
		t.Fatalf("last button ends at %v, toolbar ends at %v", end, bar.X+bar.W)
	}
}
