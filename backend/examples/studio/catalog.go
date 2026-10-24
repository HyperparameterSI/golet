// Author: HyperparameterSI
package main

// Demo describes a real, runnable component demonstration.
type Demo struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Category    string   `json:"category"`
	Description string   `json:"description"`
	Controls    []string `json:"controls"`
	Source      string   `json:"source"`
}

var demos = []Demo{
	{"Shadcn", "shadcn · Maia", "原生控件", "65 个官方目录项的纯 Go 原生演示与移植对照。", []string{"Button", "Input", "Dialog", "Select", "Tabs"}, "backend/examples/studio/shadcn.go"},
	{"BorderBeam", "边框光束", "插件", "流动色彩、呼吸光效与主题边缘辉光，纯 Go 原生绘制。", []string{"borderbeam.Box", "Medium", "Small", "Line", "PulseInner", "PulseOutside", "Paused"}, "backend/examples/studio/borderbeam.go"},
	{"Themes", "原生主题实验室", "原生控件", "切换 Maia 中性与 Liquid Glass，体验明暗外观、背景材质与真实控件交互。", []string{"ShadcnTheme", "glass.Theme", "Light", "Dark", "Material", "Button", "TextInput", "Select", "Tabs"}, "backend/examples/internal/gallerydemo/themes.go"},
	{"Overview", "基础与布局", "原生控件", "响应式布局、状态更新和主题。", []string{"Row", "Column", "Box", "Grid", "Spacer", "Text", "Button", "Badge", "Avatar", "Progress", "Spinner"}, "backend/examples/internal/gallerydemo/demo.go"},
	{"Controls", "表单与选择", "原生控件", "完整表单：输入、选择、数值、日期与颜色。", []string{"Button", "PrimaryButton", "IconButton", "Checkbox", "CheckboxGroup", "Switch", "Toggle", "ToggleGroup", "Radio", "RadioGroup", "Select", "Combobox", "Autocomplete", "SearchField", "TokenField", "Slider", "StepSlider", "RangeSlider", "NumberInput", "Stepper", "Rating", "Meter", "Calendar", "DateInput", "TimeInput", "ColorWell", "Tabs", "Segmented", "Toolbar", "Breadcrumbs", "Collapsible", "Accordion", "AccordionItem", "Field", "Fieldset", "Label"}, "backend/examples/internal/gallerydemo/demo.go"},
	{"Text", "文本与输入", "原生控件", "文本排版、编辑、选择、验证与多语言字符。", []string{"TextInput", "TextArea", "RichText", "Link", "FontFeatures", "Selection"}, "backend/examples/internal/gallerydemo/demo.go"},
	{"Rich Text", "富文本编辑器", "高级组件", "可编辑文档、格式工具栏、撤销与重做。", []string{"richtext.Editor", "Document", "Formatting"}, "backend/examples/internal/gallerydemo/richtext.go"},
	{"List", "列表、表格与树", "原生控件", "一万行虚拟列表、排序表格、树、网格与拖放。", []string{"List", "Table", "OutlineTable", "Tree", "TreeItem", "GridView", "Selection", "Drag", "Drop"}, "backend/examples/internal/gallerydemo/demo.go"},
	{"Charts", "图表", "高级组件", "折线、面积、柱状与饼图等数据可视化。", []string{"LineChart", "AreaChart", "BarChart", "PieChart", "ScatterChart", "Sparkline"}, "backend/examples/internal/gallerydemo/charts.go"},
	{"Markdown", "Markdown", "高级组件", "Markdown 文档、代码块与交互链接。", []string{"markdown.View", "CodeBlock", "Tables"}, "backend/examples/internal/gallerydemo/markdown.go"},
	{"Styling", "样式与布局", "原生控件", "网格、边框、渐变、阴影和文字装饰。", []string{"Grid", "Border", "Gradient", "Shadow", "Flex", "Scroll"}, "backend/examples/internal/gallerydemo/demo.go"},
	{"Drawing", "绘图与图像", "原生控件", "自定义绘制、路径、SVG 与位图。", []string{"Canvas", "Painter", "SVG", "Image"}, "backend/examples/internal/gallerydemo/demo.go"},
	{"Glass", "玻璃效果", "插件", "可调参数的玻璃折射与材质。", []string{"glass.Effect", "Refraction"}, "backend/examples/internal/gallerydemo/demo.go"},
	{"Overlays", "弹层与反馈", "原生控件", "弹窗、菜单、提示和原生对话框。", []string{"Modal", "AlertDialog", "Popover", "Tooltip", "MenuButton", "ContextMenu", "Toast", "Dialog"}, "backend/examples/internal/gallerydemo/demo.go"},
	{"Motion", "动画与手势", "原生控件", "转场、弹簧、拖动与动画状态。", []string{"Transition", "Spring", "Animation", "Gesture"}, "backend/examples/internal/gallerydemo/demo.go"},
	{"Languages", "导航与国际化", "原生控件", "路由历史、语言选择与从右到左布局。", []string{"Router", "Sidebar", "SidebarSection", "SidebarItem", "BackButton", "ForwardButton", "i18n", "RTL"}, "backend/examples/internal/gallerydemo/languages.go"},
	{"Advanced", "缓冲区与分栏", "原生控件", "共享文本缓冲区、Outline、颜色选择器及不同分栏方式。", []string{"TextInputBuffer", "TextAreaBuffer", "Outline", "ColorPicker", "Divider", "ScrollHorizontal", "Split", "SplitPercent", "SplitVertical", "SplitVerticalPercent"}, "backend/examples/studio/native.go"},
	{"Editor", "代码编辑与 Diff", "高级组件", "多语言高亮、补全、诊断、代码差异与编辑。", []string{"editor.View", "DiffView", "Buffer", "Diagnostics", "Completion"}, "backend/examples/internal/editordemo/demo.go"},
	{"Loading", "Loading 动画", "插件", "29 种指示器，可调整尺寸、颜色、缓动与暂停。", []string{"loading.All", "Size", "Color", "Easing", "Pause"}, "backend/examples/internal/loadingdemo/demo.go"},
	{"WebView", "嵌入式 WebView", "高级组件", "原生工具栏与系统 WebView 组合，展示离线 HTML。", []string{"ui.WebView", "NewWebView", "LoadHTML", "Reload", "Zoom"}, "backend/examples/studio/native.go"},
	{"Terminal", "终端", "插件", "真实 Shell 与原生终端控件；首次使用可能获取终端运行库。", []string{"terminal.New", "terminal.View", "PTY"}, "backend/examples/studio/native.go"},
	{"Browser", "Golet Browser", "插件", "独立 CEF 窗口：Alloy 自定义风格与 Chrome 插件兼容风格。", []string{"StyleAlloy", "StyleChrome", "Navigate", "Back", "Forward", "Reload", "NewTab", "Close"}, "backend/examples/studio/browser.go"},
}
