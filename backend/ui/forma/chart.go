// Author: HyperparameterSI
package forma

import (
	"fmt"
	"math"
	"strings"

	"github.com/HyperparameterSI/golet/backend/ui"
	"github.com/HyperparameterSI/golet/backend/ui/chart"
)

// ChartKind chooses native marks. Other geometries remain available through ChartContainer.
type ChartKind int

const (
	ChartBar ChartKind = iota
	ChartArea
	ChartLine
)

// ChartSeries contains categorical values aligned with ChartOptions.Labels.
// Missing values and non-finite numbers are gaps. A zero Color uses theme tokens.
type ChartSeries struct {
	Label  string
	Values []float64
	Color  ui.Color
}

// ChartOptions configures a responsive, themed chart. Height defaults to 260 DIPs.
// Curve is shared with ui/chart (linear, monotone and step interpolation).
type ChartOptions struct {
	Label                           string
	Kind                            ChartKind
	Labels                          []string
	Series                          []ChartSeries
	Height                          float32
	Curve                           chart.Curve
	HideLegend, HideGrid, ShowYAxis bool
	FormatValue                     func(float64) string
	OnSelect                        func(int)
}

// ChartState belongs to the view. Initialize Selected to -1 for no selection.
type ChartState struct{ Selected int }

func chartValue(s ChartSeries, i int) (float64, bool) {
	if i >= len(s.Values) {
		return 0, false
	}
	v := s.Values[i]
	return v, !math.IsNaN(v) && !math.IsInf(v, 0)
}
func chartDomain(o ChartOptions) chart.Linear {
	lo, hi := 0.0, 0.0
	for _, s := range o.Series {
		for i := range o.Labels {
			if v, ok := chartValue(s, i); ok {
				lo = min(lo, v)
				hi = max(hi, v)
			}
		}
	}
	if lo == hi {
		hi = lo + 1
	}
	return chart.NewLinear(lo, hi).Nice(4)
}

// Chart draws native marks and exposes each category through pointer and
// keyboard selection. It performs no polling and does not require a WebView.
func Chart(c *ui.Context, state *ChartState, o ChartOptions) ui.Element {
	height := o.Height
	if height <= 0 {
		height = 260
	}
	format := o.FormatValue
	if format == nil {
		format = func(v float64) string { return fmt.Sprintf("%g", v) }
	}
	colors := ChartColors(c)
	seriesColors := make([]ui.Color, len(o.Series))
	names := make([]string, len(o.Series))
	for j, s := range o.Series {
		names[j] = s.Label
		seriesColors[j] = s.Color
		if s.Color == (ui.Color{}) {
			seriesColors[j] = colors[j%len(colors)]
		}
	}
	n := len(o.Labels)
	return ui.Column(c).FillWidth().MinWidth(0).Gap(16).Children(func() {
		if n == 0 || len(o.Series) == 0 {
			state.Selected = -1
			ui.Text(c, "No chart data").TextColor(c.Theme().TextMuted)
			return
		}
		y := chartDomain(o)
		y.Format = format
		xs := make([]float32, n)
		for i := range xs {
			if o.Kind == ChartBar {
				xs[i] = (float32(i) + .5) / float32(n)
			} else if n == 1 {
				xs[i] = .5
			} else {
				xs[i] = float32(i) / float32(n-1)
			}
		}
		ticks := make([]chart.Tick, 0, min(n, 7))
		stride := max(1, (n+6)/7)
		for i := 0; i < n; i += stride {
			ticks = append(ticks, chart.Tick{At: xs[i], Label: o.Labels[i]})
		}
		axes := chart.Axes{Bottom: ticks}
		if o.ShowYAxis {
			axes.Left = y.Ticks(4)
		}
		describe := func(i int) string {
			var b strings.Builder
			b.WriteString(o.Labels[i])
			for _, s := range o.Series {
				if v, ok := chartValue(s, i); ok {
					fmt.Fprintf(&b, " · %s: %s", s.Label, format(v))
				}
			}
			return b.String()
		}
		chart.Plot(c, axes, func(area ui.Element) {
			label := o.Label
			if label == "" {
				label = "Chart plot"
			}
			area.Label(label).Description("Arrow keys select data; Escape clears selection")
			if chart.Select(c, area, &state.Selected, xs, describe) && o.OnSelect != nil {
				o.OnSelect(state.Selected)
			}
			if state.Selected < -1 {
				state.Selected = -1
			}
			border := c.Theme().Border
			area.Draw(func(p *ui.Painter, r ui.Rect) {
				if !o.HideGrid {
					for _, tick := range y.Ticks(4) {
						p.Line(r.X, chart.Y(r, tick.At), r.X+r.W, chart.Y(r, tick.At), 1, border)
					}
				}
			})
			for j, s := range o.Series {
				color := seriesColors[j]
				if o.Kind == ChartBar {
					group := float32(.8) / float32(n)
					width := group / float32(len(o.Series))
					for i := range o.Labels {
						if v, ok := chartValue(s, i); ok {
							start := xs[i] - group/2 + float32(j)*width
							chart.Rect(c, chart.Span{Lo: start + width*.04, Hi: start + width*.96}, y.Span(0, v)).Background(color).Radius(4).Label(fmt.Sprintf("%s · %s: %s", o.Labels[i], s.Label, format(v)))
						}
					}
				} else {
					ys := make([]float32, n)
					baseline := make([]float32, n)
					for i := range ys {
						baseline[i] = y.At(0)
						if v, ok := chartValue(s, i); ok {
							ys[i] = y.At(v)
						} else {
							ys[i] = float32(math.NaN())
						}
					}
					chart.Rect(c, chart.Span{Lo: 0, Hi: 1}, chart.Span{Lo: 0, Hi: 1}).Draw(func(p *ui.Painter, r ui.Rect) {
						if o.Kind == ChartArea {
							p.FillPathGradient(chart.Area(r, xs, baseline, ys, o.Curve), ui.LinearGradient{From: color.Alpha(.45), To: color.Alpha(.04), Angle: 180})
						}
						p.StrokePath(chart.Line(r, xs, ys, o.Curve), 1.5, color)
					})
				}
			}
			if i := state.Selected; i >= 0 && i < n {
				chart.RuleX(c, xs[i], 1).Background(c.Theme().TextMuted)
				for j, s := range o.Series {
					if v, ok := chartValue(s, i); ok {
						chart.Point(c, xs[i], y.At(v), 7, 7).Radius(4).Background(seriesColors[j]).Border(1, c.Theme().Surface)
					}
				}
				align := ui.Start
				if xs[i] > .5 {
					align = ui.End
				}
				chart.Anchor(c, xs[i], 1, align, ui.Start, func() {
					ui.Column(c).Padding(12).Gap(6).Radius(radius(c, 12)).Background(c.Theme().Surface).Border(1, border).Children(func() {
						ui.Text(c, o.Labels[i]).Bold().NoWrap()
						for j, s := range o.Series {
							if v, ok := chartValue(s, i); ok {
								ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
									ui.Box(c).Size(8, 8).Radius(2).Background(seriesColors[j])
									ui.Text(c, s.Label).NoWrap()
									ui.Text(c, format(v)).Bold().NoWrap()
								})
							}
						}
					}).Label("Chart tooltip")
				})
			}
		}).Height(height).Gap(10).MarginX(24).FontSize(fontSize(c, 12)).TextColor(c.Theme().TextMuted)
		if !o.HideLegend {
			ChartLegend(c, names, seriesColors).FillWidth()
		}
	})
}
