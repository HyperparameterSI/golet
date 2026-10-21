// Author: HyperparameterSI

//go:build !darwin && !(linux && (amd64 || arm64)) && !(windows && (amd64 || arm64)) && !js

package e2e

import (
	"errors"
	"testing"

	"github.com/HyperparameterSI/golet/backend"
)

// Platforms without a backend run no GUI tests.
func activateMenu(*golet.Window, ...string) error { return errors.New("unsupported") }

func pressShortcut(string, bool) (bool, bool) { return false, false }

func endSheet(*golet.Window) (bool, bool) { return false, false }

func recentDocumentRecorded(string) (recorded, supported bool) { return false, false }

func click(*golet.Window, float64, float64) bool { return false }

func appClick(*golet.Window, float64, float64) bool { return false }

func heldAppClick(*golet.Window, float64, float64) bool { return false }

func drag(*golet.Window, [][2]float64) bool { return false }

func sideButton(*golet.Window, bool) bool { return false }

func webViewAttached(*golet.Window) (bool, bool) { return false, false }

// Collection behaviors and deactivating the app are macOS's.
func collectionBehavior(*golet.Window) (uint, bool) { return 0, false }
func deactivateApp() bool                           { return false }
func appActive() (bool, bool)                       { return false, false }

func dockDevTools(*golet.Window) bool { return false }

func dockedDevToolsPlace(*golet.Window) string { return "" }

func dismissPopups() (int, bool) { return 0, false }

func setDroppedFiles(*golet.Window, []string) bool { return false }

// The progress bar is shown by the shell, out of the app's reach.
func dockTileImage() ([]byte, bool) { return nil, false }

func defaultURLHandler(string) (string, bool) { return "", false }

func dockMenu(int) ([]string, bool) { return nil, false }

// Posting key events needs the accessibility permission on macOS.
func pressCtrlShiftK() bool { return false }

// Only Linux binds global shortcuts through a desktop portal.
func usePortalShortcuts() (func(), bool) { return nil, false }

func pressKeys(...string) bool { return false }

// Only macOS windows have a toolbar.
func fullScreenHidesToolbar(*golet.Window) (bool, bool) { return false, false }

// Only macOS windows have traffic lights.
func trafficLights(*golet.Window) (float64, float64, bool) { return 0, 0, false }
func setWindowAppearance(*golet.Window, bool)              {}

func movePointer(int, int) bool                 { return false }
func pressButton(bool) bool                     { return false }
func resizeCursor(*golet.Window) (string, bool) { return "", false }

func menuBarShown(*golet.Window) (bool, bool)                { return false, false }
func activateAccelerator(*golet.Window, string) (bool, bool) { return false, false }
func enterMenuBar(*golet.Window, string) (bool, bool, bool)  { return false, false, false }
func openMenus(*golet.Window, byte) bool                     { return false }
func closeMenus()                                            {}

func titleButtons(*golet.Window) ([]string, bool) { return nil, false }
func pressTitleButton(*golet.Window, string) bool { return false }

func topNonClient(*golet.Window) (int32, int32, bool) { return 0, 0, false }

func menuBarColors(*golet.Window) (uint32, uint32, bool) { return 0, 0, false }

func titleButtonColor(*golet.Window, string) (uint8, uint8, uint8, bool) { return 0, 0, 0, false }

func composition(*golet.Window) (bool, bool, bool) { return false, false, false }
func loseComposition(bool) bool                    { return false }

func failWebViews(int) bool { return false }

// A Control-click is a secondary click on macOS only.
func controlClick(*golet.Window, float64, float64) bool { return false }

// Input methods, file drops and assistive technology are only automated
// on macOS.
func composeOver(*golet.Window, string, int, bool, int, int) bool { return false }
func inputClient(*golet.Window) ([2]int, string, bool)            { return [2]int{}, "", false }
func dropFiles(*golet.Window, float64, float64, []string) (bool, bool, bool) {
	return false, false, false
}
func accessibility(*golet.Window) ([]accessNode, bool)         { return nil, false }
func accessPerform(*golet.Window, string, string, string) bool { return false }

const roleText, roleButton, roleCheckBox, roleTextField, roleSlider = "text", "button", "check box", "text field", "slider"

// roleListItem is the role of the rows of a List.
const roleListItem = "list item"

// Typing into native UI is only automated on macOS.
func clickAndType(*golet.Window, float64, float64, string) bool { return false }
func compose(*golet.Window, string, int, bool) bool             { return false }

// Only Linux draws native UI in a GtkGLArea, nor waits to load the GPU's
// driver.
func glSurface(*golet.Window) (string, []byte, int, int, bool) { return "", nil, 0, 0, false }
func memoryUI(*testing.T) bool                                 { return false }
func surfaceOnScreen(*golet.Window) ([]byte, bool)             { return nil, false }

// Only Windows reads the screen over native UI.
func screenColor(*golet.Window, float64, float64) (uint8, uint8, uint8, bool) { return 0, 0, 0, false }

// Only macOS has key-value observing.
func observe(*golet.Window) (func(), bool) { return nil, false }

func rightClick(*golet.Window, float64, float64) bool { return false }
func popupMenus() ([][]string, bool)                  { return nil, false }
func choosePopupItem(string) bool                     { return false }

// Input methods keep the keys typed while they compose from the content
// themselves here (GTK's filtering, IMM32's VK_PROCESSKEY), which a
// composition the test makes up does not show.
func pressKey(*golet.Window, uint16, string) bool { return false }
func pressTab(*golet.Window, bool) bool           { return false }

func dragFiles(*golet.Window, float64, float64, []string) (bool, bool) { return false, false }

// AppKit's older accessibility API is macOS's.
func axAttribute(*golet.Window, string, string) (string, bool, bool, bool) {
	return "", false, false, false
}

func pressEscape(*golet.Window) bool { return false }

func documentChrome(*golet.Window) (string, string, bool, bool, bool) {
	return "", "", false, false, false
}
func sheet(*golet.Window) ([]string, []string, bool, bool) { return nil, nil, false, false }
func endSheetButton(*golet.Window, int) bool               { return false }

// Panels' window styles are Windows'.
func panelStyles(*golet.Window) (bool, bool, bool, bool) { return false, false, false, false }
