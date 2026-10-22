// Author: HyperparameterSI
//go:build windows

package chromium

import (
	"os/exec"
	"syscall"
)

// The console bootstrap needs redirected pipes, but no visible console window.
func prepareProcess(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
}
func killProcess(cmd *exec.Cmd) error { return cmd.Process.Kill() }
