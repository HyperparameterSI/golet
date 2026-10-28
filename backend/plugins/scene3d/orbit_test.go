// Author: HyperparameterSI

package scene3d

import (
	"math"
	"testing"

	"github.com/HyperparameterSI/golet/backend/ui"
)

func TestOrbitDragAndRelease(t *testing.T) {
	camera := DefaultCamera()
	original := camera
	var controls OrbitControls
	down := ui.InputEvent{Kind: ui.InputPointerDown, Button: 0, X: 30, Y: 40}
	if !controls.Handle(&camera, down, 500) {
		t.Fatal("drag not captured")
	}
	controls.Handle(&camera, ui.InputEvent{Kind: ui.InputPointerMove, Button: 0, X: 130, Y: 100}, 500)
	if camera.Position == original.Position || camera.Position.Y == original.Position.Y || camera.Target != original.Target {
		t.Fatal("drag did not orbit in both axes")
	}
	beforeDistance := original.Position.Sub(original.Target).Dot(original.Position.Sub(original.Target))
	afterDistance := camera.Position.Sub(camera.Target).Dot(camera.Position.Sub(camera.Target))
	if math.Abs(float64(beforeDistance-afterDistance)) > .001 {
		t.Fatal("orbit changed radius")
	}
	controls.Handle(&camera, ui.InputEvent{Kind: ui.InputPointerUp, Button: 0}, 500)
	held := camera
	controls.Handle(&camera, ui.InputEvent{Kind: ui.InputPointerMove, Button: -1, X: 300, Y: 400}, 500)
	if held != camera {
		t.Fatal("hover moved released camera")
	}
}
func TestOrbitPanAndZoom(t *testing.T) {
	for _, button := range []int{0, 2} {
		camera := DefaultCamera()
		original := camera
		var controls OrbitControls
		controls.Handle(&camera, ui.InputEvent{Kind: ui.InputPointerDown, Button: button, Mods: ui.Shift, X: 30, Y: 40}, 500)
		controls.Handle(&camera, ui.InputEvent{Kind: ui.InputPointerMove, Button: button, X: 60, Y: 90}, 500)
		if camera.Target == original.Target {
			t.Fatal("pan did not move pivot")
		}
		a, b := camera.Position.Sub(original.Position), camera.Target.Sub(original.Target)
		if a.Sub(b).Dot(a.Sub(b)) > 1e-10 {
			t.Fatal("pan changed viewing direction")
		}
	}
	camera := DefaultCamera()
	var controls OrbitControls
	distance := func() float32 { d := camera.Position.Sub(camera.Target); return d.Dot(d) }
	before := distance()
	controls.Handle(&camera, ui.InputEvent{Kind: ui.InputScroll, DY: -100}, 500)
	if distance() >= before {
		t.Fatal("wheel did not zoom in")
	}
	before = distance()
	controls.Handle(&camera, ui.InputEvent{Kind: ui.InputZoom, Scale: .5}, 500)
	if distance() <= before {
		t.Fatal("pinch did not zoom out")
	}
	for i := 0; i < 100; i++ {
		controls.Handle(&camera, ui.InputEvent{Kind: ui.InputScroll, DY: -10000}, 500)
	}
	if distance() < .25*.25-1e-6 || camera.validate() != nil {
		t.Fatal("zoom crossed pivot")
	}
}
func TestOrbitPolesAndInvalidInput(t *testing.T) {
	camera := DefaultCamera()
	var controls OrbitControls
	for _, pitch := range []float32{1000, -1000} {
		controls.Rotate(&camera, 0, pitch)
		if camera.validate() != nil {
			t.Fatal("orbit crossed pole")
		}
	}
	before := camera
	for _, ev := range []ui.InputEvent{{Kind: ui.InputZoom, Scale: 0}, {Kind: ui.InputScroll, DY: float32(math.NaN())}} {
		if controls.Handle(&camera, ev, 500) {
			t.Fatal("invalid event consumed")
		}
	}
	if before != camera {
		t.Fatal("invalid event changed camera")
	}
	controls.Handle(&camera, ui.InputEvent{Kind: ui.InputPointerDown, Button: 0}, 500)
	controls.Cancel()
	controls.Handle(&camera, ui.InputEvent{Kind: ui.InputPointerMove, Button: 0, X: 10, Y: 10}, 500)
	if before != camera {
		t.Fatal("cancel did not end drag")
	}
}
