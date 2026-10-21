// Author: HyperparameterSI

// Command setup builds and installs the optional CEF host without Python.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type runtimeEntry struct {
	Version string `json:"version"`
	URL     string `json:"url"`
	SHA1    string `json:"sha1"`
	Size    int64  `json:"size"`
}

type options struct{ root, sdk, cmake, output string }

func main() {
	var opts options
	flag.StringVar(&opts.root, "root", "", "repository root (defaults to nearest parent checkout)")
	flag.StringVar(&opts.sdk, "cef-dir", "", "use an extracted, matching CEF minimal SDK")
	flag.StringVar(&opts.cmake, "cmake", "cmake", "CMake executable")
	flag.StringVar(&opts.output, "output", "", "runtime installation directory")
	flag.Parse()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	if err := install(ctx, opts); err != nil {
		fmt.Fprintln(os.Stderr, "chromium setup:", err)
		os.Exit(1)
	}
}

func repositoryRoot(start string) (string, error) {
	if start == "" {
		var err error
		start, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}
	root, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(root, "backend/plugins/chromium/native/runtime.json")); err == nil {
			return root, nil
		}
		parent := filepath.Dir(root)
		if parent == root {
			return "", errors.New("cannot find Golet checkout; use --root")
		}
		root = parent
	}
}

func run(ctx context.Context, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s: %w", name, err)
	}
	return nil
}

func install(ctx context.Context, opts options) error {
	if runtime.GOOS != "darwin" {
		return errors.New("native host packaging currently supports macOS only")
	}
	root, err := repositoryRoot(opts.root)
	if err != nil {
		return err
	}
	cmake, err := findCMake(ctx, opts.cmake)
	if err != nil {
		return err
	}
	fmt.Println("Using CMake:", cmake)
	source := filepath.Join(root, "backend/plugins/chromium/native")
	data, err := os.ReadFile(filepath.Join(source, "runtime.json"))
	if err != nil {
		return err
	}
	var entries map[string]runtimeEntry
	if err = json.Unmarshal(data, &entries); err != nil {
		return err
	}
	platform := runtime.GOOS + "-" + runtime.GOARCH
	entry, ok := entries[platform]
	if !ok {
		return fmt.Errorf("unsupported platform %s", platform)
	}
	cache, err := os.UserCacheDir()
	if err != nil {
		return err
	}
	cache = filepath.Join(cache, "golet", "chromium")
	if err = os.MkdirAll(cache, 0755); err != nil {
		return err
	}
	sdk := opts.sdk
	if sdk == "" {
		sdk, err = fetchSDK(ctx, entry, cache)
	} else {
		sdk, err = filepath.Abs(sdk)
	}
	if err != nil {
		return err
	}
	header, err := os.ReadFile(filepath.Join(sdk, "include/cef_version.h"))
	if err != nil {
		return err
	}
	if !strings.Contains(string(header), strconv.Quote(entry.Version)) {
		return errors.New("CEF SDK version differs from native/runtime.json")
	}
	output := opts.output
	if output == "" {
		output = filepath.Join(cache, "host-v1", platform)
	}
	output, err = filepath.Abs(output)
	if err != nil {
		return err
	}
	if err = validateOutput(output); err != nil {
		return err
	}
	if err = os.MkdirAll(filepath.Dir(output), 0755); err != nil {
		return err
	}
	work := filepath.Join(root, ".golet/chromium-build", platform)
	arch := "arm64"
	if runtime.GOARCH == "amd64" {
		arch = "x86_64"
	}
	if err = run(ctx, cmake, "-S", source, "-B", work, "-DCEF_ROOT="+sdk, "-DPROJECT_ARCH="+arch, "-DCMAKE_BUILD_TYPE=Release"); err != nil {
		return err
	}
	if err = run(ctx, cmake, "--build", work, "--parallel", strconv.Itoa(min(runtime.NumCPU(), 6))); err != nil {
		return err
	}
	staging, err := os.MkdirTemp(filepath.Dir(output), ".golet-chromium-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(staging)
	prepared := filepath.Join(staging, "runtime")
	if err = packageApp(ctx, root, sdk, work, prepared, staging); err != nil {
		return err
	}
	marker, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return err
	}
	if err = os.WriteFile(filepath.Join(prepared, "golet-chromium-runtime.json"), append(marker, '\n'), 0644); err != nil {
		return err
	}
	if err = replaceRuntime(prepared, output); err != nil {
		return err
	}
	fmt.Println("Installed:", filepath.Join(output, "Golet Browser.app"))
	fmt.Println("Run: pnpm dev:browser")
	return nil
}

func validateOutput(output string) error {
	info, err := os.Lstat(output)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return errors.New("output must be a directory, not a symlink or file")
	}
	marker, err := os.Lstat(filepath.Join(output, "golet-chromium-runtime.json"))
	if err != nil || !marker.Mode().IsRegular() {
		return errors.New("refusing to replace a directory not installed by this builder")
	}
	return nil
}

// Keep the old installation outside staging so even a failed rollback cannot
// be deleted by staging cleanup. Publishing uses renames on the same volume.
func replaceRuntime(prepared, output string) error {
	if err := validateOutput(output); err != nil {
		return err
	}
	if _, err := os.Lstat(output); os.IsNotExist(err) {
		return os.Rename(prepared, output)
	}
	backup, err := os.MkdirTemp(filepath.Dir(output), ".golet-chromium-backup-")
	if err != nil {
		return err
	}
	old := filepath.Join(backup, "runtime")
	if err = os.Rename(output, old); err != nil {
		os.Remove(backup)
		return err
	}
	if err = os.Rename(prepared, output); err != nil {
		if restoreErr := os.Rename(old, output); restoreErr != nil {
			return fmt.Errorf("publish failed: %v; restore failed: %v; previous runtime retained at %s", err, restoreErr, old)
		}
		os.Remove(backup)
		return err
	}
	return os.RemoveAll(backup)
}
