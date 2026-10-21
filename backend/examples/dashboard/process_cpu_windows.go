// Author: HyperparameterSI

package main

import "syscall"

func processCPUSeconds() (float64, error) {
	process, err := syscall.GetCurrentProcess()
	if err != nil {
		return 0, err
	}
	var created, exited, kernel, user syscall.Filetime
	if err = syscall.GetProcessTimes(process, &created, &exited, &kernel, &user); err != nil {
		return 0, err
	}
	ticks := func(t syscall.Filetime) float64 { return float64(uint64(t.HighDateTime)<<32 | uint64(t.LowDateTime)) }
	return (ticks(kernel) + ticks(user)) / 1e7, nil
}
