// Author: HyperparameterSI

//go:build linux && (amd64 || arm64)

package golet

import (
	"github.com/HyperparameterSI/golet/backend/internal/linux"
	"github.com/HyperparameterSI/golet/backend/internal/platform"
)

func newBackend() platform.Backend { return linux.New() }
