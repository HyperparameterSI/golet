// Author: HyperparameterSI

//go:build darwin

package e2e

import (
	"log"
	"os/exec"
	"slices"
	"testing"

	"github.com/HyperparameterSI/golet/backend"
	"github.com/HyperparameterSI/golet/backend/internal/darwin"
	"github.com/HyperparameterSI/golet/backend/internal/wait"
)

// recentDocumentRecorded reports whether path is among the app's recent
// documents.
func recentDocumentRecorded(path string) (recorded, supported bool) {
	var paths []string
	golet.RunOnMain(func() { paths = darwin.TestRecentDocuments() })
	return slices.Contains(paths, path), true
}

func activateMenu(w *golet.Window, path ...string) (err error) {
	golet.RunOnMain(func() { err = darwin.TestPerformMenuItem(path...) })
	return err
}

func pressShortcut(key string, shift bool) (handled bool, supported bool) {
	golet.RunOnMain(func() { handled = darwin.TestPerformKeyEquivalent(key, shift) })
	return handled, true
}

func endSheet(w *golet.Window) (ok bool, supported bool) {
	golet.RunOnMain(func() { ok = darwin.TestEndSheet(w.NativeHandle()) })
	return ok, true
}

func click(w *golet.Window, x, y float64) bool {
	waitOnScreen(w)
	golet.RunOnMain(func() { darwin.TestClick(w.NativeHandle(), x, y) })
	return true
}

// appClick clicks through the application, whose event monitors see it.
func appClick(w *golet.Window, x, y float64) bool {
	waitOnScreen(w)
	golet.RunOnMain(func() { darwin.TestAppClick(w.NativeHandle(), x, y) })
	return true
}

// heldAppClick is appClick with the button held down, as a user clicks:
// what the press did runs before the release, the work it left for the
// main thread too.
func heldAppClick(w *golet.Window, x, y float64) bool {
	waitOnScreen(w)
	golet.RunOnMain(func() { darwin.TestAppPress(w.NativeHandle(), x, y, true) })
	golet.RunOnMain(func() {}) // after what the press posted
	golet.RunOnMain(func() { darwin.TestAppPress(w.NativeHandle(), x, y, false) })
	return true
}

// waitOnScreen waits for a window to show on screen, as a user sees a
// window before clicking it: a window shows a while after it is ordered
// front and made key, longer on a busy runner, and AppKit puts a click
// made for it before then elsewhere in it, as outside its content, where
// it clicks nothing. wait.Guard is a hang guard: it says so if the window
// did not show by then, before the click that will miss.
func waitOnScreen(w *golet.Window) {
	if !wait.Poll(func() (on bool) {
		golet.RunOnMain(func() { on = darwin.TestOnScreen(w.NativeHandle()) })
		return on
	}) {
		log.Printf("e2e: the window %q did not show on screen in %v: a click made for it misses", w.Title(), wait.Guard)
	}
}

func sideButton(w *golet.Window, back bool) (ok bool) {
	golet.RunOnMain(func() { ok = darwin.TestSideButton(w.NativeHandle(), back) })
	return ok
}

func drag(w *golet.Window, points [][2]float64) bool {
	golet.RunOnMain(func() { darwin.TestDrag(w.NativeHandle(), points) })
	return true
}

func webViewAttached(w *golet.Window) (attached bool, supported bool) {
	golet.RunOnMain(func() { attached = darwin.TestWebViewAttached(w.NativeHandle()) })
	return attached, true
}

// deactivateApp activates the Finder, as the user would.
func deactivateApp() bool {
	return exec.Command("open", "-a", "Finder").Run() == nil
}

func appActive() (active bool, supported bool) {
	golet.RunOnMain(func() { active = darwin.TestIsFrontmost() })
	return active, true
}

func collectionBehavior(w *golet.Window) (behavior uint, supported bool) {
	golet.RunOnMain(func() { behavior = darwin.TestCollectionBehavior(w.NativeHandle()) })
	return behavior, true
}

func dockDevTools(w *golet.Window) (ok bool) {
	golet.RunOnMain(func() { ok = darwin.TestDockInspector(w.NativeHandle()) })
	return ok
}

func dockedDevToolsPlace(w *golet.Window) (place string) {
	golet.RunOnMain(func() { place = darwin.TestInspectorPlace(w.NativeHandle()) })
	return place
}

// Context menus track the mouse in a modal loop; not automated.
func dismissPopups() (int, bool) { return 0, false }

func setDroppedFiles(w *golet.Window, paths []string) bool {
	golet.RunOnMain(func() { darwin.TestSetDroppedFiles(w.NativeHandle(), paths) })
	return true
}

func dockTileImage() (png []byte, supported bool) {
	golet.RunOnMain(func() { png = darwin.TestDockTileImage() })
	return png, true
}

func defaultURLHandler(string) (string, bool) { return "", false }

func dockMenu(click int) (titles []string, supported bool) {
	golet.RunOnMain(func() { titles = darwin.TestDockMenu(click) })
	return titles, true
}

// Posting key events needs the accessibility permission on macOS.
func pressCtrlShiftK() bool { return false }

// Only Linux binds global shortcuts through a desktop portal.
func usePortalShortcuts() (func(), bool) { return nil, false }

func pressKeys(...string) bool { return false }

func fullScreenHidesToolbar(w *golet.Window) (hides bool, supported bool) {
	golet.RunOnMain(func() { hides = darwin.TestFullScreenHidesToolbar(w.NativeHandle()) })
	return hides, true
}

func trafficLights(w *golet.Window) (x, y float64, supported bool) {
	golet.RunOnMain(func() { x, y = darwin.TestTrafficLights(w.NativeHandle()) })
	return x, y, true
}

func setWindowAppearance(w *golet.Window, dark bool) {
	golet.RunOnMain(func() { darwin.TestSetWindowAppearance(w.NativeHandle(), dark) })
}

// Frameless windows keep native resize borders here.
func movePointer(int, int) bool                 { return false }
func pressButton(bool) bool                     { return false }
func resizeCursor(*golet.Window) (string, bool) { return "", false }

// The menu bar belongs to the application on macOS.
func menuBarShown(*golet.Window) (bool, bool)                { return false, false }
func activateAccelerator(*golet.Window, string) (bool, bool) { return false, false }
func enterMenuBar(*golet.Window, string) (bool, bool, bool)  { return false, false, false }
func openMenus(*golet.Window, byte) bool                     { return false }
func closeMenus()                                            {}

// The traffic lights are AppKit's own.
func titleButtons(*golet.Window) ([]string, bool) { return nil, false }
func pressTitleButton(*golet.Window, string) bool { return false }

func topNonClient(*golet.Window) (int32, int32, bool) { return 0, 0, false }

func menuBarColors(*golet.Window) (uint32, uint32, bool) { return 0, 0, false }

func titleButtonColor(*golet.Window, string) (uint8, uint8, uint8, bool) { return 0, 0, 0, false }

// AppKit composes the window: there is no redirection bitmap to go without.
func composition(*golet.Window) (bool, bool, bool) { return false, false, false }
func loseComposition(bool) bool                    { return false }

// WebKit creates web views with their windows.
func failWebViews(int) bool { return false }

func clickAndType(w *golet.Window, x, y float64, text string) (ok bool) {
	golet.RunOnMain(func() { ok = darwin.TestClickAndType(w.NativeHandle(), x, y, text) })
	return ok
}

// pressKey presses the key of a virtual key code in a window of native UI,
// through the input method.
func pressKey(w *golet.Window, code uint16, chars string) (ok bool) {
	golet.RunOnMain(func() { ok = darwin.TestKey(w.NativeHandle(), code, chars) })
	return ok
}

// AppKit routes a drag to the view under the pointer, which
// dropFiles, dropping on the content's, skips.
func dragFiles(*golet.Window, float64, float64, []string) (bool, bool) { return false, false }

// pressTab presses Tab in a window of native UI, with Shift when back, as
// the keyboard does: in the content or in the page of a web view.
func pressTab(w *golet.Window, back bool) (ok bool) {
	golet.RunOnMain(func() { ok = darwin.TestKeyShift(w.NativeHandle(), 48, "\t", back) })
	return ok
}

func compose(w *golet.Window, text string, caret int, commit bool) (ok bool) {
	golet.RunOnMain(func() { ok = darwin.TestCompose(w.NativeHandle(), text, caret, commit) })
	return ok
}

func controlClick(w *golet.Window, x, y float64) bool {
	golet.RunOnMain(func() { darwin.TestControlClick(w.NativeHandle(), x, y) })
	return true
}

func composeOver(w *golet.Window, text string, caret int, commit bool, from, length int) (ok bool) {
	golet.RunOnMain(func() { ok = darwin.TestComposeOver(w.NativeHandle(), text, caret, commit, from, length) })
	return ok
}

func inputClient(w *golet.Window) (selected [2]int, document string, ok bool) {
	golet.RunOnMain(func() { selected, document = darwin.TestInputClient(w.NativeHandle()) })
	return selected, document, true
}

func dropFiles(w *golet.Window, x, y float64, paths []string) (over, dropped, ok bool) {
	golet.RunOnMain(func() { over, dropped = darwin.TestDropFiles(w.NativeHandle(), x, y, paths) })
	return over, dropped, true
}

// The roles of elements in the accessibility API of AppKit.
const roleText, roleButton, roleCheckBox, roleTextField, roleSlider = "AXStaticText", "AXButton", "AXCheckBox", "AXTextField", "AXSlider"

// roleListItem is the role of the rows of a List.
const roleListItem = "AXRow/AXTableRow"

func accessibility(w *golet.Window) (nodes []accessNode, ok bool) {
	golet.RunOnMain(func() {
		for _, n := range darwin.TestAccessibility(w.NativeHandle()) {
			role := n.Role
			if n.Subrole != "" {
				role += "/" + n.Subrole
			}
			nodes = append(nodes, accessNode{role, n.Label, n.Value})
		}
	})
	return nodes, true
}

// axAttribute reads an attribute of an element through AppKit's older
// accessibility API, which only macOS has.
func axAttribute(w *golet.Window, label, attr string) (value string, settable, named, ok bool) {
	golet.RunOnMain(func() { value, settable, named, ok = darwin.TestAccessibilityAttribute(w.NativeHandle(), label, attr) })
	return value, settable, named, ok
}

func accessPerform(w *golet.Window, label, action, value string) (ok bool) {
	golet.RunOnMain(func() { ok = darwin.TestAccessibilityPerform(w.NativeHandle(), label, action, value) })
	return ok
}

// observe has key-value observing watch the view of a window showing
// native UI and its elements, until stop.
func observe(w *golet.Window) (stop func(), supported bool) {
	var end func()
	golet.RunOnMain(func() { end = darwin.TestObserve(w.NativeHandle()) })
	return func() { golet.RunOnMain(end) }, true
}

// Only Linux draws native UI in a GtkGLArea, nor waits to load the GPU's
// driver.
func glSurface(*golet.Window) (string, []byte, int, int, bool) { return "", nil, 0, 0, false }
func memoryUI(*testing.T) bool                                 { return false }
func surfaceOnScreen(*golet.Window) ([]byte, bool)             { return nil, false }

// Only Windows reads the screen over native UI.
func screenColor(*golet.Window, float64, float64) (uint8, uint8, uint8, bool) { return 0, 0, 0, false }

// Context menus are not automated on macOS: one shown waits for the user.
func rightClick(*golet.Window, float64, float64) bool { return false }
func popupMenus() ([][]string, bool)                  { return nil, false }
func choosePopupItem(string) bool                     { return false }

// pressEscape presses Escape in the key window, as the keyboard does.
func pressEscape(*golet.Window) bool {
	golet.RunOnMain(func() { darwin.TestAppKey(53, "\x1b") })
	return true
}

func documentChrome(w *golet.Window) (title, file string, icon, edited, supported bool) {
	golet.RunOnMain(func() { title, file, icon, edited = darwin.TestDocumentChrome(w.NativeHandle()) })
	return title, file, icon, edited, true
}

func sheet(w *golet.Window) (texts, buttons []string, ok, supported bool) {
	golet.RunOnMain(func() { texts, buttons, ok = darwin.TestSheet(w.NativeHandle()) })
	return texts, buttons, ok, true
}

func endSheetButton(w *golet.Window, i int) (ok bool) {
	golet.RunOnMain(func() { ok = darwin.TestEndSheetButton(w.NativeHandle(), i) })
	return ok
}

// Panels' window styles are Windows'.
func panelStyles(*golet.Window) (bool, bool, bool, bool) { return false, false, false, false }
