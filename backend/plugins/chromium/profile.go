// Author: HyperparameterSI
package chromium

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
)

// prepareProfile never reuses persistent data with the development keychain.
func prepareProfile(opts Options) (string, []string, func(), error) {
	cleanup := func() {}
	if opts.EphemeralDemo {
		if opts.ProfileDir != "" {
			return "", nil, cleanup, errors.New("chromium: EphemeralDemo cannot use a persistent ProfileDir")
		}
		dir, err := os.MkdirTemp("", "golet-browser-demo-")
		if err != nil {
			return "", nil, cleanup, err
		}
		var args []string
		if runtime.GOOS == "darwin" {
			args = append(args, "--use-mock-keychain")
		}
		return dir, args, func() { _ = os.RemoveAll(dir) }, nil
	}
	dir := opts.ProfileDir
	if dir == "" {
		base, err := os.UserConfigDir()
		if err != nil {
			return "", nil, cleanup, err
		}
		dir = filepath.Join(base, "golet", "chromium", "default")
	}
	dir, err := filepath.Abs(dir)
	if err != nil {
		return "", nil, cleanup, err
	}
	if err = os.MkdirAll(dir, 0700); err != nil {
		return "", nil, cleanup, err
	}
	return dir, nil, cleanup, nil
}
