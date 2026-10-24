# shadcn / Base UI · Maia 原生组件

Golet 的可选 Go 组件层，导入 `github.com/HyperparameterSI/golet/backend/ui/shadcn`。所有运行时 UI 都由 Golet 绘制；不使用 WebView、React、CSS 解释器或 JavaScript。原生的焦点、输入法、无障碍、剪贴板及窗口生命周期由 `ui` 提供。

以 [shadcn 官方 Base UI Maia 注册表](https://ui.shadcn.com/r/styles/base-maia/registry.json) 2026-10-24 的 **65 个目录项**为固定基准。源码、每份文件的 SHA-256 和上游 MIT 署名保存在 `reference/`；这些文件不嵌入应用。其中 Form 是上游无源码的占位项，Golet 提供原生表单组合入口。

## 使用

```go
import (
    "github.com/HyperparameterSI/golet/backend/ui"
    "github.com/HyperparameterSI/golet/backend/ui/shadcn"
)

// 每个窗口保存一次主题及表单数据，在 ui.View 中构建控件。
theme := shadcn.Theme(false)
email := ""
view := ui.View(func(c *ui.Context) {
    c.SetTheme(theme)
    shadcn.Card(c).Children(func() {
        shadcn.CardHeader(c).Children(func() {
            shadcn.CardTitle(c, "Create workspace")
            shadcn.CardDescription(c, "Build something native.")
        })
        shadcn.CardContent(c).Children(func() {
            shadcn.Field(c, "Email", func() {
                shadcn.Input(c, &email, shadcn.InputOptions{
                    Placeholder: "you@example.com",
                })
            })
        })
        shadcn.CardFooter(c).Children(func() {
            shadcn.Button(c, "Create").OnClick(func() {
                c.Toast("Created")
            })
        })
    })
    shadcn.Sonner(c)
})
_ = view // 传给 golet.WindowOptions.Content。
```

采用 Maia 的胶囊按钮/输入框、16 DIP 卡片、36 DIP 默认控件高度、24/32/40 DIP 按钮尺寸、14 DIP 正文字号和 3 DIP 焦点环。按钮和徽章支持 Default、Outline、Secondary、Ghost、Destructive、Link。`Theme(true)` 提供暗色外观；自定义主题仍可使用，但改变颜色或材质后不再是 Maia 的视觉基准。

Studio → 原生控件 → **shadcn · Maia** 可选择全部 65 个目录项。默认使用参考中性色，支持 Dark 和“跟随 Studio 主题”；后者用于验证主题适配。

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
CGO_ENABLED=0 go test ./backend/ui/shadcn ./backend/examples/studio
CGO_ENABLED=0 go vet ./backend/ui/shadcn ./backend/examples/studio
CGO_ENABLED=0 GOOS=js GOARCH=wasm go build ./backend/ui/shadcn
```

测试覆盖 65 个 Studio 页面在明暗主题下构建，以及按钮尺寸和禁用、键盘选择、日历导航、原生右键菜单、对话框与确认框、粘贴过滤、命令只执行一次、轮播和减少动态效果。交互测试使用 `ui.Tester` 的事件与虚拟时钟，不使用 sleep。
