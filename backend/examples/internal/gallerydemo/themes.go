// Author: HyperparameterSI

package gallerydemo

import (
	"github.com/HyperparameterSI/golet/backend/examples/internal/demotheme"
	"github.com/HyperparameterSI/golet/backend/plugins/glass"
	"github.com/HyperparameterSI/golet/backend/ui"
)

// publishTheme shares a user's selection with other native demo windows.
func (g *gallery) publishTheme() {
	if g.sharedTheme {
		demotheme.Set(g.liquidGlass, g.appearance)
	}
}

// applyTheme resolves appearance before controls are built without resetting
// form state. Standalone test previews keep their own cached presets.
func (g *gallery) applyTheme(c *ui.Context) {
	if g.sharedTheme {
		g.liquidGlass, g.appearance = demotheme.Selection()
		demotheme.Apply(c)
		return
	}
	dark := c.Theme().Dark
	if g.appearance == 1 {
		dark = false
	}
	if g.appearance == 2 {
		dark = true
	}
	if g.neutralLight == nil {
		g.neutralLight, g.neutralDark = ui.FormaTheme(false), ui.FormaTheme(true)
	}
	theme := g.neutralLight
	if dark {
		theme = g.neutralDark
	}
	if g.liquidGlass {
		if g.glassLight == nil {
			g.glassLight, g.glassDark = glass.Theme(false), glass.Theme(true)
		}
		theme = g.glassLight
		if dark {
			theme = g.glassDark
		}
	}
	c.SetTheme(theme)
}

func (g *gallery) themes(c *ui.Context) {
	t := c.Theme()
	ui.Text(c, "选择外观，体验同一套原生控件。").TextColor(t.TextMuted)
	ui.Row(c).Gap(12).Wrap().AlignItems(ui.Center).Children(func() {
		style := 0
		if g.liquidGlass {
			style = 1
		}
		if ui.Segmented(c, &style, "Maia 中性", "Liquid Glass").Label("控件风格").Changed() {
			g.liquidGlass = style == 1
			g.publishTheme()
		}
		if ui.Segmented(c, &g.appearance, "跟随系统", "浅色", "深色").Label("主题外观").Changed() {
			g.publishTheme()
		}
		ui.Checkbox(c, &g.themePlain, "纯色背景")
	})
	a, b := ui.Hex("#e6edf3"), ui.Hex("#eeebf1")
	if t.Dark {
		a, b = ui.Hex("#1c2733"), ui.Hex("#2a2732")
	}
	stage := ui.Column(c).Padding(20).Gap(16).Radius(16).Clip()
	if g.themePlain || c.Preferences().HighContrast {
		stage.Background(t.Background).Border(1, t.Border)
	} else {
		stage.Gradient(a, b, 135)
		if g.liquidGlass {
			// Static wallpaper shapes reveal the lens along each glass edge.
			// No background animation or additional timer is required.
			stage.Draw(func(p *ui.Painter, r ui.Rect) {
				p.Clip(r, 16, func() {
					p.Fill(ui.Rect{X: r.X + r.W*0.40, Y: r.Y - r.H*0.15, W: r.W * 0.9, H: r.H * 0.85}, ui.Hex("#679bcf").Alpha(0.32), r.W*0.45)
					p.Fill(ui.Rect{X: r.X - r.W*0.25, Y: r.Y + r.H*0.45, W: r.W * 0.95, H: r.H}, ui.Hex("#91b6c6").Alpha(0.24), r.W*0.48)
				})
			})
		}
	}
	stage.Children(func() {
		ui.Row(c).AlignItems(ui.Center).Children(func() {
			ui.Text(c, "工作空间").FontSize(20).FontWeight(600)
			ui.Spacer(c)
			style := "Maia"
			if g.liquidGlass {
				style = "Glass"
			}
			ui.Text(c, style).FontSize(12).TextColor(t.TextMuted)
		})
		ui.Tabs(c, &g.themeTab, "账户", "通知", "关于").Label("预览页签")
		panel := ui.Column(c).Padding(20).Gap(20).Radius(12).Background(t.Surface).Border(1, t.Border)
		if g.liquidGlass && !c.Preferences().HighContrast {
			panel.Material(t.PanelMaterial).
				Border(1, ui.RGBA(255, 255, 255, 0.16)).
				Shadow(0, 8, 24, -6, ui.RGBA(10, 24, 44, 0.16))
		}
		panel.Children(func() {
			switch g.themeTab {
			case 0:
				ui.Text(c, "基本信息").FontSize(14).FontWeight(600)
				ui.Column(c).Gap(6).Children(func() {
					ui.Text(c, "名称").FontSize(12).TextColor(t.TextMuted)
					ui.TextInput(c, &g.themeName).Placeholder("例如：设计工作室").Label("工作空间名称")
				})
				ui.Column(c).Gap(6).Children(func() {
					ui.Text(c, "类型").FontSize(12).TextColor(t.TextMuted)
					ui.Select(c, &g.themePlan, []string{"个人空间", "团队空间", "共享空间"}).Label("空间类型")
				})
				ui.Checkbox(c, &g.themeChecked, "同步工作空间偏好")
			case 1:
				ui.Text(c, "通知偏好").FontSize(14).FontWeight(600)
				ui.Row(c).Gap(12).AlignItems(ui.Center).Children(func() {
					ui.Switch(c, &g.themeEnabled).Label("桌面通知")
					ui.Text(c, "开启桌面通知")
				})
				ui.Text(c, "这里只演示控件状态，不会发送系统通知。").TextColor(t.TextMuted)
			case 2:
				ui.Text(c, "由 Golet 原生绘制").FontSize(14).FontWeight(600)
				ui.Text(c, "玻璃采样当前窗口内的背景。高对比度使用实色，减少动画遵循系统偏好。").TextColor(t.TextMuted)
			}
			ui.Divider(c)
			ui.Row(c).Gap(8).Wrap().Children(func() {
				if ui.PrimaryButton(c, "保存设置").Clicked() {
					g.themeNotice = "已保存当前演示设置"
				}
				if ui.Button(c, "重置").Clicked() {
					g.themeName, g.themePlan, g.themeNotice = "", "", "已重置演示表单"
					g.themeEnabled, g.themeChecked = false, false
				}
			})
			if g.themeNotice != "" {
				ui.Text(c, g.themeNotice).TextColor(t.TextMuted).Label("操作结果")
			}
		})
	})
}
