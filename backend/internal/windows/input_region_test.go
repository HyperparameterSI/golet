// Author: HyperparameterSI

//go:build windows && (amd64 || arm64)

package windows

import "github.com/HyperparameterSI/golet/backend/internal/platform"

// The core reaches a window's input region through the interface alone.
var _ platform.InputRegionWindow = (*window)(nil)
