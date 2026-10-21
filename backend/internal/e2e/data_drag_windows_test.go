// Author: HyperparameterSI

//go:build windows && (amd64 || arm64)

package e2e

import (
	"github.com/HyperparameterSI/golet/backend"
	"github.com/HyperparameterSI/golet/backend/internal/windows"
	"github.com/HyperparameterSI/golet/backend/transfer"
)

func dropNativeData(w *golet.Window, x, y float64, d transfer.Data, ops transfer.Operation) (operation transfer.Operation, dropped, supported bool) {
	golet.RunOnMain(func() { operation, dropped = windows.TestDropData(w.NativeHandle(), x, y, d, ops) })
	return operation, dropped, true
}
