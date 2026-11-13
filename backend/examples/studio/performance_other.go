//go:build !darwin

// Author: HyperparameterSI
package main

type unsupportedSource struct{}

func newSource() source { return unsupportedSource{} }
func (unsupportedSource) Read() raw {
	return raw{errors: []string{"此平台暂仅提供 Go runtime 指标；系统与进程指标尚未接入。"}}
}
func (unsupportedSource) Close() {}
