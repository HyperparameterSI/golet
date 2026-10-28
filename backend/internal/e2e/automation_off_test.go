// Author: HyperparameterSI

//go:build golet_noinspector && !js

package e2e

// Builds without the inspector leave plugins/automation out.
func automationApp() {}
