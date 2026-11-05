// Author: HyperparameterSI

//go:build windows && (amd64 || arm64)

package e2e

import (
	"sync/atomic"
	"testing"

	"github.com/HyperparameterSI/golet/backend"
	win "github.com/HyperparameterSI/golet/backend/internal/windows"
	"github.com/HyperparameterSI/golet/backend/ui"
)

// The items Golet draws of a context menu are as wide as they need, at the
// window's DPI, with long labels, checks, shortcuts and submenus: Windows
// widens them by as much as it doesn't tell (#922). The first menu may be
// wider than it needs; the next, measured knowing that, are no wider.
func TestContextMenuItemWidths(t *testing.T) {
	var frames atomic.Int32
	view := func(c *ui.Context) {
		frames.Add(1)
		ui.Box(c).Fill().ContextMenu(func(m *ui.Menu) {
			m.Item("复制歌名和歌手")
			m.Separator()
			m.Item("查看歌手：房东的猫")
			m.Item("查看专辑：云烟成雨")
			m.Item("Show the artist: The Landlord's Cat").Checked(true)
			m.Item("Copy the song's name").Shortcut(ui.Ctrl|ui.Shift, ui.KeyC)
			m.Submenu("添加到歌单", func(m *ui.Menu) {
				m.Item("日常")
			})
		})
	}
	w := newWindow(t, golet.WindowOptions{Title: "Menu widths", Width: 400, Height: 300, Content: ui.View(view)})
	eventually(t, "a frame", func() bool { return frames.Load() > 0 })
	for _, pass := range []string{"the first menu", "the next"} {
		if !rightClick(w, 100, 50) {
			t.Fatalf("%s: right-clicking failed", pass)
		}
		eventually(t, pass+" to show", func() bool { menus, _ := popupMenus(); return len(menus) == 1 })
		// The menu shows; the next waits for Windows' widening learned
		// from its items drawn.
		var items []win.PopupItemWidth
		eventually(t, pass+" to draw", func() bool {
			var known bool
			golet.RunOnMain(func() { items, known = win.TestPopupItemWidths(w.NativeHandle()) })
			return known || len(items) == 0
		})
		closeMenus()
		eventually(t, pass+" to close", func() bool { menus, _ := popupMenus(); return len(menus) == 0 })
		if len(items) == 0 {
			t.Skip("Windows draws the items of menus here")
		}
		var needs int32
		for _, it := range items {
			if it.Shown < it.Needs {
				t.Errorf("%s: %q shows %d pixels wide, needs %d", pass, it.Label, it.Shown, it.Needs)
			}
			needs = max(needs, it.Needs)
		}
		if pass == "the next" && items[0].Shown != needs {
			t.Errorf("%s: the items show %d pixels wide, the widest needs %d", pass, items[0].Shown, needs)
		}
	}
}
