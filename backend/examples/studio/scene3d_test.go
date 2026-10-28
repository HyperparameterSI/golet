// Author: HyperparameterSI

package main

import (
	"bytes"
	"go/parser"
	"go/token"
	"testing"
	"time"

	"github.com/HyperparameterSI/golet/backend/examples/internal/demotheme"
	"github.com/HyperparameterSI/golet/backend/plugins/scene3d"
	"github.com/HyperparameterSI/golet/backend/ui"
)

func TestScene3DStudioInteractions(t *testing.T) {
	d, err := newScene3DDemo()
	if err != nil {
		t.Fatal(err)
	}
	tt := ui.NewTester(d.view, 1180, 900)
	tt.Frame()
	if d.renderError != nil || d.stats.ShadedPixels == 0 {
		t.Fatalf("initial render: %v / %+v", d.renderError, d.stats)
	}
	before := d.scene.Instances[0].Transform
	tt.Advance(time.Second / 10)
	if before == d.scene.Instances[0].Transform {
		t.Fatal("visible viewport did not animate")
	}
	click := func(label string) {
		t.Helper()
		if err := tt.Click(label); err != nil {
			t.Fatal(err)
		}
	}
	click("暂停动画")
	held := d.scene.Instances[0].Transform
	tt.Advance(time.Second)
	if held != d.scene.Instances[0].Transform || tt.FramePending() {
		t.Fatal("paused viewport keeps animating")
	}
	for i, label := range []string{"立方体", "OBJ 金字塔", "共享网格 × 100", "深度遮挡"} {
		click(label)
		if d.mode != i || d.renderError != nil || d.stats.ShadedPixels == 0 {
			t.Fatalf("mode %s failed: %v", label, d.renderError)
		}
		if i == 1 && d.scene.Instances[0].Mesh != d.pyramid {
			t.Fatal("OBJ mesh not selected")
		}
		if i == 2 && (d.stats.Instances != 100 || d.scene.Instances[0].Mesh != d.scene.Instances[99].Mesh) {
			t.Fatal("instances do not share geometry")
		}
	}
	click("Scene3D 效果视口")
	camera := d.camera
	tt.Key(0, ui.KeyRight)
	if camera == d.camera {
		t.Fatal("viewport keyboard focus failed")
	}
	click("Scene3D 分辨率")
	tt.Key(0, ui.KeyHome)
	tt.Key(0, ui.KeyEnter)
	if d.resolution != "480 px" {
		t.Fatalf("resolution unchanged: %s", d.resolution)
	}
	click("无光照材质")
	if !d.scene.Instances[0].Material.Unlit {
		t.Fatal("material toggle ignored")
	}
	click("继续播放")
	tt.Advance(time.Second / 10)
	click("使用示例")
	if !tt.HasText("01 · 最小可运行应用") {
		t.Fatal("usage page missing")
	}
	held = d.scene.Instances[0].Transform
	tt.Advance(time.Second)
	if held != d.scene.Instances[0].Transform || tt.FramePending() {
		t.Fatal("hidden viewport keeps animating")
	}
	click("效果演示")
	tt.Advance(time.Second / 10)
	if held == d.scene.Instances[0].Transform {
		t.Fatal("return to demo did not resume")
	}
}
func TestScene3DStudioThemeAndReducedMotion(t *testing.T) {
	glass, appearance := demotheme.Selection()
	defer demotheme.Set(glass, appearance)
	demotheme.Set(false, 1)
	d, err := newScene3DDemo()
	if err != nil {
		t.Fatal(err)
	}
	tt := ui.NewTester(d.view, 960, 760)
	tt.Frame()
	light := d.scene.Background
	demotheme.Set(true, 2)
	tt.Frame()
	if light == d.scene.Background {
		t.Fatal("shared theme did not reach viewport")
	}
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	tt.Frame()
	held := d.scene.Instances[0].Transform
	tt.Advance(time.Second)
	if held != d.scene.Instances[0].Transform || tt.FramePending() {
		t.Fatal("reduced motion ignored")
	}
}
func TestScene3DUsageSyntax(t *testing.T) {
	if _, err := parser.ParseFile(token.NewFileSet(), "main.go", scene3DUsage, parser.AllErrors); err != nil {
		t.Fatal(err)
	}
}

// Exercise pointer capture through Golet's real headless input router, including
// dragging outside the viewport and synthesised release on focus loss.
func TestScene3DStudioMouseOrbit(t *testing.T) {
	d, err := newScene3DDemo()
	if err != nil {
		t.Fatal(err)
	}
	d.playing = false
	d.selectMode(2)
	tt := ui.NewTester(d.view, 1180, 900)
	tt.Frame()
	r, ok := tt.Find("Scene3D 效果视口")
	if !ok {
		t.Fatal("viewport missing")
	}
	x, y := r.X+r.W/2, r.Y+r.H/2
	before := d.camera
	imageBefore := append([]byte(nil), tt.Image().Pix...)
	tt.Press(x, y)
	tt.Move(x+80, y+60)
	if before.Position.Y == d.camera.Position.Y || before.Position.X == d.camera.Position.X {
		t.Fatal("drag did not change both viewing angles")
	}
	if bytes.Equal(imageBefore, tt.Image().Pix) {
		t.Fatal("camera moved but displayed frame did not change")
	}
	before = d.camera
	tt.Move(-20, y)
	if before == d.camera {
		t.Fatal("drag outside viewport lost capture")
	}
	tt.Release(-20, y)
	before = d.camera
	tt.Move(x, y)
	if before != d.camera {
		t.Fatal("camera followed hover after release")
	}
	tt.Scroll(x, y, 0, -120)
	if before.Position == d.camera.Position {
		t.Fatal("wheel did not dolly camera")
	}
	tt.Press(x, y)
	tt.Move(x+30, y+20)
	tt.SetFocused(false)
	before = d.camera
	tt.SetFocused(true)
	tt.Move(x+100, y+100)
	if before != d.camera {
		t.Fatal("drag stuck after focus loss")
	}
	if tt.FramePending() {
		t.Fatal("mouse controls leave an animation loop running")
	}
	if err := tt.Click("重置视角"); err != nil {
		t.Fatal(err)
	}
	if d.camera.Target != (scene3d.Vec3{}) {
		t.Fatal("reset did not restore pivot")
	}
}
