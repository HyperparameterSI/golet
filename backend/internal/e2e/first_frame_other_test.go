// Author: HyperparameterSI

//go:build !js && !(windows && (amd64 || arm64))

package e2e

import "github.com/HyperparameterSI/golet/backend"

// firstFrameRevealed reports whether w shows its first frame once shown:
// only Windows keeps a window out of sight after it showed (cloaked).
func firstFrameRevealed(*golet.Window) bool { return true }
