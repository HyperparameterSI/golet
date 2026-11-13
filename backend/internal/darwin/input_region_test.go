// Author: HyperparameterSI

//go:build darwin

package darwin

import "github.com/HyperparameterSI/golet/backend/internal/platform"

// The core reaches a window's input region through the interface alone.
var _ platform.InputRegionWindow = (*window)(nil)
