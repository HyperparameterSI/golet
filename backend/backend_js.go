// Author: HyperparameterSI

//go:build js

package golet

import (
	"github.com/HyperparameterSI/golet/backend/internal/platform"
	"github.com/HyperparameterSI/golet/backend/internal/web"
)

func newBackend() platform.Backend { return web.New() }
