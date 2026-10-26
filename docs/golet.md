# Golet 迁移与验证

## 当前目录重构（2026-10-21）

Go 实现在 `backend/`，Web 源码在 `frontend/`；根 `go.mod` 保留。公共 Go 导入路径已变为 `github.com/HyperparameterSI/golet/backend`。本页下方初始迁移数据是重构前记录，当前路径以根 README 和 architecture.md 为准。

- Studio 与默认 shadcn 模板使用 TanStack Start SPA、React、Vite 和 Tailwind/shadcn；默认模板生成 `backend/main.go` 与 `frontend/src/`。
- Start 在 `frontend` 工作目录构建，输出 `dist/client/index.html`；Go 只嵌入 client 产物。
- 完整 Go 测试、201 项 Bun 测试、根与文档站类型检查通过；三平台编译、独立新脚手架构建、文档站全量预渲染通过。
- macOS 打包应用在 `golet://localhost/` 加载 TanStack 路由，实际 Go 调用成功；Windows/Linux 未执行 GUI 实机验收。
- 移动公共包后发布示例体积：counter-native 9,232,354 B，hello 9,337,362 B，各增加 33,024 B（约 0.36%）。Mach-O 检查显示原生示例 `__text` 保持 3,820,820 B，主要增加在 `__gopclntab`（+31,472 B）及类型元数据，来自更长的包路径；未引入新的运行时代码。


验证日期：2026-10-21。宿主：macOS / Apple M5 arm64，Go 1.27.1，Bun 1.4.3。

## 来源与范围

基于本地 `../mygo` 的已跟踪文件迁移，参考提交 `5fb60098a44a4c629cbdeac5f0ff80ce2f843c2f`。保留框架核心、全部平台后端、原生 UI、插件、CLI、文档、测试和示例。未复制源仓库 Git 历史或依赖缓存。

Go module、包名、CLI 路径、配置文件、环境变量、build tags、IPC 协议与生成客户端统一为 golet。npm 包使用 `golet-runtime`、`golet-cli`、`@golet/cli-*` 和 `@golet-plugins/*`。这是新的命名空间，旧项目应重新生成客户端并更新配置，不能混用两种桥接运行时。

必要保留项：原作者许可证、NOTICE、来自独立工具的 DMG 测试夹具名称、可选原生库的上游下载地址，以及这些已发布 SQLite 二进制导出的 `mygo_sqlite_*` ABI。改变 ABI 名称但继续下载旧二进制会使插件无法加载，因此这些名称不是公共 Golet API。

## 新增与改进

- 默认 `init` 模板为 `shadcn`，保留 `web` 和 `native` 模板。
- React + TypeScript + Tailwind v4 + Vite；本地 Button、Card、Input，CVA、Radix Slot、可扩展的 `components.json`、`@/` 别名。
- 明暗主题、键盘焦点、禁用/加载状态、错误提示、ARIA live region，以及实际 Go 绑定的调用示例。
- shadcn 模板不再创建一秒一次的时钟 IPC 事件。
- 原生 `ui.FormaTheme(bool)`，native 模板缓存两套主题，按系统明暗模式使用；不引入 WebView 或新的运行时依赖。
- 监听器移除、快照保留列表原地压缩，清除尾部函数引用，保持快照独立、once 语义、可重入回调和并发安全。
- 删除继承的网站第三方统计脚本。网站 canonical URL 为部署配置，当前未发布网站。
- CI 增加 Studio 前端绑定生成与生产构建。

## 性能证据

同机 `BenchmarkListenerSnapshot`，16 个常驻监听器，三次运行：

| 指标 | 迁移基线 | 优化后 |
| --- | --- | --- |
| ns/op | 162.9 / 161.1 / 158.8 | 51.60 / 53.45 / 48.18 |
| B/op | 872 | 128 |
| allocs/op | 6 | 1 |

分配字节减少约 85%，分配次数减少约 83%。这只是事件监听器分发路径的微基准；不能推导为整个桌面应用性能提升。耗时随机器负载变化，内存分配计数更稳定。

发布参数：`CGO_ENABLED=0 go build -trimpath -tags golet_noinspector -ldflags '-s -w -X github.com/HyperparameterSI/golet/backend.production=1'`。

| 示例 | 修改前字节 | 修改后字节 |
| --- | ---: | ---: |
| counter-native | 9,199,330 | 9,199,330 |
| hello（WebView） | 9,304,338 | 9,304,338 |

未使用 shadcn 主题的应用体积没有增长。

## 验证

已通过：

- 完整 `CGO_ENABLED=0 go test ./...`（包括核心、CLI、UI 和插件）。
- macOS、Linux、Windows 的 `CGO_ENABLED=0 go build ./...` 与 `go vet ./...`；生产 build tag 的 vet。
- WebAssembly（js/wasm）全部非 CLI 包编译通过。
- `go test -race -run TestListener .`，快照修改、重入 once 和并发订阅回归测试。
- Bun 201 项测试全部通过；工作区 TypeScript 类型检查；bridge/runtime/plugin 构建。
- `golet init -golet <checkout> -template shadcn` 生成真实临时项目、安装依赖、生成 Go 客户端、TypeScript 检查和 Vite 生产构建。
- Studio 经 CLI 构建为独立 `.app`，实际窗口加载 `golet://localhost/`，显示 Go 运行时信息，输入 Ada 返回 `Hello, Ada! This message comes from Go.`，暗色切换正常。

Linux / Windows 这里验证的是交叉编译和静态检查，未在这两种系统上执行 GUI 验收；原有 CI 保留相应原生测试。未执行完整性能/耗电/长期稳定性评估，不声称所有平台和功能均达到生产认证。

## 本地运行与发布

参见根 README。当前 npm 包、GitHub Releases、网站均未由此任务发布；第一次发布前须配置包所有权和发布凭证，确认网站部署地址。可选插件二进制仍从 MyGo 的固定版本下载并检查 SHA-256，后续自建产物时必须同时更新下载地址、摘要和 ABI。

`examples/studio/build/` 与 npm `dist/` 都是可重建产物，不提交。Golet 尚无 Git 提交时，Go 的 VCS 探测会打印 `does not have any commits yet`，CLI 已能继续构建；建立首个提交后该提示消失。

## 2026-10-21 代码检查与第二轮优化

Author: HyperparameterSI

本轮对全仓库执行格式、编译、静态检查和测试，重点人工检查事件分发、Eval 取消、Dashboard 采样与轮询路径；不代表逐行证明所有平台实现都没有缺陷。保留公开框架 API、许可证和第三方归属声明。

### 代码规范与结构

- 统一 Go 的 gofmt 和前端 Prettier 格式，添加 `.editorconfig`、`.prettierrc.json` 与忽略规则。`bun run format:check` 已加入 Linux CI；生成的绑定与路由文件不手工格式化。
- Dashboard 的 `Monitor` / `Sample` 改为 `ProcessMonitor` / `ProcessSample`，拆分 `process_monitor.go`、`process_metrics.go`、`sample_history.go` 和各平台 `process_cpu_*.go`。这是示例绑定名称调整，框架公开 API 保持不变，生成的 TS 客户端已更新。
- 前端拆出采样 hook、历史合并、数值格式化、图表和指标卡。新增代码保留 `Author: HyperparameterSI`，注释说明并发所有权、单位、生命周期及测量边界。
- `threadRecords` 明确表示 Go 线程 profile 的累计记录数，不误称当前 OS 线程数；运行时内存明确不等于 RSS。
- 修复 Eval 取消回调读取可变测试钩子的竞态；异步执行前捕获钩子，防止上一轮回调与下一测试修改重叠。

### 减少持续工作量

监听器在订阅不变时复用不可变快照；增删或清空订阅时使缓存失效。保留 once、重入和并发语义，代价是每个监听器列表多一个 slice 字段及稳定订阅期间保留的快照。

Dashboard 使用固定 300 点环形缓冲区，写入不再搬移全部历史；GC P95 只在 GC 次数变化时计算，使用栈上缓冲区。首次读取历史，之后按序号增量读取；前端保留有界历史，请求结束后才安排下一轮，暂停和隐藏时停止轮询。后端继续每秒采样，`runtime.ReadMemStats` 本身仍有成本。

### 同机微基准

Apple M5、darwin/arm64、Go 1.27.1。事件基准取 5 次中位数，Dashboard 取 3 次中位数；耗时受系统负载影响，以下不是整个应用 CPU 使用率的下降百分比。

| 路径 | 本轮修改前 | 本轮修改后 |
| --- | ---: | ---: |
| 16 个稳定监听器快照 | 65.10 ns/op，128 B/op，1 alloc/op | 13.57 ns/op，0 B/op，0 alloc/op |
| 满历史采样 | 29,866 ns/op，2,104 B/op，3 allocs/op | 21,997 ns/op，0 B/op，0 allocs/op |
| 完整历史 JSON | 152,251 ns/op，157,342 B/op | 133,361 ns/op，169,879 B/op |
| 稳态轮询 JSON（原 300 点 → 新增 1 点） | 152,251 ns/op，157,342 B/op | 852.6 ns/op，1,008 B/op |

B/op 是该基准的内存分配量，**不是网络包大小**。最后一行比较的是实际轮询策略变化，处理的数据量不同；首次加载或完整导出仍需全量数据。完整 JSON 分配略增来自新增序号等字段，不能声称所有操作都减少内存。

复现：

```sh
go test ./backend -run '^$' -bench BenchmarkListenerSnapshot -benchmem -count 5
go test ./backend/examples/dashboard -run '^$' -bench . -benchmem -count 3
```

### 发布体积与空闲内存

沿用上文发布参数，对 counter-native 和 hello 构建前后版本：

| 示例 | 修改前 | 修改后 | 增量 |
| --- | ---: | ---: | ---: |
| 原生二进制 | 9,232,354 B | 9,232,482 B | 128 B |
| WebView 二进制 | 9,337,362 B | 9,354,226 B | 16,864 B（约 0.18%） |

使用 `-ldflags '-dumpdep …'` 对照链接依赖，新增依赖边集中在监听器缓存 slice 的类型信息、GC 指针布局及 Eval 闭包捕获；没有引入新的第三方运行时依赖。这是以少量体积和缓存保留换取重复分发时更少分配的取舍。

macOS `internal/idlemem`：每个发布程序运行 3 次，每次等待空闲 6 秒，取物理内存占用中位数（十进制 MB）：

| 示例 | 修改前 | 修改后 |
| --- | ---: | ---: |
| 原生应用 | 36.34 MB | 36.85 MB |
| WebView 主进程 | 34.72 MB | 34.74 MB |
| WebView 含子进程 | 69.88 MB | 70.09 MB |

本轮没有证明空闲内存下降；变化小且受系统与 WebView 状态影响。主要收益是高频路径的分配和轮询工作量减少，不应将它们换算成常驻内存或全应用 CPU 降幅。尚未进行长期负载、耗电或各系统真实桌面基准。

### 本轮验证

- Go 全量测试通过；核心与 Dashboard 的 race 检查通过；Eval/监听器重点 race 测试重复 3 次通过。
- Bun 204 项测试通过，包含增量历史合并边界；工作区、Studio、Dashboard 和网站 TypeScript 检查通过。
- macOS / Linux / Windows 编译与 vet、生产标签检查，以及 js/wasm 非 CLI 包 vet 通过。Linux / Windows GUI 未在本机实测。
- 前端格式检查、bridge/runtime/plugins 构建通过；Dashboard 经 CLI 构建独立 macOS app 和 dmg 成功。
