// Author: HyperparameterSI

//go:build !darwin && !(windows && (amd64 || arm64)) && !js

package e2e

import (
	"github.com/HyperparameterSI/golet/backend"
	"github.com/HyperparameterSI/golet/backend/transfer"
)

func dropNativeData(*golet.Window, float64, float64, transfer.Data, transfer.Operation) (transfer.Operation, bool, bool) {
	return transfer.None, false, false
}
