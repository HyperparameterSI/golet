// Author: HyperparameterSI

package main

const scene3DUsage = `// Author: HyperparameterSI
package main

import (
    "image/color"
    "log"

    "github.com/HyperparameterSI/golet/backend"
    "github.com/HyperparameterSI/golet/backend/plugins/scene3d"
    "github.com/HyperparameterSI/golet/backend/ui"
    "github.com/HyperparameterSI/golet/backend/ui/forma"
)

func main() {
    world := &scene3d.Scene{
        Background: color.RGBA{18, 24, 38, 255},
        Ambient: .25,
        Lights: []scene3d.DirectionalLight{{
            Direction: scene3d.Vec3{X: 2, Y: 4, Z: 3}, Intensity: .8,
        }},
        Instances: []scene3d.Instance{{
            Mesh: scene3d.Cube(),
            Transform: scene3d.Identity(),
            Material: scene3d.Material{Color: color.RGBA{99, 158, 241, 255}},
        }},
    }
    camera := scene3d.DefaultCamera()
    var revision uint64
    golet.App.WhenReady(func() {
        golet.NewWindow(golet.WindowOptions{
            Title: "Scene3D", Width: 1000, Height: 760,
            TitleBarStyle: golet.TitleBarHiddenInset,
            TrafficLightPosition: &golet.Point{X: 12, Y: 11},
            Content: ui.View(func(c *ui.Context) {
                ui.Column(c).Fill().Padding(20).Gap(12).Children(func() {
                    bar := c.TitleBar()
                    ui.Box(c).Height(max(bar.Height, 36)).DragWindow()
                    if forma.Button(c, "Rotate").Clicked() {
                        instance := &world.Instances[0]
                        instance.Transform = instance.Transform.Mul(scene3d.RotationY(.3))
                        revision++
                    }
                    scene3d.View(c, scene3d.ViewOptions{
                        Scene: world, Camera: camera, Revision: revision,
                        MaxDimension: 720,
                        OnFrame: func(_ scene3d.Stats, err error) {
                            if err != nil { log.Printf("scene3d: %v", err) }
                        },
                    }).Fill()
                })
            }),
        })
    })
    if err := golet.App.Run(); err != nil { log.Fatal(err) }
}`

const scene3DOBJUsage = `// Author: HyperparameterSI
// 以下代码在已有 world / revision 的 UI 事件中使用；补充导入 os。
// 大模型应在后台解析，再通过 Window.Update 交付到 UI 线程。
f, err := os.Open("model.obj")
if err != nil { log.Printf("open model: %v", err); return }
mesh, err := scene3d.LoadOBJ(f, scene3d.OBJLimits{
    Bytes: 16 << 20, Vertices: 250000, Triangles: 500000,
})
f.Close()
if err != nil { log.Printf("load OBJ: %v", err); return }
world.Instances[0].Mesh = mesh
revision++`

const scene3DAnimationUsage = `// Author: HyperparameterSI
// 补充导入 time；last 在窗口状态中保存，不要每帧重建。
var last time.Time
animate := func(now time.Time) {
    if !last.IsZero() {
        dt := float32(max(0, min(now.Sub(last).Seconds(), .1)))
        for i := range world.Instances {
            instance := &world.Instances[i]
            instance.Transform = instance.Transform.Mul(scene3d.RotationY(dt * .65))
        }
    }
    last = now
}
// 将 animate 传入 ViewOptions.Animate。暂停时传 nil，恢复前重置 last。
// 独立游戏逻辑需要固定步长；此回调只演示模型旋转。`

const scene3DOrbitUsage = `// Author: HyperparameterSI
// controls、camera 和 height 保存在窗口状态中，不要每帧重新创建。
var controls scene3d.OrbitControls
var height float32

// 在 ui.View 构建函数中：
scene3d.View(c, scene3d.ViewOptions{
    Scene: world, Camera: camera, Revision: revision,
    OnInput: func(ev ui.InputEvent) bool {
        return controls.Handle(&camera, ev, height)
    },
}).Fill().DrawOver(func(_ *ui.Painter, r ui.Rect) {
    height = r.H
})
// 重置视角或移除视口时调用 controls.Cancel()。
// 左键环绕；Shift+左键/中键平移；滚轮/捏合缩放。`
