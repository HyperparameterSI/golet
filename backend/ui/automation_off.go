// Author: HyperparameterSI

//go:build golet_noinspector

package ui

// Production builds leave plugins/automation's hook out, as they do the
// inspector.
func (h *windowHost) automate() {}
