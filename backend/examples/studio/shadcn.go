// Author: HyperparameterSI
package main

import (
	"fmt"
	"github.com/HyperparameterSI/golet/backend"
	"github.com/HyperparameterSI/golet/backend/examples/internal/demotheme"
	"github.com/HyperparameterSI/golet/backend/ui"
	"github.com/HyperparameterSI/golet/backend/ui/chart"
	s "github.com/HyperparameterSI/golet/backend/ui/shadcn"
	"time"
)

func openShadcn() *golet.Window {
	return golet.NewWindow(golet.WindowOptions{Title: "Studio · shadcn / Maia", Width: 1040, Height: 840, Content: ui.View(newShadcnView())})
}
func newShadcnView(initial ...string) func(*ui.Context) {
	selected := "button"
	if len(initial) > 0 {
		selected = initial[0]
	}
	dark, follow := false, false
	lightTheme, darkTheme := s.Theme(false), s.Theme(true)
	text, choice, answer := "", "Go", ""
	on, expanded, opened := false, false, false
	value := float64(40)
	tab, page, slide, question := 0, 1, 0, 0
	date := time.Date(2026, 10, 24, 0, 0, 0, 0, time.UTC)
	calendarMode := "Single"
	var dates []time.Time
	dateRange := s.DateRange{}
	split := float32(180)
	list := ui.ListState{}
	scroll := ui.ScrollState{}
	return func(c *ui.Context) {
		demotheme.Apply(c)
		if !follow {
			theme := lightTheme
			if dark {
				theme = darkTheme
			}
			c.SetTheme(theme)
		}
		ui.Column(c).Fill().Background(c.Theme().Background).Padding(32).Gap(24).Children(func() {
			ui.Row(c).Gap(16).AlignItems(ui.Center).Children(func() {
				ui.Text(c, "shadcn / Maia").FontSize(26).Bold().Grow(1)
				ui.Checkbox(c, &dark, "Dark")
				ui.Checkbox(c, &follow, "跟随 Studio 主题")
			})
			ui.Text(c, "65 个官方目录项 · 纯 Go 原生组件 · 基准源码已固定").TextColor(c.Theme().TextMuted)
			s.Select(c, &selected, s.Components()).Width(280).Label("Component")
			ui.Scroll(c).Grow(1).Children(func() {
				s.Card(c).Children(func() {
					s.CardHeader(c).Children(func() {
						s.CardTitle(c, selected)
						s.CardDescription(c, "交互演示；可切换明暗外观。完整对齐状态见组件移植清单。")
					})
					s.CardContent(c).MinHeight(230).Children(func() {
						menu := []s.MenuItem{{Label: "Profile", Action: func() { c.Toast("Profile selected") }}, {Label: "Notifications", Checked: &on}, {Separator: true}, {Label: "Disabled action", Disabled: true}}
						switch selected {
						case "button":
							ui.Row(c).Gap(8).Children(func() {
								for i, name := range []string{"Default", "Outline", "Secondary", "Ghost", "Destructive", "Link"} {
									s.Button(c, name, s.ButtonOptions{Variant: s.Variant(i)}).OnClick(func() { c.Toast(name) })
								}
							})
							ui.Row(c).Gap(8).Children(func() {
								for _, size := range []s.Size{s.XS, s.Small, s.Normal, s.Large} {
									s.Button(c, "Button", s.ButtonOptions{Size: size})
								}
								s.Button(c, "Disabled", s.ButtonOptions{Disabled: true})
							})
						case "badge":
							ui.Row(c).Gap(8).Children(func() {
								for i, name := range []string{"Default", "Outline", "Secondary", "Ghost", "Destructive", "Link"} {
									s.Badge(c, name, s.Variant(i))
								}
							})
						case "card":
							s.Card(c).Children(func() {
								s.CardHeader(c).Children(func() { s.CardTitle(c, "Create project"); s.CardDescription(c, "Start building with Golet.") })
								s.CardContent(c).Children(func() { s.Input(c, &text, s.InputOptions{Placeholder: "Project name"}) })
								s.CardFooter(c).Children(func() { s.Button(c, "Create") })
							})
						case "input":
							s.Input(c, &text, s.InputOptions{Placeholder: "you@example.com"})
							s.Input(c, &text, s.InputOptions{Invalid: true})
							s.Input(c, &text, s.InputOptions{Disabled: true})
						case "textarea":
							s.Textarea(c, &text, s.InputOptions{Placeholder: "Write a message…"}).Height(120)
						case "checkbox":
							s.Checkbox(c, &on, "Accept terms and conditions")
						case "switch":
							s.Switch(c, &on).Label("Airplane mode")
						case "radio-group":
							s.RadioGroup(c, func() {
								for _, option := range []string{"Go", "TypeScript", "Rust"} {
									s.Radio(c, &choice, option, option)
								}
							})
						case "select":
							s.Select(c, &choice, []string{"Go", "TypeScript", "Rust"})
						case "native-select":
							s.NativeSelect(c, &choice, []string{"Go", "TypeScript", "Rust"})
						case "combobox":
							s.Combobox(c, &choice, []string{"Go", "TypeScript", "Rust"})
						case "tabs":
							s.Tabs(c, &tab, "Account", "Password", "Notifications")
							s.CardDescription(c, fmt.Sprintf("Selected panel: %d", tab+1))
						case "toggle":
							s.Toggle(c, &on, "Bold")
						case "toggle-group":
							s.ToggleGroup(c, func() { s.Toggle(c, &on, "Bold"); s.Toggle(c, &expanded, "Italic") })
						case "slider":
							s.Slider(c, &value, 0, 100).Width(320)
							ui.Textf(c, "%.0f", value)
						case "progress":
							s.Progress(c, value/100).Width(320)
							s.Button(c, "Advance").OnClick(func() { value = float64((int(value) + 10) % 110) })
						case "calendar":
							s.Select(c, &calendarMode, []string{"Single", "Multiple", "Range"}).Width(200)
							options := s.CalendarOptions{}
							if calendarMode == "Multiple" {
								options.Multiple = &dates
							}
							if calendarMode == "Range" {
								options.Range = &dateRange
							}
							s.Calendar(c, &date, options)
						case "accordion":
							s.Accordion(c, func() {
								s.AccordionItem(c, "Is it native?", &expanded, func() { ui.Text(c, "Yes. Every element is rendered by Golet, without a webview.") })
								s.AccordionItem(c, "Does it support dark mode?", &on, func() { ui.Text(c, "Use the Dark switch above.") })
							})
						case "collapsible":
							s.Collapsible(c, "Dependencies", &expanded, func() { ui.Text(c, "golet/backend/ui/shadcn") })
						case "dialog", "alert-dialog", "sheet", "drawer":
							s.Button(c, "Open "+selected).OnClick(func() { opened = true })
							switch selected {
							case "dialog":
								s.Dialog(c, &opened, "Edit profile", func() {
									s.Input(c, &text, s.InputOptions{Placeholder: "Name"})
									s.Button(c, "Save").OnClick(func() { opened = false })
								})
							case "alert-dialog":
								s.AlertDialog(c, &opened, "Are you sure?", "This is a demonstration; no data is removed.", func() { c.Toast("Confirmed") })
							case "sheet":
								s.Sheet(c, &opened, s.Right, "Settings", func() { s.Input(c, &text) })
							case "drawer":
								s.Drawer(c, &opened, "Move goal", func() { s.Slider(c, &value, 0, 100); s.Button(c, "Done").OnClick(func() { opened = false }) })
							}
						case "popover":
							b := s.Button(c, "Open popover")
							b.OnClick(func() { opened = !opened })
							s.Popover(c, b, &opened, func() { s.Label(c, "Dimensions"); s.Input(c, &text) })
						case "tooltip":
							b := s.Button(c, "Hover or focus")
							s.Tooltip(c, b, "Pure Go tooltip")
						case "hover-card":
							b := s.Button(c, "@golet", s.ButtonOptions{Variant: s.Link})
							s.HoverCard(c, b, func() { s.CardTitle(c, "Golet"); s.CardDescription(c, "Build native apps with Go.") })
						case "dropdown-menu":
							s.DropdownMenu(c, "Open menu", menu)
						case "context-menu":
							s.ContextMenu(c, ui.Box(c).Height(160).Border(1, c.Theme().Border).Radius(16).Center().Children(func() { ui.Text(c, "Right click or Shift+F10") }), menu)
						case "menubar":
							s.Menubar(c, func() { s.DropdownMenu(c, "File", menu); s.DropdownMenu(c, "Edit", menu) })
						case "navigation-menu":
							s.NavigationMenu(c, func() {
								s.DropdownMenu(c, "Products", menu)
								s.Button(c, "Documentation", s.ButtonOptions{Variant: s.Ghost})
							})
						case "command":
							s.Command(c, &text, []s.CommandItem{{ID: "calendar", Label: "Calendar", Group: "Suggestions", Run: func() { c.Toast("Calendar") }}, {ID: "settings", Label: "Settings", Group: "Preferences", Run: func() { c.Toast("Settings") }}})
						case "pagination":
							s.Pagination(c, &page, 12)
							ui.Textf(c, "Page %d", page)
						case "sidebar":
							s.Sidebar(c, &choice, func() {
								s.SidebarItem(c, "Go", "Overview")
								s.SidebarItem(c, "TypeScript", "Projects")
								s.SidebarItem(c, "Rust", "Settings")
							}).Height(200)
						case "breadcrumb":
							s.Breadcrumb(c, []string{"Home", "Components", "Maia"}, &tab)
						case "button-group":
							s.ButtonGroup(c, func() {
								s.GroupItem(s.Button(c, "Previous", s.ButtonOptions{Variant: s.Outline}), 0, 2)
								s.GroupItem(s.Button(c, "Next", s.ButtonOptions{Variant: s.Outline}), 1, 2)
							})
						case "aspect-ratio":
							s.AspectRatio(c, 16.0/9).Width(320).Background(c.Theme().SurfaceHover).Center().Children(func() { ui.Text(c, "16 : 9") })
						case "avatar":
							s.Avatar(c, "Golet", nil)
						case "alert":
							s.Alert(c, "Heads up!", "You can add components to your app using Go.")
							s.Alert(c, "Something went wrong", "Please try again.", true)
						case "empty":
							s.Empty(c, "No projects yet", "Create your first project to get started.", func() { s.Button(c, "Create project") })
						case "item":
							s.Item(c, "Team workspace", "Manage your project and team.", func() { s.Button(c, "Open", s.ButtonOptions{Variant: s.Outline, Size: s.Small}) })
						case "field", "form":
							s.Form(c, func() {
								s.Field(c, "Email", func() { s.Input(c, &text, s.InputOptions{Placeholder: "you@example.com"}) }).Description("We will never share your email.")
								s.Checkbox(c, &on, "Subscribe")
								s.Button(c, "Save")
							})
						case "label":
							s.Label(c, "Email address")
						case "input-group":
							s.InputGroup(c, func() {
								ui.Text(c, "https://")
								s.InputGroupInput(c, &text)
								s.Button(c, "Go", s.ButtonOptions{Size: s.XS})
							})
						case "input-otp":
							s.InputOTP(c, &text, 6)
							s.CardDescription(c, "Enter or paste six digits.")
						case "separator":
							ui.Text(c, "Golet")
							s.Separator(c)
							ui.Text(c, "Native UI")
						case "kbd":
							ui.Row(c).Gap(4).Children(func() { s.Kbd(c, "⌘"); s.Kbd(c, "K") })
						case "resizable":
							s.Resizable(c, &split, func() { ui.Text(c, "Left panel") }, func() { ui.Text(c, "Right panel") }).Height(180)
						case "scroll-area":
							s.ScrollArea(c).Height(180).Children(func() {
								for i := 0; i < 30; i++ {
									ui.Textf(c, "Item %d", i+1).Padding(8)
									s.Separator(c)
								}
							})
						case "skeleton":
							s.Skeleton(c).Size(48, 48)
							s.Skeleton(c).Size(260, 16)
							s.Skeleton(c).Size(200, 16)
						case "spinner":
							s.Spinner(c)
						case "carousel":
							s.Carousel(c, &slide, 5, func(i int) { s.Card(c).Height(160).Center().Children(func() { ui.Textf(c, "%d", i+1).FontSize(32) }) })
						case "table":
							s.Table(c, &list, []ui.TableColumn{{Title: "Name", Width: 180}, {Title: "Status", Width: 140}}, 5, func(row, col int) {
								if col == 0 {
									ui.Textf(c, "Project %d", row+1)
								} else {
									s.Badge(c, "Active", s.Secondary)
								}
							}).Height(220)
						case "chart":
							s.ChartContainer(c, chart.Axes{}, func(area ui.Element) {
								area.Draw(func(p *ui.Painter, r ui.Rect) {
									xs := []float32{0, .2, .4, .6, .8, 1}
									ys := []float32{.2, .5, .3, .8, .6, .9}
									p.StrokePath(chart.Line(r, xs, ys, chart.CurveMonotone), 2, c.Theme().Accent)
								})
							}).Height(180)
							s.ChartLegend(c, []string{"Activity"}, []ui.Color{c.Theme().Accent})
						case "toast", "sonner":
							s.Button(c, "Show notification").OnClick(func() {
								s.Toast(c, ui.Toast{Title: "Changes saved", Description: "Your profile has been updated.", Action: "Undo", OnAction: func() { c.Toast("Undone") }})
							})
						case "direction":
							s.Direction(c, ui.RightToLeft, func() { s.Button(c, "مرحبا"); s.Input(c, &text) })
						case "attachment":
							s.Attachment(c, "report.pdf", "2.4 MB", "done", func() { c.Toast("Remove clicked") })
							s.Attachment(c, "image.png", "Uploading…", "uploading", nil)
						case "bubble":
							s.Bubble(c, "A native message bubble", s.Default)
							s.Bubble(c, "Secondary reply", s.Secondary)
						case "message":
							s.Message(c, "You", true, func() { s.Bubble(c, "Build something native.", s.Default) })
							s.Message(c, "Assistant", false, func() { s.Bubble(c, "Let's build it with Golet.", s.Secondary) })
						case "message-scroller":
							s.MessageScroller(c, &scroll, func() {
								for i := 0; i < 20; i++ {
									s.Message(c, "Assistant", false, func() { s.Bubble(c, fmt.Sprintf("Message %d", i+1), s.Secondary) })
								}
							}).Height(220)
						case "marker":
							s.Marker(c, "Today")
						case "trace":
							s.Trace(c, "Read project files", &expanded, func() { ui.Text(c, "main.go"); ui.Text(c, "go.mod") })
						case "code-block":
							s.CodeBlock(c, "Go", "func main() {\n    println(\"Hello, Golet\")\n}")
						case "questionnaire":
							s.Questionnaire(c, &question, []s.Question{{Title: "Choose a language", Description: "Select one option.", Choices: []string{"Go", "TypeScript", "Rust"}, Answer: &answer, Required: true}}, func() { c.Toast("Submitted: " + answer) })
						}
					})
				})
			})
		})
		s.Sonner(c)
	}
}
