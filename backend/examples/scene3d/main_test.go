// Author: HyperparameterSI

package main

import (
	"os"
	"testing"
	"time"

	"github.com/HyperparameterSI/golet/backend/plugins/scene3d"
	"github.com/HyperparameterSI/golet/backend/ui"
)

func TestOBJExample(t *testing.T) {
	f, err := os.Open("assets/pyramid.obj")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	mesh, err := scene3d.LoadOBJ(f, scene3d.OBJLimits{})
	if err != nil {
		t.Fatal(err)
	}
	d := newDemo(mesh)
	tt := ui.NewTester(d.view, 860, 620)
	tt.Frame()
	if mesh.TriangleCount() != 6 || d.renderError != nil || d.stats.ShadedPixels == 0 {
		t.Fatalf("OBJ demo failed: %v", d.renderError)
	}
}

func TestFormaControlsAndAnimation(t *testing.T) {
	d := newDemo(scene3d.Cube())
	tt := ui.NewTester(d.view, 1200, 820)
	tt.Frame()
	if d.renderError != nil || d.stats.ShadedPixels == 0 {
		t.Fatalf("initial frame: %+v / %v", d.stats, d.renderError)
	}
	if tt.FramePending() {
		t.Fatal("static demo animates")
	}
	camera := d.camera
	tt.Click("Orbit left")
	if camera == d.camera {
		t.Fatal("Forma orbit button did not update camera")
	}
	tt.Click("Light / dark")
	if d.dark {
		t.Fatal("theme button did not switch")
	}
	tt.Click("Play / pause")
	before := d.scene.Instances[0].Transform
	tt.Advance(time.Second / 30)
	if before == d.scene.Instances[0].Transform {
		t.Fatal("animation did not advance")
	}
	tt.Click("Play / pause")
	tt.Frame()
	before = d.scene.Instances[0].Transform
	tt.Advance(time.Second)
	if before != d.scene.Instances[0].Transform || tt.FramePending() {
		t.Fatal("paused animation keeps running")
	}
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	tt.Click("Play / pause")
	tt.Frame()
	before = d.scene.Instances[0].Transform
	tt.Advance(time.Second)
	if before != d.scene.Instances[0].Transform || tt.FramePending() {
		t.Fatal("reduced motion ignored")
	}
}

func TestViewportKeyboardFocus(t *testing.T) {
	d := newDemo(scene3d.Cube())
	tt := ui.NewTester(d.view, 1200, 820)
	tt.Frame()
	if err := tt.Click("3D viewport"); err != nil {
		t.Fatal(err)
	}
	before := d.camera
	tt.Key(0, ui.KeyLeft)
	if before == d.camera {
		t.Fatal("viewport did not receive arrow key")
	}
}
