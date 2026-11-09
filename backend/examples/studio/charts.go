// Author: HyperparameterSI
package main

import (
	"math"
	"time"

	"github.com/HyperparameterSI/golet/backend/ui"
	"github.com/HyperparameterSI/golet/backend/ui/chart"
	"github.com/HyperparameterSI/golet/backend/ui/forma"
)

func newChartDemo() func(*ui.Context) {
	bars, area := forma.ChartState{Selected: -1}, forma.ChartState{Selected: -1}
	period, curveName := "Last 3 months", "Smooth"
	labels := make([]string, 90)
	desktop, mobile := make([]float64, 90), make([]float64, 90)
	start := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)
	for i := range labels {
		labels[i] = start.AddDate(0, 0, i).Format("Jan 2")
		desktop[i] = 80 + 170*math.Pow(math.Sin(float64(i)*1.13), 2) + 60*math.Pow(math.Sin(float64(i)*.32), 2)
		mobile[i] = 25 + 105*math.Pow(math.Sin(float64(i)*1.13+.7), 2)
	}
	return func(c *ui.Context) {
		ui.Column(c).FillWidth().Gap(24).Children(func() {
			chartDemoCard(c, "Visitors by month", "Desktop and mobile · January–June", func() {
				forma.Chart(c, &bars, forma.ChartOptions{Label: "Monthly visitors", Kind: forma.ChartBar, Labels: []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"}, Series: []forma.ChartSeries{{Label: "Desktop", Values: []float64{186, 305, 237, 73, 209, 214}, Color: ui.Hex("#2563eb")}, {Label: "Mobile", Values: []float64{80, 200, 120, 190, 130, 140}, Color: ui.Hex("#60a5fa")}}, Height: 240})
			})
			chartDemoCard(c, "Total Visitors", "示例数据 · 切换范围查看访问趋势", func() {
				ui.Row(c).Gap(8).Children(func() {
					for _, value := range []string{"Last 3 months", "Last 30 days", "Last 7 days"} {
						variant := forma.Outline
						if period == value {
							variant = forma.Secondary
						}
						if forma.Button(c, value, forma.ButtonOptions{Variant: variant, Size: forma.Small}).Clicked() {
							period = value
							area.Selected = -1
						}
					}
				})
				forma.Select(c, &curveName, []string{"Smooth", "Linear", "Step"}).Width(160).Label("Chart interpolation")
				from := 0
				if period == "Last 30 days" {
					from = 60
				}
				if period == "Last 7 days" {
					from = 83
				}
				curve := chart.CurveMonotone
				if curveName == "Linear" {
					curve = chart.CurveLinear
				}
				if curveName == "Step" {
					curve = chart.CurveStep
				}
				forma.Chart(c, &area, forma.ChartOptions{Label: "Visitor trends", Kind: forma.ChartArea, Labels: labels[from:], Series: []forma.ChartSeries{{Label: "Desktop", Values: desktop[from:]}, {Label: "Mobile", Values: mobile[from:]}}, Curve: curve, Height: 260})
			})
			ui.Text(c, "悬停查看数值 · Tab 聚焦图表后使用方向键 · Esc 清除选择").TextColor(c.Theme().TextMuted)
		})
	}
}

func chartDemoCard(c *ui.Context, title, description string, build func()) {
	forma.Card(c).FillWidth().Children(func() {
		forma.CardHeader(c).Children(func() { forma.CardTitle(c, title); forma.CardDescription(c, description) })
		forma.CardContent(c).Children(build)
	})
}
