// Author: HyperparameterSI

//go:build !(linux && (amd64 || arm64))

package golet

func linkDMABUFProbe() {}
