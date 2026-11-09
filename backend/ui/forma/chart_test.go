// Author: HyperparameterSI
package forma

import (
	"github.com/HyperparameterSI/golet/backend/ui"
	"math"
	"testing"
)

func TestChartDomainAndMissingValues(t *testing.T) {
	o := ChartOptions{Labels: []string{"a", "b", "c"}, Series: []ChartSeries{{Values: []float64{-5, math.Inf(1), 10}}, {Values: []float64{math.NaN()}}}}
	d := chartDomain(o)
	if d.Lo > -5 || d.Hi < 10 {
		t.Fatalf("domain: %+v", d)
	}
	for _, i := range []int{1, 3} {
		if _, ok := chartValue(o.Series[0], i); ok {
			t.Fatal("invalid value accepted")
		}
	}
	d = chartDomain(ChartOptions{})
	if d.Hi <= d.Lo {
		t.Fatal("degenerate empty domain")
	}
}
func TestChartKeyboardSelection(t *testing.T) {
	state := ChartState{Selected: -1}
	selected := -1
	tt := ui.NewTester(func(c *ui.Context) {
		Chart(c, &state, ChartOptions{Labels: []string{"Jan", "Feb"}, Series: []ChartSeries{{Label: "Desktop", Values: []float64{12, 18}}}, OnSelect: func(i int) { selected = i }})
	}, 600, 400)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnd)
	if state.Selected != 1 || selected != 1 || !tt.HasText("18") {
		t.Fatalf("selection=%d callback=%d texts=%v", state.Selected, selected, tt.Texts())
	}
	tt.Key(0, ui.KeyEscape)
	if state.Selected != -1 {
		t.Fatal("selection not cleared")
	}
}
func TestChartEmptyAndShortSeries(t *testing.T) {
	for _, kind := range []ChartKind{ChartBar, ChartArea, ChartLine} {
		s := ChartState{Selected: 99}
		tt := ui.NewTester(func(c *ui.Context) {
			Chart(c, &s, ChartOptions{Kind: kind, Labels: []string{"a", "b"}, Series: []ChartSeries{{Values: []float64{2}}}})
		}, 600, 400)
		tt.Frame()
		if s.Selected != -1 {
			t.Fatal("out-of-range selection retained")
		}
	}
}
