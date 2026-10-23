// Author: HyperparameterSI

//go:build darwin

package golet

import (
	"github.com/HyperparameterSI/golet/backend/internal/darwin"
	"github.com/HyperparameterSI/golet/backend/internal/platform"
)

func newBackend() platform.Backend { return darwin.New() }

func enableNotificationSenders() { darwin.EnableNotificationSenders() }

func linkContentRules() { darwin.LinkContentRules() }
