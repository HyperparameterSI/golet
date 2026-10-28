// Author: HyperparameterSI

//go:build windows && (amd64 || arm64)

package e2e

import (
	"log"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/HyperparameterSI/golet/backend"
	win "github.com/HyperparameterSI/golet/backend/internal/windows"
)

func activateMenu(w *golet.Window, path ...string) (err error) {
	golet.RunOnMain(func() { err = win.TestActivateMenuItem(w.NativeHandle(), path...) })
	return err
}

// Keyboard, dialog and popup automation are not wired on Windows.
func pressShortcut(string, bool) (bool, bool) { return false, false }

func endSheet(*golet.Window) (bool, bool) { return false, false }

// recentDocumentRecorded reports whether path is in the Recent category of
// the app's jump list.
func recentDocumentRecorded(path string) (recorded, supported bool) {
	var paths []string
	golet.RunOnMain(func() { paths = win.TestRecentDocuments() })
	return slices.ContainsFunc(paths, func(p string) bool { return strings.EqualFold(p, path) }), true
}

// Clicks reach windows showing native UI only: WebView2 takes the mouse in
// windows of its own process.
func click(w *golet.Window, x, y float64) (ok bool) {
	golet.RunOnMain(func() { ok = win.TestClickSurface(w.NativeHandle(), x, y) })
	return ok
}

// Clicks reach windows showing native UI only: a web view's page is
// not clicked.
func clickPage(*golet.Window, float64, float64) bool { return false }

// appClick clicks with the system's input, which activates windows, and
// returns once the app handled the click.
func appClick(w *golet.Window, x, y float64) (ok bool) {
	var missed string
	golet.RunOnMain(func() { ok, missed = win.TestInputClick(w.NativeHandle(), x, y) })
	if missed != "" {
		log.Printf("e2e: the click at (%v, %v) did not reach the app; the window there: %s", x, y, missed)
	}
	return ok
}

// heldAppClick is appClick: the owner of a focusable flyout eats the
// click that closes it, however long it is held.
func heldAppClick(w *golet.Window, x, y float64) bool { return appClick(w, x, y) }

func sideButton(w *golet.Window, back bool) (ok bool) {
	golet.RunOnMain(func() { ok = win.TestSideButton(w.NativeHandle(), back) })
	return ok
}

// pinch pinches the touchpad at (x, y) as Windows sends it, scrolling
// with Ctrl held, which rounds scale.
func pinch(w *golet.Window, x, y, scale float64) (got float64, ok bool) {
	golet.RunOnMain(func() { got, ok = win.TestPinchSurface(w.NativeHandle(), x, y, scale) })
	return got, ok
}

func drag(w *golet.Window, points [][2]float64) (ok bool) {
	golet.RunOnMain(func() { ok = win.TestDragSurface(w.NativeHandle(), points) })
	return ok
}

func webViewAttached(*golet.Window) (bool, bool) { return false, false }

// Collection behaviors and deactivating the app are macOS's.
func collectionBehavior(*golet.Window) (uint, bool)               { return 0, false }
func deactivateApp() bool                                         { return false }
func material(*golet.Window) (string, int, float64, bool, bool)   { return "", 0, 0, false, false }
func roundedBacking(*golet.Window) (bool, float64, float64, bool) { return false, 0, 0, false }
func appActive() (bool, bool)                                     { return false, false }

// WebView2 opens DevTools in a window of its own.
func dockDevTools(*golet.Window) bool { return false }

func dockedDevToolsPlace(*golet.Window) string { return "" }

func dismissPopups() (int, bool) { return 0, false }

func setDroppedFiles(w *golet.Window, paths []string) bool {
	golet.RunOnMain(func() { win.TestSetDroppedFiles(w.NativeHandle(), paths) })
	return true
}

// The progress bar is shown by the shell, out of the app's reach.
func dockTileImage() ([]byte, bool) { return nil, false }

func defaultURLHandler(string) (string, bool) { return "", false }

func dockMenu(int) ([]string, bool) { return nil, false }

func pressCtrlShiftK() bool {
	golet.RunOnMain(func() { win.TestPressKeys(0x11, 0x10, 'K') }) // VK_CONTROL, VK_SHIFT
	return true
}

// Only Linux binds global shortcuts through a desktop portal.
func usePortalShortcuts() (func(), bool) { return nil, false }

func pressKeys(...string) bool { return false }

// Only macOS windows have a toolbar.
func fullScreenHidesToolbar(*golet.Window) (bool, bool) { return false, false }

// Only macOS windows have traffic lights.
func trafficLights(*golet.Window) (float64, float64, bool) { return 0, 0, false }
func setWindowAppearance(*golet.Window, bool)              {}

// Frameless windows keep native resize borders here.
func movePointer(int, int) bool                 { return false }
func pressButton(bool) bool                     { return false }
func resizeCursor(*golet.Window) (string, bool) { return "", false }

func menuBarShown(w *golet.Window) (shown, supported bool) {
	golet.RunOnMain(func() { shown = win.TestMenuBarShown(w.NativeHandle()) })
	return shown, true
}

func activateAccelerator(w *golet.Window, accel string) (handled, supported bool) {
	golet.RunOnMain(func() { handled = win.TestActivateAccelerator(w.NativeHandle(), accel) })
	return handled, true
}

// enterMenuBar sends the window the SC_KEYMENU that Alt and F10 become, and
// leaves the menus from inside the menu loop, which runs the app's work.
func enterMenuBar(w *golet.Window, _ string) (during, after, supported bool) {
	hwnd := w.NativeHandle()
	inLoop := make(chan bool, 1)
	go func() {
		time.Sleep(300 * time.Millisecond)
		golet.RunOnMain(func() {
			inLoop <- win.TestMenuBarShown(hwnd)
			win.TestEndMenu()
		})
	}()
	golet.RunOnMain(func() { supported = win.TestEnterMenuBar(hwnd) })
	if !supported {
		return false, false, false
	}
	during = <-inLoop
	golet.RunOnMain(func() { after = win.TestMenuBarShown(hwnd) })
	return during, after, true
}

// openMenus opens the window's menus from the keyboard: Alt alone, or Alt
// and letter (a virtual key, such as 'A'). Their loop runs the app's work
// until closeMenus.
func openMenus(w *golet.Window, letter byte) (supported bool) {
	hwnd := w.NativeHandle()
	golet.RunOnMain(func() { supported = win.TestForeground(hwnd) })
	if !supported {
		return false
	}
	if letter == 0 {
		go golet.RunOnMain(func() { win.TestEnterMenuBar(hwnd) }) // returns once the menus close
	} else {
		golet.RunOnMain(func() { win.TestPressKeys(0x12, letter) }) // VK_MENU
	}
	return true
}

// closeMenus leaves the menus, as Escape does.
func closeMenus() { golet.RunOnMain(win.TestEndMenu) }

// titleButtons returns the window controls of a hidden title bar, named
// after the hit-test codes over them.
func titleButtons(w *golet.Window) (names []string, supported bool) {
	golet.RunOnMain(func() { names = win.TestCaptionButtons(w.NativeHandle()) })
	return names, true
}

func pressTitleButton(w *golet.Window, name string) (ok bool) {
	golet.RunOnMain(func() { ok = win.TestPressCaptionButton(w.NativeHandle(), name) })
	return ok
}

// menuBarColors returns the colors of a window's menu bar and of the line
// below it, as 0xBBGGRR.
func menuBarColors(w *golet.Window) (bar, line uint32, supported bool) {
	golet.RunOnMain(func() { bar, line, supported = win.TestMenuBarColors(w.NativeHandle()) })
	return bar, line, supported
}

// topNonClient returns how many pixels at the top of a window are not its
// page's, and how many a hidden title bar keeps there.
func topNonClient(w *golet.Window) (px, want int32, supported bool) {
	golet.RunOnMain(func() { px, want = win.TestTopNonClient(w.NativeHandle()) })
	return px, want, true
}

// titleButtonColor puts the pointer over a window control of a hidden
// title bar and returns the color on the screen beside its glyph.
func titleButtonColor(w *golet.Window, name string) (r, g, b uint8, ok bool) {
	golet.RunOnMain(func() { r, g, b, ok = win.TestCaptionColor(w.NativeHandle(), name) })
	return r, g, b, ok
}

// composition reports whether a window has no redirection bitmap, for the
// material behind its page, and whether the controls of its hidden title
// bar show through DirectComposition.
func composition(w *golet.Window) (noRedirect, composed, supported bool) {
	golet.RunOnMain(func() { noRedirect, composed = win.TestComposed(w.NativeHandle()) })
	return noRedirect, composed, true
}

// canCompose reports whether DirectComposition is available.
func canCompose() (ok bool) {
	golet.RunOnMain(func() { ok = win.TestCanCompose() })
	return ok
}

// loseComposition takes the device of DirectComposition away, as a draw
// that finds it invalid does, or, with removed, as a removed GPU does.
func loseComposition(removed bool) (ok bool) {
	golet.RunOnMain(func() { ok = win.TestLoseComposition(removed) })
	return ok
}

// failWebViews makes WebView2 fail to create the next n webviews.
func failWebViews(n int) bool {
	golet.RunOnMain(func() { win.TestFailWebViews(n) })
	return true
}

// A Control-click is a secondary click on macOS only.
func controlClick(*golet.Window, float64, float64) bool { return false }

// composeOver does what an input method that converts typed text again
// does through IMM32: it settles the range it reconverts, then composes or
// commits.
func composeOver(w *golet.Window, text string, caret int, commit bool, from, length int) (ok bool) {
	golet.RunOnMain(func() { ok = win.TestComposeOver(w.NativeHandle(), text, caret, commit, from, length) })
	return ok
}

// inputClient returns the selection and the text around it that input
// methods get with IMR_DOCUMENTFEED, which holds no composition.
func inputClient(w *golet.Window) (selected [2]int, document string, ok bool) {
	golet.RunOnMain(func() {
		var start, length int
		document, start, length, ok = win.TestDocumentFeed(w.NativeHandle())
		selected = [2]int{start, length}
	})
	return selected, document, ok
}

func dropFiles(w *golet.Window, x, y float64, paths []string) (over, dropped, ok bool) {
	golet.RunOnMain(func() { over, dropped = win.TestDropFiles(w.NativeHandle(), x, y, paths) })
	return over, dropped, true
}

// The control types of elements in UI Automation.
const roleText, roleButton, roleCheckBox, roleTextField, roleSlider = "Text", "Button", "CheckBox", "Edit", "Slider"

// roleListItem is the role of the rows of a List.
const roleListItem = "ListItem"

func accessibility(w *golet.Window) (nodes []accessNode, ok bool) {
	golet.RunOnMain(func() {
		var list []win.TestAccessNode
		list, ok = win.TestAccessibility(w.NativeHandle())
		for _, n := range list {
			nodes = append(nodes, accessNode{n.Role, n.Label, n.Value})
		}
	})
	return nodes, ok
}

func accessPerform(w *golet.Window, label, action, value string) (ok bool) {
	golet.RunOnMain(func() { ok = win.TestAccessibilityPerform(w.NativeHandle(), label, action, value) })
	return ok
}

// Typing into native UI is only automated on macOS.
func clickAndType(*golet.Window, float64, float64, string) bool { return false }
func compose(*golet.Window, string, int, bool) bool             { return false }

// Only Linux draws native UI in a GtkGLArea, nor waits to load the GPU's
// driver.
func glSurface(*golet.Window) (string, []byte, int, int, bool) { return "", nil, 0, 0, false }
func memoryUI(*testing.T) bool                                 { return false }
func surfaceOnScreen(*golet.Window) ([]byte, bool)             { return nil, false }

// screenColor returns the color the screen shows at (x, y), in DIPs, in a
// window showing native UI, or false when the screen cannot be read.
func screenColor(w *golet.Window, x, y float64) (r, g, b uint8, ok bool) {
	golet.RunOnMain(func() { r, g, b, ok = win.TestSurfacePixel(w.NativeHandle(), x, y) })
	return r, g, b, ok
}

// Only macOS has key-value observing.
func observe(*golet.Window) (func(), bool) { return nil, false }

// rightClick clicks (x, y) in a window showing native UI with the
// secondary button, with the messages a mouse sends.
func rightClick(w *golet.Window, x, y float64) (ok bool) {
	golet.RunOnMain(func() { ok = win.TestRightClickSurface(w.NativeHandle(), x, y) })
	return ok
}

// popupMenus returns the labels of the items of the context menus shown.
// The modal loop of a menu runs the main thread's work.
func popupMenus() (menus [][]string, supported bool) {
	golet.RunOnMain(func() { menus = win.TestPopups() })
	return menus, true
}

func choosePopupItem(label string) (ok bool) {
	golet.RunOnMain(func() { ok = win.TestChoosePopupItem(label) })
	return ok
}

// Input methods keep the keys typed while they compose from the content
// themselves here (GTK's filtering, IMM32's VK_PROCESSKEY), which a
// composition the test makes up does not show.
func pressKey(*golet.Window, uint16, string) bool { return false }

// dragFiles drags files from another application to (x, y) of a window of
// native UI and drops them, as OLE does, waiting for the target to answer.
func dragFiles(w *golet.Window, x, y float64, paths []string) (dropped, ok bool) {
	golet.RunOnMain(func() { dropped = win.TestDragFiles(w.NativeHandle(), x, y, paths) })
	return dropped, true
}

// pressTab presses Tab, with Shift when back, on the keyboard, in the
// window brought to the front.
func pressTab(w *golet.Window, back bool) (ok bool) {
	golet.RunOnMain(func() {
		if ok = win.TestForeground(w.NativeHandle()); !ok {
			return
		}
		if back {
			win.TestPressKeys(0x10, 0x09) // VK_SHIFT, VK_TAB
		} else {
			win.TestPressKeys(0x09)
		}
	})
	return ok
}

// AppKit's older accessibility API is macOS's.
func axAttribute(*golet.Window, string, string) (string, bool, bool, bool) {
	return "", false, false, false
}

// pressEscape presses Escape, as the keyboard does.
func pressEscape(*golet.Window) bool {
	golet.RunOnMain(func() { win.TestPressKeys(0x1B) }) // VK_ESCAPE
	return true
}

// Windows shows no document but the title.
func documentChrome(w *golet.Window) (title, file string, icon, edited, supported bool) {
	golet.RunOnMain(func() { title = win.TestWindowTitle(w.NativeHandle()) })
	return title, "", false, false, true
}

func sheet(*golet.Window) ([]string, []string, bool, bool) { return nil, nil, false, false }
func endSheetButton(*golet.Window, int) bool               { return false }

func panelStyles(w *golet.Window) (tool, topmost, owned, supported bool) {
	golet.RunOnMain(func() { tool, topmost, owned = win.TestPanelStyles(w.NativeHandle()) })
	return tool, topmost, owned, true
}

// throttlingOff reports whether the page of w, or of its first web view,
// is set up to keep running while it can't be seen, and how many of the
// settings that take the platform has.
func throttlingOff(w *golet.Window, webView bool) (off bool, settings int) {
	golet.RunOnMain(func() { off, settings = win.TestBackgroundThrottlingOff(w.NativeHandle(), webView) })
	return off, settings
}

// WebView2 tells no capture.
func mockCapture(*golet.Window) bool { return false }

// pageGestures reports whether a swipe navigates and a pinch zooms the page
// of a window, or with webView 0 or more, of its web view of that index.
func pageGestures(w *golet.Window, webView int) (swipe, pinch, ok bool) {
	golet.RunOnMain(func() { swipe, pinch, ok = win.TestPageGestures(w.NativeHandle(), webView) })
	return swipe, pinch, ok
}

// WebView2 takes touch input from Windows' pointer messages, which a test
// can't send it.
func pinchPage(*golet.Window, int, float64, float64, ...float64) bool { return false }

// pageMenusHooked tells that the tests can right-click a page and read
// the menu it opens, through rightClickPage, pageMenu, choosePageMenuItem
// and dismissPageMenu. A menu WebView2 draws itself can't be read: the
// tests change theirs, which Golet draws.
const pageMenusHooked = true

// rightClickPage right-clicks (x, y) in a window's page, with the
// DevTools protocol's mouse events.
func rightClickPage(w *golet.Window, x, y float64) (ok bool) {
	golet.RunOnMain(func() { ok = win.TestRightClickPage(w.NativeHandle(), x, y) })
	return ok
}

// pageMenu returns the labels of the items of the context menu Golet
// draws over a page, "-" for separators, and whether one is shown. The
// modal loop of a menu runs the main thread's work.
func pageMenu() (labels []string, shown bool) {
	golet.RunOnMain(func() { labels, shown = win.TestPageMenu() })
	return labels, shown
}

// choosePageMenuItem chooses an item of that menu by its label.
func choosePageMenuItem(label string) (ok bool) {
	golet.RunOnMain(func() { ok = win.TestChoosePageMenuItem(label) })
	return ok
}

func dismissPageMenu() { golet.RunOnMain(win.TestEndMenu) }

// takeForeground brings a window to the foreground, after Alt pressed and
// released when Windows refuses it: a process that sent the last input may
// take the foreground.
func takeForeground(w *golet.Window) {
	golet.RunOnMain(func() {
		if !win.TestForeground(w.NativeHandle()) {
			win.TestPressKeys(0x12) // VK_MENU
			win.TestForeground(w.NativeHandle())
		}
	})
}
