# Scene3D（实验性）

Author: HyperparameterSI

纯 Go、可选的三维视口，导入路径：
`github.com/HyperparameterSI/golet/backend/plugins/scene3d`。
不依赖浏览器、前端代码、cgo 或第三方引擎运行时。

**当前是软件光栅器，不是已完成的 GPU 游戏引擎。** Golet 的 GPU 合成器负责
显示最终图像；三角形变换、裁剪、光照和深度测试在 CPU 上执行。
没有“比其他引擎快 10～100 倍”的对比证据。

## 已实现

- 右手坐标系，透视摄像机，仿射模型变换，索引网格和共享几何实例。
- 六个视锥平面的三角形裁剪、包围盒视锥剔除、背面剔除和逐像素深度测试。
- 不透明纯色材质、双面/无光照模式、环境光和方向光（平面 Lambert 着色）。
- 有文件大小/顶点数/三角形数限制的 OBJ 加载；只接受已三角化的面。
- `scene3d.View` 参与原生 UI 的布局、裁剪、输入和图像合成，可与 Forma 混排。
- 可选动画回调，遵循减少动态效果偏好；默认无持续动画或后台计时器。
- 稳定尺寸下复用颜色、深度、顶点工作缓冲和 UI 位图；静态帧复用。
- 可替换 `Renderer`（当前接口输出 CPU 图像），可查询实际工作计数和缓冲容量。

## 运行示例

Studio 已集成「插件 → Scene3D · 三维视口」：运行 `pnpm dev`，打开该项即可体验
四种场景、参数调节和可复制的 Go 使用示例。其界面跟随 Studio 原生主题设置。

在仓库根目录：

```sh
go run ./backend/examples/scene3d
go run ./backend/examples/scene3d backend/examples/scene3d/assets/pyramid.obj
# 可选：加载自己的已三角化 OBJ；不会自动缩放或居中模型。
go run ./backend/examples/scene3d /path/to/model.obj
```

示例提供 Forma 按钮：环绕摄像机、旋转模型、播放/暂停、切换明暗主题。
点击三维视口获得键盘焦点，左右方向键环绕。仅模型与背景颜色跟随示例设置；
材质不会被组件主题隐式改写。默认窗口沿用 HiddenInset 和原生交通灯。

## 接入

场景、网格和摄像机只创建一次，在应用状态中保存；不要在每帧创建网格。

```go
// Author: HyperparameterSI
mesh := scene3d.Cube()
scene := &scene3d.Scene{
    Background: color.RGBA{18, 24, 38, 255},
    Ambient: .2,
    Lights: []scene3d.DirectionalLight{
        {Direction: scene3d.Vec3{X: 2, Y: 4, Z: 3}, Intensity: .8},
    },
    Instances: []scene3d.Instance{{
        Mesh: mesh,
        Material: scene3d.Material{Color: color.RGBA{100, 160, 240, 255}},
    }},
}
camera := scene3d.DefaultCamera()
var revision uint64

// 在 ui.View 的构建函数中，与 forma.Button 等组件放在同一布局中：
scene3d.View(c, scene3d.ViewOptions{
    Scene: scene, Camera: camera, Revision: revision,
    MaxDimension: 960,
    OnFrame: func(stats scene3d.Stats, err error) {
        if err != nil { log.Printf("scene3d: render failed: %v", err) }
    },
}).Fill()
```

修改实例、材质、灯光或背景后必须增加 `Revision`，通过正常的窗口更新流程
重绘。摄像机与尺寸自动参与缓存键。`Animate func(time.Time)` 在可见绘制帧中
更新场景并触发重绘，不需要手动增加 Revision；使用传入时间进行固定步长或
有界增量更新。暂停时将 Animate 设为 nil。示例演示了变换动画，不包含骨骼动画。

`Scene` 是调用者拥有的数据，渲染期间不能并发修改。`View`、位图与动画回调
在 UI 线程运行；后台加载完成后通过 `Window.Update` 交付场景。独立软件渲染可在
工作 goroutine 上进行，但每个工作线程要有自己的 `Target`，交付图像时要保证
其生命周期，不能同时改写正在显示的像素。

## 资源与扩展边界

### 鼠标摄像机控制

`OrbitControls` 是可复用的 Y-up 环绕控制器，不产生后台动画。
左键拖动控制水平与俯仰角（限制在极点以内），Shift＋左键或中键平移观察中心，
滚轮/捏合控制距离。支持自定义旋转/缩放速度与最小/最大观察距离。
在窗口状态中保留一个控制器，并在构建视口时接入：

```go
// Author: HyperparameterSI
// controls、camera、height 都保存在窗口状态中。
scene3d.View(c, scene3d.ViewOptions{
    Scene: world, Camera: camera, Revision: revision,
    OnInput: func(ev ui.InputEvent) bool {
        return controls.Handle(&camera, ev, height)
    },
}).Fill().DrawOver(func(_ *ui.Painter, r ui.Rect) { height = r.H })
```

控制器类型为 `scene3d.OrbitControls`，`height` 为 float32，使用显示高度而非渲染
像素，因此切换分辨率不改变平移灵敏度。Golet 将捕获后的移动、松开事件发送到
视口，即使指针已拖出视口；失焦也会发送松开。移除视口、切换页签或重置摄像机
时调用 `controls.Cancel()` 清除拖动状态。Studio 的使用示例页包含完整接入片段。

### 缓冲与渲染器

- `NewMesh` 复制输入并保持几何不可变；多个 Instance 可共享 Mesh。
- 每个像素需要 4 字节 RGBA 和 4 字节深度。640×360 两个缓冲合计约 1.76 MiB。
- `MaxPixels` 限制单目标为 4,194,304 像素（颜色+深度 32 MiB）。变换顶点另计。
- `Scale` 默认为每 DIP 一个渲染像素，不会自动按 Retina 放大。`MaxDimension`
  默认 1024、上限 2048，等比例限制分辨率；这是质量/资源的显式权衡。
- `Stats.BufferBytes` 不包含模型、Go 对象头、UI/GPU 纹理或驱动内存。
- `Renderer` 可替换图像生成算法；`NewRenderer` 每次挂载仅调用一次，更换需重挂载。
  它**不是**原生 GPU 纹理/命令队列接口。
- 直接使用 `Software.Render` 时，`Target.Image` 为借用数据，下次 Render 会覆盖；
  `Target.Release` 释放 CPU 引用。视口移除后其状态随 UI 生命周期回收。

尚未实现：原生 GPU 三维管线、GPU 实例化、纹理采样、PBR、透明排序、阴影、
抗锯齿、glTF/GLB、骨骼动画、物理、编辑器。OBJ 的 UV、法线和材质声明被忽略，
不打开 MTL 或纹理文件；不支持的几何记录和非三角面明确报错。

## 验证与性能

```sh
CGO_ENABLED=0 go test ./backend/plugins/scene3d ./backend/examples/scene3d
CGO_ENABLED=0 go test ./backend/plugins/scene3d -run '^$' -bench BenchmarkSoftware -benchmem -count=3
```

2026-10-28，Apple M5 / macOS arm64 / Go 1.27.1，三次中位数：

| 软件光栅器，640×360 | 每帧时间 | 预热后分配 |
| --- | ---: | ---: |
| 1 个共享网格实例（立方体） | 0.181 ms | 0 B / 0 次 |
| 100 个共享网格实例（立方体） | 0.369 ms | 0 B / 0 次 |

场景由基准源码固定，包含视锥剔除，并非所有三角形都覆盖屏幕。
只测 Render；不含 OBJ 加载、UI 构建、纹理上传、呈现或真实游戏逻辑。
不能换算成游戏 FPS，也不能据此声称超过其他引擎。GUI 交互由 ui.Tester 检查，
并不等于已经完成各平台真实显卡验证。

原始结果见 [基准输出](../../../docs/performance/scene3d-software-2026-10-28.txt)。
研究来源和 GPU 阶段设计见 [设计记录](../../../docs/plugins/scene3d.md)。
