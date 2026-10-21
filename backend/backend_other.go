// Author: HyperparameterSI

//go:build !darwin && !(linux && (amd64 || arm64)) && !(windows && (amd64 || arm64)) && !js

package golet

import (
	"github.com/HyperparameterSI/golet/backend/internal/platform"
	"github.com/HyperparameterSI/golet/backend/internal/unsupported"
)

func newBackend() platform.Backend { return unsupported.New() }
