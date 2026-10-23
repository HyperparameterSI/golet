# 框架诊断日志

<!-- Author: HyperparameterSI -->

Golet 核心通过 `SetLogHandler` 提供结构化诊断。默认使用 Go 标准库 `log` 输出警告和错误；应用、窗口生命周期和正常 IPC 调用仅在开启详细日志后输出。不创建后台日志线程、轮询器或无界缓存。

在 `App.Run()` 前设置处理器。例如在应用中接入标准库 `log/slog`：

```go
package main

import (
    "log/slog"

    golet "github.com/HyperparameterSI/golet/backend"
)

func configureLogging() {
    golet.SetLogHandler(golet.LogInfo, func(r golet.LogRecord) {
        level := slog.LevelInfo
        switch r.Level {
        case golet.LogDebug:
            level = slog.LevelDebug
        case golet.LogWarn:
            level = slog.LevelWarn
        case golet.LogError:
            level = slog.LevelError
        }
        slog.Log(nil, level, r.Event,
            "window_id", r.WindowID,
            "method", r.Method,
            "error", r.Err,
            "stack", r.Stack,
        )
    })
}
```

`configureLogging` 在启动应用前调用。排查 IPC 时将最低级别改为 `golet.LogDebug`，同时开启目标日志系统的 Debug 级别。`SetLogHandler(golet.LogWarn, nil)` 恢复默认输出；替换处理器是并发安全的，正在执行的旧处理器可能继续完成。可以通过 `SetLogHandler(golet.LogDebug, func(golet.LogRecord) {})` 禁用核心诊断输出。

## 记录范围

| 事件 | 级别 | 含义 |
| --- | --- | --- |
| `app.starting` / `app.ready` / `app.stopped` | Info | 应用生命周期 |
| `app.init_failed` / `app.run_failed` | Error | 平台初始化或事件循环失败 |
| `window.created` / `window.closed` | Info | 窗口生命周期 |
| `window.load_finished` | Debug | 页面加载完成 |
| `window.load_failed` | Warn | 页面加载失败，含错误码和说明 |
| `window.create_failed` / `window.render_process_gone` | Error | 创建失败或渲染进程退出 |
| `window.open_external_failed` / `window.open_options_invalid` | Warn | 默认外链打开失败或新窗口选项无效 |
| `ipc.started` / `ipc.completed` | Debug | IPC 调用开始或成功完成 |
| `ipc.canceled` | Debug | 调用被取消（例如页面关闭），默认不刷警告 |
| `ipc.malformed` / `ipc.failed` | Warn | 请求无法解析或调用返回错误 |
| `ipc.panic` | Error | 服务方法 panic，含堆栈；该调用随后也产生 `ipc.failed` |
| `ipc.encode_failed` | Error | 返回值无法序列化（无窗口和方法上下文） |
| `protocol.handler_missing` / `protocol.request_invalid` | Warn | 自定义协议缺少处理器或请求无效 |
| `protocol.panic` | Error | 自定义协议处理器 panic，含堆栈 |

`LogRecord` 包含事件时间、级别、事件名、窗口 ID、方法名、错误和堆栈。不主动记录参数、返回值、桥接令牌或请求 URL，但错误文本和堆栈可能包含业务信息；导出前按业务需要脱敏。生命周期事件不会回放，处理器应在启动前安装。IPC 和自定义协议诊断适用于桌面 WebView。

处理器在产生事件的 goroutine 上同步执行，可能是 UI 主线程。必须并发安全、快速返回，不要调用 UI 方法或触发 panic。网络上传、慢磁盘写入应由应用的日志系统通过有界队列处理，并明确丢弃/背压策略，避免阻塞界面。文件轮转和日志保留由应用配置。

此接口覆盖核心生命周期、窗口、IPC 和自定义协议；CLI、平台底层和独立插件仍可能使用各自的日志渠道，不受 `SetLogHandler` 统一过滤。
