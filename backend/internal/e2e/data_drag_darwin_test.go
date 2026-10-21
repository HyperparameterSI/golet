// Author: HyperparameterSI

//go:build darwin

package e2e

import (
	"github.com/HyperparameterSI/golet/backend"
	"github.com/HyperparameterSI/golet/backend/internal/darwin"
	"github.com/HyperparameterSI/golet/backend/transfer"
)

func dropNativeData(w *golet.Window, x, y float64, d transfer.Data, ops transfer.Operation) (operation transfer.Operation, dropped, supported bool) {
	golet.RunOnMain(func() { operation, dropped = darwin.TestDropData(w.NativeHandle(), x, y, d, ops) })
	return operation, dropped, true
}
