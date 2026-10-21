// Author: HyperparameterSI
package chromium

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestDemoProfileIsIsolated(t *testing.T) {
	dir, args, cleanup, err := prepareProfile(Options{EphemeralDemo: true})
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	if runtime.GOOS == "darwin" && (len(args) != 1 || args[0] != "--use-mock-keychain") {
		t.Fatalf("missing demo keychain flag: %v", args)
	}
	other, _, remove, err := prepareProfile(Options{EphemeralDemo: true})
	if err != nil {
		t.Fatal(err)
	}
	defer remove()
	if dir == other {
		t.Fatal("profiles reused")
	}
	cleanup()
	if _, err := os.Stat(dir); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("profile survived cleanup: %v", err)
	}
}
func TestDemoCannotReusePersistentProfile(t *testing.T) {
	dir := t.TempDir()
	marker := filepath.Join(dir, "existing-data")
	os.WriteFile(marker, []byte("keep"), 0600)
	if _, _, _, err := prepareProfile(Options{EphemeralDemo: true, ProfileDir: dir}); err == nil {
		t.Fatal("persistent demo profile accepted")
	}
	_, args, cleanup, err := prepareProfile(Options{ProfileDir: dir})
	if err != nil {
		t.Fatal(err)
	}
	cleanup()
	if len(args) != 0 {
		t.Fatal("normal profile weakened")
	}
	if _, err := os.Stat(marker); err != nil {
		t.Fatal("existing profile removed")
	}
}
