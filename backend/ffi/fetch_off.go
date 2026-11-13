// Author: HyperparameterSI

//go:build golet_nodownload

package ffi

import "github.com/HyperparameterSI/golet/backend/internal/nativelib"

var fetch func(f nativelib.File, path string) error
