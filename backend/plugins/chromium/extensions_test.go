// Author: HyperparameterSI
package chromium

import (
	"os"
	"path/filepath"
	"testing"
)

func TestChromeExtensionArguments(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "manifest.json"), []byte(`{"manifest_version":3,"name":"test"}`), 0600)
	if _, err := extensionArguments(Options{Extensions: []string{dir}}); err == nil {
		t.Fatal("Alloy accepted Chrome extension")
	}
	args, err := extensionArguments(Options{Style: StyleChrome, Extensions: []string{dir}})
	if err != nil || len(args) != 2 {
		t.Fatalf("%v %v", args, err)
	}
	os.WriteFile(filepath.Join(dir, "manifest.json"), []byte(`{"manifest_version":2,"name":"test"}`), 0600)
	if _, err := extensionArguments(Options{Style: StyleChrome, Extensions: []string{dir}}); err == nil {
		t.Fatal("MV2 accepted")
	}
}
