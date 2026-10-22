// Author: HyperparameterSI
package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

var cmakeVersionPattern = regexp.MustCompile(`(?m)^cmake version (\d+)\.(\d+)(?:\.|\s|$)`)

func validateCMakeVersion(output string) error {
	match := cmakeVersionPattern.FindStringSubmatch(output)
	if match == nil {
		return fmt.Errorf("cannot identify CMake version: %s", strings.TrimSpace(output))
	}
	major, _ := strconv.Atoi(match[1])
	minor, _ := strconv.Atoi(match[2])
	if major < 3 || (major == 3 && minor < 21) {
		return fmt.Errorf("CMake %s.%s is too old; 3.21+ is required", match[1], match[2])
	}
	return nil
}

// findCMake honors explicit configuration and checks standard installation paths
// when an IDE or terminal has not inherited the tool's PATH.
func findCMake(ctx context.Context, requested string) (string, error) {
	explicit := requested != "" && requested != "cmake"
	candidates := []string{requested}
	if !explicit {
		if configured := os.Getenv("CMAKE"); configured != "" {
			candidates = []string{configured}
			explicit = true
		} else {
			candidates = []string{"cmake", "/opt/homebrew/bin/cmake", "/usr/local/bin/cmake", "/Applications/CMake.app/Contents/bin/cmake"}
			if runtime.GOOS == "windows" {
				candidates = []string{"cmake"}
				if dir := os.Getenv("ProgramFiles"); dir != "" {
					candidates = append(candidates, filepath.Join(dir, "CMake", "bin", "cmake.exe"))
				}
			}
		}
	}
	var problems []string
	for _, candidate := range candidates {
		path, err := exec.LookPath(candidate)
		if err != nil {
			if explicit {
				problems = append(problems, err.Error())
			}
			continue
		}
		output, err := exec.CommandContext(ctx, path, "--version").CombinedOutput()
		if err == nil {
			err = validateCMakeVersion(string(output))
		}
		if err == nil {
			return path, nil
		}
		problems = append(problems, fmt.Sprintf("%s: %v", path, err))
	}
	detail := ""
	if len(problems) > 0 {
		detail = "\n" + strings.Join(problems, "\n")
	}
	install := "brew install cmake"
	if runtime.GOOS == "windows" {
		install = "winget install Kitware.CMake (also install Visual Studio 2022 Desktop development with C++ and Windows SDK)"
	}
	return "", fmt.Errorf("CMake 3.21+ was not found or could not run%s\nInstall it with: %s\nThen retry: pnpm setup:chromium\nOr specify: pnpm setup:chromium --cmake /absolute/path/to/cmake", detail, install)
}
