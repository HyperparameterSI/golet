// Author: HyperparameterSI
package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPackageWindowsRuntime(t *testing.T) {
	root := t.TempDir()
	put := func(name, value string) {
		t.Helper()
		p := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(value), 0644); err != nil {
			t.Fatal(err)
		}
	}
	for _, file := range []string{"Release/libcef.dll", "Release/bootstrapc.exe", "Release/libcef.lib", "Resources/icudtl.dat", "Resources/resources.pak", "Resources/locales/en-US.pak", "LICENSE.txt"} {
		put("sdk/"+file, file)
	}
	put("build/Release/golet-chromium.dll", "client")
	put("frontend/examples/browser/newtab/index.html", "home")
	put("frontend/examples/browser/newtab/branding.js", "default")
	out := filepath.Join(root, "runtime")
	if err := packageWindows(context.Background(), root, filepath.Join(root, "sdk"), filepath.Join(root, "build"), out, options{appName: "Test Browser"}); err != nil {
		t.Fatal(err)
	}
	for _, file := range []string{"golet-chromium.exe", "golet-chromium.dll", "libcef.dll", "icudtl.dat", "resources.pak", "locales/en-US.pak", "CEF-LICENSE.txt", "newtab/index.html"} {
		if _, err := os.Stat(filepath.Join(out, file)); err != nil {
			t.Fatal(err)
		}
	}
	data, err := os.ReadFile(filepath.Join(out, "golet-chromium.exe"))
	if err != nil || string(data) != "Release/bootstrapc.exe" {
		t.Fatalf("sandbox bootstrap not preserved: %q %v", data, err)
	}
	data, err = os.ReadFile(filepath.Join(out, "newtab/branding.js"))
	if err != nil || !strings.Contains(string(data), "Test Browser") {
		t.Fatalf("branding: %q %v", data, err)
	}
	if _, err := os.Stat(filepath.Join(out, "libcef.lib")); !os.IsNotExist(err) {
		t.Fatal("import library was packaged")
	}
	if err := os.Remove(filepath.Join(root, "sdk/Release/bootstrapc.exe")); err != nil {
		t.Fatal(err)
	}
	if err := packageWindows(context.Background(), root, filepath.Join(root, "sdk"), filepath.Join(root, "build"), filepath.Join(root, "incomplete"), options{}); err == nil {
		t.Fatal("accepted SDK without sandbox bootstrap")
	}
}
