# MyGo v0.8.3 → v0.9.0 同步记录

作者：HyperparameterSI。同步日期：2026-11-12。

基线为 `v0.8.3`（`5add3817`），目标为 [v0.9.0](https://github.com/egoist/mygo/releases/tag/v0.9.0)（`4811a286`）。按两个标签之间的差异迁移，未混入目标标签之后的 main 提交。Golet 的发布版本保持 `0.1.17`，本次代码升级不创建发布标签。

## 功能与修复

- **原生模块（实验性）**：新增 `backend/native`、`backend/ffi`，通过 Go 接口及生成器调用 Swift、Rust、Zig；增加 CLI 的生成、检查、构建与分发流程，以及 `ui.NativeView`。Go 应用仍使用 `CGO_ENABLED=0`；自行开发原生模块需要对应语言的编译器，使用已发布模块可下载校验后的预编译库。
- **视频与音频**：新增 `ui/video`，接入 AVFoundation、Media Foundation、GStreamer 和浏览器视频后端；加入视频帧、YUV/NV12 渲染、播放控制及相关生命周期处理。音频新增 AAC、FLAC、MP4/M4A（含 ALAC）及 DASH 支持。
- **原生 UI**：新增 SplitButton、PrimarySplitButton、ListState.Choose、对数/对称对数图表轴，更新滚动条、焦点滚入视口、输入、无障碍、布局及绘制逻辑。分体按钮适配 Golet 的 36 像素默认高度、主题材质和中性焦点环，切换窗口活动状态不改变其尺寸。
- **编辑器与文档**：扩充语言支持、预览、patch diff、行内控件、富文本扩展，以及 Markdown CodeLink 等能力；共享的 Studio/Gallery 演示同步更新。
- **窗口与系统**：同步带参数重新启动、清理窗口状态、选定语言、退出钩子、KeepRunning、输入区域、敏感剪贴板、默认应用/Quick Look、客户端证书等接口和平台实现。
- **插件与构建**：同步 updater、fetch 等修复，包括拒绝带用户名密码的请求 URL；更新开发签名、开发者身份、后台模式和 RPM 打包。SQLite 改用原生模块实现。
- **测试与 CI**：增加原生模块、视频、编辑器、图表、布局等回归测试；CI 加入 Zig/SQLite 原生构建、原生模块往返测试以及 Linux GStreamer 依赖。

## Golet 适配

保留 `backend/`、`frontend/` 分离结构、Forma、Liquid Glass、Chromium、Scene3D、BorderBeam、Omnix 集成、结构化日志以及原有 Studio 性能面板。新增包使用 Golet 模块路径与命令名；保留上游版权及真实发布资源地址。

Studio 新增「原生模块 · MapKit」入口，与 Go 绘制的按钮共处一个窗口。该演示需要 macOS，使用 Golet CLI 构建 Swift 模块；其他平台会返回明确的不可用错误。普通应用不依赖该演示。

SQLite 继续使用上游固定 SHA-256 的发布二进制。生成的 Golet 原生库使用 `golet_` ABI 前缀，加载器兼容上游二进制的 `mygo_` 前缀；没有修改发布资源校验和来假装资源已重新构建。

## 升级注意事项

1. 原生模块接口目前为实验性。修改声明后通过 `golet generate` 或 `golet dev` 重新生成并构建，不要手改生成文件。
2. SQLite 的旧 `LibraryPath` / `GOLET_SQLITE_LIBRARY` 配置迁移到 `GOLET_NATIVE_SQLITECORE`，指向新原生模块库；系统的普通 libsqlite3 不能替代该模块。
3. Linux 视频播放需要 GStreamer 及所用格式的解码插件；支持哪些格式取决于平台可用的解码器。
4. 自行构建 SQLite 原生模块需要 Zig 0.16.0；普通应用通过 CLI 获取已发布库，无需安装 Zig。
5. Swift MapKit 示例运行：`go run ./backend/cmd/golet dev backend/examples/native-module`。Studio 可运行 `pnpm dev`，首次构建新增 Swift 模块会耗时。

## 体积检查

同一工具链、`CGO_ENABLED=0`、`golet_noinspector`、`-trimpath`、`-s -w`：

| macOS arm64 示例 | 升级前（字节） | 升级后（字节） | 增幅 |
| --- | ---: | ---: | ---: |
| counter-native | 9,847,346 | 9,950,242 | 1.04% |
| hello | 9,763,618 | 9,814,178 | 0.52% |

`-dumpdep` 检查未发现普通原生示例链接 `native`、`ffi`、`ui/video`、音频解码器或插件包。增长来自公共平台/渲染/UI 路径（如窗口证书、输入区域、滚动/焦点、视频帧合成和 NV12 光栅支持）。这是二进制体积对比，不能据此推断 CPU、内存或帧率提升；本次没有宣称倍数性能提升。

## 验证记录

- `CGO_ENABLED=0 go test ./...` 全量通过；生产模式 `go test -tags golet_noinspector ./backend/ui ./backend` 通过。
- macOS、Linux、Windows 的 `go build ./...` 和 `go vet ./...` 通过；`golet_noinspector` 及 WebAssembly 静态检查通过。
- 前端 258 项测试、TypeScript 检查、格式检查、所有包构建及 Studio 生产构建通过。
- Swift Language/Motion/Places 示例构建通过；实际 Go/Swift 原生模块往返及跨模块测试通过，无跳过。
- 使用 SHA-256 校验后的 macOS arm64 SQLite 发布库，实际插件测试通过，无跳过。
- Linux/Windows 本机 GUI、视频播放与 Rust/Zig 原生模块运行时未在本机执行；相关依赖和验证步骤已加入 CI。未以交叉编译代替平台运行验证。
