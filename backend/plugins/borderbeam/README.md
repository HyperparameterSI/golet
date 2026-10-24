# Border Beam — 原生 Go

将 [Border Beam](https://border-beam.frad.me/) 的五种模式移植到 Golet 原生 UI，不使用 React、WebView、JavaScript 或运行时 WGSL。

```go
borderbeam.Box(c, borderbeam.Options{
    Preset: borderbeam.Medium,
    Palette: borderbeam.Colorful,
    Strength: 1,
    Radius: 16,
}).Padding(24).Background(c.Theme().Surface).Children(func() {
    ui.Text(c, "Give your ideas an edge.").Bold()
    ui.Button(c, "继续")
})
```

| Preset | 动画 | 默认周期 |
| --- | --- | --- |
| `Medium` | 局部色场、旋转锥形遮罩、内辉光 | 1.96 秒 |
| `Small` | 小半径色斑与独立遮罩，适合紧凑按钮 | 1.96 秒 |
| `Line` | 沿底边移动的光束、两端淡出、光柱 | 3.1 秒 |
| `PulseInner` | 色斑位移、宽高呼吸、分区透明度变化 | 2.3 秒 |
| `PulseOutside` | 边框呼吸与向外扩散的双层辉光 | 2.3 秒 |

`Palette` 支持 `Colorful`、`Ocean`、`Sunset`、`Mono`，以及跟随主题的默认 `Theme`；`Colors` 可以覆盖色板。`Duration` 调节运动周期；色相使用独立周期。`StaticColors` 只固定色相，不停止运动。`Brightness`、`Saturation`、`HueRange` 可调节颜色。

动画按 30 FPS 请求定时重绘，不创建后台 goroutine。`Paused` 保留相位，恢复后继续；`Disabled` 隐藏并停止请求帧；系统减少动态效果自动静止；不可见的控件不请求动画帧。装饰不拦截内容的点击。`PulseOutside` 应预留约 60 DIP 外部空间，父容器若裁剪会截掉外发光。

Studio → 插件 → 边框光束提供五种模式切换、四组配色、暂停/继续、隐藏、周期、强度和内辉光范围，并跟随 Studio 明暗主题。

## 来源与实现边界

`reference/beam-color.wgsl` 与 `reference/border-beam.css` 保存上游参考快照；`reference/LICENSE` 保留 Jakub Antalik 与 Frad 的 MIT 署名。`go generate ./backend/plugins/borderbeam` 用 Go 生成器从 WGSL 提取原始色斑与色板数值。参考文件不嵌入应用。

`field.go` 移植二维径向色场和脉冲振荡；`render.go` 实现各模式遮罩及分层混合；`blur.go` 在复用缓冲区中执行降采样高斯模糊。细边框单独保留高分辨率采样。应用不引入 CEF 或额外运行库。

这是原生渲染移植，不承诺与浏览器逐像素一致：颜色滤镜、降采样模糊和浏览器合成存在差异；当前没有移植 DOM 激活/消失时的 CSS 淡入淡出。原生窗口 GPU 端的实机视觉对比仍需单独验证。测试覆盖五种动画的运行、暂停、恢复、减少动态效果、离屏停止和内容交互。
