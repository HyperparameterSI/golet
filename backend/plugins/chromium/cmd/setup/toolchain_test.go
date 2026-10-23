// Author: HyperparameterSI
package main

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestCMakeVersion(t *testing.T) {
	for _, tc := range []struct {
		version string
		valid   bool
	}{{"cmake version 3.21.0", true}, {"cmake version 4.2.1", true}, {"cmake version 3.20.9", false}, {"cmake version 2.99.0", false}, {"not cmake", false}} {
		if err := validateCMakeVersion(tc.version); (err == nil) != tc.valid {
			t.Errorf("%s: %v", tc.version, err)
		}
	}
}
func TestExplicitCMake(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cmake.exe")
	source := filepath.Join(dir, "main.go")
	// A native executable exercises process discovery on Windows as well as Unix.
	if err := os.WriteFile(source, []byte("package main\nimport \"fmt\"\nfunc main() { fmt.Println(\"cmake version 3.21.0\") }\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if out, err := exec.Command("go", "build", "-o", path, source).CombinedOutput(); err != nil {
		t.Fatalf("build CMake fixture: %v\n%s", err, out)
	}
	got, err := findCMake(context.Background(), path)
	if err != nil || got != path {
		t.Fatalf("%q %v", got, err)
	}
	t.Setenv("CMAKE", path)
	got, err = findCMake(context.Background(), "")
	if err != nil || got != path {
		t.Fatalf("environment: %q %v", got, err)
	}
	_, err = findCMake(context.Background(), filepath.Join(dir, "missing"))
	install := "brew install cmake"
	if runtime.GOOS == "windows" {
		install = "winget install Kitware.CMake"
	}
	if err == nil || !strings.Contains(err.Error(), install) {
		t.Fatalf("missing actionable error: %v", err)
	}
}
