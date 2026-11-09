# Golet Forma 原生组件

Golet 的可选 Go 组件层，导入 `github.com/HyperparameterSI/golet/backend/ui/forma`。所有运行时 UI 都由 Golet 绘制；不使用 WebView、React、CSS 解释器或 JavaScript。原生的焦点、输入法、无障碍、剪贴板及窗口生命周期由 `ui` 提供。

以 [shadcn 官方 Base UI Maia 注册表](https://ui.shadcn.com/r/styles/base-maia/registry.json) 2026-10-24 的 **65 个目录项**为固定基准。源码、每份文件的 SHA-256 和上游 MIT 署名保存在 [`docs/references/shadcn-base-maia`](../../../docs/references/shadcn-base-maia)；这些文件不嵌入应用。其中 Form 是上游无源码的占位项，Golet 提供原生表单组合入口。

## 使用

### 原生 CodeBlock

沿用 `CodeBlock(c, title, code)`，可通过可选参数开启扩展功能：

```go
// Author: HyperparameterSI
// 在 View 外创建一次；不传 Highlight 时不加载语言高亮器。
highlight := forma.HighlightCode("go")
expanded := false
// 在 View 内：
forma.CodeBlock(c, "main.go", source, forma.CodeBlockOptions{
    LineNumbers: true, StartLine: 98,
    Highlight: highlight, HighlightKey: "go",
    Wrap: false, Expanded: &expanded, CollapsedHeight: 120,
})
```

支持默认/Small 尺寸、Destructive 错误样式、自定义 Actions/Footer、自动换行、横纵滚动、键盘方向键/PageUp/PageDown、展开收起，以及按零起始行索引设置 `Lines`（`CodeLineAdded`、`CodeLineRemoved`、`CodeLineHighlighted`）。代码区始终从左到右，标题跟随界面方向。`MaxHeight` 默认 320，负数取消高度上限。

复制按钮保留原始空行、尾换行和缩进，不附带行号及差异标记；成功/失败反馈持续两秒。语法高亮按行缓存，追加内容只重新处理新增或变化的后缀；修改高亮器配置时需同步更新 `HighlightKey`。普通文本使用单个文本节点。Studio → Forma → code-block 展示主题、高亮、差异、错误及展开效果。

参考 [官方 Base Code Block](https://ui.shadcn.com/docs/components/base/code-block)，原生实现不依赖 Shiki/React。目前无标题操作栏保持可见，尚未移植网页的悬停浮动操作、滚动边缘渐隐及 ANSI 终端转义着色；不宣称像素级完全一致。

### 基础组合

```go
import (
    "github.com/HyperparameterSI/golet/backend/ui"
    "github.com/HyperparameterSI/golet/backend/ui/forma"
)

// 每个窗口保存一次主题及表单数据，在 ui.View 中构建控件。
theme := forma.Theme(false)
email := ""
view := ui.View(func(c *ui.Context) {
    c.SetTheme(theme)
    forma.Card(c).Children(func() {
        forma.CardHeader(c).Children(func() {
            forma.CardTitle(c, "Create workspace")
            forma.CardDescription(c, "Build something native.")
        })
        forma.CardContent(c).Children(func() {
            forma.Field(c, "Email", func() {
                forma.Input(c, &email, forma.InputOptions{
                    Placeholder: "you@example.com",
                })
            })
        })
        forma.CardFooter(c).Children(func() {
            forma.Button(c, "Create").OnClick(func() {
                c.Toast("Created")
            })
        })
    })
    forma.Sonner(c)
})
_ = view // 传给 golet.WindowOptions.Content。
```

采用 Maia 的胶囊按钮/输入框、16 DIP 卡片、36 DIP 默认控件高度、24/32/40 DIP 按钮尺寸、14 DIP 正文字号和 3 DIP 焦点环。按钮和徽章支持 Default、Outline、Secondary、Ghost、Destructive、Link。`Theme(true)` 提供暗色外观；自定义主题仍可使用，但改变颜色或材质后不再是 Maia 的视觉基准。

Studio → 原生控件 → **Forma · Maia** 可选择全部 65 个目录项。默认使用参考中性色，支持 Dark 和“跟随 Studio 主题”；后者用于验证主题适配。

## 移植范围与验收

当前完成的是 65 个目录项的原生实现入口与可运行演示，**不是全部 React API、所有变体和逐像素一致性的完成声明**。以下差异仍需继续补齐：

- 全组件逐状态的浏览器/原生截图对比尚未完成。已检查基础表单的明暗原生截图；字体栅格化、部分图标和阴影仍可能不同。
- Calendar 支持单选、多选、日期范围及禁用日期；方向键移动焦点、Enter 确认，Home/End 跳到周首尾。年份下拉、多月视图和范围长度限制还未提供。Table 支持虚拟化行与选择、基础排序状态；列调整和单元格块选择使用 `ui.Table`。
- Drawer 当前为可拖拽关闭的底部面板；未提供上游全部方向、嵌套堆叠和 snap points API。Carousel 支持分页、按钮、键盘与拖拽；未实现 Embla 的全部选项。
- Command 提供搜索、分组和命令激活；Dropdown/ContextMenu 提供勾选、子菜单和键盘焦点组。上游所有子菜单悬停延迟、复杂多级焦点返回仍需更完整的验证。
- Chart 提供原生 `ui/chart` 容器、图例和 tooltip；不兼容 Recharts 的 React API。Form/Field 使用 Go 数据绑定和错误描述，不包含 React Hook Form 适配层。
- Message、Bubble、Attachment、Questionnaire、Trace 提供基础原生组合；上游全部尺寸/状态、反应条、流式展示和复杂问卷 API 尚未逐项对齐。
- ButtonGroup 用 `GroupItem(element, index, count)` 连接子项圆角并消除重复边框；InputGroup 用 `InputGroupInput` 共享外轮廓与焦点环。附加区域仍需要由调用方显式组合。NativeSelect 使用原生 Golet 选择器，不是 HTML `<select>`。

这些项目没有被隐藏在“已移植”计数里；后续应按“外观 → 键盘/焦点 → 边界状态 → 明暗/RTL → 截图”逐项验收。完整目录映射见 [PORTING.md](PORTING.md)。

## 验证

```sh
CGO_ENABLED=0 go test ./backend/ui/forma ./backend/examples/studio
CGO_ENABLED=0 go vet ./backend/ui/forma ./backend/examples/studio
CGO_ENABLED=0 GOOS=js GOARCH=wasm go build ./backend/ui/forma
```

测试覆盖 65 个 Studio 页面在明暗主题下构建，以及按钮尺寸和禁用、键盘选择、日历导航、原生右键菜单、对话框与确认框、粘贴过滤、命令只执行一次、轮播和减少动态效果。交互测试使用 `ui.Tester` 的事件与虚拟时钟，不使用 sleep。

## 原生主题与内置配色

`Theme(false)` 保持旧版中性 Maia 用法。使用 `NewTheme` 可选择 **Neutral、Zinc、Stone、Blue、Violet、Emerald、Rose**，每套都支持明暗模式。这些是 Golet 参考 Tailwind 色系设计的预设，不是 shadcn 官方所有 preset 的完整复制。

```go
radius := float32(8) // 默认 16；8 更柔和紧凑；0 为直角
light := forma.NewTheme(forma.ThemeOptions{
    Preset: forma.Violet,
    Radius: &radius,
    FontSize: 14,
    Tokens: map[string]ui.Color{
        "primary": ui.Hex("#6d28d9"),
        "primary-foreground": ui.Hex("#ffffff"),
        "card": ui.Hex("#faf8ff"),
        "popover": ui.Hex("#ffffff"),
        "ring": ui.Hex("#8b5cf6"),
    },
})
dark := forma.NewTheme(forma.ThemeOptions{Preset: forma.Violet, Dark: true})
// 在窗口状态中保存 light/dark，仅在用户修改配置时重建。
// 构建视图时 c.SetTheme(light) 或 c.SetTheme(dark)。
```

遵循 [shadcn 的语义颜色约定](https://ui.shadcn.com/docs/theming)，背景和文字颜色分别配置。主色更改后，请同时选择合适的 `primary-foreground`；自定义颜色不会自动计算对比度。`Tokens` 在构造时复制，允许显式透明颜色，不会因为零值而覆盖回默认色。

| Token | 原生用途 |
| --- | --- |
| `background` / `foreground` | 窗口背景、默认文字 |
| `card` / `card-foreground` | 卡片背景、卡片内容 |
| `popover` / `popover-foreground` | 选择菜单、弹出菜单、对话面板 |
| `primary` / `primary-foreground` | 主按钮、选中控件、主徽章 |
| `secondary` / `secondary-foreground` | 辅助按钮和徽章 |
| `muted` / `muted-foreground` | 次级背景、描述文字 |
| `accent` / `accent-foreground` | Ghost、菜单高亮、Toggle |
| `border` / `input` / `ring` | 边框、输入区域、键盘焦点 |
| `destructive` | 错误和危险操作 |
| `sidebar` / `sidebar-foreground` / `sidebar-primary` / `sidebar-primary-foreground` / `sidebar-accent` / `sidebar-border` / `sidebar-ring` | 导航区域及其交互状态 |
| `chart-1` … `chart-5` | `ChartColors(c)`，图例未传颜色时使用 |

`forma.Color(c, "primary")` 可在自定义绘制中读取配置，`ChartColors(c)` 返回独立色板切片。自定义曲线需将色板传给 `ui/chart`；主题不会修改调用方显式指定的曲线颜色。

参考 [Tailwind 的设计变量机制](https://tailwindcss.com/docs/theme)，我们将配色与组件解耦，但提供的是 Go 配置 API，并不解析 CSS、class 字符串或执行 Tailwind。圆角按现有 Maia 原生比例缩放，保留默认视觉尺寸；没有照搬最新网页的全部 `radius-*` 比例。圆形头像、单选圆点、滑块手柄和开关保留其功能形状。`GroupItem` 的显式连接圆角仍由调用方控制。

`Font` 可设置系统可用字体，`FontSize` 按比例调整组件文字；不会自动扩大整个布局，较大字体应配合 `Large` 按钮和显式尺寸。主题只负责配色与度量，不会自动启用 Liquid Glass 材质。

Studio 的 **Forma / Maia** 页面可切换配色、Maia / Soft / Sharp 圆角、14 / 16 DIP 字号和 Dark。选择“跟随 Studio 主题”时禁用这些独立配置，使用 Studio 当前主题。

## 从旧包迁移到 Forma

将导入路径 `github.com/HyperparameterSI/golet/backend/ui/shadcn` 改为 `github.com/HyperparameterSI/golet/backend/ui/forma`，调用前缀改为 `forma`。组件 API 和主题配置保持原有结构；主题内部命名空间改为 `forma.`，通过 `ThemeOptions.Tokens` 配置时仍使用 `primary` 等语义名称，无需加前缀。

Studio 示例已迁移为 `backend/examples/studio/forma.go`，入口为 **Forma · Maia**。基础主题函数更名为 `ui.FormaTheme`；`ui.ShadcnTheme` 作为弃用别名保留，以兼容已使用该基础函数的应用。原生组件旧包路径不再保留。

本目录只有 Go 实现、测试、文档和许可证。TSX 上游参考快照已单独归档到 `docs/references/shadcn-base-maia`，不参与原生组件运行。前端项目仍可使用真正的 shadcn React 组件，两者互不依赖。

### 可选浏览器主题适配

`github.com/HyperparameterSI/golet/backend/ui/forma/browser` 提供 `Launch`、
`Theme` 和 `SetTheme`，把同一份 `*ui.Theme` 传给独立 CEF 浏览器窗口。
目前支持 macOS Alloy 的配色、明暗外观、字体与圆角；不会影响网页内容，
也不是嵌入式 `ui.Element`。普通 Forma 应用不依赖此包或 CEF。
用法与平台限制见 [Chromium 插件文档](../../plugins/chromium/README.md#跟随-forma-主题macos-alloy)。

## AI Message 与 Model Selector

<!-- Author: HyperparameterSI -->

参考官方 [AI Message](https://ui.shadcn.com/docs/ai-elements/base/message) 与
[Model Selector](https://ui.shadcn.com/docs/ai-elements/base/model-selector)，提供纯 Go 原生组合；无 React、WebView 或模型 SDK 依赖。

```go
// Author: HyperparameterSI
// 窗口状态，在 View 外保存。
feedback := forma.FeedbackNone
model, effort := "reasoner", "medium"
models := []forma.ModelOption{
    {ID: "reasoner", Name: "Reasoner", Provider: "My provider",
        Efforts: []forma.ModelEffort{
            {ID: "low", Name: "Low"}, {ID: "medium", Name: "Medium"},
        }},
    {ID: "chat", Name: "Chat", Provider: "My provider"},
}
// View 内：
forma.AIMessage(c, forma.MessageUser, "帮我创建一个窗口")
forma.AIMessage(c, forma.MessageAssistant, reply, forma.AIMessageOptions{
    Streaming: generating,
    Feedback: &feedback,
    OnRetry: retry,
})
forma.ModelSelector(c, &model, &effort, models)
```

`AIMessage` 将用户文本置于末端的浅色气泡，助手内容使用原生 Markdown，围栏代码交给 `CodeBlock`；`Streaming` 使用增量 Markdown 渲染并禁用消息操作。支持复制原文及两秒结果反馈、互斥且可取消的好/差评价、重试、编辑和自定义 Footer。`OnEdit` 收到原始文本。`MessageContent`、`MessageHeader`、`MessageFooter` 和 `MessageResponse` 可用于自定义组合，旧 `Message(c, author, end, build)` 保持兼容。底部操作使用复制、点赞、点踩、重试和编辑线框图标，带悬停提示、无障碍名称及选中状态。用户气泡与代码块使用主题的 muted 色，正文保持无边框。

`ModelSelector` 的两个菜单分别绑定模型 ID 与思考强度 ID，支持提供商分组、描述、禁用条目、当前选择标记和原生键盘导航，选中后关闭。模型不提供 `Efforts` 时隐藏强度菜单；切换模型时保留仍受支持的强度，否则使用第一个可用选项或空值。模型 ID 应唯一且非空。`OnChange` 在选择改变时收到模型和强度。整个控件可通过 `Disabled` 禁用。

控件不请求模型、不保存密钥。应用需自行把选择映射到已支持的后端模型参数；此 UI 不会自动使 Assistant 插件或供应商支持任意思考强度。重试、编辑和评价也仅调用应用提供的回调。

Studio → Forma 的 `ai-message` 和 `model-selector` 有交互演示；基础 `message` 页面也显示 AI 消息组合。`AIComponents()` 单独列出 AI 组合入口，保持 `Components()` 对固定基础注册表的映射不变。

CodeBlock 的默认外观采用浅灰（深色主题为深灰）面板、大圆角、终端图标标题及图标复制按钮。传入 `WrapToggle: &wrap` 可显示自动换行切换；传入 `OnDownload: func(source string) { ... }` 可显示下载按钮。下载回调收到原始文本，保存位置和错误处理由应用决定。Studio 的 Terminal 示例使用系统保存对话框导出文本。图标保留无障碍标签、键盘焦点和 Tooltip，复制反馈显示勾选图标。

### 模型输入面板

`ModelSelector` 默认采用带细边框的胶囊按钮；`ModelComposer` 把它与无边框多行输入、添加和发送图标组合为浅灰圆角面板，自动跟随明暗主题：

```go
// Author: HyperparameterSI
forma.ModelComposer(c, &prompt, &model, &effort, models, forma.ModelComposerOptions{
    Placeholder: "让助手帮你修改代码…",
    OnAdd: openAttachmentPicker,
    OnSubmit: func(text, model, effort string) {
        // 将本次输入快照交给应用的后台任务。
    },
})
```

`OnAdd` 或 `OnSubmit` 未设置时，对应按钮禁用；空白输入、无效或禁用模型不能发送。`Disabled` 禁用输入及整组操作。Enter 保持多行输入行为，点击发送才提交。发送不自动清空文本，应用可在请求成功后自行清空。`OnModelChange` 接收模型或强度变更。Studio 的按钮连接演示回调，不自动发起 AI 请求或附件上传。

### Prompt Input

参考 [官方 Prompt Input](https://ui.shadcn.com/docs/ai-elements/base/prompt-input)，`PromptInput` 提供独立的纯 Go 消息输入控件：

```go
// Author: HyperparameterSI
// 在窗口状态中保存，不能每帧重新创建。
input := forma.PromptInputState{Model: "reasoner", Effort: "medium"}
// View 内：
forma.PromptInput(c, &input, forma.PromptInputOptions{
    Models: models, // 可省略：不显示模型菜单，也不要求选择模型。
    OnAttach: openAttachmentPicker,
    OnSubmit: func(request forma.PromptSubmission) {
        // 启动后台任务，将文本、附件元数据及模型选择交给应用处理。
    },
})
```

- Enter 发送，Shift+Enter 换行，空白且无附件时禁止发送；只有附件也可提交。
- `Attachments` 是应用持有的元数据，可展示文件名、描述并逐项移除；`OnRemove` 通知应用清理资源。`OnDropFiles` 接收拖入的路径，不读取文件。
- `OnAttach` 由应用接入原生文件选择器；文件读取、大小/格式校验、上传以及错误处理由应用负责，控件不做隐式 I/O。
- `PromptSubmitted` / `PromptStreaming` 显示等待/生成状态并替换为停止按钮，`OnStop` 由应用连接取消操作；`Disabled` 禁用交互，`Error` 显示错误。`Tools` 可插入自定义工具按钮。
- 提交回调收到文本、模型、强度及独立的附件列表快照；不会自动清空草稿。应用可在提交成功后清空，并通过 `Window.Update` 更新 UI 状态。
- 原有 `ModelComposer` 继续保持 Enter 换行行为，迁移时请明确选择所需键盘语义。

Studio → Forma → `prompt-input` 演示附件、文件拖入、发送/停止及明暗主题。添加按钮生成示例附件，拖入文件只保存元数据，发送使用模拟响应，均不会上传文件或调用模型。

See [Native AI Elements](AI_ELEMENTS.md) for the complete 21-entry mapping, usage, and differences from the web examples.

## Native charts

`Chart(c, &state, options)` adds themed grouped bars (`ChartBar`), gradient areas
(`ChartArea`) and lines (`ChartLine`) on top of `ui/chart`. No React, Recharts,
JavaScript or WebView is required. Studio → Forma → chart demonstrates grouped
monthly bars and a 90/30/7-day area chart with smooth, linear and step curves.

```go
// Author: HyperparameterSI
state := forma.ChartState{Selected: -1} // Keep outside the view callback.
// Inside the view:
forma.Chart(c, &state, forma.ChartOptions{
    Kind: forma.ChartBar,
    Labels: []string{"Jan", "Feb", "Mar"},
    Series: []forma.ChartSeries{
        {Label: "Desktop", Values: []float64{186, 305, 237}},
        {Label: "Mobile", Values: []float64{80, 200, 120}},
    },
    Height: 260,
})
```

Zero series colors inherit `ChartColors(c)`; explicit colors override them.
`Curve`, `ShowYAxis`, `HideGrid`, `HideLegend` and `FormatValue` customize the
presentation. Hover or focus the plot and use Left/Right/Home/End to inspect
values; Escape clears selection. `OnSelect` receives the category index, or -1
when cleared. Missing/NaN/infinite values produce gaps; negative values are
included in the automatic zero-based domain. Categories are evenly spaced;
filter your data before passing it to Chart to change the visible date range.
Keep supplied slices immutable for the duration of a frame.

This is a native adaptation of the [shadcn Chart](https://ui.shadcn.com/docs/components/base/chart)
and [area examples](https://ui.shadcn.com/charts/area), not a Recharts API clone.
Areas overlap rather than stack; pie/radar/radial, animation and stacking are
not exposed by this convenience API. Existing `ChartContainer`, `ChartLegend`,
`ChartTooltip` and `ui/chart` remain available for custom compositions.
