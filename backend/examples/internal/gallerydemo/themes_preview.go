// Author: HyperparameterSI

package gallerydemo

import "github.com/HyperparameterSI/golet/backend/ui"

// precisionCard demonstrates exact decimal values without converting them to floats.
func (g *gallery) precisionCard(c *ui.Context) {
	card(c, "DecimalInput · 精确小数", func() {
		host := ui.Column(c).Gap(8)
		amount := ui.Local(host, "amount", func() string { return "19.99" })
		host.Children(func() {
			ui.Text(c, "以字符串保存金额；步长 0.01，超过两位小数时提示错误，不自动舍入。").TextColor(c.Theme().TextMuted)
			ui.DecimalInput(c, amount, "0", "999999999999999999.99", "0.01").Label("Exact amount")
			ui.Text(c, "当前值："+*amount)
		})
	})
}

// themePreviewCard keeps preview themes local to the subtree and reuses them across frames.
func (g *gallery) themePreviewCard(c *ui.Context) {
	card(c, "Themed · 局部主题", func() {
		row := ui.Row(c).Gap(12).Wrap()
		themes := ui.Local(row, "themes", func() [2]*ui.Theme { return [2]*ui.Theme{ui.LightTheme(), ui.DarkTheme()} })
		row.Children(func() {
			for i, label := range []string{"浅色预览", "深色预览"} {
				ui.Themed(c.Key(label), themes[i]).Width(260).Padding(16).Gap(12).Radius(12).Children(func() {
					ui.Text(c, label).Bold()
					ui.Text(c, "这个区域拥有独立主题，外部控件继续使用窗口主题。").TextColor(c.Theme().TextMuted)
					if ui.Button(c, "预览按钮").Clicked() {
						c.Toast(label)
					}
				})
			}
		})
	})
}
