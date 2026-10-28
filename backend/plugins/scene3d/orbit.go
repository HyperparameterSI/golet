// Author: HyperparameterSI

package scene3d

import (
	"math"

	"github.com/HyperparameterSI/golet/backend/ui"
)

// OrbitControls provides event-driven, Y-up camera navigation. Keep one per
// viewport on its UI thread. Left drag orbits, Shift+left or middle drag pans,
// wheel and touchpad pinch dolly. It never schedules animation or allocates
// resources. Camera.Target is the orbit pivot and moves together with panning.
type OrbitControls struct {
	// RotationSpeed is radians per DIP, default .006.
	RotationSpeed float32
	// ZoomSpeed is the exponential wheel factor per DIP, default .002.
	ZoomSpeed float32
	// Distance bounds default to max(.25, Near*2) and Far*.9 respectively.
	MinDistance, MaxDistance float32
	dragging, panning        bool
	button                   int
	x, y                     float32
}

// Cancel releases drag state when removing a viewport or resetting its camera.
// Golet also delivers pointer-up when a captured window loses focus.
func (o *OrbitControls) Cancel() { o.dragging = false; o.panning = false }

// Rotate changes yaw and pitch in radians while retaining target and distance.
// Pitch stops just short of the poles, preventing camera flips. Returns false
// for an invalid camera or non-finite deltas. Navigation assumes world Y is up.
func (o *OrbitControls) Rotate(camera *Camera, yaw, pitch float32) bool {
	if camera == nil || camera.validate() != nil || !finite(yaw) || !finite(pitch) {
		return false
	}
	offset := camera.Position.Sub(camera.Target)
	radius := float32(math.Sqrt(float64(offset.Dot(offset))))
	azimuth := math.Atan2(float64(offset.X), float64(offset.Z)) + float64(yaw)
	elevation := math.Asin(float64(max(float32(-1), min(float32(1), offset.Y/radius)))) + float64(pitch)
	elevation = max(-math.Pi/2+.01, min(math.Pi/2-.01, elevation))
	horizontal := float64(radius) * math.Cos(elevation)
	camera.Position = camera.Target.Add(Vec3{float32(horizontal * math.Sin(azimuth)), radius * float32(math.Sin(elevation)), float32(horizontal * math.Cos(azimuth))})
	camera.Up = Vec3{Y: 1}
	return true
}
func (o *OrbitControls) dolly(camera *Camera, factor float32) bool {
	if !finite(factor) || factor <= 0 {
		return false
	}
	offset := camera.Position.Sub(camera.Target)
	distance := float32(math.Sqrt(float64(offset.Dot(offset))))
	lo, hi := o.MinDistance, o.MaxDistance
	if !finite(lo) || lo <= 0 {
		lo = max(float32(.25), camera.Near*2)
	}
	if !finite(hi) || hi <= 0 {
		hi = camera.Far * .9
	}
	if hi < lo {
		return false
	}
	wanted := max(lo, min(hi, distance*factor))
	camera.Position = camera.Target.Add(offset.Scale(wanted / distance))
	return true
}

// Handle consumes navigation events and updates camera. viewportHeight is the
// displayed height in DIPs, used to make panning independent of render scale.
// Forward events from ViewOptions.OnInput; Golet captures accepted drags and
// delivers moves/releases outside the viewport. Other inputs return false.
func (o *OrbitControls) Handle(camera *Camera, ev ui.InputEvent, viewportHeight float32) bool {
	if ev.Kind == ui.InputPointerUp {
		if o.dragging && ev.Button == o.button {
			o.Cancel()
			return true
		}
		return false
	}
	if camera == nil || camera.validate() != nil {
		o.Cancel()
		return false
	}
	switch ev.Kind {
	case ui.InputPointerDown:
		if (ev.Button != 0 && ev.Button != 2) || !finite(ev.X) || !finite(ev.Y) {
			return false
		}
		o.dragging = true
		o.panning = ev.Button == 2 || ev.Mods&ui.Shift != 0
		o.button = ev.Button
		o.x, o.y = ev.X, ev.Y
		return true
	case ui.InputPointerMove:
		if !o.dragging {
			return false
		}
		if ev.Button != o.button || !finite(ev.X) || !finite(ev.Y) {
			o.Cancel()
			return false
		}
		dx, dy := ev.X-o.x, ev.Y-o.y
		o.x, o.y = ev.X, ev.Y
		if !finite(dx) || !finite(dy) {
			o.Cancel()
			return false
		}
		if o.panning {
			if !finite(viewportHeight) || viewportHeight <= 0 {
				return false
			}
			forward := camera.Target.Sub(camera.Position)
			distance := float32(math.Sqrt(float64(forward.Dot(forward))))
			forward = forward.Unit()
			right := forward.Cross(camera.Up).Unit()
			up := right.Cross(forward)
			units := 2 * distance * float32(math.Tan(float64(camera.FOVY)/2)) / viewportHeight
			delta := right.Scale(-dx * units).Add(up.Scale(dy * units))
			position, target := camera.Position.Add(delta), camera.Target.Add(delta)
			if !position.valid() || !target.valid() {
				return false
			}
			camera.Position, camera.Target = position, target
			return true
		}
		speed := o.RotationSpeed
		if !finite(speed) || speed <= 0 {
			speed = .006
		}
		return o.Rotate(camera, -dx*speed, dy*speed)
	case ui.InputScroll:
		if !finite(ev.DY) || ev.DY == 0 {
			return false
		}
		speed := o.ZoomSpeed
		if !finite(speed) || speed <= 0 {
			speed = .002
		}
		return o.dolly(camera, float32(math.Exp(float64(max(float32(-20), min(float32(20), ev.DY*speed))))))
	case ui.InputZoom:
		if !finite(ev.Scale) || ev.Scale <= 0 {
			return false
		}
		return o.dolly(camera, 1/ev.Scale)
	}
	return false
}
