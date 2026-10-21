// Author: HyperparameterSI
//go:build !darwin && !linux

package chromium

import "os/exec"

// Native CEF host packaging is currently available only for macOS.
func prepareProcess(cmd *exec.Cmd)    {}
func killProcess(cmd *exec.Cmd) error { return cmd.Process.Kill() }
