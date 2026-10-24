// Author: HyperparameterSI
package main

import (
	"time"

	"github.com/HyperparameterSI/golet/backend"
	"github.com/HyperparameterSI/golet/backend/examples/internal/demotheme"
	"github.com/HyperparameterSI/golet/backend/plugins/borderbeam"
	"github.com/HyperparameterSI/golet/backend/ui"
)

func openBorderBeam() *golet.Window {
	return golet.NewWindow(golet.WindowOptions{Title: "Studio · Border Beam", Width: 960, Height: 820, Content: ui.View(newBorderBeamView())})
}

// newBorderBeamView keeps demo state separate from the window for deterministic UI tests.
func newBorderBeamView() func(*ui.Context) {
	paused, disabled := false, false
	preset := borderbeam.Medium
	speed := float64(1.96)
	strength, glow := float64(1), float64(28)
	return func(c *ui.Context) {
		demotheme.Apply(c)
		t := c.Theme()
		ui.Scroll(c).Fill().Padding(28).Gap(20).Children(func() {
			ui.Text(c, "Border Beam").FontSize(28).Bold()
			ui.Text(c, "五种原生动画 · 参考色场与遮罩 · 全局同步调节")
			names := []string{"Medium", "Small", "Line", "Pulse Inner", "Pulse Outside"}
			status := "播放中 · " + names[preset]

			if c.Preferences().ReduceMotion {
				status = "静态展示 · 系统已开启减少动态效果"
			}
			if paused {
				status = "已暂停 · 取消勾选暂停以继续"
			}
			if disabled {
				status = "效果已隐藏 · 取消勾选隐藏效果以显示"
			}
			ui.Text(c, status).FontSize(13).TextColor(t.TextMuted)
			ui.Row(c).Gap(16).Children(func() {
				ui.Checkbox(c, &paused, "暂停")
				ui.Checkbox(c, &disabled, "隐藏效果")
			})
			ui.Row(c).Gap(8).Children(func() {
				for i, name := range []string{"Medium", "Small", "Line", "Pulse Inner", "Pulse Outside"} {
					ui.Button(c, name).OnClick(func() {
						preset = borderbeam.Preset(i)
						speed = 1.96
						if preset == borderbeam.Line {
							speed = 3.1
						}
						if preset >= borderbeam.PulseInner {
							speed = 2.3
						}
					})
				}
			})
			ui.Grid(c).Columns(3).Gap(20).Children(func() {
				ui.Column(c).Gap(8).Children(func() { ui.Text(c, "周期（秒）"); ui.Slider(c, &speed, 1, 8) })
				ui.Column(c).Gap(8).Children(func() { ui.Text(c, "光效强度"); ui.Slider(c, &strength, .1, 1) })
				ui.Column(c).Gap(8).Children(func() { ui.Text(c, "内辉光范围"); ui.Slider(c, &glow, 8, 48) })
			})
			o := borderbeam.Options{Strength: float32(strength), Glow: float32(glow), Paused: paused, Disabled: disabled, Preset: preset, Duration: time.Duration(float64(speed) * float64(time.Second))}
			borderbeam.Box(c, o).Padding(28).Gap(16).Background(t.Surface).Children(func() {
				ui.Text(c, "Give your ideas an edge.").FontSize(24).Bold()
				ui.Text(c, "默认配色跟随 Studio 主题。切换 Maia / Liquid Glass 或明暗外观后再次体验。")
				label := "暂停动画"
				if paused {
					label = "继续播放"
				}
				ui.Button(c, label).OnClick(func() { paused = !paused })
			})
			ui.Text(c, "COLOR STUDIES / 彩色光效").FontSize(12).TextColor(t.TextMuted)
			ui.Grid(c).Columns(2).Gap(16).Children(func() {
				for _, sample := range []struct {
					name, description, accent string
					palette                   borderbeam.Palette
					colors                    []ui.Color
				}{
					{"Prism · 多彩棱镜", "红、蓝、绿、紫交织，展示局部流光。", "#a78bfa", borderbeam.Colorful, nil},
					{"Ocean · 海洋蓝", "青色与深蓝渐隐，适合工具与工作台。", "#38bdf8", borderbeam.Ocean, nil},
					{"Sunset · 落日余晖", "珊瑚红、暖橙与粉紫，柔和而明亮。", "#fb923c", borderbeam.Sunset, nil},
					{"Aurora · 极光", "自定义薄荷绿、冰蓝和紫色光斑。", "#34d399", borderbeam.Theme, []ui.Color{ui.Hex("#34d399"), ui.Hex("#67e8f9"), ui.Hex("#a78bfa"), ui.Hex("#22d3ee")}},
				} {
					sampleOptions := o
					sampleOptions.Palette = sample.palette
					sampleOptions.Colors = sample.colors
					accent := ui.Hex(sample.accent)
					tint := float32(.035)
					if t.Dark {
						tint = .07
					}
					borderbeam.Box(c, sampleOptions).Label(sample.name + " card").Padding(24).Gap(12).Background(t.Surface.Mix(accent, tint)).Children(func() {
						ui.Row(c).Gap(6).Children(func() {
							for _, color := range []ui.Color{accent, accent.Mix(ui.RGB(255, 255, 255), .4), accent.Mix(t.Surface, .4)} {
								ui.Box(c).Size(7, 7).Radius(4).Background(color)
							}
						})
						ui.Text(c, sample.name).FontSize(19).Bold()
						ui.Text(c, sample.description).TextColor(t.TextMuted)
						ui.Text(c, "NATIVE COLOR FIELD").FontSize(10).TextColor(accent.Mix(t.Text, .3))
					})
				}
			})
			ui.Text(c, "系统开启减少动态效果时自动静止；不可见时不持续请求动画帧。")
		})
	}
}
