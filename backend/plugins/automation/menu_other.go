// Author: HyperparameterSI

//go:build !darwin && !(linux && (amd64 || arm64)) && !(windows && (amd64 || arm64)) && !golet_noinspector

package automation

import "github.com/HyperparameterSI/golet/backend"

// chooseMenuItem has no menus to choose from on other platforms.
func chooseMenuItem(*golet.Window, []string) error {
	return errorf(codeUnsupported, "this platform has no menus to choose from")
}
