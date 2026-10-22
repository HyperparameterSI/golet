// Author: HyperparameterSI
//go:build !darwin && !linux && !windows

package chromium

import "os/exec"

func prepareProcess(cmd *exec.Cmd)    {}
func killProcess(cmd *exec.Cmd) error { return cmd.Process.Kill() }
