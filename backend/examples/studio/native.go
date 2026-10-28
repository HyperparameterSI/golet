// Author: HyperparameterSI
package main

import (
	"errors"
	"github.com/HyperparameterSI/golet/backend"
	"github.com/HyperparameterSI/golet/backend/examples/internal/demotheme"
	"github.com/HyperparameterSI/golet/backend/examples/internal/editordemo"
	"github.com/HyperparameterSI/golet/backend/examples/internal/gallerydemo"
	"github.com/HyperparameterSI/golet/backend/examples/internal/loadingdemo"
	"github.com/HyperparameterSI/golet/backend/plugins/terminal"
	"github.com/HyperparameterSI/golet/backend/ui"
)

// Showcase exposes the catalog and opens native demos in the Studio process.
// windows is accessed exclusively on the main UI thread.
type Showcase struct{ windows map[string]*golet.Window }

func (s *Showcase) Catalog() []Demo { return demos }
func (s *Showcase) Open(id string) error {
	valid := false
	for _, demo := range demos {
		if demo.ID == id && id != "Browser" {
			valid = true
			break
		}
	}
	if !valid {
		return errors.New("unknown native demo")
	}
	// Fetch optional terminal dependencies on the bound-method goroutine.
	if id == "Terminal" {
		if err := terminal.Load(); err != nil {
			return err
		}
	}
	var err error
	golet.RunOnMain(func() {
		if s.windows == nil {
			s.windows = map[string]*golet.Window{}
		}
		if win := s.windows[id]; win != nil {
			win.Show()
			win.Focus()
			return
		}
		var win *golet.Window
		switch id {
		case "Scene3D":
			win, err = openScene3D()
		case "Editor":
			win = editordemo.Open()
		case "Forma":
			win = openForma()
		case "BorderBeam":
			win = openBorderBeam()
		case "Loading":
			win = loadingdemo.Open()
		case "Advanced":
			win = openAdvanced()
		case "WebView":
			win, err = openWebView()
		case "Terminal":
			win, err = openTerminal()
		default:
			win = gallerydemo.Open(id)
		}
		if win != nil {
			demotheme.Register(win)
			s.windows[id] = win
			win.OnClosed(func() { delete(s.windows, id) })
		}
	})
	return err
}

func openWebView() (*golet.Window, error) {
	var view *golet.WebView
	zoom := 1.0
	win := golet.NewWindow(golet.WindowOptions{Title: "Studio · WebView", Width: 950, Height: 680, Content: ui.View(func(c *ui.Context) {
		demotheme.Apply(c)
		ui.Column(c).Fill().Children(func() {
			ui.Row(c).Padding(12).Gap(12).Children(func() {
				ui.Text(c, "原生工具栏 + 系统 WebView")
				if ui.Button(c, "加载示例").Clicked() && view != nil {
					view.Page().LoadHTML(webViewHTML, "")
				}
				if ui.Button(c, "放大").Clicked() && view != nil {
					zoom += 0.1
					view.Page().SetZoomFactor(zoom)
				}
				if ui.Button(c, "重置缩放").Clicked() && view != nil {
					zoom = 1
					view.Page().SetZoomFactor(zoom)
				}
			})
			if view != nil {
				ui.WebView(c, view).Grow(1)
			}
		})
	})})
	var err error
	view, err = win.NewWebView(golet.WebViewOptions{})
	if err != nil {
		win.Destroy()
		return nil, err
	}
	view.Page().LoadHTML(webViewHTML, "")
	win.Invalidate()
	return win, nil
}

const webViewHTML = `<!doctype html><html lang="zh"><meta charset="utf-8"><style>body{font:18px system-ui;padding:60px;color:#27483c;background:#f3f7f1}button{padding:14px 24px;border:0;border-radius:12px;background:#315f49;color:white;cursor:pointer}h1{font-size:44px}p{line-height:1.8}</style><h1>Web 与原生，共处一窗。</h1><p>上方是 Go 绘制的控件，下方由系统 WebView 渲染。无需下载 Chromium。</p><button onclick="this.textContent='已点击 '+(++window.count)+' 次'">试试页面交互</button><script>window.count=0</script></html>`

func openTerminal() (*golet.Window, error) {
	var win *golet.Window
	term, err := terminal.New(terminal.Options{OnTitle: func(title string) {
		golet.RunOnMain(func() {
			if win != nil {
				win.SetTitle(title)
			}
		})
	}, OnExit: func(int) {
		golet.RunOnMain(func() {
			if win != nil {
				win.Close()
			}
		})
	}})
	if err != nil {
		return nil, err
	}
	win = golet.NewWindow(golet.WindowOptions{Title: "Studio · Terminal", Width: 900, Height: 600, Content: ui.View(func(c *ui.Context) {
		demotheme.Apply(c)
		terminal.View(c, term).Fill().AutoFocus()
	})})
	win.OnClosed(func() { term.Close() })
	return win, nil
}
func openAdvanced() *golet.Window {
	input := ui.NewTextBuffer("共享缓冲区")
	area := ui.NewTextBuffer("这是 TextAreaBuffer。\n编辑文字，尝试撤销与选择。")
	color := ui.Hex("#527d66")
	var tree ui.OutlineState[string]
	mode := 0
	size := float32(200)
	ratio := float32(0.5)
	return golet.NewWindow(golet.WindowOptions{Title: "Studio · 缓冲区与分栏", Width: 1000, Height: 760, Content: ui.View(func(c *ui.Context) {
		demotheme.Apply(c)
		ui.Scroll(c).Fill().Padding(24).Gap(16).Children(func() {
			ui.Text(c, "缓冲区与布局组件").FontSize(26).Bold()
			ui.TextInputBuffer(c, input)
			ui.TextInputBuffer(c, input)
			ui.TextAreaBuffer(c, area).Height(130)
			ui.Divider(c)
			ui.ColorPicker(c, &color)
			ui.Outline(c, &tree, []string{"Golet"}, func(key string) []string {
				if key == "Golet" {
					return []string{"backend", "frontend"}
				}
				return nil
			}, func(key string) { ui.Text(c, key) }).Height(100)
			ui.ScrollHorizontal(c).Height(64).Children(func() {
				ui.Row(c).Gap(12).Children(func() {
					for _, label := range []string{"Row", "Column", "Grid", "Box", "Spacer", "Scroll", "ScrollBoth", "ScrollHorizontal", "Split", "Outline"} {
						ui.Button(c, label).Width(130)
					}
				})
			})
			ui.Segmented(c, &mode, "Split", "SplitPercent", "SplitVertical", "SplitVerticalPercent")
			left := func() {
				ui.Column(c).Fill().Padding(20).Background(c.Theme().Surface).Children(func() { ui.Text(c, "面板 A · 拖动分隔线") })
			}
			right := func() { ui.Column(c).Fill().Padding(20).Children(func() { ui.Text(c, "面板 B") }) }
			switch mode {
			case 0:
				ui.Split(c, &size, left, right).Height(220)
			case 1:
				ui.SplitPercent(c, &ratio, left, right).Height(220)
			case 2:
				ui.SplitVertical(c, &size, left, right).Height(300)
			case 3:
				ui.SplitVerticalPercent(c, &ratio, left, right).Height(300)
			}
		})
	})})
}
