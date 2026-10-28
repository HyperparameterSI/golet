// Author: HyperparameterSI

//go:build windows && (amd64 || arm64)

package golet

import (
	"github.com/HyperparameterSI/golet/backend/internal/platform"
	win "github.com/HyperparameterSI/golet/backend/internal/windows"
)

func newBackend() platform.Backend { return win.New() }

func linkContentRules() { win.LinkContentRules() }

func linkNowPlaying() { win.LinkNowPlaying() }

func linkMenuIcons() { win.LinkMenuIcons() }

func linkContextMenus() { win.LinkContextMenus() }
