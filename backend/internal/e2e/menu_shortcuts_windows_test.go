// Author: HyperparameterSI

//go:build windows

package e2e

import (
	"strings"
	"testing"
	"time"

	"github.com/HyperparameterSI/golet/backend"
	win "github.com/HyperparameterSI/golet/backend/internal/windows"
	"github.com/HyperparameterSI/golet/backend/transfer"
	"github.com/HyperparameterSI/golet/backend/ui"
)

// The shortcuts of the app's menu items work in a window of native UI, as
// in a page: before the content, which gets the rest. The roles' keys, as
// Edit › Copy's Ctrl+C, go to the content first, as a terminal's Ctrl+C.
func TestNativeMenuShortcuts(t *testing.T) {
	var settings, split, copied, typed int
	prev := golet.App.Menu()
	defer golet.App.SetMenu(prev)
	golet.App.SetMenu(golet.NewMenu([]*golet.MenuItem{
		{Label: "App", Submenu: []*golet.MenuItem{
			{Label: "Settings", Accelerator: "Ctrl+,", Click: func(*golet.MenuItem, *golet.Window) { settings++ }},
			{Label: "Split", Accelerator: "Alt+Shift+=", Click: func(*golet.MenuItem, *golet.Window) { split++ }},
		}},
		{Role: golet.RoleEditMenu},
	}))
	w := newWindow(t, golet.WindowOptions{Width: 300, Height: 200, Content: ui.View(func(c *ui.Context) {
		c.OnShortcut(ui.Ctrl, ui.KeyC, func() { copied++ })
		c.OnShortcut(ui.Ctrl, ui.KeyComma, func() { typed++ })
		ui.Text(c, "shortcuts")
	})})
	const vkControl, vkShift, vkMenu, vkComma, vkEqual = 0x11, 0x10, 0x12, 0xBC, 0xBB
	key := func(vk byte, mods ...byte) {
		t.Helper()
		var ok bool
		golet.RunOnMain(func() { ok = win.TestSurfaceKey(w.NativeHandle(), vk, mods...) })
		if !ok {
			t.Fatal("no window of native UI to press keys in")
		}
	}
	counts := func() (s, p, c, ty int) {
		golet.RunOnMain(func() { s, p, c, ty = settings, split, copied, typed })
		return
	}
	key(vkComma, vkControl)
	key(vkEqual, vkMenu, vkShift)
	key('C', vkControl)
	eventually(t, "the content's Ctrl+C", func() bool { _, _, c, _ := counts(); return c == 1 })
	if s, p, _, ty := counts(); s != 1 || p != 1 || ty != 0 {
		t.Errorf("Ctrl+, ran Settings %d times and reached the content %d, Alt+Shift+= ran Split %d times; want once, never, once", s, ty, p)
	}
}

// In a page, the edit roles' keys are the page's, as in a browser and as in
// native UI: WebView2 pastes with the page's own paste event and every
// format of the clipboard, where Edit › Paste would type its text alone.
func TestPageEditKeys(t *testing.T) {
	prev := golet.App.Menu()
	defer golet.App.SetMenu(prev)
	golet.App.SetMenu(golet.NewMenu([]*golet.MenuItem{{Role: golet.RoleEditMenu}}))
	w := newWindow(t, golet.WindowOptions{Width: 400, Height: 300})
	w.Page().LoadHTML(`<textarea id=t></textarea><script>
addEventListener("paste", (e) => { window.pasted = e.clipboardData.types.join(" ") })
</script>`, "")
	waitFor(t, w, "document.getElementById('t') !== null")
	if err := golet.Clipboard.Write(transfer.New(transfer.NewItem(
		transfer.Bytes(transfer.Text, []byte("Golet")), transfer.Bytes(transfer.HTML, []byte("<b>Golet</b>"))))); err != nil {
		t.Fatal(err)
	}
	var front bool
	golet.RunOnMain(func() { front = win.TestForeground(w.NativeHandle()) })
	if !front {
		t.Fatal("the window did not come to the front")
	}
	if _, err := w.Page().Eval("document.getElementById('t').focus()"); err != nil {
		t.Fatal(err)
	}
	golet.RunOnMain(func() { win.TestPressKeys(0x11, 'V') }) // VK_CONTROL
	waitFor(t, w, "window.pasted !== undefined")
	if types, _ := golet.EvalAs[string](w.Page(), "window.pasted"); !strings.Contains(types, "text/html") {
		t.Errorf("the page's paste had %q, want text/html among them", types)
	}
}

// Ctrl+V in a text field pastes once: the field takes the key, and Edit ›
// Paste, whose key it is too, leaves it.
func TestNativeFieldPastesOnce(t *testing.T) {
	prev := golet.App.Menu()
	defer golet.App.SetMenu(prev)
	golet.App.SetMenu(golet.NewMenu([]*golet.MenuItem{{Role: golet.RoleEditMenu}}))
	golet.RunOnMain(func() { golet.Clipboard.WriteText("host.ts.net") })
	var value string
	w := newWindow(t, golet.WindowOptions{Width: 300, Height: 200, Content: ui.View(func(c *ui.Context) {
		ui.TextInput(c, &value).AutoFocus()
	})})
	read := func() (v string) { golet.RunOnMain(func() { v = value }); return }
	time.Sleep(500 * time.Millisecond) // the field shows and takes the focus
	const vkControl = 0x11
	golet.RunOnMain(func() { win.TestSurfaceKey(w.NativeHandle(), 'V', vkControl) })
	eventually(t, "the paste", func() bool { return read() != "" })
	time.Sleep(300 * time.Millisecond)
	if got := read(); got != "host.ts.net" {
		t.Errorf("after Ctrl+V the field has %q, want %q", got, "host.ts.net")
	}
}
