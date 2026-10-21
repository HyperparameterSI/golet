# Golet Observatory

Author: HyperparameterSI

Go 进程性能监控示例。后端位于本目录，React + TanStack Start + shadcn 前端位于 `frontend/examples/dashboard`。

从仓库根目录运行：

```sh
pnpm install
pnpm dev:dashboard
pnpm build:dashboard
```

开发端口为 5174，不占用 Studio 的 5173。开发 `.app` 需要通过 `golet dev` 启动；直接打开独立应用请使用 `build/darwin-arm64/Golet Observatory.app`。

## 指标与口径

- 每秒后台采样，最多保留 300 点；UI 可选 1/2/5 秒刷新、1/3/5 分钟窗口、暂停显示和 JSON 导出。隐藏窗口停止前端轮询，后端仍保留最新历史。没有模拟数据。
- CPU：macOS/Linux 使用 `getrusage`，Windows 使用 `GetProcessTimes`，计算相邻样本的用户态 + 内核态 CPU 时间差 / 真实时间差。单核满载为 100%，多核可能超过 100%。首个样本及不支持的平台显示不可用。
- 内存：Go `runtime.ReadMemStats` 的活跃堆、占用堆页、栈、已归还 OS 的内存、对象数、累计分配及分配速率。运行时保留内存为 `Sys - HeapReleased`，**不是进程 RSS，也不含 WebView 子进程或原生库的全部内存**。
- GC：累计次数、最近暂停、最近至多 256 次暂停 P95、累计 STW 暂停、生命周期 GC CPU 占比及下次目标堆。事件按秒聚合，暂停值是该秒内最后一次 GC，不是整秒总和。
- Goroutine、运行时线程 profile 记录数（不是 OS 当前线程数）、Go 版本、PID、OS/架构、逻辑 CPU 数、GOMAXPROCS、运行时长、采样开销与 IPC 往返耗时。
- 健康检查阈值：CPU 80%（单核口径）、GC P95 10 ms、Goroutine 1,000、采样耗时 10 ms。只表示观察阈值，不能据此断言内存泄漏或性能故障。

采样在独立 goroutine 中运行，固定容量环形缓冲区以 O(1) 代价保存新样本。首次请求返回历史，之后通过单调序号仅传输新增样本；前端自行合并并限制为 300 点，请求完成后才调度下一次轮询，避免慢请求叠加。GC P95 只在发生 GC 后重新计算。采样器自身仍会产生 CPU 开销，不能当作零成本监控。导出包含进程元信息和最多 300 个样本，不包含堆对象内容。
