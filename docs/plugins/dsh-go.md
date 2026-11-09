# DSH-Go AI 助手（纯 Go）

<!-- Author: HyperparameterSI -->

`backend/plugins/dsh-go` 是独立可选插件，集成
[dsh-go](https://github.com/HyperparameterSI/dsh-go) 的模型接口、工具 schema 校验、工具执行管线和会话事件。
`dsh-go/native` 提供跟随 Forma 主题的原生聊天面板，包含 Markdown、流式输出、发送、停止、清空及逐次工具授权。
实现和示例均为 Go，不需要 Node、Python、浏览器或 cgo。

## 直接运行

```sh
CGO_ENABLED=0 go run ./backend/examples/assistant
```

没有 API key 时，这是**明确标注的离线交互演示**：请求授权读取当前时间，并展示工具结果，不生成真实模型回答。
在本机环境中设置 `DEEPSEEK_API_KEY` 后重新运行即可连接 DeepSeek；可选 `DEEPSEEK_MODEL`（默认 `deepseek-chat`）与 `DEEPSEEK_BASE_URL`（默认 `https://api.deepseek.com`）。
密钥只留在 Go 进程，不写入 UI、会话记录或日志。真实模型调用需要网络及供应商账户。

## 集成到应用

```go
// Author: HyperparameterSI
adapter, err := dshgo.NewDeepSeek(dshgo.DeepSeekOptions{
    APIKey: os.Getenv("DEEPSEEK_API_KEY"),
})
if err != nil { return err }
panel, err := native.New(dshgo.Options{
    Adapter: adapter,
    Model: "deepseek-chat",
    System: "你是一个简洁友好的桌面助手。",
})
if err != nil { return err }
window := golet.NewWindow(golet.WindowOptions{
    Title: "Assistant", Width: 1000, Height: 760,
    TitleBarStyle: golet.TitleBarHiddenInset,
    TrafficLightPosition: &golet.Point{X: 12, Y: 11},
    Content: ui.View(panel.View),
})
panel.SetWindow(window)
```

以上窗口创建代码放在 `golet.App.WhenReady` 中；在进程入口调用 `golet.App.Run()`。
导入路径为 `github.com/HyperparameterSI/golet/backend/plugins/dsh-go`、其 `/native` 子包及 Golet 的 `backend`、`backend/ui`。
完整可运行版本见 [示例](../../backend/examples/assistant/main.go)。面板继承外围 UI 主题；独立示例使用 `forma.Theme(c.Theme().Dark)`。

不需要聊天 UI 时，使用 `dshgo.New(options)` 创建会话，在工作 goroutine 调用
`conversation.Run(ctx, input, func(event dshgo.Event) { ... })`。
回调同步执行；更新原生界面需使用 `window.Update`。不要从回调调用 `History`。
`Cancel` 停止当前轮，`Reset` 清空空闲会话，`Close` 阻止新请求并取消正在进行的工作。

## 注册 Go 工具

```go
// Author: HyperparameterSI
tool := dshgo.Tool{
    Name: "current_time",
    Description: "读取本机当前时间",
    RequiresApproval: true,
    Execute: func(ctx context.Context, input map[string]any) (any, error) {
        return time.Now().Format(time.RFC3339), ctx.Err()
    },
}
// 在 Options.Tools 中传入 []dshgo.Tool{tool}。
```

`Parameters` 接受 dsh-go 支持的 JSON Schema；省略时为空对象。
框架先校验参数，再授权，再调用 Go 函数。原生面板会展示工具名和完整可滚动参数，由用户单次允许或拒绝。
使用无 UI 的会话时，受保护工具必须配置 `Options.Approve`，否则拒绝执行。
默认不注册 shell、文件系统、网络、浏览器操作或桌面控制工具。由应用显式接入所需能力，并自行限制路径、域名和权限。
工具及自定义模型适配器必须遵守 context 取消；框架不能强行终止任意 Go 函数。

## 限额与会话

默认每轮最多 8 次模型请求、16 次工具调用、2048 输出 token/请求、256 KiB 内容，超时 2 分钟。
历史最多保留 20 个完整成功轮次、1 MiB；超过上下文请求预算返回 `ErrLimit`，可清空历史后重试。
工具输出超限会结束当前轮，避免将大结果继续交给模型。相同会话并发调用返回 `ErrBusy`。
失败或取消的轮次不进入历史，但已经执行的工具副作用不会回滚。
`Result` 提供用量、步骤数及 dsh-go 事件记录；`History()` 返回独立快照。没有默认磁盘持久化或遥测。
这些记录可能含用户输入和工具结果，保存时由应用决定访问权限及保留策略。

## 兼容范围与来源

固定使用 dsh-go 提交 `c1023d186802c2fc0df87eb57c1ff98bde238728`（MIT，Copyright 2026 JopenChen）。
依赖来源为 `github.com/HyperparameterSI/dsh-go`。该仓库目前仍声明
`module github.com/JopenChen/dsh-go`，因此通过根目录 `go.mod` 的 `replace`
切换下载来源，Go 导入路径保留其模块声明，避免产生两套不兼容的类型。
外部应用使用本插件时也需在自己的 `go.mod` 添加相同替换规则（Go 不继承依赖模块的 `replace`）：

```go
replace github.com/JopenChen/dsh-go => github.com/HyperparameterSI/dsh-go v0.0.0-20260911155656-c1023d186802
```

dsh-go 自身是 [DeepSeek Harness](https://github.com/deepseek-ai/deepseek-harness) 的 Go 参考实现。
本插件复用其公开接口与执行管线，但没有直接使用其尚不完整的 agent 循环及 DeepSeek provider：Golet 实现了有界循环和标准 HTTP/SSE 工具消息协议，正确发送 `role=tool`、工具参数及模型名称。

当前支持文本、推理内容转发、串行工具调用和可替换的 `llm.LLMAdapter`。
尚未实现完整官方 Harness 的技能、MCP、沙箱、子代理、检查点恢复、多模态或本地模型推理引擎。
可通过 `BaseURL` 连接兼容的本地推理服务（仅 loopback 允许 HTTP），或实现 dsh-go 的 `LLMAdapter` 接入其他供应商。
HTTP 重定向被拒绝，不自动重试模型或工具调用。没有引入官方 Harness 的 Node 运行时。

该插件只有在应用导入时才进入依赖图；普通原生 UI/webview 应用不启用任何 AI 后台线程或服务。

## Studio 演示与命名迁移

插件目录由 `backend/plugins/assistant` 改为 `backend/plugins/dsh-go`，包名为 `dshgo`。旧目录不保留兼容包装；应用需更新导入路径及包名前缀。底层依赖仍通过 `go.mod replace` 使用 HyperparameterSI/dsh-go。

Studio 的「插件 → DSH-Go · AI Assistant」打开真实原生聊天窗口，复用 `native.Panel`。聊天使用 `forma.AIMessage`，输入使用 `forma.PromptInput`（Enter 发送，Shift+Enter 换行），代码围栏使用 `forma.CodeBlock`，模型标签/菜单使用 `forma.ModelSelector`。跟随 Studio 的明暗和材质主题。支持流式回复、复制、评价、将问题放回输入框编辑、工具单次授权/拒绝、停止、清空以及步骤/token 统计。关闭窗口自动取消请求。

未设置 `DEEPSEEK_API_KEY` 时使用离线适配器，真实经过工具校验和授权链路，但不调用模型，也不虚构 token 用量。设置 key 后重新打开演示窗口，使用 `DEEPSEEK_MODEL` 和 `DEEPSEEK_BASE_URL` 配置的模型。模型菜单仅展示当前实际配置的模型，不提供未接入的假选项。此文本演示没有附件上传和供应商思考强度支持，因此附件按钮禁用、强度菜单不显示。评价仅保存在本次界面状态中。

### 在 Studio 内设置 API Key

打开「DSH-Go · AI Assistant」后点击右上角 **模型设置**，输入 DeepSeek API Key（遮罩输入）、API 地址及模型，点击 **应用并新建会话** 即刻切换在线模式，无需重启。配置只保存在当前窗口的 Go 内存中，不写文件、不经过 Web IPC；关闭窗口后丢弃。已有环境变量会作为初始值。应用设置会清空当前会话；正在回复时须先停止。密钥有效性在首次发送请求时验证。点击 **离线演示** 可清除窗口中的密钥并切回离线适配器。
