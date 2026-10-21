// Author: HyperparameterSI

//go:build !windows && !darwin && !(linux && (amd64 || arm64)) && !js

package ui

import (
	"errors"

	"github.com/HyperparameterSI/golet/backend/internal/platform"
)

func newGPURenderer(platform.SurfaceNative) (gpuRenderer, error) {
	return nil, errors.New("no GPU renderer for this platform yet")
}
