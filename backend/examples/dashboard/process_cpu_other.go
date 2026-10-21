// Author: HyperparameterSI

//go:build !darwin && !linux && !windows

package main

import "errors"

func processCPUSeconds() (float64, error) {
	return 0, errors.New("process CPU metrics unavailable on this platform")
}
