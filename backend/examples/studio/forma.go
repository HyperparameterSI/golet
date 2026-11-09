// Author: HyperparameterSI
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/HyperparameterSI/golet/backend"
	"github.com/HyperparameterSI/golet/backend/examples/internal/demotheme"
	"github.com/HyperparameterSI/golet/backend/ui"

	"github.com/HyperparameterSI/golet/backend/ui/forma"
)

func openForma() *golet.Window {
	return golet.NewWindow(golet.WindowOptions{Title: "Studio · Forma / Maia", Width: 1040, Height: 840, Content: ui.View(newFormaView())})
}
func newFormaView(initial ...string) func(*ui.Context) {
	aiPatterns := newAIPatternDemo()
	chartDemo := newChartDemo()
	selected := "button"
	if len(initial) > 0 {
		selected = initial[0]
	}
	dark, follow := false, false
	preset, corners := "Neutral", "Maia"
	fontScale := "14"
	type themeKey struct {
		preset, corners, font string
		dark                  bool
	}
	themes := make(map[themeKey]*ui.Theme)
	presetNames := make([]string, 0, len(forma.Presets()))
	for _, name := range forma.Presets() {
		presetNames = append(presetNames, string(name))
	}
	text, choice, answer := "", "Go", ""
	on, expanded, opened := false, false, false
	value := float64(40)
	tab, page, slide, question := 0, 1, 0, 0
	date := time.Date(2026, 10, 24, 0, 0, 0, 0, time.UTC)
	calendarMode := "Single"
	var dates []time.Time
	dateRange := forma.DateRange{}
	split := float32(180)
	list := ui.ListState{}
	scroll := ui.ScrollState{}
	codeWrap, codeNumbers, codeExpanded := false, true, false
	terminalWrap := true
	codeHighlight := forma.HighlightCode("go")
	messageFeedback := forma.FeedbackNone
	messageStreaming := false
	messageStatus := ""
	modelID, modelEffort := "demo-reasoner", "medium"
	prompt := forma.PromptInputState{Model: "demo-reasoner", Effort: "medium"}
	promptStatus := forma.PromptReady
	promptResult := ""
	models := []forma.ModelOption{
		{ID: "demo-reasoner", Name: "Reasoner", Provider: "Golet Demo", Description: "支持思考强度的示例模型", Efforts: []forma.ModelEffort{{ID: "low", Name: "Low"}, {ID: "medium", Name: "Medium"}, {ID: "high", Name: "High"}}},
		{ID: "demo-chat", Name: "Chat", Provider: "Golet Demo", Description: "不提供思考强度选项"},
		{ID: "local", Name: "Local model", Provider: "Local", Description: "尚未配置", Disabled: true},
	}
	return func(c *ui.Context) {
		demotheme.Apply(c)
		if !follow {
			key := themeKey{preset, corners, fontScale, dark}
			theme := themes[key]
			if theme == nil {
				radius := float32(16)
				if corners == "Soft" {
					radius = 8
				}
				if corners == "Sharp" {
					radius = 0
				}
				font := float32(14)
				if fontScale == "16" {
					font = 16
				}
				theme = forma.NewTheme(forma.ThemeOptions{Preset: forma.Preset(preset), Dark: dark, Radius: &radius, FontSize: font})
				themes[key] = theme
			}
			c.SetTheme(theme)
		}
		ui.Column(c).Fill().Background(c.Theme().Background).Padding(32).Gap(24).Children(func() {
			ui.Row(c).Gap(16).AlignItems(ui.Center).Children(func() {
				ui.Text(c, "Forma / Maia").FontSize(26).Bold().Grow(1)
				ui.Checkbox(c, &dark, "Dark")
				ui.Checkbox(c, &follow, "跟随 Studio 主题")
			})
			ui.Text(c, "65 个基础目录项 + AI 组合 · 纯 Go 原生组件").TextColor(c.Theme().TextMuted)
			ui.Row(c).Gap(12).AlignItems(ui.Center).Children(func() {
				forma.Select(c, &selected, append(forma.Components(), forma.AIComponents()...)).Width(240).Label("Component")
				forma.Select(c, &preset, presetNames).Width(160).Label("Theme palette").Disabled(follow)
				forma.Select(c, &corners, []string{"Maia", "Soft", "Sharp"}).Width(140).Label("Corner style").Disabled(follow)
				forma.Select(c, &fontScale, []string{"14", "16"}).Width(80).Label("Font size").Disabled(follow)
				for _, name := range []string{"primary", "secondary", "accent", "ring"} {
					ui.Box(c).Size(16, 16).Radius(8).Background(forma.Color(c, name)).Label(name)
				}
			})
			ui.Scroll(c).Grow(1).Children(func() {
				forma.Card(c).Children(func() {
					forma.CardHeader(c).Children(func() {
						forma.CardTitle(c, selected)
						forma.CardDescription(c, "交互演示；可切换明暗外观。完整对齐状态见组件移植清单。")
					})
					forma.CardContent(c).MinHeight(230).Children(func() {
						menu := []forma.MenuItem{{Label: "Profile", Action: func() { c.Toast("Profile selected") }}, {Label: "Notifications", Checked: &on}, {Separator: true}, {Label: "Disabled action", Disabled: true}}
						switch selected {
						case "button":
							ui.Row(c).Gap(8).Children(func() {
								for i, name := range []string{"Default", "Outline", "Secondary", "Ghost", "Destructive", "Link"} {
									forma.Button(c, name, forma.ButtonOptions{Variant: forma.Variant(i)}).OnClick(func() { c.Toast(name) })
								}
							})
							ui.Row(c).Gap(8).Children(func() {
								for _, size := range []forma.Size{forma.XS, forma.Small, forma.Normal, forma.Large} {
									forma.Button(c, "Button", forma.ButtonOptions{Size: size})
								}
								forma.Button(c, "Disabled", forma.ButtonOptions{Disabled: true})
							})
						case "badge":
							ui.Row(c).Gap(8).Children(func() {
								for i, name := range []string{"Default", "Outline", "Secondary", "Ghost", "Destructive", "Link"} {
									forma.Badge(c, name, forma.Variant(i))
								}
							})
						case "card":
							forma.Card(c).Children(func() {
								forma.CardHeader(c).Children(func() { forma.CardTitle(c, "Create project"); forma.CardDescription(c, "Start building with Golet.") })
								forma.CardContent(c).Children(func() { forma.Input(c, &text, forma.InputOptions{Placeholder: "Project name"}) })
								forma.CardFooter(c).Children(func() { forma.Button(c, "Create") })
							})
						case "input":
							forma.Input(c, &text, forma.InputOptions{Placeholder: "you@example.com"})
							forma.Input(c, &text, forma.InputOptions{Invalid: true})
							forma.Input(c, &text, forma.InputOptions{Disabled: true})
						case "textarea":
							forma.Textarea(c, &text, forma.InputOptions{Placeholder: "Write a message…"}).Height(120)
						case "checkbox":
							forma.Checkbox(c, &on, "Accept terms and conditions")
						case "switch":
							forma.Switch(c, &on).Label("Airplane mode")
						case "radio-group":
							forma.RadioGroup(c, func() {
								for _, option := range []string{"Go", "TypeScript", "Rust"} {
									forma.Radio(c, &choice, option, option)
								}
							})
						case "select":
							forma.Select(c, &choice, []string{"Go", "TypeScript", "Rust"})
						case "native-select":
							forma.NativeSelect(c, &choice, []string{"Go", "TypeScript", "Rust"})
						case "combobox":
							forma.Combobox(c, &choice, []string{"Go", "TypeScript", "Rust"})
						case "tabs":
							forma.Tabs(c, &tab, "Account", "Password", "Notifications")
							forma.CardDescription(c, fmt.Sprintf("Selected panel: %d", tab+1))
						case "toggle":
							forma.Toggle(c, &on, "Bold")
						case "toggle-group":
							forma.ToggleGroup(c, func() { forma.Toggle(c, &on, "Bold"); forma.Toggle(c, &expanded, "Italic") })
						case "slider":
							forma.Slider(c, &value, 0, 100).Width(320)
							ui.Textf(c, "%.0f", value)
						case "progress":
							forma.Progress(c, value/100).Width(320)
							forma.Button(c, "Advance").OnClick(func() { value = float64((int(value) + 10) % 110) })
						case "calendar":
							forma.Select(c, &calendarMode, []string{"Single", "Multiple", "Range"}).Width(200)
							options := forma.CalendarOptions{}
							if calendarMode == "Multiple" {
								options.Multiple = &dates
							}
							if calendarMode == "Range" {
								options.Range = &dateRange
							}
							forma.Calendar(c, &date, options)
						case "accordion":
							forma.Accordion(c, func() {
								forma.AccordionItem(c, "Is it native?", &expanded, func() { ui.Text(c, "Yes. Every element is rendered by Golet, without a webview.") })
								forma.AccordionItem(c, "Does it support dark mode?", &on, func() { ui.Text(c, "Use the Dark switch above.") })
							})
						case "collapsible":
							forma.Collapsible(c, "Dependencies", &expanded, func() { ui.Text(c, "golet/backend/ui/forma") })
						case "dialog", "alert-dialog", "sheet", "drawer":
							forma.Button(c, "Open "+selected).OnClick(func() { opened = true })
							switch selected {
							case "dialog":
								forma.Dialog(c, &opened, "Edit profile", func() {
									forma.Input(c, &text, forma.InputOptions{Placeholder: "Name"})
									forma.Button(c, "Save").OnClick(func() { opened = false })
								})
							case "alert-dialog":
								forma.AlertDialog(c, &opened, "Are you sure?", "This is a demonstration; no data is removed.", func() { c.Toast("Confirmed") })
							case "sheet":
								forma.Sheet(c, &opened, forma.Right, "Settings", func() { forma.Input(c, &text) })
							case "drawer":
								forma.Drawer(c, &opened, "Move goal", func() { forma.Slider(c, &value, 0, 100); forma.Button(c, "Done").OnClick(func() { opened = false }) })
							}
						case "popover":
							b := forma.Button(c, "Open popover")
							b.OnClick(func() { opened = !opened })
							forma.Popover(c, b, &opened, func() { forma.Label(c, "Dimensions"); forma.Input(c, &text) })
						case "tooltip":
							b := forma.Button(c, "Hover or focus")
							forma.Tooltip(c, b, "Pure Go tooltip")
						case "hover-card":
							b := forma.Button(c, "@golet", forma.ButtonOptions{Variant: forma.Link})
							forma.HoverCard(c, b, func() { forma.CardTitle(c, "Golet"); forma.CardDescription(c, "Build native apps with Go.") })
						case "dropdown-menu":
							forma.DropdownMenu(c, "Open menu", menu)
						case "context-menu":
							forma.ContextMenu(c, ui.Box(c).Height(160).Border(1, c.Theme().Border).Radius(16).Center().Children(func() { ui.Text(c, "Right click or Shift+F10") }), menu)
						case "menubar":
							forma.Menubar(c, func() { forma.DropdownMenu(c, "File", menu); forma.DropdownMenu(c, "Edit", menu) })
						case "navigation-menu":
							forma.NavigationMenu(c, func() {
								forma.DropdownMenu(c, "Products", menu)
								forma.Button(c, "Documentation", forma.ButtonOptions{Variant: forma.Ghost})
							})
						case "command":
							forma.Command(c, &text, []forma.CommandItem{{ID: "calendar", Label: "Calendar", Group: "Suggestions", Run: func() { c.Toast("Calendar") }}, {ID: "settings", Label: "Settings", Group: "Preferences", Run: func() { c.Toast("Settings") }}})
						case "pagination":
							forma.Pagination(c, &page, 12)
							ui.Textf(c, "Page %d", page)
						case "sidebar":
							forma.Sidebar(c, &choice, func() {
								forma.SidebarItem(c, "Go", "Overview")
								forma.SidebarItem(c, "TypeScript", "Projects")
								forma.SidebarItem(c, "Rust", "Settings")
							}).Height(200)
						case "breadcrumb":
							forma.Breadcrumb(c, []string{"Home", "Components", "Maia"}, &tab)
						case "button-group":
							forma.ButtonGroup(c, func() {
								forma.GroupItem(forma.Button(c, "Previous", forma.ButtonOptions{Variant: forma.Outline}), 0, 2)
								forma.GroupItem(forma.Button(c, "Next", forma.ButtonOptions{Variant: forma.Outline}), 1, 2)
							})
						case "aspect-ratio":
							forma.AspectRatio(c, 16.0/9).Width(320).Background(c.Theme().SurfaceHover).Center().Children(func() { ui.Text(c, "16 : 9") })
						case "avatar":
							forma.Avatar(c, "Golet", nil)
						case "alert":
							forma.Alert(c, "Heads up!", "You can add components to your app using Go.")
							forma.Alert(c, "Something went wrong", "Please try again.", true)
						case "empty":
							forma.Empty(c, "No projects yet", "Create your first project to get started.", func() { forma.Button(c, "Create project") })
						case "item":
							forma.Item(c, "Team workspace", "Manage your project and team.", func() { forma.Button(c, "Open", forma.ButtonOptions{Variant: forma.Outline, Size: forma.Small}) })
						case "field", "form":
							forma.Form(c, func() {
								forma.Field(c, "Email", func() { forma.Input(c, &text, forma.InputOptions{Placeholder: "you@example.com"}) }).Description("We will never share your email.")
								forma.Checkbox(c, &on, "Subscribe")
								forma.Button(c, "Save")
							})
						case "label":
							forma.Label(c, "Email address")
						case "input-group":
							forma.InputGroup(c, func() {
								ui.Text(c, "https://")
								forma.InputGroupInput(c, &text)
								forma.Button(c, "Go", forma.ButtonOptions{Size: forma.XS})
							})
						case "input-otp":
							forma.InputOTP(c, &text, 6)
							forma.CardDescription(c, "Enter or paste six digits.")
						case "separator":
							ui.Text(c, "Golet")
							forma.Separator(c)
							ui.Text(c, "Native UI")
						case "kbd":
							ui.Row(c).Gap(4).Children(func() { forma.Kbd(c, "⌘"); forma.Kbd(c, "K") })
						case "resizable":
							forma.Resizable(c, &split, func() { ui.Text(c, "Left panel") }, func() { ui.Text(c, "Right panel") }).Height(180)
						case "scroll-area":
							forma.ScrollArea(c).Height(180).Children(func() {
								for i := 0; i < 30; i++ {
									ui.Textf(c, "Item %d", i+1).Padding(8)
									forma.Separator(c)
								}
							})
						case "skeleton":
							forma.Skeleton(c).Size(48, 48)
							forma.Skeleton(c).Size(260, 16)
							forma.Skeleton(c).Size(200, 16)
						case "spinner":
							forma.Spinner(c)
						case "carousel":
							forma.Carousel(c, &slide, 5, func(i int) {
								forma.Card(c).Height(160).Center().Children(func() { ui.Textf(c, "%d", i+1).FontSize(32) })
							})
						case "table":
							forma.Table(c, &list, []ui.TableColumn{{Title: "Name", Width: 180}, {Title: "Status", Width: 140}}, 5, func(row, col int) {
								if col == 0 {
									ui.Textf(c, "Project %d", row+1)
								} else {
									forma.Badge(c, "Active", forma.Secondary)
								}
							}).Height(220)
						case "chart":
							chartDemo(c)
						case "toast", "sonner":
							forma.Button(c, "Show notification").OnClick(func() {
								forma.Toast(c, ui.Toast{Title: "Changes saved", Description: "Your profile has been updated.", Action: "Undo", OnAction: func() { c.Toast("Undone") }})
							})
						case "direction":
							forma.Direction(c, ui.RightToLeft, func() { forma.Button(c, "مرحبا"); forma.Input(c, &text) })
						case "attachment":
							forma.Attachment(c, "report.pdf", "2.4 MB", "done", func() { c.Toast("Remove clicked") })
							forma.Attachment(c, "image.png", "Uploading…", "uploading", nil)
						case "bubble":
							forma.Bubble(c, "A native message bubble", forma.Default)
							forma.Bubble(c, "Secondary reply", forma.Secondary)
						case "message", "ai-message":
							ui.Column(c).FillWidth().Gap(28).Children(func() {
								forma.Checkbox(c, &messageStreaming, "模拟生成中（禁用操作）")
								forma.AIMessage(c, forma.MessageUser, "如何创建一个原生窗口？", forma.AIMessageOptions{OnEdit: func(value string) { text = value; messageStatus = "已将消息填入编辑框" }})
								forma.AIMessage(c, forma.MessageAssistant, "使用 **Golet** 创建原生窗口：\n\n```go\nwindow := golet.NewWindow(options)\n```\n\n界面和交互均由 Go 实现。", forma.AIMessageOptions{Streaming: messageStreaming, Feedback: &messageFeedback, OnRetry: func() { messageStatus = "已触发重试回调（演示不调用模型）" }, OnFeedback: func(v forma.MessageFeedback) { messageStatus = "反馈：" + string(v) }})
								forma.Input(c, &text, forma.InputOptions{Placeholder: "编辑消息…"})
								ui.Text(c, messageStatus).TextColor(c.Theme().TextMuted)
							})
						case "chain-of-thought", "chat", "checkpoint", "confirmation", "context", "file-tree", "inline-citation", "plan", "queue", "reasoning", "sources", "stack-trace", "subagent", "suggestion", "task", "terminal", "test-results", "tool":
							aiPatterns(c, selected)
						case "prompt-input":
							ui.Column(c).FillWidth().Gap(16).Children(func() {
								forma.PromptInput(c, &prompt, forma.PromptInputOptions{Models: models, Status: promptStatus, OnAttach: func() {
									prompt.Attachments = append(prompt.Attachments, forma.PromptAttachment{ID: fmt.Sprint(len(prompt.Attachments)), Name: "notes.txt", Description: "演示附件 · 未读取或上传文件"})
								}, OnDropFiles: func(paths []string) {
									for _, path := range paths {
										prompt.Attachments = append(prompt.Attachments, forma.PromptAttachment{ID: path, Path: path, Name: filepath.Base(path), Description: "本地文件 · 未上传"})
									}
								}, OnSubmit: func(p forma.PromptSubmission) {
									promptResult = fmt.Sprintf("提交快照：%s · %s · %d 个附件 · %s", p.Model, p.Effort, len(p.Attachments), p.Text)
									promptStatus = forma.PromptStreaming
								}, OnStop: func() { promptStatus = forma.PromptReady }})
								ui.Text(c, "Enter 发送 · Shift+Enter 换行 · 可拖入文件。添加按钮使用示例附件；不调用模型、不上传文件。").TextColor(c.Theme().TextMuted)
								ui.Text(c, promptResult)
								forma.Button(c, "完成模拟响应", forma.ButtonOptions{Variant: forma.Outline}).OnClick(func() { promptStatus = forma.PromptReady })
							})
						case "model-selector":
							ui.Column(c).FillWidth().Gap(16).Children(func() {
								forma.ModelComposer(c, &text, &modelID, &modelEffort, models, forma.ModelComposerOptions{Placeholder: "让助手帮你修改代码…", OnAdd: func() { messageStatus = "已触发添加回调（可由应用接入附件选择）" }, OnSubmit: func(prompt, model, effort string) {
									messageStatus = "已提交给演示回调：" + model + " / " + effort + " · " + prompt
								}})
								ui.Text(c, messageStatus).TextColor(c.Theme().TextMuted)
								ui.Text(c, "Model: "+modelID+" · Effort: "+modelEffort)
								ui.Text(c, "列表由应用提供；此演示不会发起模型请求。").TextColor(c.Theme().TextMuted)
							})
						case "message-scroller":
							forma.MessageScroller(c, &scroll, func() {
								for i := 0; i < 20; i++ {
									forma.Message(c, "Assistant", false, func() { forma.Bubble(c, fmt.Sprintf("Message %d", i+1), forma.Secondary) })
								}
							}).Height(220)
						case "marker":
							forma.Marker(c, "Today")
						case "trace":
							forma.Trace(c, "Read project files", &expanded, func() { ui.Text(c, "main.go"); ui.Text(c, "go.mod") })
						case "code-block":
							ui.Column(c).FillWidth().Gap(16).Children(func() {
								forma.CodeBlock(c, "Terminal", "$ pnpm lint\napps/web/components/message-list.tsx\n  18:22  error  'messages' is possibly 'undefined'. Add a guard before calling map  @typescript-eslint/no-unsafe-call\n  42:7  warning  React Hook React.useEffect has a missing dependency: 'scrollToBottom'  react-hooks/exhaustive-deps\n\n✖ 2 problems (1 error, 1 warning)", forma.CodeBlockOptions{WrapToggle: &terminalWrap, MaxHeight: -1, OnDownload: func(source string) {
									go func() {
										path, err := golet.Dialog.Save(golet.SaveDialogOptions{Title: "保存终端输出", DefaultPath: "terminal.txt"})
										if err == nil && path != "" {
											err = os.WriteFile(path, []byte(source), 0600)
										}
										if err != nil {
											log.Printf("studio: save terminal output: %v", err)
										}
									}()
								}})
								ui.Row(c).Gap(16).Children(func() {
									forma.Checkbox(c, &codeWrap, "自动换行")
									forma.Checkbox(c, &codeNumbers, "行号")
								})
								forma.CodeBlock(c, "main.go", "package main\n\nfunc main() {\n    println(\"Hello, Golet\")\n    println(\"Forma code blocks are rendered entirely in Go, with optional syntax highlighting.\")\n}\n", forma.CodeBlockOptions{LineNumbers: codeNumbers, StartLine: 98, Wrap: codeWrap, Highlight: codeHighlight, HighlightKey: "go", Expanded: &codeExpanded, CollapsedHeight: 100})
								forma.CodeBlock(c, "修改预览", "func greeting() string {\n    return \"Hello\"\n    return \"Hello, Forma\"\n}", forma.CodeBlockOptions{LineNumbers: true, Lines: map[int]forma.CodeLineVariant{1: forma.CodeLineRemoved, 2: forma.CodeLineAdded, 0: forma.CodeLineHighlighted}})
								forma.CodeBlock(c, "构建错误 · Small", "Error: configuration is missing.\nPlease set the application name before building.", forma.CodeBlockOptions{Size: forma.Small, Variant: forma.Destructive, Wrap: true})
								forma.CodeBlock(c, "", "$ go build ./...\n\nBuild completed successfully.\n", forma.CodeBlockOptions{Size: forma.Small})
							})
						case "questionnaire":
							forma.Questionnaire(c, &question, []forma.Question{{Title: "Choose a language", Description: "Select one option.", Choices: []string{"Go", "TypeScript", "Rust"}, Answer: &answer, Required: true}}, func() { c.Toast("Submitted: " + answer) })
						}
					})
				})
			})
		})
		forma.Sonner(c)
	}
}
