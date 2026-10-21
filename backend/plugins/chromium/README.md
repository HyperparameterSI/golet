# Chromium（可选、实验性）

Author: HyperparameterSI

真正的 CEF / Chromium 浏览器窗口，Go 代码通过私有 stdin/stdout 管道控制单独的原生宿主进程。普通 Golet 应用不会导入或下载它，系统 WebView 保持原样。Go 侧 `CGO_ENABLED=0`；**首次构建 CEF 宿主仍需要 C++ / Objective-C++ 工具链**，这并非纯 Go 实现的 Chromium。

## 当前范围

- macOS 原生宿主已实现，Apple Silicon 已通过真实窗口测试；Intel 有固定 SDK 和构建配置，未在 Intel 设备实测。
- 自定义 Cocoa 标签栏、圆角地址栏和导航按钮，CEF Alloy 负责网页渲染；离线新标签页提供搜索与快捷入口；没有向任意网页注入 Golet 的 Go 方法桥。
- Go API 支持导航、后退、前进、刷新、停止、新增、切换、关闭标签，获取标签列表和焦点页状态，以及关闭窗口。
- **它目前是独立浏览器窗口，不是可放进 `ui.Row` 的嵌入式 `ui.Chromium` 控件。** Golet 布局内的 CEF 子视图、Windows/Linux 原生宿主和对应打包尚未实现。Go 包在这些平台可编译不代表原生浏览器可用。
- 不提供 Chrome 账号同步、DRM 或商业编解码器可用性的保证；网页兼容性以所固定的 CEF 发行版为准。

## 安装和运行

macOS：需要 Go、CMake 3.21+、兼容 CEF SDK 的 Xcode 命令行工具（SDK 要求见其 README）。在仓库根目录执行：

```sh
brew install cmake  # 首次构建所需；已安装可跳过
pnpm setup:chromium
pnpm dev:browser
# 指定初始网址：
pnpm dev:browser -url https://go.dev
```

安装脚本会检查 CMake 3.21+，自动查找 PATH、Homebrew（Apple Silicon / Intel）和 `/Applications/CMake.app`。也可设置 `CMAKE` 环境变量。CMake 只用于编译 CEF 原生宿主，运行浏览器无需 CMake，也不需要 Python。

如果 CMake 安装在其他位置：

```sh
go run ./backend/plugins/chromium/cmd/setup --cmake /absolute/path/to/cmake
```

安装脚本按 `native/runtime.json` 下载固定版本的 CEF minimal SDK，检查文件大小及 CEF 索引公布的 SHA-1，构建宿主并生成带 Golet 图标的 `.app`。不会下载“最新版”或运行远程脚本。运行应用本身不会自动下载、安装或更新 CEF。

当前 macOS 固定 CEF `154.0.34+g14c5a08+chromium-154.0.8037.98`。Apple Silicon SDK 压缩包约 132 MB，实测宿主 `.app` 约 332 MiB；这是独立的 Chromium 运行环境，不能套用系统 WebView 示例的体积或内存数据。

默认安装位置：

```text
~/Library/Caches/golet/chromium/host-v1/darwin-arm64/Golet Browser.app
```

支持 `--cef-dir /path/to/extracted/sdk` 复用同版本 SDK，`--output /path/to/runtime` 指定安装目录。设置 `GOLET_CHROMIUM_HOST` 或 `Options.HostPath` 时必须指向 `.app/Contents/MacOS/Golet Browser` 可执行文件。

## Go 调用

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

browser, err := chromium.Launch(ctx, chromium.Options{
    URL: "https://go.dev",
    // 每个并发运行的宿主必须使用不同的 ProfileDir。
    ProfileDir: "/absolute/path/to/browser-profile",
})
if err != nil {
    return err
}
// Navigate 返回代表命令已接收；网页完成加载需要观察 Events。
if err := browser.Navigate(ctx, "https://go.dev/doc/"); err != nil {
    return err
}
for state := range browser.Events() {
    log.Printf("%s %s loading=%v", state.Title, state.URL, state.Loading)
}
return browser.Wait()
```

导入路径为 `github.com/HyperparameterSI/golet/backend/plugins/chromium`。也可调用 `Back(ctx)`、`Forward(ctx)`、`Reload(ctx)`、`Stop(ctx)`。`Close(ctx)` 关闭全部标签页；超时会强制结束进程组，正常关闭则让 CEF 刷新用户配置。生命周期 context 取消时也会关闭宿主。

`Events()` 是容量为 1 的最新状态流，消费者较慢时跳过中间状态；不是完整事件日志。没有轮询线程、调试 HTTP 端口或 WebSocket 控制端口。公开导航接口接受 HTTP、HTTPS 和 `about:blank`，自定义地址栏支持网址与 Google 搜索。

## 应用窗口

设置 `Options{AppMode: true, URL: appURL, ProfileDir: appProfile}` 可使用只有内容和系统标题栏的独立应用窗口。地址栏、导航按钮、标签栏及对应菜单不会创建，网页标题作为窗口标题。新增标签命令返回错误，网页弹窗被拒绝；刷新、编辑、窗口缩放和关闭仍然可用。

AppMode 只支持默认的 Alloy，不可与 Chrome 样式、扩展或 AssistantURL 同用。这是窗口呈现模式，不是页面权限或导航隔离机制。更新 Go 包后需重新运行 Chromium setup 构建匹配的原生宿主。

## 自定义界面

- 首页：`frontend/examples/browser/newtab/` 中的 HTML、CSS、JavaScript 和 Logo；无需前端运行时服务。
- 标签栏、地址栏、窗口配色：`backend/plugins/chromium/native/shell_mac.mm`。
- 修改后重新运行 `pnpm setup:chromium`，再启动 `pnpm dev:browser`。
- 快捷键：⌘L 地址栏、⌘T 新标签、⌘W 关闭当前标签、⌘R 刷新。
- `NewTab(ctx, url)` 新增页面（空 URL 打开首页）；`SelectTab(ctx, id)` 切换；`CloseTab(ctx, id)` 关闭。`State.Tabs` 提供标签 ID、URL 和标题。关闭最后一个标签会退出窗口。

## 分发

将完整的 `Golet Browser.app` 随产品携带，通过 `HostPath` 指定它；只复制主可执行文件不能运行。必须保留 Frameworks、各角色 Helper、Resources 中的 CEF 许可证和第三方清单。构建脚本只做本地 ad-hoc 签名；正式发布需要开发者签名和公证，当前 `golet build` 不会自动把该宿主打入其他应用。

## 验证

```sh
CGO_ENABLED=0 go test ./backend/plugins/chromium
# 打开真实窗口，使用本地 HTTP 页面验证导航和生命周期：
GOLET_CHROMIUM_E2E=1 go test -race -v ./backend/plugins/chromium
```

已在当前 macOS/arm64 环境通过真实页面加载、多标签新增/切换/关闭、离线首页、前进后退、刷新、停止、关闭和 race 检查；协议单元测试使用假宿主验证并发命令和取消，不用它代替 CEF 渲染验证。快速加载/退出测试中偶见 CEF GPU mailbox / FrameWidgetHost 诊断，功能断言仍通过；这些上游诊断尚未定位，未完成长时间浏览稳定性测试。

CEF 宿主按照 [CEF 应用结构文档](https://github.com/chromiumembedded/cef/blob/master/docs/tutorial.md) 接入；macOS Helper 按 [CEF sandbox 文档](https://github.com/chromiumembedded/cef/blob/master/docs/sandbox_setup.md) 初始化沙箱，未使用 `--no-sandbox`。CEF 许可证见 `native/CEF-LICENSE.txt`，原始分发的许可证与完整 credits 会一起打包。

## AI 与扩展 API

`DevTools(ctx, tabID, method, params)` 通过宿主私有管道调用 CDP；`Evaluate(ctx, tabID, expression, result)` 执行可信 JavaScript 并解析返回值，无需远程调试端口。这两个 API 拥有页面控制权限，不能直接暴露给网页或不可信模型输出。

`assistant` 子包提供可选的模型客户端、页面快照、受限操作、回环助手服务和用户脚本扩展。浏览器示例中可打开助手进行问答、审核操作和运行扩展，具体配置见 [示例说明](../../examples/browser/README.md#ai-助手与页面操作)。未导入该子包的应用不会启动 HTTP 服务或连接模型。

使用标准 Chrome 扩展时设置 `Options.Style: chromium.StyleChrome`，通过 `Options.Extensions` 传入可信的 Manifest V3 解压目录。此模式采用标准 Chrome 窗口，无法同时使用当前 Alloy 自定义外观。已实测 content scripts；其他扩展 API 兼容性需逐项验证。

## 在 Go 代码中选择窗口风格

```go
// Author: HyperparameterSI
browser, err := chromium.Launch(ctx, chromium.Options{
    Style: chromium.StyleAlloy, // Golet 自定义窗口（默认）
    URL: "https://go.dev",
})
if err != nil {
    return err
}
// 使用 browser；退出时调用 Close 或取消 ctx。
```

将 `Style` 改为 `chromium.StyleChrome` 即使用标准 Chrome 窗口，需要 Chrome 插件时再设置 `Extensions: []string{"/absolute/path/to/extension"}`。两种风格均支持 Go 导航、状态订阅和 DevTools 操作；AI 助手与用户脚本由应用按需接入。风格在 `Launch` 时确定，运行中切换需要关闭后重新启动宿主。

`Style` 不填写时默认为 Alloy。无效风格或 Alloy 搭配 Chrome 插件会在启动前报错。旧的 `ChromeStyle` 字段和 `-chrome-style` 参数保留兼容，推荐迁移至 `Style` 和 `-style`；显式 Alloy 与旧字段 `true` 同时设置会报冲突。

```sh
pnpm dev:browser -style alloy
pnpm dev:browser -style chrome
```

## 开发演示与钥匙串提示

默认使用系统安全存储。macOS 上，临时签名的宿主重新构建后可能重新出现 Safe Storage 授权提示。仅用于控件演示时，可设置 `Options{EphemeralDemo: true}`：创建独立临时 profile，并使用 Chromium 的 `--use-mock-keychain` 开发选项。正常关闭会清理临时目录；异常杀死 Go 父进程可能留下临时文件。

该模式不提供系统钥匙串保护，不应登录真实账号或存储敏感信息，不能与 `ProfileDir` 同时设置。默认值为 `false`，不会改变正式应用的安全存储；不会访问或删除原有持久 profile。Studio Browser 演示默认启用。独立示例可运行 `pnpm dev:browser -ephemeral-demo`。

参考：[Chromium 开发构建权限提示说明](https://chromium.googlesource.com/chromium/src/+/main/docs/mac_build_instructions.md#avoiding-system-permissions-dialogs-after-each-build)。
