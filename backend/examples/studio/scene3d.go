// Author: HyperparameterSI

package main

import (
	"fmt"
	"image/color"
	"math"
	"strings"
	"time"

	"github.com/HyperparameterSI/golet/backend"
	"github.com/HyperparameterSI/golet/backend/examples/internal/demotheme"
	"github.com/HyperparameterSI/golet/backend/plugins/scene3d"
	"github.com/HyperparameterSI/golet/backend/ui"
	"github.com/HyperparameterSI/golet/backend/ui/forma"
)

// scene3DDemo owns its meshes once per window. All changes run on the UI thread.
type scene3DDemo struct {
	scene          scene3d.Scene
	camera         scene3d.Camera
	cube, pyramid  *scene3d.Mesh
	tab, mode      int
	resolution     string
	playing, unlit bool
	light          float64
	controls       scene3d.OrbitControls
	viewportHeight float32
	revision       uint64
	last           time.Time
	stats          scene3d.Stats
	renderError    error
	background     color.RGBA
}

func openScene3D() (*golet.Window, error) {
	d, err := newScene3DDemo()
	if err != nil {
		return nil, err
	}
	return golet.NewWindow(golet.WindowOptions{
		Title: "Studio · Scene3D", Width: 1180, Height: 900, MinWidth: 960, MinHeight: 760,
		TitleBarStyle: golet.TitleBarHiddenInset, TrafficLightPosition: &golet.Point{X: 12, Y: 11},
		Content: ui.View(d.view),
	}), nil
}
func newScene3DDemo() (*scene3DDemo, error) {
	pyramid, err := scene3d.LoadOBJ(strings.NewReader(scene3DPyramidOBJ), scene3d.OBJLimits{})
	if err != nil {
		return nil, fmt.Errorf("Studio Scene3D: load bundled OBJ: %w", err)
	}
	d := &scene3DDemo{cube: scene3d.Cube(), pyramid: pyramid, playing: true, light: .85, resolution: "720 px"}
	d.scene.Ambient = .2
	d.scene.Lights = []scene3d.DirectionalLight{{Direction: scene3d.Vec3{X: 2, Y: 4, Z: 3}, Intensity: .85}}
	d.selectMode(0)
	return d, nil
}
func (d *scene3DDemo) selectMode(mode int) {
	d.mode = mode
	d.last = time.Time{}
	d.scene.Instances = nil
	mesh := d.cube
	if mode == 1 {
		mesh = d.pyramid
	}
	count := 3
	if mode == 2 {
		count = 100
	}
	colors := []color.RGBA{{99, 158, 241, 255}, {169, 125, 232, 255}, {48, 187, 166, 255}}
	for i := 0; i < count; i++ {
		position := scene3d.Vec3{X: float32(i-1) * 1.5}
		scale := float32(1)
		if mode == 2 {
			position = scene3d.Vec3{X: float32(i%10) - 4.5, Z: float32(i/10) - 4.5}
			scale = .55
		}
		if mode == 3 {
			position = scene3d.Vec3{X: float32(i-1) * .35, Z: float32(i-1) * .9}
			scale = 1.4
		}
		d.scene.Instances = append(d.scene.Instances, scene3d.Instance{Mesh: mesh, Transform: scene3d.Translation(position).Mul(scene3d.Scaling(scene3d.Vec3{X: scale, Y: scale, Z: scale})), Material: scene3d.Material{Color: colors[i%len(colors)], Unlit: d.unlit}})
	}
	d.resetCamera()
	d.revision++
}
func (d *scene3DDemo) resetCamera() {
	d.camera = scene3d.DefaultCamera()
	d.controls.Cancel()
	radius, height := float32(6), float32(2.8)
	if d.mode == 2 {
		radius, height = 13, 9
	}
	d.camera.Position = scene3d.Vec3{X: float32(math.Sin(.5)) * radius, Y: height, Z: float32(math.Cos(.5)) * radius}
}
func (d *scene3DDemo) orbit(delta float32) {
	d.controls.Rotate(&d.camera, delta, 0)
}
func (d *scene3DDemo) animate(now time.Time) {
	if !d.last.IsZero() {
		dt := float32(max(0, min(now.Sub(d.last).Seconds(), .1)))
		rotation := scene3d.RotationY(dt * .65)
		for i := range d.scene.Instances {
			in := &d.scene.Instances[i]
			in.Transform = in.Transform.Mul(rotation)
		}
	}
	d.last = now
}
func (d *scene3DDemo) view(c *ui.Context) {
	demotheme.Apply(c)
	t := c.Theme()
	stage := t.Surface.Mix(t.Text, .035)
	bg := color.RGBA{stage.R, stage.G, stage.B, 255}
	if bg != d.background {
		d.background = bg
		d.scene.Background = bg
		d.revision++
	}
	ui.Column(c).Fill().Padding(20).Gap(12).Children(func() {
		bar := c.TitleBar()
		ui.Row(c).Height(max(bar.Height, 36)).Padding(0, bar.Right, 0, bar.Left).DragWindow().Children(func() { ui.Text(c, "Scene3D · 三维视口").FontSize(25).Bold() })
		ui.Text(c, "纯 Go 软件渲染 · 真实深度遮挡 · Forma 同窗口交互").TextColor(t.TextMuted)
		forma.Tabs(c, &d.tab, "效果演示", "使用示例")
		if d.tab == 1 {
			d.controls.Cancel()
			d.last = time.Time{}
			d.usage(c)
			return
		}
		ui.Row(c).Gap(8).Children(func() {
			for i, name := range []string{"立方体", "OBJ 金字塔", "共享网格 × 100", "深度遮挡"} {
				variant := forma.Outline
				if i == d.mode {
					variant = forma.Default
				}
				if forma.Button(c, name, forma.ButtonOptions{Variant: variant}).Clicked() {
					d.selectMode(i)
				}
			}
		})
		ui.Row(c).Gap(8).Children(func() {
			label := "暂停动画"
			if !d.playing {
				label = "继续播放"
			}
			if forma.Button(c, label).Clicked() {
				d.playing = !d.playing
				d.last = time.Time{}
			}
			if forma.Button(c, "向左环绕").Clicked() {
				d.orbit(-.25)
			}
			if forma.Button(c, "向右环绕").Clicked() {
				d.orbit(.25)
			}
			if forma.Button(c, "重置视角").Clicked() {
				d.resetCamera()
			}
			ui.Text(c, "分辨率上限").TextColor(t.TextMuted)
			forma.Select(c, &d.resolution, []string{"480 px", "720 px", "960 px"}).Width(120).Label("Scene3D 分辨率")
		})
		ui.Row(c).Gap(12).AlignItems(ui.Center).Children(func() {
			ui.Text(c, "方向光")
			before := d.light
			forma.Slider(c, &d.light, 0, 1.5).Width(180).Label("方向光强度")
			if before != d.light {
				d.scene.Lights[0].Intensity = float32(d.light)
				d.revision++
			}
			was := d.unlit
			forma.Switch(c, &d.unlit).Label("无光照材质")
			ui.Text(c, "无光照材质")
			if was != d.unlit {
				for i := range d.scene.Instances {
					d.scene.Instances[i].Material.Unlit = d.unlit
				}
				d.revision++
			}
			if c.Preferences().ReduceMotion {
				ui.Text(c, "减少动态效果：静态展示").TextColor(t.TextMuted)
			}
		})
		var animate func(time.Time)
		if d.playing && !c.Preferences().ReduceMotion {
			animate = d.animate
		} else {
			d.last = time.Time{}
		}
		limit := 720
		if d.resolution == "480 px" {
			limit = 480
		} else if d.resolution == "960 px" {
			limit = 960
		}
		scene3d.View(c.Key("scene3d-viewport"), scene3d.ViewOptions{
			Scene: &d.scene, Camera: d.camera, Revision: d.revision, MaxDimension: limit, Animate: animate,
			OnInput: func(ev ui.InputEvent) bool {
				if d.controls.Handle(&d.camera, ev, d.viewportHeight) {
					return true
				}
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
			},
			OnFrame: func(stats scene3d.Stats, err error) { d.stats, d.renderError = stats, err },
		}).Fill().MinHeight(220).Radius(t.Radius).Clip().Label("Scene3D 效果视口").DrawOver(func(_ *ui.Painter, r ui.Rect) { d.viewportHeight = r.H })
		// Paint-time text stays current when animation redraws without rebuilding UI.
		ui.Box(c).Height(24).Draw(func(p *ui.Painter, r ui.Rect) {
			text := fmt.Sprintf("实例 %d  ·  绘制三角形 %d  ·  剔除实例 %d  ·  CPU 缓冲 %.2f MiB", d.stats.Instances, d.stats.RasterizedTriangles, d.stats.CulledInstances, float64(d.stats.BufferBytes)/(1024*1024))
			if d.renderError != nil {
				text = "渲染失败：" + d.renderError.Error()
			}
			p.Text(r.X, r.Y, text, 13, t.TextMuted)
		})
		ui.Text(c, "左键拖动：环绕 / 俯仰 · Shift＋左键或中键：平移 · 滚轮 / 捏合：缩放 · ← / →：环绕").FontSize(12).TextColor(t.TextMuted)
	})
}
func (d *scene3DDemo) usage(c *ui.Context) {
	ui.Scroll(c).Fill().Gap(16).Children(func() {
		ui.Text(c, "01 · 最小可运行应用").FontSize(18).Bold()
		ui.Text(c, "保存为 main.go，在已依赖 Golet 的 Go 模块中运行 go run .。网格与场景在窗口外创建，按钮通过 Revision 更新视口。")
		forma.CodeBlock(c, "Go · 场景 + Forma + 视口", scene3DUsage)
		ui.Text(c, "02 · OBJ 模型加载").FontSize(18).Bold()
		ui.Text(c, "效果页的金字塔通过 LoadOBJ 解析下方模型。仅支持三角面；不会加载 MTL、纹理或外部文件。")
		forma.CodeBlock(c, "Go · 替换已有实例的网格", scene3DOBJUsage)
		forma.CodeBlock(c, "OBJ · 内置金字塔", scene3DPyramidOBJ)
		ui.Text(c, "03 · 动画与资源预算").FontSize(18).Bold()
		ui.Text(c, "Animate 在绘制帧执行，设为 nil 即暂停；系统减少动态效果会自动停用。静态场景修改材质、灯光或实例后增加 Revision。摄像机、尺寸变化自动刷新。")
		forma.CodeBlock(c, "Go · ViewOptions.Animate 回调", scene3DAnimationUsage)
		ui.Text(c, "04 · 鼠标控制摄像机").FontSize(18).Bold()
		ui.Text(c, "OrbitControls 保存在窗口状态中，处理鼠标捕获、俯仰限制、缩放和平移。平移使用显示高度，与渲染分辨率无关。")
		forma.CodeBlock(c, "Go · OrbitControls", scene3DOrbitUsage)
		ui.Text(c, "MaxDimension 控制渲染尺寸；共享 Mesh 避免复制几何。Renderer 可替换 CPU 图像生成器，目前没有原生 GPU 三维管线，不能将这些统计当作与其他引擎的性能对比。")
		ui.Text(c, "完整实现：backend/examples/studio/scene3d.go\n接口说明：backend/plugins/scene3d/README.md").Selectable()
	})
}

const scene3DPyramidOBJ = `# Author: HyperparameterSI
v -0.7 0 -0.7
v 0.7 0 -0.7
v 0.7 0 0.7
v -0.7 0 0.7
v 0 1 0
f 1 2 3
f 1 3 4
f 1 5 2
f 2 5 3
f 3 5 4
f 4 5 1
`
