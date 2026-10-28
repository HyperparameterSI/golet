// Author: HyperparameterSI

// Scene3d demonstrates opaque 3D geometry and Forma controls in one native
// window, with no browser or frontend runtime. Run from the repository root:
//
//	go run ./backend/examples/scene3d
package main

import (
	"image/color"
	"log"
	"math"
	"os"
	"time"

	"github.com/HyperparameterSI/golet/backend"
	"github.com/HyperparameterSI/golet/backend/plugins/scene3d"
	"github.com/HyperparameterSI/golet/backend/ui"
	"github.com/HyperparameterSI/golet/backend/ui/forma"
)

type demo struct {
	scene       scene3d.Scene
	camera      scene3d.Camera
	revision    uint64
	angle       float32
	dark        bool
	animate     bool
	last        time.Time
	stats       scene3d.Stats
	renderError error
}

func newDemo(mesh *scene3d.Mesh) *demo {
	d := &demo{camera: scene3d.DefaultCamera(), dark: true}
	d.scene = scene3d.Scene{Background: color.RGBA{18, 24, 38, 255}, Ambient: .2, Lights: []scene3d.DirectionalLight{{Direction: scene3d.Vec3{X: 2, Y: 4, Z: 3}, Intensity: .85}}}
	for i, col := range []color.RGBA{{104, 161, 244, 255}, {173, 132, 239, 255}, {62, 199, 178, 255}} {
		d.scene.Instances = append(d.scene.Instances, scene3d.Instance{Mesh: mesh, Transform: scene3d.Translation(scene3d.Vec3{X: float32(i-1) * 1.4}), Material: scene3d.Material{Color: col}})
	}
	d.camera.Position = scene3d.Vec3{X: 4, Y: 3, Z: 6}
	return d
}
func (d *demo) orbit(delta float32) {
	d.angle += delta
	d.camera.Position = scene3d.Vec3{X: float32(math.Sin(float64(d.angle))) * 7, Y: 3, Z: float32(math.Cos(float64(d.angle))) * 7}
}
func (d *demo) view(c *ui.Context) {
	c.SetTheme(forma.Theme(d.dark))
	ui.Column(c).Fill().Padding(20).Gap(12).Children(func() {
		bar := c.TitleBar()
		ui.Row(c).Height(max(bar.Height, 36)).Padding(0, bar.Right, 0, bar.Left).DragWindow().Children(func() {
			ui.Text(c, "Scene3D · Native lab").FontSize(24).Bold()
		})
		ui.Text(c, "Pure Go · Software reference renderer · Arrow keys orbit").TextColor(c.Theme().TextMuted)
		ui.Row(c).Gap(8).Children(func() {
			if forma.Button(c, "Orbit left").Clicked() {
				d.orbit(-.25)
			}
			if forma.Button(c, "Orbit right").Clicked() {
				d.orbit(.25)
			}
			if forma.Button(c, "Rotate meshes").Clicked() {
				for i := range d.scene.Instances {
					in := &d.scene.Instances[i]
					in.Transform = in.Transform.Mul(scene3d.RotationY(.3))
				}
				d.revision++
			}
			if forma.Button(c, "Play / pause").Clicked() {
				d.animate = !d.animate
				d.last = time.Time{}
			}
			if forma.Button(c, "Light / dark").Clicked() {
				d.dark = !d.dark
				if d.dark {
					d.scene.Background = color.RGBA{18, 24, 38, 255}
				} else {
					d.scene.Background = color.RGBA{235, 239, 246, 255}
				}
				d.revision++
			}
		})
		var animate func(time.Time)
		if d.animate && !c.Preferences().ReduceMotion {
			animate = func(now time.Time) {
				if !d.last.IsZero() {
					dt := max(0, min(now.Sub(d.last).Seconds(), .1))
					for i := range d.scene.Instances {
						in := &d.scene.Instances[i]
						in.Transform = in.Transform.Mul(scene3d.RotationY(float32(dt) * .6))
					}
				}
				d.last = now
			}
		} else {
			d.last = time.Time{}
		}
		scene3d.View(c, scene3d.ViewOptions{Scene: &d.scene, Camera: d.camera, Revision: d.revision, Animate: animate, MaxDimension: 960, OnInput: func(ev ui.InputEvent) bool {
			if ev.Kind != ui.InputKeyDown {
				return false
			}
			switch ev.Key {
			case ui.KeyLeft:
				d.orbit(-.1)
			case ui.KeyRight:
				d.orbit(.1)
			default:
				return false
			}
			return true
		}, OnFrame: func(s scene3d.Stats, err error) {
			d.stats, d.renderError = s, err
			if err != nil {
				log.Printf("scene3d: render failed: %v", err)
			}
		}}).Fill().MinHeight(240)
		ui.Textf(c, "Previous frame · %d triangles · %d culled instances · %.2f MiB CPU buffers", d.stats.RasterizedTriangles, d.stats.CulledInstances, float64(d.stats.BufferBytes)/(1024*1024)).FontSize(12)
		if d.renderError != nil {
			ui.Text(c, d.renderError.Error())
		}
	})
}
func main() {
	mesh := scene3d.Cube()
	if len(os.Args) > 1 {
		f, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		mesh, err = scene3d.LoadOBJ(f, scene3d.OBJLimits{})
		f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}
	d := newDemo(mesh)
	golet.App.WhenReady(func() {
		golet.NewWindow(golet.WindowOptions{Title: "Golet Scene3D", Width: 1200, Height: 820, MinWidth: 860, MinHeight: 620, TitleBarStyle: golet.TitleBarHiddenInset, TrafficLightPosition: &golet.Point{X: 12, Y: 11}, Content: ui.View(d.view)})
	})
	if err := golet.App.Run(); err != nil {
		log.Fatal(err)
	}
}
