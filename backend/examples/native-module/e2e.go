// Author: HyperparameterSI

//go:build e2e

package main

// Builds with -tags e2e can be driven through the automation plugin:
// GOLET_AUTOMATION=<socket> ./app, then golet-automation.
import _ "github.com/HyperparameterSI/golet/backend/plugins/automation"
