// Author: HyperparameterSI
package main

import (
	"testing"

	"github.com/HyperparameterSI/golet/backend/ui"
	"github.com/HyperparameterSI/golet/backend/ui/forma"
)

func TestMaiaCatalogRenders(t *testing.T) {
	for _, name := range forma.Components() {
		t.Run(name, func(t *testing.T) {
			tt := ui.NewTester(newFormaView(name), 1040, 840)
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

func TestMaiaThemePicker(t *testing.T) {
	tt := ui.NewTester(newFormaView("button"), 1040, 840)
	if err := tt.Click("Theme palette"); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 3; i++ {
		tt.Key(0, ui.KeyDown)
	}
	tt.Key(0, ui.KeyEnter)
	if !tt.HasText("Blue") {
		t.Fatal("palette did not change")
	}
	if err := tt.Click("Corner style"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyEnd)
	tt.Key(0, ui.KeyEnter)
	if !tt.HasText("Sharp") {
		t.Fatal("corner style did not change")
	}
	if err := tt.Click("Dark"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if !tt.HasText("Default") {
		t.Fatal("lost component after theme switch")
	}
}
