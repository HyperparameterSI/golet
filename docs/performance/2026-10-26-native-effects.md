# 原生动画性能优化记录 · 2026-10-26

作者：HyperparameterSI

## 范围与结论

本次针对可选的 `plugins/borderbeam`；没有证据支持“整个 Golet 快 10 倍”。
暂停动画的重复窗口重绘超过 10 倍；实际动态帧主要改善 Line 和 PulseInner。
Small 和部分 Outside 场景没有稳定收益，表中的轻微下降如实保留。
不降低采样倍率、动画帧率、模糊精度，也不改变外观或停用动画。

## 测量方法

Apple M5、darwin/arm64、CGO_ENABLED=0；基线为 v0.1.11（19c1ca0）。
基线源码通过 Go `-overlay` 映射，使用相同的新基准测试文件。
每项固定 30 帧、重复 3 次，取中位数；动画每次推进 1/30 秒，
避免自动调整 b.N 导致两组测量经历不同动画阶段。
card 为 480×160，window 为 1200×800。使用 ui.NewTester 的 CPU 渲染器，
包含生成纹理、更新及合成；不代表 Metal、D3D、GPU 上传时间或整机 CPU 占用率。
PausedFrame 预热后重复绘制同一暂停效果，测量真实重复工作消除；
不是活动动画加速比。ns/op 转换为下面的 ms/frame。
这些是短时本机基准，存在系统负载波动，不能视作跨设备 SLA。

| 场景 | 优化前 ms/frame | 优化后 ms/frame | 前/后 | B/op 前 → 后 |
| --- | ---: | ---: | ---: | ---: |
| card/medium | 9.69549 | 9.05493 | 1.07× | 4320 → 4321 |
| card/small | 14.30649 | 14.67338 | 0.97× | 2272 → 2273 |
| card/line | 15.05144 | 6.80767 | 2.21× | 5568 → 5569 |
| card/inner | 15.80991 | 9.15528 | 1.73× | 5328 → 5329 |
| card/outside | 26.01655 | 26.61997 | 0.98× | 3392 → 3393 |
| window/medium | 32.16980 | 27.81585 | 1.16× | 4339 → 4340 |
| window/small | 82.69158 | 84.04425 | 0.98× | 2295 → 2292 |
| window/line | 76.02803 | 24.95838 | 3.05× | 5584 → 5585 |
| window/inner | 44.23305 | 27.51043 | 1.61× | 5351 → 5348 |
| window/outside | 75.90892 | 66.93010 | 1.13× | 3443 → 3441 |
| PausedFrame | 9.72741 | 0.00184 | 5278.03× | 4480 → 160 |

暂停场景分配次数由 15 次降到 1 次，分配字节由 4480 降到 160（28 倍减少）。
缓存复用已有的三层纹理，只额外保留最后一帧参数及颜色快照；不缓存动画历史帧。
活动动画的纹理内存规模未改变。B/op 是累计分配量，并非 RSS 或常驻内存。

## 实现

- 相同尺寸、时间、主题和选项时复用已有纹理，避免暂停/减少动态效果时重新生成。
- 颜色保存独立快照，调用方原地修改切片也会失效。
- 将 Line 相位相关参数移出逐像素循环。
- Line 仅扫描两个径向遮罩覆盖区域的保守外接矩形；透明区域保留为零。
- Medium / PulseInner 提前跳过不可能贡献像素的内部区域。
- 删除随后必定被覆盖的高光颜色场采样。

## 验证与复现

```sh
CGO_ENABLED=0 go test ./backend/plugins/borderbeam
CGO_ENABLED=0 go test ./backend/plugins/borderbeam -run '^$' \
  -bench 'Benchmark(AnimationFrame|PausedFrame)' -benchtime=30x -count=3 -benchmem
```

`render_reference_test.go` 保留发布前的渲染算法作为测试基准，不进入应用二进制。
逐像素比较覆盖全部五种模式、明暗主题、多个相位、普通及小数尺寸、大圆角，
同时比较锐利层与模糊缓冲区。缓存失效测试覆盖动画推进、缩放、明暗和选项变化。

## 其他发现

兼容 `TextArea` 的 `*string` 接口必须把编辑结果回写完整字符串；8 MiB 文本
连续 20 次小删除会产生大量整段复制。框架现有的 `TextBuffer` / `TextAreaBuffer`
可以避免此类复制，但这属于已有接口的选型收益，未计入本次优化。
后续整体性能目标应分别定义启动耗时、IPC、活动帧、空闲 CPU 和 RSS，
以实际应用端到端测量，不能用一个暂停动画的加速比替代。

## 未使用插件的应用体积

使用 `-trimpath -tags golet_noinspector -ldflags "-s -w -X github.com/HyperparameterSI/golet/backend.production=1"`：

| 应用 | 优化前 bytes | 优化后 bytes |
| --- | ---: | ---: |
| counter-native | 9266386 | 9266386 |
| hello（WebView） | 9422258 | 9422258 |

优化仍然是按需链接；未使用 Border Beam 的应用没有体积增长。
