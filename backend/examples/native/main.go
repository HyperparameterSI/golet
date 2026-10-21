// Author: HyperparameterSI

// Native tours the desktop integration: application and context menus, a
// tray icon, dialogs, notifications, a global shortcut, the clipboard and
// dark mode.
//
//	go run ./examples/native
package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"strings"
	"time"

	"github.com/HyperparameterSI/golet/backend"
)

// Native is called from the page.
type Native struct{}

// OpenFiles lets the user pick files and returns their paths.
func (Native) OpenFiles(ctx context.Context) ([]string, error) {
	return golet.Dialog.Open(golet.OpenDialogOptions{
		Parent:   golet.CallerWindow(ctx),
		Title:    "Pick some files",
		Multiple: true,
	})
}

// Confirm asks a yes/no question.
func (Native) Confirm(ctx context.Context, question string) (bool, error) {
	res, err := golet.Dialog.Message(golet.MessageOptions{
		Parent:  golet.CallerWindow(ctx),
		Type:    golet.MessageQuestion,
		Message: question,
		Detail:  "This dialog is native.",
		Buttons: []string{"Yes", "No"},
	})
	return err == nil && res.Button == 0, err
}

// Copy puts text on the clipboard.
func (Native) Copy(text string) { golet.Clipboard.WriteText(text) }

// Paste reads text from the clipboard.
func (Native) Paste() string { return golet.Clipboard.ReadText() }

// Notify shows a desktop notification.
func (Native) Notify(title, body string) error {
	if !golet.NotificationsSupported() {
		return fmt.Errorf("notifications need a packaged app (golet dev or golet build)")
	}
	return golet.NewNotification(golet.NotificationOptions{Title: title, Body: body}).Show()
}

// ShowContextMenu pops up a native context menu at the mouse position.
func (Native) ShowContextMenu(ctx context.Context) {
	win := golet.CallerWindow(ctx)
	golet.NewMenu([]*golet.MenuItem{
		{Label: "Reload", Click: func(*golet.MenuItem, *golet.Window) { win.Page().Reload() }},
		{Label: "Toggle Developer Tools", Click: func(*golet.MenuItem, *golet.Window) { win.Page().ToggleDevTools() }},
		golet.Separator(),
		{Role: golet.RoleCopy},
		{Role: golet.RolePaste},
	}).Popup(win)
}

// SetDark switches between the light and dark appearance.
func (Native) SetDark(dark bool) {
	if dark {
		golet.Theme.SetSource(golet.ThemeDark)
	} else {
		golet.Theme.SetSource(golet.ThemeLight)
	}
}

// Log is sent to the page for anything worth showing.
var Log = golet.NewEvent[string]("log")

const page = `<!doctype html>
<html>
<head>
  <title>Native</title>
  <style>
    :root { color-scheme: light dark; font-family: system-ui, sans-serif; }
    body { margin: 24px; }
    button { font: inherit; margin: 0 6px 8px 0; padding: 6px 12px; }
    pre { background: rgba(128,128,128,.12); padding: 12px; border-radius: 8px; min-height: 120px; white-space: pre-wrap; }
  </style>
</head>
<body>
  <h2>Native integration</h2>
  <p>Right-click anywhere for a native context menu. Try the tray icon and ⌘⇧Space / Ctrl+Shift+Space.</p>
  <button id="open">Open files…</button>
  <button id="confirm">Ask a question</button>
  <button id="copy">Copy text</button>
  <button id="paste">Paste text</button>
  <button id="notify">Notify</button>
  <button id="dark">Toggle dark mode</button>
  <pre id="log"></pre>
  <script>
    const log = (s) => (document.getElementById("log").textContent += s + "\n");
    const call = (m, ...a) => golet.call("Native." + m, ...a).catch((e) => log("error: " + e.message));
    golet.on("log", log);
    let dark = matchMedia("(prefers-color-scheme: dark)").matches;
    document.getElementById("open").onclick = async () => log("picked: " + JSON.stringify(await call("OpenFiles")));
    document.getElementById("confirm").onclick = async () => log("answer: " + await call("Confirm", "Do you like Golet?"));
    document.getElementById("copy").onclick = () => call("Copy", "Copied from Golet at " + new Date().toLocaleTimeString()).then(() => log("copied"));
    document.getElementById("paste").onclick = async () => log("clipboard: " + await call("Paste"));
    document.getElementById("notify").onclick = () => call("Notify", "Hello", "From Golet");
    document.getElementById("dark").onclick = () => call("SetDark", (dark = !dark));
    addEventListener("contextmenu", (e) => { e.preventDefault(); call("ShowContextMenu"); });
  </script>
</body>
</html>`

// trayIcon draws a 32x32 template image: a filled circle.
func trayIcon() []byte {
	img := image.NewNRGBA(image.Rect(0, 0, 32, 32))
	for y := range 32 {
		for x := range 32 {
			if dx, dy := x-16, y-16; dx*dx+dy*dy <= 12*12 {
				img.Set(x, y, color.Black)
			}
		}
	}
	var b bytes.Buffer
	_ = png.Encode(&b, img)
	return b.Bytes()
}

func main() {
	golet.Bind(Native{})
	app := golet.App
	app.WhenReady(func() {
		var win *golet.Window
		show := func() {
			if win == nil || win.IsDestroyed() {
				win = golet.NewWindow(golet.WindowOptions{Title: "Native", Width: 720, Height: 520})
				win.Page().LoadHTML(page, "")
			}
			win.Show()
		}
		logf := func(format string, args ...any) {
			if win != nil {
				_ = Log.Emit(win, fmt.Sprintf(format, args...))
			}
		}

		app.SetMenu(golet.NewMenu([]*golet.MenuItem{
			{Role: golet.RoleAppMenu},
			{Role: golet.RoleEditMenu},
			{Label: "Demo", Submenu: []*golet.MenuItem{
				{Label: "Say Hi", Accelerator: "CmdOrCtrl+Shift+H", Click: func(*golet.MenuItem, *golet.Window) { logf("Hi from the menu!") }},
				{ID: "always-on-top", Label: "Always on Top", Type: golet.MenuItemCheckbox, Click: func(item *golet.MenuItem, w *golet.Window) {
					if w != nil {
						w.SetAlwaysOnTop(item.IsChecked())
					}
				}},
				golet.Separator(),
				{Label: "Small", Type: golet.MenuItemRadio, Checked: true, Click: func(_ *golet.MenuItem, w *golet.Window) { w.SetSize(720, 520) }},
				{Label: "Large", Type: golet.MenuItemRadio, Click: func(_ *golet.MenuItem, w *golet.Window) { w.SetSize(1000, 720) }},
			}},
			{Role: golet.RoleViewMenu},
			{Role: golet.RoleWindowMenu},
		}))

		tray, err := golet.NewTray(golet.TrayOptions{
			Icon:           trayIcon(),
			IconIsTemplate: true,
			ToolTip:        "Golet Native",
			Menu: golet.NewMenu([]*golet.MenuItem{
				{Label: "Show Window", Click: func(*golet.MenuItem, *golet.Window) { show() }},
				{Label: "What time is it?", Click: func(*golet.MenuItem, *golet.Window) { logf("It's %s", time.Now().Format(time.Kitchen)) }},
				golet.Separator(),
				{Role: golet.RoleQuit},
			}),
		})
		if err != nil {
			log.Println("tray:", err)
		} else {
			_ = tray
		}

		if err := golet.GlobalShortcut.Register("CmdOrCtrl+Shift+Space", func() {
			show()
			logf("global shortcut pressed")
		}); err != nil {
			log.Println("shortcut:", err)
		}

		golet.Theme.OnUpdated(func() {
			logf("appearance is now %s", map[bool]string{true: "dark", false: "light"}[golet.Theme.IsDark()])
		})
		show()
		d := golet.Screen.PrimaryDisplay()
		logf("primary display: %s %dx%d @%gx", strings.TrimSpace(d.Label), d.Bounds.Width, d.Bounds.Height, d.ScaleFactor)
	})
	// Keep running in the tray after the window is closed.
	app.OnWindowAllClosed(func() {})
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
