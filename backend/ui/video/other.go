// Author: HyperparameterSI

//go:build !(linux && (amd64 || arm64)) && !darwin && !(windows && (amd64 || arm64)) && !js

package video

import "github.com/HyperparameterSI/golet/backend/internal/media"

func newEngine(source) (media.Engine, error) { return nil, ErrUnsupported }
