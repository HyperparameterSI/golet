// Author: HyperparameterSI

//go:build darwin && !golet_noinspector

package automation

import (
	"github.com/HyperparameterSI/golet/backend"
	"github.com/HyperparameterSI/golet/backend/internal/darwin"
)

// chooseMenuItem chooses an item of the application menu by the titles
// along its path, as the menu bar does. Main thread only.
func chooseMenuItem(_ *golet.Window, path []string) error {
	return darwin.TestPerformMenuItem(path...)
}
