// Author: HyperparameterSI

//go:build windows && (amd64 || arm64)

package golet

import (
	"github.com/HyperparameterSI/golet/backend/internal/platform"
	win "github.com/HyperparameterSI/golet/backend/internal/windows"
)

func newBackend() platform.Backend { return win.New() }
