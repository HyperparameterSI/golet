# Scene3D 设计与性能验收

Author: HyperparameterSI

## 定位

为 Golet 增加独立可选、纯 Go 的三维组件。第一阶段提供真实三角形、透视摄像机、
裁剪与深度遮挡，以及 Forma 共用窗口的可运行基础；当前交付为 **CPU 软件参考
渲染器**。它适合简单模型预览和正确性验证，不能作为“高性能 GPU 引擎已完成”的证明。

[使用方法和当前能力](../../backend/plugins/scene3d/README.md)

## 开源设计调研

研究了以下官方资料；本实现没有复制这些项目的源码，也没有引入其运行时依赖。

| 项目 | 官方设计 | 对 Golet 的启发与限制 |
| --- | --- | --- |
| [Filament](https://github.com/google/filament) | 跨平台实时物理渲染；[PBR 说明](https://github.com/google/filament/blob/main/docs/Filament.md.html) | 材质和光照应有明确模型。完整 PBR、资源管线和原生运行时不能当作一个简单控件直接搬入；本阶段仅平面 Lambert。 |
| [bgfx](https://bkaradzic.github.io/bgfx/overview.html) | 渲染后端与应用框架解耦，覆盖 Metal、D3D、OpenGL 等 | 将场景描述与提交/资源后端分开。它提供 C/C++ API；直接接入不符合“实现全为 Go”的目标，现阶段只借鉴边界设计。 |
| [Babylon.js thin instances](https://doc.babylonjs.com/features/featuresDeepDive/mesh/copies/thinInstances/) | 共享网格和紧凑实例数据，同时明确批次可见性与动态增删的权衡 | 当前使用共享不可变 Mesh + 连续 Instance 切片，并逐实例包围盒剔除；这还不是 GPU instanced draw。 |

## 当前数据流

```text
调用者的 Scene / Camera / shared Mesh
  → instance bounds culling
  → world / camera vertex transforms (reused scratch)
  → back-face rejection / flat directional light
  → six-plane clipping
  → triangle coverage + reciprocal-depth test
  → reusable RGBA target
  → ui.Bitmap → Golet compositor → native window
                          ↑
                 Forma controls / overlays
```

插件仅依赖公开 ui API，无 init 注册、无核心 ui 导入插件，无 cgo、无默认后台线程。
未使用 scene3d 的应用不会链接其实现。独立软件渲染可并行使用不同 Target；UI 的
位图更新和显示仍遵循 Golet 主线程约束。错误由 Render 返回或由 View.OnFrame 交付，
框架不强制全局日志配置，也不吞掉失败。

## 后续原生 GPU 阶段（尚未实现）

1. 增加显式可释放的三维资源/视口接口：顶点、索引、实例缓冲、颜色与深度附件，
   可查询后端能力和内存统计。保持注册在插件构造路径，避免增加普通应用依赖。
2. 先实现一个平台并与软件参考帧比对，再扩展 Metal、D3D11、OpenGL/WebGL2。
   加入二维合成器共享的离屏纹理与资源屏障；避免逐帧 GPU→CPU→GPU 回读。
3. GPU 实例化、分批/状态排序、LOD、空间分区与遮挡优化。决定批次粒度时同时
   测量多余绘制与 CPU 剔除成本，不把“合并成一个 draw call”当作普遍最优。
4. 分离资源包：纹理/PBR、glTF、骨骼与动画、阴影按需启用，定义销毁与预算。
5. 多窗口/不同分辨率、设备丢失、缩放、关闭时资源释放，以及 Linux/Windows
   真机和浏览器上的像素、输入、性能验证。

目前 Renderer 接口只支持输出 CPU Image，不能绕过上述 GPU 合成工作。不能通过
把 CPU 运算移到后台、减少分辨率、关闭阴影等办法，声称同等画质性能提升。

## 性能目标的验收方式

“快 10 倍/100 倍”只可针对约定的场景与基线表达，不能对所有市售引擎承诺。
对比前固定引擎版本、硬件、驱动、电源状态、模型/纹理、镜头、分辨率、画质、
可见物体数量和预热流程；同一场景同等效果下记录：

- CPU update / prepare / submit 的中位数与 p95。
- GPU timestamp 帧时、draw call、实际可见三角形数量。
- Go 堆分配、进程常驻内存、GPU 资源与峰值内存。
- 静态空闲场景与持续动画场景分别报告；不能把“不重绘”与“全帧渲染”相除。
- 保存源码、原始结果和图像对照；测试断言检查工作计数/分配，不依赖墙钟阈值。

已提交基准只建立 Golet 软件实现自身的可复现基线；原生 GPU 消耗与外部引擎对比
尚未测量。实际实施顺序应先保正确性与资源边界，再用剖析选择可证实的热点。
