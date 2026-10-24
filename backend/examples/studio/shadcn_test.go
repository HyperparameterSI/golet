// Author: HyperparameterSI
package main

import (
	"github.com/HyperparameterSI/golet/backend/ui"
	"github.com/HyperparameterSI/golet/backend/ui/shadcn"
	"testing"
)

func TestMaiaCatalogRenders(t *testing.T) {
	for _, name := range shadcn.Components() {
		t.Run(name, func(t *testing.T) {
			tt := ui.NewTester(newShadcnView(name), 1040, 840)
			tt.Frame()
			if !tt.HasText(name) {
				t.Fatal("missing demo heading")
			}
			if err := tt.Click("Dark"); err != nil {
				t.Fatal(err)
			}
			tt.Frame()
		})
	}
}
