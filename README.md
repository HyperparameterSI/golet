# Golet

**Go 跨平台桌面应用框架，兼顾原生 UI 和 Web 前端。**

Golet 基于 [MyGo](https://github.com/egoist/mygo) 的完整实现迁移，保留 MIT 许可证和来源声明。Go 模块仍位于根目录；公共 Go 包的导入路径为 `github.com/HyperparameterSI/golet/backend`。

- **原生 UI**：纯 Go 编写，GPU 绘制，支持布局、输入法、无障碍、虚拟列表、编辑器、图表和动画。原生控件由框架绘制，并非操作系统原生控件的封装。
- **Web 前端**：使用系统 WebView（macOS WKWebView、Linux WebKitGTK、Windows WebView2），不捆绑 Chromium。
- **shadcn 模板**：TanStack Start、React、TypeScript、Vite、Tailwind CSS，提供组件源码、明暗主题、类型化 Go 调用和错误处理。可使用 shadcn CLI 继续添加组件。
- **Go ↔ TypeScript**：生成类型安全客户端，支持事件、流式通道和插件。
- **桌面能力**：窗口、菜单、托盘、剪贴板、对话框、快捷键、文件监听、自动更新及跨平台打包。
- **无 cgo**：三平台可交叉编译；原生 UI 也支持 WebAssembly。

## 从零构建一个跨平台项目

Golet 有两种使用方式：通过 CLI 创建独立应用，或在已有 Go 项目中导入框架。以下流程使用 **本地源码 + pnpm**，不依赖 Golet 已发布到 Go module proxy 或 npm。应用与框架是两个独立目录，不要把业务项目放进 Golet 仓库。

### 1. 准备环境与框架

- Go：当前源码要求 **1.27.1+**，以根目录 `go.mod` 为准。
- Web 项目：Node.js **22.12+**（推荐 24）及 pnpm。纯 Go 原生 UI 项目不需要 Node.js、pnpm 或 Bun。
- WebView 运行环境：macOS 使用系统 WKWebView；Windows 需要 WebView2 Runtime；Linux 需要 WebKitGTK（Debian/Ubuntu 可安装 `libwebkit2gtk-4.1-0`）。原生 UI 的平台依赖见[原生 UI 文档](docs/native.md)。

先将 Golet 源码放到本地 `golet/`。已有本仓库可直接使用，不必重新克隆；下面的示例工作目录是 `workspace/`：

```text
workspace/
├── golet/          # 框架源码
└── my-desktop/     # 下一步由 CLI 创建
```

在框架仓库根目录执行：

```sh
cd golet
pnpm install
pnpm --dir frontend/packages/runtime build
```

这里只需要构建供业务项目使用的 `golet-runtime`；不用执行依赖 Bun 的框架完整构建/测试脚本。

### 2. 一条命令初始化 Web 桌面应用

仍在 `golet/` 根目录执行（目标目录必须不存在或为空）：

```sh
go run ./backend/cmd/golet init -golet . -template shadcn -name "My Desktop" -module example.com/my-desktop ../my-desktop
```

`-golet .` 使用当前框架源码，自动写入 Go 的 `replace`、前端 runtime 的本地 `file:` 依赖和本地 CLI 调用。`-module` 是新业务项目自己的 Go 模块名，可以换成你自己的仓库路径。

CLI 会生成：

```text
my-desktop/
├── backend/main.go          # Go 服务、窗口和应用生命周期
├── frontend/
│   ├── src/components/      # React 页面、可编辑 shadcn 组件
│   ├── src/routes/          # TanStack 文件路由
│   ├── src/golet.ts         # 自动生成的 Go 调用客户端
│   ├── src/style.css        # Tailwind 与主题样式
│   ├── vite.config.ts
│   └── components.json
├── resources/icon.png       # 应用图标
├── golet.config.ts          # 名称、标识符、版本、前端路径、打包配置
├── go.mod                   # 独立业务模块，本地 replace 指向 golet
└── package.json             # 从此目录运行 pnpm 命令
```

### 3. 安装、启动与开发

```sh
cd ../my-desktop
pnpm install
pnpm dev
```

`pnpm dev` 同时启动 Vite 和 Go 桌面窗口，支持前端热更新与 Go 代码变更后的重建。它不是单独打开浏览器：网页要在 Golet 窗口里才有桌面 IPC 能力。

```sh
pnpm generate       # Go 服务方法变更后，重新生成 TypeScript 客户端
pnpm typecheck      # 检查前端类型
pnpm build          # 构建当前平台的独立桌面应用
```

新模板内部通过 `npm run dev:web` / `npm run build:web` 启动前端；npm 随 Node.js 提供，外层可以继续用 pnpm。若此前已生成旧项目，请把 `golet.config.ts` 中对应的 `bun run ...` 改成 `npm run ...`。

本地模式下不要移动或删除 `golet/`：生成的 Go `replace`、runtime `file:` 依赖及配置导入都引用它。团队协作时可约定同级目录，或在有正式发布版本后改用版本化依赖。

### 4. 添加 Go 业务能力并在 React 中调用

模板已包含 `Greeter`，在 `backend/main.go` 中给它增加方法即可：

```go
// Author: HyperparameterSI

// Add adds two numbers; the method is exposed through the existing binding.
func (Greeter) Add(a, b int) int {
    return a + b
}
```

模板中的 `golet.Bind(Greeter{})` 已注册此服务，不需要重复注册。新建服务时，使用导出的结构体与方法，并在 `App.Run()` 前调用 `golet.Bind(YourService{})`。

运行 `pnpm generate` 后，在 React 组件中调用生成的方法：

```tsx
// Author: HyperparameterSI

import { useState } from "react";
import { Greeter } from "@/golet";
import { Button } from "@/components/ui/button";

export function AddExample() {
  const [message, setMessage] = useState("");
  const [pending, setPending] = useState(false);
  async function calculate() {
    setPending(true);
    try { setMessage(`计算结果：${await Greeter.add(20, 22)}`); }
    catch (error) { setMessage(error instanceof Error ? error.message : String(error)); }
    finally { setPending(false); }
  }
  return <div>
    <Button disabled={pending} onClick={calculate}>调用 Go</Button>
    <p role="status">{message}</p>
  </div>;
}
```

将该组件放入 `frontend/src/components/`，并在现有页面中引入。`frontend/src/golet.ts` 会由 CLI 覆盖，不要手工修改。耗时计算或 I/O 应放入后台任务，避免阻塞 UI 主线程；参见 [IPC 文档](docs/frontend.md)。

### 5. 配置与跨平台打包

先编辑业务项目根目录的 `golet.config.ts`：

```ts
// 保留模板已有的 defineConfig 导入和其他字段，按产品修改：
name: "My Desktop",
identifier: "com.yourcompany.mydesktop",
version: "0.1.0",
main: "./backend",
frontendDist: "frontend/dist/client",
bindings: "frontend/src/golet.ts",
out: "build",
```

这是配置对象的字段片段，不是完整文件。默认应用图标采用 蓝紫渐变圆角底与白色圆环标识；替换 `resources/icon.png` 可设置自己的应用图标。TanStack Start 仅在构建时预渲染，最终应用内嵌 `frontend/dist/client`，不附带 Node.js 服务。

在业务项目根目录运行：

```sh
pnpm build

# 指定平台；下面直接调用本地 replace 指向的 Golet CLI
go run github.com/HyperparameterSI/golet/backend/cmd/golet build -platform windows/amd64
go run github.com/HyperparameterSI/golet/backend/cmd/golet build -platform linux/amd64

# 在 macOS 上生成同时支持 Apple Silicon / Intel 的应用
go run github.com/HyperparameterSI/golet/backend/cmd/golet build -platform darwin/universal
```

| 目标 | 产物位置 / 说明 |
| --- | --- |
| macOS | `build/darwin-arm64/`、`build/darwin-amd64/` 或 `build/darwin-universal/`，包含 `.app` 和 `.dmg` |
| Windows | `build/windows-amd64/` 等目录，包含 `.exe`；生成安装器还需要 NSIS |
| Linux | `build/linux-amd64/` 等目录，包含可执行文件、Debian 包和安装归档 |

Golet 本身不依赖 cgo，但**交叉编译可执行文件不等于可在任意宿主完成所有签名和安装器步骤**：macOS 签名、DMG 和公证需要 macOS；Windows 安装器需要 NSIS；业务自行引入 cgo 后还需对应工具链。建议分别在 macOS、Windows、Linux CI 上构建、测试和发布。正式分发的签名、公证、安装器说明见[分发文档](docs/distribution.md)。

### 6. 只用 Go：创建原生 UI 项目

在框架仓库根目录执行，不需要准备前端依赖：

```sh
go run ./backend/cmd/golet init -golet . -template native -name "My Native App" -module example.com/my-native-app ../my-native-app
cd ../my-native-app
go tool golet dev
go test ./...
go tool golet build
```

原生模板目前把 `main.go`、视图测试和 `golet.json` 放在业务项目根目录，不生成 `frontend/`。视图通过 `github.com/HyperparameterSI/golet/backend/ui` 构建；配置文件是 JSON，因此运行它也不需要 JavaScript 工具链。

### 7. 在已有 Go 项目中引入框架

下面从一个新 Go 模块演示；已有模块跳过 `go mod init`。先回到 `workspace/`（框架 `golet/` 的父目录），再执行：

```sh
mkdir my-native-lib
cd my-native-lib
go mod init example.com/my-native-lib
go mod edit -go=1.27.1
go mod edit -replace=github.com/HyperparameterSI/golet=../golet
go mod edit -require=github.com/HyperparameterSI/golet@v0.0.0
go get -tool github.com/HyperparameterSI/golet/backend/cmd/golet
```

创建 `main.go`：

```go
// Author: HyperparameterSI

package main

import (
    "log"
    "github.com/HyperparameterSI/golet/backend"
    "github.com/HyperparameterSI/golet/backend/ui"
)

func main() {
    golet.App.WhenReady(func() {
        golet.NewWindow(golet.WindowOptions{
            Title: "Hello Golet", Width: 640, Height: 420,
            Content: ui.View(func(c *ui.Context) {
                ui.Column(c).Fill().Center().Children(func() {
                    ui.Text(c, "Hello, Golet!").FontSize(28)
                })
            }),
        })
    })
    if err := golet.App.Run(); err != nil { log.Fatal(err) }
}
```

创建 `golet.json`：

```json
{
  "name": "Hello Golet",
  "identifier": "com.example.hellogolet",
  "version": "0.1.0",
  "out": "build"
}
```

然后运行：

```sh
go mod tidy
go run .
# 或启用开发重载 / 应用打包
go tool golet dev
go tool golet build
```

引入的是 Go 模块 `github.com/HyperparameterSI/golet`，Go 代码中导入的公共包是 **`github.com/HyperparameterSI/golet/backend`**，两者不要混淆。

### 发布后的安装方式

当前文档以本地源码模式为可用路径，不假定 `v0.1.0` 或 npm 包已发布。**待对应 Go tag、`golet-cli` 和 `golet-runtime` 均发布后**，才可以省去本地 `-golet`：

```sh
# 仅在 v0.1.0 已正式发布后执行
go install github.com/HyperparameterSI/golet/backend/cmd/golet@v0.1.0
golet init -template shadcn -module example.com/my-desktop my-desktop
```

`golet` 需位于 PATH（通常来自 `go env GOPATH` 下的 `bin` 目录）。没有正式发布版本时，请使用前面的本地初始化流程。

## 项目结构

```text
golet/
├── backend/                 # 全部 Go 实现
│   ├── *.go                 # 框架公共 API、IPC、窗口与生命周期
│   ├── cmd/golet/           # Go CLI（含可发布的嵌入式脚手架资源）
│   ├── internal/            # macOS / Linux / Windows 后端及渲染
│   ├── ui/                  # 纯 Go 原生 UI
│   ├── plugins/             # 插件的 Go 实现
│   └── examples/            # 示例 Go 入口与桌面构建配置
├── frontend/                # 全部 Web 应用及 TypeScript 包
│   ├── examples/studio/     # TanStack Start + shadcn + React + Vite
│   ├── examples/dashboard/ # 性能监控 Dashboard
│   ├── examples/todo/       # 轻量 TypeScript 示例
│   ├── packages/            # bridge / runtime / CLI npm 包
│   ├── plugins/             # 插件的 TypeScript 客户端
│   └── website/             # TanStack Start 文档站
├── docs/                    # 使用文档与架构
├── scripts/                 # 仓库构建、版本和性能工具
├── go.mod                   # 单一 Go 模块，根目录没有 Go 实现文件
└── package.json             # 前端工作区及统一命令
```

Studio 的视图在 `frontend/examples/studio/src/components/app.tsx`，文件路由在 `src/routes/`；Go 服务在 `backend/examples/studio/main.go`。类型化 IPC 客户端生成到前端的 `src/golet.ts`。

TanStack Start 使用 [SPA 模式](https://tanstack.com/start/latest/docs/framework/react/guide/spa-mode)，将外壳预渲染为 `dist/client/index.html`，由 Go 嵌入并提供服务；`dist/server` 只用于构建时预渲染，不随桌面应用分发。业务后端仍使用 Go，不使用 Start server functions。

根目录快捷命令：`bun run dev`、`bun run build:app`、`bun run build:web`、`bun run test:go`。`bun run build` 构建 bridge/runtime/plugin 包。

## 本地开始

需要 Go 1.27.1+ 和 Node.js 22.12+（推荐 Node.js 24）。Studio 支持 pnpm，无需安装 Bun；框架自身的测试、bridge/plugin 打包和发布脚本仍使用 Bun。最终应用不依赖 Node.js 或 Bun。Go 的自动工具链机制可下载 `go.mod` 指定的版本。

```sh
pnpm install
pnpm dev         # 启动 Studio
pnpm build:web   # 构建 Web 静态资源
pnpm build:app   # 打包桌面应用
pnpm build:app:mac:intel      # Intel Mac：darwin/amd64
pnpm build:app:mac:arm64      # Apple Silicon：darwin/arm64
pnpm build:app:mac:universal  # 双架构安装包：Intel + Apple Silicon
```

Intel 和 Universal 打包命令可在 Apple Silicon 或 Intel Mac 上运行；签名与生成 DMG 需要 macOS。Studio 产物分别位于 `backend/examples/studio/build/darwin-amd64/` 和 `darwin-universal/`。Universal 同时包含 x86_64 与 arm64。可选 Chromium 宿主仍需在目标机器单独安装，不会随 Studio 自动打包。

新生成的 shadcn 应用提供对应的 `pnpm build:mac:intel`、`pnpm build:mac:arm64` 和 `pnpm build:mac:universal`。已有应用也可直接执行 `golet build -platform darwin/amd64`。

上述 Studio 命令会自动编译本地 `golet-runtime`。开发框架自身时可继续使用 Bun：

```sh
bun install
bun run build

# shadcn 桌面示例（Go 后端 + React WebView）
go run ./backend/cmd/golet dev backend/examples/studio

# 纯 Go 原生 UI
go run ./backend/examples/counter-native
go run ./backend/examples/gallery
```

从本地源码生成应用（目标目录必须为空）：

```sh
# 默认是 shadcn；显式指定 -template web 可生成轻量原生 TypeScript 前端
go run ./backend/cmd/golet init -golet "$PWD" -template shadcn ../my-desktop-app

# 纯 Go 原生 UI，默认使用中性 shadcn 风格主题
go run ./backend/cmd/golet init -golet "$PWD" -template native ../my-native-app
```

新建的默认 Web 项目也将 Go 入口放入 `backend/`、Web 源码与 Vite 配置放入 `frontend/`。Web 项目进入根目录执行 `pnpm dev` / `pnpm build`；原生项目执行 `go tool golet dev` / `go tool golet build`。`-golet` 指向本地仓库，不需要等待 Golet npm 包发布。

**尚未发布到 npm 或 GitHub Releases**：当前使用本地源码工作流。仓库中的发布配置是待配置的基础设施，不代表包已经存在。可选 SQLite / terminal 插件及 NSIS 使用原项目固定版本、SHA-256 校验的二进制下载；SQLite 的底层 ABI 符号保留原名以保持兼容。

## 使用 AI 生成 Golet 应用

Web 桌面应用默认使用 **TanStack Start + React + TypeScript + Vite + shadcn/ui + Tailwind CSS**。`golet init` 默认选择 `shadcn` 模板，也可以显式指定：

```sh
golet init -template shadcn -module example.com/my-app my-app
cd my-app
pnpm install
pnpm dev
```

新生成的 shadcn 项目会携带以下 AI 开发指引：

| 文件 | 用途 |
| --- | --- |
| `AGENTS.md` | 技术栈、目录边界、IPC、SPA 打包、代码注释与验证规则的统一来源 |
| `CLAUDE.md` | 引导 Claude Code 读取项目规则 |
| `.github/copilot-instructions.md` | 引导 GitHub Copilot 读取项目规则 |
| `README.md` | 启动、构建及可直接复制的 AI 提示词 |
| `.agents/skills/` | CLI 随项目安装的 Golet 维护技能与参考资料 |

可以把下面这段话交给 AI：

> 使用 Golet 构建跨平台桌面应用。先读取项目 AGENTS.md 和 golet.config.ts。前端使用现有 TanStack Start + React + TypeScript + Vite + shadcn/ui，保持 SPA 模式；前端源码放 frontend/，Go 业务与系统能力放 backend/，通过生成的 Golet 客户端通信。复用 components/ui 和主题变量，不另建框架或 Node 服务。完成后运行相关类型检查、测试与构建，并说明未验证的平台。我要实现的功能是：……

README 负责说明，项目规则负责提示 AI，模板里的依赖、路由和构建配置提供实际默认实现。AI 工具是否自动加载规则取决于其支持情况；未自动读取时，应显式附上 `AGENTS.md`。这些文件不能强制所有模型遵循，也不能替代类型检查和构建验证。用户明确要求其他技术栈时，以用户要求为准。

已有项目可以按实际技术栈合并这些规则，不要覆盖原有项目约定。`golet install-skills .` 仅更新内置技能，不会自动补写或覆盖项目根目录的 `AGENTS.md`、`CLAUDE.md` 和 README。纯 Go 原生 UI 仍使用 `-template native`，轻量 TypeScript 前端使用 `-template web`。

## shadcn 与原生主题

Web 模板的 `components.json`、`@/` 路径别名和 CSS 语义变量已配置：

```sh
# 在生成的 Web 项目的 frontend 目录内
bunx shadcn@latest add dialog
```

原生应用可在初始化时创建主题，在视图中设置：

```go
var theme = ui.ShadcnTheme(false) // true 为暗色

func view(c *ui.Context) {
    c.SetTheme(theme)
    ui.Column(c).Fill().Center().Children(func() {
        ui.Text(c, "Hello, Golet").FontSize(28).Bold()
        if ui.PrimaryButton(c, "Continue").Clicked() {
            // 应用逻辑
        }
    })
}
```

实现参考 [shadcn Vite 配置](https://ui.shadcn.com/docs/installation/vite) 和 [主题变量](https://ui.shadcn.com/docs/theming)。原生主题使用相同中性色语义，实际绘制仍由 Go 完成。

## 验证与性能

```sh
CGO_ENABLED=0 go test ./...
CGO_ENABLED=0 go build ./...
CGO_ENABLED=0 GOOS=linux go build ./...
CGO_ENABLED=0 GOOS=windows go build ./...
bun run test
bun run typecheck
bun run format:check

go test -run '^$' -bench '^BenchmarkListenerSnapshot$' -benchmem ./backend
```

事件监听器复用稳定订阅的不可变快照：16 个监听器的本轮基准由 **128 B / 1 alloc** 降为 **0 B / 0 alloc**，同机耗时中位数从 65.10 ns 降为 13.57 ns。Dashboard 使用固定容量环形历史与增量 IPC，暂停或隐藏时停止前端轮询。这些是局部工作量与临时分配的改进，实测空闲内存没有明显下降，不能推导为整个应用 CPU 或内存同比降低。shadcn 模板无定时轮询，空闲时不发送时钟 IPC。

Go 源码使用 `gofmt`；前端通过 `bun run format` 格式化、`bun run format:check` 检查（已加入 CI），生成文件除外。新增注释注明作者 `HyperparameterSI`，并优先解释单位、并发约束和设计原因。

详细验证与限制见 [Golet 迁移说明](docs/golet.md)，框架 API 见 [文档](docs/README.md)，内部实现见 [架构](docs/architecture.md)。

## License

MIT。保留 [LICENSE](LICENSE) 中的原作者版权，来源版本见 [NOTICE](NOTICE)。

## 性能监控 Dashboard

```sh
pnpm dev:dashboard     # Golet Observatory，Web 开发端口 5174
pnpm build:dashboard   # 打包独立监控示例
```

Go 后端在 `backend/examples/dashboard`，Web 前端在 `frontend/examples/dashboard`。
展示真实进程 CPU、Go 堆/栈内存、分配速率、Goroutine、GC 暂停与事件、环境信息和健康检查，支持趋势窗口、暂停和 JSON 导出。
指标不包含 WebView 子进程；Go 运行时内存不等于 RSS。完整口径见 [Dashboard 说明](backend/examples/dashboard/README.md)。

## 可选 Chromium 浏览器（实验性）

```sh
pnpm setup:chromium  # 首次构建 CEF 宿主，需要 Go、CMake 和 Xcode 命令行工具
pnpm dev:browser     # 已安装后直接运行
```

提供自定义原生标签栏、圆角地址栏和离线搜索首页，使用 CEF Alloy 渲染网页；Go 侧支持多标签管理、导航、历史、刷新、停止、状态订阅和生命周期管理。首页可在 `frontend/examples/browser/newtab/` 自定义，修改后重新运行安装命令。只有选择该插件才需携带 CEF；当前 macOS 宿主的 Apple Silicon 已实测；Windows amd64/arm64 已提供实验性宿主与安装流程，尚待 Windows 原生编译和实机验证。它是独立浏览器窗口，嵌入 Golet 布局、Linux 宿主和自动集成 `golet build` 尚未实现。详情见 [Chromium 插件](backend/plugins/chromium/README.md) 与 [浏览器示例](backend/examples/browser/README.md)。

Golet Browser 还提供可选 AI 助手、确认后执行的页面操作、按网站授权的用户脚本，以及 Chrome 风格窗口下的本地 Manifest V3 插件加载。模型配置、扩展格式和兼容范围见 [浏览器使用说明](backend/examples/browser/README.md#ai-助手与页面操作)。
