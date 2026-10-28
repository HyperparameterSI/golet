// Author: HyperparameterSI

//go:build linux && (amd64 || arm64) && !golet_noinspector

package automation

import (
	"github.com/HyperparameterSI/golet/backend"
	"github.com/HyperparameterSI/golet/backend/internal/linux"
)

// chooseMenuItem chooses an item of the window's menu bar, which shows
// the application menu, by the labels along its path. Main thread only.
func chooseMenuItem(w *golet.Window, path []string) error {
	return linux.TestActivateMenuItem(w.NativeHandle(), path...)
}
