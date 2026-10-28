# Golet Studio

Author: HyperparameterSI

Studio 是框架的交互式组件目录，使用 TanStack Start + React + shadcn 作为入口，通过类型化 Go 绑定打开真实原生控件窗口。不是以 HTML 仿制原生控件。

```sh
pnpm dev         # 开发运行
pnpm build:app   # 构建 Studio 桌面应用
pnpm build:app:mac:intel      # Intel Mac，输出 build/darwin-amd64/
pnpm build:app:mac:arm64      # Apple Silicon，输出 build/darwin-arm64/
pnpm build:app:mac:universal  # 双架构，输出 build/darwin-universal/
```

以上命令从仓库根目录运行，产物目录位于 `backend/examples/studio/build/`。Intel 与 Universal 包均可在 Apple Silicon Mac 上构建；签名和 DMG 生成需要 macOS。Chromium 是独立可选宿主，需要在目标 Mac 上运行 `pnpm setup:chromium` 安装匹配架构的版本。

## 演示范围

- 基础布局：Row、Column、Grid、Box、Spacer、滚动、分栏与主题。
- 表单控件：按钮、输入、复选、开关、单选、选择器、组合框、自动补全、搜索、标签输入、数值、滑块、评分、进度、日期、时间和颜色选择。
- 数据与导航：虚拟列表、表格、树、Outline、网格、Tabs、Sidebar、Breadcrumbs、Router、拖放和选择。
- 反馈：对话框、弹层、菜单、Tooltip、Toast、折叠面板和手风琴。
- 高级内容：文本缓冲区、富文本、Markdown、代码编辑器、Diff、图表、绘图、SVG、动画和国际化。
- 插件：Scene3D 三维视口、边框光束、玻璃效果、全部 29 种 Loading 指示器、真实终端和 Golet Browser。
- Web：shadcn Button / Input / Card 与 Go IPC；原生窗口中的系统 WebView。

可按名称、API 名称和分类搜索演示分组。选择卡片后点击“打开原生演示”。再次打开同一分组会聚焦原窗口。关闭后再次打开重置表单状态。

在“原生主题实验室”切换 Liquid Glass 或 Maia，并选择跟随系统、浅色或深色外观。设置在当前 Studio 会话中共享：已打开的原生演示窗口同步更新，之后打开的窗口继承选择。控件、导航和 Gallery 卡片统一使用主题材质，文字、图表和颜色样本保留可读性及语义颜色。系统 WebView 的网页内容与独立 Browser 界面由各自的样式控制。

## Scene3D

在「插件 → Scene3D · 三维视口」点击「打开原生演示」，或搜索 `Scene3D`、`LoadOBJ`。
窗口包含两个页签，全部由 Go 原生 UI 与 Forma 构建：

- **效果演示**：立方体、通过 OBJ 解析的金字塔、100 个共享网格实例、深度遮挡。
  支持播放/暂停、摄像机环绕/重置、方向光强度、无光照材质、480/720/960 像素上限。
  左键拖动环绕/俯仰，Shift＋左键或中键拖动平移，滚轮或触控板捏合缩放；
  点击视口后也可用左右方向键。重置视角恢复原始位置和观察中心。
  底部实时显示实例、绘制三角形、剔除数和 CPU 缓冲容量。
- **使用示例**：可复制的完整 Go 应用、OBJ 加载片段与内置模型、动画回调、缓存更新和 OrbitControls 接入说明。

窗口跟随 Studio 原生主题实验室的 Maia / Liquid Glass 与明暗设置，3D 背景同步调整。
玻璃材质作用于 UI；3D 模型仍使用 Scene3D 的纯色材质。减少动态效果、暂停以及
切换至使用页时不会持续请求 3D 动画帧。再次点击目录会聚焦已有演示窗口。

当前为 **CPU 软件光栅器**，不是原生 GPU 三维引擎。统计不包含 UI/GPU 纹理，不等于
FPS 或其他引擎的性能比较。接口、限制和独立运行方式见 [Scene3D 说明](../../plugins/scene3d/README.md)。

## Browser

首次构建浏览器宿主：

```sh
brew install cmake  # 已有 CMake 3.21+ 可跳过
pnpm setup:chromium
```

在 Studio 的 Browser 分组选择 Alloy 或 Chrome，输入 HTTP(S) 地址并启动。面板支持导航、前进、后退、刷新、停止、新建标签、状态检查与关闭。切换风格前先关闭当前实例。

Browser 使用独立 CEF 窗口，目前宿主仅支持 macOS。它不是嵌入 Studio 的 WebView。Studio 启动时不自动下载或运行 Chromium；没有宿主时会显示安装提示，其他控件仍可用。Studio 使用 `EphemeralDemo: true`，每次启动创建独立临时配置，macOS 下不访问系统钥匙串，正常关闭后删除临时配置，不保留登录状态。该开发模式不提供钥匙串加密保护，请勿用于真实账号；原有持久配置不会被修改或删除。Studio 进程退出会关闭其浏览器实例。完整 AI 助手和扩展管理演示仍见 `backend/examples/browser`。

Terminal 运行真实 Shell，输入命令会影响本机；运行库仅在打开终端时加载。系统 WebView 示例使用离线 HTML，不依赖外部网站。

## 代码结构与验证

- `backend/examples/studio/`：组件目录、Go 服务、Browser 生命周期及补充原生演示。
- `frontend/examples/studio/`：React 页面、shadcn 组件及生成的 TypeScript 绑定。
- `backend/examples/internal/gallerydemo/`：Gallery 与 Studio 共用的原生组件示例。
- `backend/examples/internal/editordemo/`、`loadingdemo/`：共用编辑器和 Loading 示例。

原 `go run ./backend/examples/gallery`、`code-editor` 和 `loading` 命令保留。共享 Gallery 的定时采样会在窗口关闭后停止。

```sh
go test ./backend/examples/studio ./backend/examples/internal/...
go vet ./backend/examples/studio ./backend/examples/internal/...
```

覆盖检查会检查生成的公开 UI 函数是否出现在演示中，排除非可视辅助函数和与样式控件对应的 Base 构件；它不代表每个参数组合、每个平台或第三方 Chrome 插件都经过验证。

### Border Beam

「插件 → 边框光束」打开纯 Go 原生演示窗口，支持 Medium、Small、Line、Pulse Inner、Pulse Outside 五种动画，以及四组彩色卡片、暂停/继续、隐藏、周期、强度、内辉光范围调节。外观跟随 Studio 的明暗主题；减少动态效果时自动静止。实现与上游差异见 [插件说明](../../plugins/borderbeam/README.md)。

### Forma · Maia

「原生控件 → Forma · Maia」展示固定官方目录的 65 个条目，使用纯 Go 原生组件。下拉框选择组件，Dark 切换明暗，“跟随 Studio 主题”验证全局主题适配。高级能力与视觉验收状态见 [原生组件说明](../../ui/forma/README.md)。

Maia 原生演示现支持 Neutral、Zinc、Stone、Blue、Violet、Emerald、Rose 七套配色，搭配 Dark、Maia/Soft/Sharp 圆角和 14/16 DIP 字号。主题按选择缓存，不在动画帧中重复构造。自定义方式见 `backend/ui/forma/README.md` 的 `NewTheme` 示例。
