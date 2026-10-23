// Author: HyperparameterSI

package gallerydemo

import (
	"github.com/HyperparameterSI/golet/backend/ui"
	"testing"
)

func TestThemeShowcasePreservesForm(t *testing.T) {
	g := &gallery{}
	var dark bool
	tt := ui.NewTester(func(c *ui.Context) {
		g.applyTheme(c)
		dark = c.Theme().Dark
		ui.Column(c).Fill().Padding(16).Gap(12).Children(func() { g.themes(c) })
	}, 900, 900)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	click := func(s string) {
		t.Helper()
		if err := tt.Click(s); err != nil {
			t.Fatal(err)
		}
		tt.Frame()
	}
	click("工作空间名称")
	tt.Type("Golet Design")
	click("Liquid Glass")
	click("深色")
	if !g.liquidGlass || !dark || g.themeName != "Golet Design" {
		t.Fatalf("theme switch lost state: glass=%v dark=%v name=%q", g.liquidGlass, dark, g.themeName)
	}
	click("空间类型")
	click("团队空间")
	click("保存设置")
	if g.themePlan != "团队空间" || g.themeNotice != "已保存当前演示设置" || !tt.HasText("操作结果") {
		t.Fatalf("form interaction failed: plan=%q notice=%q texts=%q", g.themePlan, g.themeNotice, tt.Texts())
	}
	click("Maia 中性")
	click("浅色")
	if g.liquidGlass || dark || g.themePlan != "团队空间" {
		t.Fatal("returning to neutral lost state")
	}
	click("重置")
	if g.themeName != "" || g.themePlan != "" {
		t.Fatal("reset did not clear form")
	}
}
