// Author: HyperparameterSI

//go:build windows && (amd64 || arm64) && !golet_noinspector

package automation

import (
	"github.com/HyperparameterSI/golet/backend"
	"github.com/HyperparameterSI/golet/backend/internal/windows"
)

// chooseMenuItem chooses an item of the window's menu bar, which shows
// the application menu, by the labels along its path. Main thread only.
func chooseMenuItem(w *golet.Window, path []string) error {
	return windows.TestActivateMenuItem(w.NativeHandle(), path...)
}
