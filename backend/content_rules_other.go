// Author: HyperparameterSI

//go:build !darwin && !(linux && (amd64 || arm64)) && !(windows && (amd64 || arm64))

package golet

func linkContentRules() {}
