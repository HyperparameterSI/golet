// Author: HyperparameterSI
package main

import (
	"context"
	"os"
	"path/filepath"
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
	path := filepath.Join(dir, "cmake")
	if err := os.WriteFile(path, []byte("#!/bin/sh\necho 'cmake version 3.21.0'\n"), 0700); err != nil {
		t.Fatal(err)
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
	if err == nil || !strings.Contains(err.Error(), "brew install cmake") {
		t.Fatalf("missing actionable error: %v", err)
	}
}
