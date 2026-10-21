// Author: HyperparameterSI
//go:build darwin || linux

package chromium

import (
	"os"
	"os/exec"
	"syscall"
)

// A private process group lets forced shutdown also reap CEF helpers. Normal
// shutdown always goes through the protocol to let Chromium flush its profile.
func prepareProcess(cmd *exec.Cmd) { cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true} }
func killProcess(cmd *exec.Cmd) error {
	err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	if err == syscall.ESRCH {
		return os.ErrProcessDone
	}
	return err
}
