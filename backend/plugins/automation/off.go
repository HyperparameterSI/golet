// Author: HyperparameterSI

//go:build golet_noinspector

package automation

import "os"

// Builds without the inspector leave the plugin out, as they are the
// app's release: it says so rather than leave a test runner waiting.
func init() {
	if os.Getenv("GOLET_AUTOMATION") != "" {
		os.Stderr.WriteString("golet: GOLET_AUTOMATION is ignored: this build leaves automation out (golet_noinspector)\n")
	}
}
