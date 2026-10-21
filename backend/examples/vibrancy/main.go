// Author: HyperparameterSI

// Vibrancy shows a window in the style of a native macOS app, in native UI:
// a sidebar over a translucent material, under a hidden title bar whose
// window controls sit over the view, the traffic lights inset over the
// sidebar on macOS. Pick a material in the sidebar to see it behind the
// window.
//
//	go run ./examples/vibrancy
package main

import (
	"log"
	"time"

	"github.com/HyperparameterSI/golet/backend"
	"github.com/HyperparameterSI/golet/backend/ui"
)

// material is a Vibrancy the sidebar offers: its constant in package golet,
// its name, and what it is for.
type material struct {
	value           golet.Vibrancy
	constant, label string
	about           string
}

var sections = []struct {
	title string
	icon  *ui.SVG
	items []material
}{
	{"Window areas", windowIcon, []material{
		{golet.VibrancySidebar, "VibrancySidebar", "Sidebar", "The background of sidebars, like this one."},
		{golet.VibrancyTitlebar, "VibrancyTitlebar", "Title bar", "The background of title bars and toolbars."},
		{golet.VibrancyHeader, "VibrancyHeader", "Header", "Headers and footers inside content."},
		{golet.VibrancyWindow, "VibrancyWindow", "Window", "The background of opaque windows."},
		{golet.VibrancyContent, "VibrancyContent", "Content", "The background of opaque content."},
		{golet.VibrancyUnderWindow, "VibrancyUnderWindow", "Under window", "What shows under the background of a window."},
		{golet.VibrancyUnderPage, "VibrancyUnderPage", "Under page", "The area behind the pages of a document."},
	}},
	{"Transient UI", transientIcon, []material{
		{golet.VibrancyMenu, "VibrancyMenu", "Menu", "The background of menus."},
		{golet.VibrancyPopover, "VibrancyPopover", "Popover", "The background of popovers."},
		{golet.VibrancyHUD, "VibrancyHUD", "HUD", "The background of heads-up displays."},
		{golet.VibrancySheet, "VibrancySheet", "Sheet", "The background of sheets."},
		{golet.VibrancyTooltip, "VibrancyTooltip", "Tooltip", "The background of tooltips."},
		{golet.VibrancySelection, "VibrancySelection", "Selection", "The highlight of selected content."},
		{golet.VibrancyFullScreenUI, "VibrancyFullScreenUI", "Full screen UI", "The background of full screen modal interfaces."},
	}},
	{"Windows 11", backdropIcon, []material{
		{golet.VibrancyMica, "VibrancyMica", "Mica", "Tints the window with the desktop wallpaper. Under window on macOS."},
		{golet.VibrancyAcrylic, "VibrancyAcrylic", "Acrylic", "Blurs what is behind the window. HUD on macOS."},
		{golet.VibrancyTabbed, "VibrancyTabbed", "Tabbed", "Mica for windows with tabs in the title bar. Under window on macOS."},
	}},
}

func icon(shapes string) *ui.SVG {
	return ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.3" stroke-linejoin="round">` + shapes + `</svg>`))
}

var (
	windowIcon    = icon(`<rect x="1.5" y="2.5" width="13" height="11" rx="2"/><path d="M1.5 5.5h13"/>`)
	transientIcon = icon(`<path d="M3 2.5h10A1.5 1.5 0 0 1 14.5 4v6a1.5 1.5 0 0 1-1.5 1.5H7.5L4.5 14v-2.5H3A1.5 1.5 0 0 1 1.5 10V4A1.5 1.5 0 0 1 3 2.5z"/>`)
	backdropIcon  = icon(`<path d="M8 1.5l1.7 4.8 4.8 1.7-4.8 1.7L8 14.5l-1.7-4.8L1.5 8l4.8-1.7z"/>`)
	sidebarIcon   = icon(`<rect x="1.5" y="2.5" width="13" height="11" rx="2"/><path d="M6 2.5v11"/>`)
)

const (
	titleBarHeight = 52 // the height of an inset title bar
	sidebarWidth   = 220
)

// app is the state the window shows.
type app struct {
	win        *golet.Window
	chosen     string // the Vibrancy chosen in the sidebar
	collapsed  bool   // the sidebar is hidden
	seeThrough bool   // the material shows behind the pane too
}

func (a *app) find() material {
	for _, s := range sections {
		for _, m := range s.items {
			if string(m.value) == a.chosen {
				return m
			}
		}
	}
	return sections[0].items[0]
}

func (a *app) view(c *ui.Context) {
	t := c.Theme()
	// The material shows wherever the view draws no background, where the
	// window shows one: the root draws none then, nor does the sidebar.
	vibrant := c.Vibrancy()
	if vibrant {
		c.Root().Background(ui.Transparent)
	}
	bar := c.TitleBar()
	row := ui.Row(c).Fill().AlignItems(ui.Stretch)
	shown := row.Animate("sidebar", b2f(!a.collapsed), 250*time.Millisecond)
	row.Children(func() {
		// The sidebar slides out to the left as it hides.
		w := sidebarWidth * shown
		ui.Row(c).Width(w).Shrink(0).AlignItems(ui.Stretch).Clip().Children(func() {
			side := ui.Column(c).Width(sidebarWidth).Shrink(0).Margin(0, 0, 0, w-sidebarWidth)
			if !vibrant {
				side.Background(t.Surface)
			}
			side.Children(func() { a.sidebar(c, vibrant) })
		})
		if w > 0 {
			ui.Box(c).Width(1).Background(t.Border)
		}
		a.pane(c, vibrant, 20+(1-shown)*(bar.Left+36))
	})

	// Next to the window controls on the left, whether the sidebar shows or
	// not.
	label := "Hide Sidebar"
	if a.collapsed {
		label = "Show Sidebar"
	}
	toggle := ui.ButtonBase(c.Key("toggle")).Absolute().Left(bar.Left+8).Top((titleBarHeight-28)/2).Size(32, 28).
		Radius(6).Center().TextColor(t.TextMuted).Label(label).Tooltip(label)
	if toggle.Hovered() {
		toggle.Background(t.SurfaceHover)
	}
	toggle.Children(func() { ui.Icon(c, sidebarIcon).Size(16, 16) })
	if toggle.Clicked() {
		a.collapsed = !a.collapsed
	}
}

// sidebar is the materials to choose from, under the room the window
// controls take.
func (a *app) sidebar(c *ui.Context, vibrant bool) {
	ui.Box(c).Height(titleBarHeight).Shrink(0).DragWindow()
	list := ui.Sidebar(c, &a.chosen, func() {
		for _, s := range sections {
			ui.SidebarSection(c, s.title, nil, func() {
				for _, m := range s.items {
					ui.SidebarItem(c, string(m.value), s.icon, m.label).Tooltip(m.about)
				}
			})
		}
	}).Grow(1).Label("Materials")
	if vibrant {
		list.Background(ui.Transparent)
	}
	if list.Changed() && a.win != nil {
		a.win.SetVibrancy(golet.Vibrancy(a.chosen))
	}
}

// pane is the material chosen, and how to choose it in code; its title
// starts left of the window controls when the sidebar hides.
func (a *app) pane(c *ui.Context, vibrant bool, titleLeft float32) {
	t := c.Theme()
	m := a.find()
	bar := c.TitleBar()
	pane := ui.Column(c).Grow(1).MinWidth(0)
	if !a.seeThrough || !vibrant {
		pane.Background(t.Background)
	}
	pane.Children(func() {
		ui.Row(c).Height(titleBarHeight).Shrink(0).AlignItems(ui.Center).Padding(0, bar.Right+20, 0, titleLeft).DragWindow().Children(func() {
			ui.Text(c, m.label).FontSize(15).Bold().SingleLine()
		})
		ui.Divider(c)
		ui.Scroll(c).Grow(1).Children(func() {
			ui.Column(c).Padding(24, 28).Gap(16).MaxWidth(560).Children(func() {
				ui.Text(c, m.about).FontSize(15)
				ui.Text(c, "win.SetVibrancy(golet."+m.constant+")").Font("monospace").FontSize(12).
					Padding(12, 14).Radius(8).Background(t.Surface).Selectable()
				ui.Text(c, "A material shows wherever the view draws no background: the sidebar has none, this pane has one.")
				ui.Checkbox(c, &a.seeThrough, "Show the material behind this pane too").Disabled(!vibrant)
				note := "Windows 11 shows the closest of its backdrops: Acrylic for transient UI, Mica otherwise. Linux has no materials."
				if !vibrant {
					note = "This window shows no material here, so the sidebar draws a background of its own: " +
						"materials show on macOS, and on Windows 11 22H2 and later."
				}
				ui.Text(c, note).FontSize(12).TextColor(t.TextMuted).Margin(8, 0, 0, 0)
			})
		})
	})
}

func b2f(b bool) float32 {
	if b {
		return 1
	}
	return 0
}

func main() {
	a := &app{chosen: string(golet.VibrancySidebar)}
	golet.App.WhenReady(func() {
		a.win = golet.NewWindow(golet.WindowOptions{
			Title:     "Vibrancy",
			Width:     880,
			Height:    640,
			MinWidth:  560,
			MinHeight: 400,
			// Hide the title bar and inset the traffic lights over the
			// sidebar (macOS). Elsewhere the window controls sit in the
			// view's title bar, 52 pixels tall.
			TitleBarStyle:  golet.TitleBarHiddenInset,
			TitleBarHeight: titleBarHeight,
			// The material shows wherever the view draws no background
			// (macOS, Windows 11 22H2); ui.Context.Vibrancy tells where it
			// does.
			Vibrancy: golet.VibrancySidebar,
			Content:  ui.View(a.view),
		})
	})
	if err := golet.App.Run(); err != nil {
		log.Fatal(err)
	}
}
