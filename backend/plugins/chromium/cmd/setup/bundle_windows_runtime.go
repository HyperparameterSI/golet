// Author: HyperparameterSI
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// packageWindows keeps the official sandbox bootstrap beside the client DLL.
// Renaming both to the same stem lets CEF locate the DLL without extra switches.
func packageWindows(ctx context.Context, root, sdk, work, destination string, opts options) error {
	if err := os.MkdirAll(destination, 0755); err != nil {
		return err
	}
	for _, file := range []string{"libcef.dll", "bootstrapc.exe"} {
		if info, err := os.Stat(filepath.Join(sdk, "Release", file)); err != nil || !info.Mode().IsRegular() {
			return fmt.Errorf("Windows CEF SDK missing Release/%s", file)
		}
	}
	for _, file := range []string{"icudtl.dat", "resources.pak"} {
		if info, err := os.Stat(filepath.Join(sdk, "Resources", file)); err != nil || !info.Mode().IsRegular() {
			return fmt.Errorf("Windows CEF SDK missing Resources/%s", file)
		}
	}
	entries, err := os.ReadDir(filepath.Join(sdk, "Release"))
	if err != nil {
		return err
	}
	for _, entry := range entries {
		// Import libraries and alternative bootstrap launchers are build inputs.
		if strings.HasSuffix(entry.Name(), ".lib") || strings.HasPrefix(entry.Name(), "bootstrap") {
			continue
		}
		source := filepath.Join(sdk, "Release", entry.Name())
		target := filepath.Join(destination, entry.Name())
		if entry.IsDir() {
			err = copyTree(ctx, source, target)
		} else {
			err = copyFile(source, target)
		}
		if err != nil {
			return err
		}
	}
	if err = copyFile(filepath.Join(sdk, "Release/bootstrapc.exe"), filepath.Join(destination, "golet-chromium.exe")); err != nil {
		return err
	}
	if err = copyFile(filepath.Join(work, "Release/golet-chromium.dll"), filepath.Join(destination, "golet-chromium.dll")); err != nil {
		return err
	}
	if err = copyTree(ctx, filepath.Join(sdk, "Resources"), destination); err != nil {
		return err
	}
	newTab := opts.newTabDir
	if newTab == "" {
		newTab = filepath.Join(root, "frontend/examples/browser/newtab")
	}
	if info, err := os.Stat(filepath.Join(newTab, "index.html")); err != nil || !info.Mode().IsRegular() {
		return fmt.Errorf("new-tab directory must contain index.html")
	}
	target := filepath.Join(destination, "newtab")
	if err = copyTree(ctx, newTab, target); err != nil {
		return err
	}
	if opts.newTabDir == "" {
		logo := ""
		if opts.icon != "" {
			logo = "product-icon.png"
			if err = copyFile(opts.icon, filepath.Join(target, logo)); err != nil {
				return err
			}
		}
		name := opts.appName
		if name == "" {
			name = "Golet Browser"
		}
		branding, err := json.Marshal(map[string]string{"name": name, "icon": logo})
		if err != nil {
			return err
		}
		file := filepath.Join(target, "branding.js")
		if err = os.Remove(file); err != nil && !os.IsNotExist(err) {
			return err
		}
		if err = os.WriteFile(file, append(append([]byte("window.browserBranding = "), branding...), ';'), 0644); err != nil {
			return err
		}
	}
	return copyFile(filepath.Join(sdk, "LICENSE.txt"), filepath.Join(destination, "CEF-LICENSE.txt"))
}
