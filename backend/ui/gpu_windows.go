// Author: HyperparameterSI

package ui

import (
	"errors"

	"github.com/HyperparameterSI/golet/backend/internal/gpu/d3d11"
	"github.com/HyperparameterSI/golet/backend/internal/platform"
)

func newGPURenderer(n platform.SurfaceNative) (gpuRenderer, error) {
	if n.HWND == 0 {
		return nil, errors.New("the surface has no window")
	}
	if n.Composed {
		// The window shows its material where the frames are transparent.
		return d3d11.NewComposed(n.HWND)
	}
	return d3d11.New(n.HWND)
}
