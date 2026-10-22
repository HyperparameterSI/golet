# Golet Browser

Author: HyperparameterSI

使用独立 Chromium 插件启动真正的 CEF 浏览器窗口，包括自定义原生标签栏、圆角地址栏和导航按钮，搭配浅色搜索首页。网页由 CEF Alloy 渲染。

首页资源独立放在 `frontend/examples/browser/newtab/`，可修改 HTML、CSS、JavaScript、Logo 和快捷入口；原生工具栏在 `backend/plugins/chromium/native/shell_mac.mm`。修改后重新运行安装命令即可生效。

支持 ⌘L 定位地址栏、⌘T 新建标签、⌘W 关闭当前标签、⌘R 刷新。首页与地址栏均支持网址和 Google 搜索。

在仓库根目录运行：

```sh
pnpm setup:chromium  # 首次安装；macOS，或实验性 Windows amd64/arm64
pnpm dev:browser
pnpm dev:browser -url https://go.dev/doc/
```

可通过 `-host` 指定宿主可执行文件，通过 `-profile` 隔离用户数据。关闭浏览器窗口或按 Ctrl+C 结束示例。

安装依赖、Go API、打包方式和当前限制见 [Chromium 插件](../../plugins/chromium/README.md)。这仍是实验性独立窗口组件，不是 Golet 布局内的嵌入式控件。

## AI 助手与页面操作

点击工具栏的 ✦ 打开助手。可直接提问，或选择目标标签并勾选共享页面，让模型总结内容、建议点击、填写和滚动。模型只生成建议；点击“确认执行这一步”后执行一次，随后重新提问以获取最新页面状态。切换网址或目标元素变化会拒绝旧操作。当前只处理主文档中的可见控件，不支持 iframe、Shadow DOM、文件上传或全自动多步任务。

配置支持 Chat Completions 请求格式的服务（完整 endpoint，不自动补路径）：

```sh
export GOLET_AI_ENDPOINT='http://127.0.0.1:11434/v1/chat/completions'
export GOLET_AI_MODEL='your-installed-model'
# 云端服务使用 HTTPS，并通过环境变量 GOLET_AI_API_KEY 提供密钥。
pnpm dev:browser
```

未配置服务时会显示“模型尚未配置”，不会伪造 AI 回复。密钥仅留在 Go 后端；不写入前端或浏览器存储。页面共享默认关闭，不上传截图、Cookie 或表单值；共享文本仍可能包含私人内容，应由用户决定是否发送。模型服务自行配置和运行，本示例不自动下载模型。每次请求独立，不持久化聊天记录。

助手使用随机端口的回环 HTTP 服务和每次启动生成的 Bearer token，无 CORS 权限，不向第三方网页提供控制桥。`DevTools` / `Evaluate` 是供可信 Go 代码使用的高权限接口，AI 只能调用固定的操作集合。应用退出即关闭服务。

前端在 `frontend/examples/browser/assistant/`，Go 服务在 `backend/plugins/chromium/assistant/`。从仓库根目录运行；分发示例时同时携带前端文件，通过 `-assistant-assets` 指定位置。

## 用户扩展

使用可信的本地 Golet 扩展：

```sh
pnpm dev:browser -url https://go.dev/doc/ \
  -user-extension frontend/examples/browser/extensions/reading-mode
```

助手的“用户扩展”页可启用、停用并手动运行扩展；状态仅保留到本次进程结束。脚本有页面读取和修改权限，只运行在 manifest 声明的精确 origin。示例再次运行可取消阅读样式。扩展目录包含：

```json
{
  "name": "阅读模式",
  "description": "改善阅读排版",
  "origins": ["https://go.dev"],
  "script": "content.js"
}
```

保存为 `golet-extension.json`，再添加 `content.js`。可重复指定 `-user-extension`。这种用户脚本不提供 `chrome.*` API，也不自动在每个页面执行。

## Chrome 浏览器插件

CEF 的 Chrome 扩展 API 要求 Chrome 风格窗口。兼容模式会使用标准浏览器窗口，替换自定义标签栏；助手会在启动时自动打开。

```sh
pnpm dev:browser -style chrome -url https://go.dev/doc/ \
  -extension frontend/examples/browser/extensions/chrome-reading-mode
```

支持加载用户明确指定的本地、已解压 Manifest V3 扩展目录，可重复指定 `-extension`。已通过真实 CEF content script 测试；不承诺所有 Chrome API 或商店插件兼容。目前没有 CRX 导入、Chrome Web Store 一键安装、插件自动更新或完整扩展管理 UI。移除启动参数即可停止加载该目录。扩展权限由 Chromium 按 manifest 执行，请仅加载信任的代码。

参考：[CEF Chrome style 扩展限制](https://github.com/chromiumembedded/cef/issues/3685)、[CEF DevTools 私有接口](https://github.com/chromiumembedded/cef/blob/master/include/cef_browser.h)。

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
