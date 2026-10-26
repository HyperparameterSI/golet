// Author: HyperparameterSI
package forma

import (
	"fmt"
	"strings"

	"github.com/HyperparameterSI/golet/backend/ui"
)

func panel(c *ui.Context, e ui.Element) ui.Element {
	return e.Radius(radius(c, 16)).Padding(4).Gap(2).Background(token(c, "popover", c.Theme().Surface)).TextColor(token(c, "popover-foreground", c.Theme().Text)).FontSize(fontSize(c, 14)).Border(1, c.Theme().Border).Shadow(0, 12, 24, 0, ui.RGBA(0, 0, 0, .18))
}
func Select(c *ui.Context, value *string, options []string) ui.Element {
	s := ui.SelectBase(c, value)
	inputStyle(c, s.Trigger.Height(36).Padding(8, 12).Gap(6), radius(c, 32), InputOptions{})
	s.Trigger.Children(func() { ui.Text(c, *value).Grow(1); ui.Text(c, "⌄").TextColor(c.Theme().TextMuted) })
	s.Popup(func(p ui.Element) {
		panel(c, p).MinWidth(144).MaxHeight(320)
		ui.Scroll(c).Children(func() {
			for _, option := range options {
				e := s.Item(option).Padding(8, 12).Radius(radius(c, 12)).Gap(10)
				if e.Highlighted() {
					e.Background(token(c, "accent", c.Theme().SurfaceHover)).TextColor(token(c, "accent-foreground", c.Theme().Text))
				}
				e.Children(func() {
					ui.Text(c, option).Grow(1)
					if option == *value {
						ui.Text(c, "✓")
					}
				})
			}
		})
	})
	return s.Trigger
}

// NativeSelect uses the same Go selection semantics with Maia trigger geometry.
func NativeSelect(c *ui.Context, value *string, options []string) ui.Element {
	return Select(c, value, options)
}
func Combobox(c *ui.Context, value *string, options []string) ui.Element {
	s := ui.ComboboxBase(c, value)
	inputStyle(c, s.Input.Height(36).Padding(4, 12), radius(c, 32), InputOptions{})
	s.Popup(func(p ui.Element) {
		panel(c, p).MaxHeight(320)
		ui.Scroll(c).Children(func() {
			count := 0
			for _, option := range options {
				if s.Filtering() && !strings.Contains(strings.ToLower(option), strings.ToLower(*value)) {
					continue
				}
				count++
				e := s.Item(option).Padding(8, 12).Radius(radius(c, 12))
				if e.Highlighted() {
					e.Background(token(c, "accent", c.Theme().SurfaceHover)).TextColor(token(c, "accent-foreground", c.Theme().Text))
				}
				e.Children(func() { ui.Text(c, option) })
			}
			if count == 0 {
				CardDescription(c, "No results found.").Padding(12)
			}
		})
	})
	return s.Input
}
func Tabs(c *ui.Context, selected *int, labels ...string) ui.Element {
	parts := ui.TabsBase(c, selected, len(labels))
	parts.List.Height(36).Padding(3).Radius(radius(c, 32)).Background(c.Theme().SurfaceHover).Children(func() {
		for i, label := range labels {
			e := parts.Tab(i).Height(30).Padding(4, 8).Radius(radius(c, 12)).FontSize(fontSize(c, 14)).FontWeight(500)
			if *selected == i {
				e.Background(c.Theme().Background).TextColor(c.Theme().Text)
			} else {
				e.TextColor(c.Theme().TextMuted)
			}
			focus(c, e, radius(c, 12), false)
			e.Children(func() { ui.Text(c, label) })
		}
	})
	return parts.List
}
func Collapsible(c *ui.Context, title string, open *bool, build func()) ui.Element {
	parts := ui.CollapsibleBase(c, open)
	parts.Trigger.Padding(16).Gap(24).FontSize(fontSize(c, 14)).FontWeight(500).Children(func() {
		ui.Text(c, title).Grow(1)
		symbol := "⌄"
		if *open {
			symbol = "⌃"
		}
		ui.Text(c, symbol)
	})
	parts.Panel(func() { ui.Column(c).Padding(0, 16, 16).Children(build) })
	return parts.Trigger
}
func Accordion(c *ui.Context, build func()) ui.Element {
	return ui.Column(c).Radius(radius(c, 16)).Border(1, c.Theme().Border).Clip().Children(build)
}
func AccordionItem(c *ui.Context, title string, open *bool, build func()) ui.Element {
	e := ui.Column(c)
	if *open {
		e.Background(c.Theme().SurfaceHover.Alpha(.5))
	}
	e.Children(func() { Collapsible(c, title, open, build); Separator(c) })
	return e
}
func Pagination(c *ui.Context, page *int, total int) ui.Element {
	total = max(1, total)
	*page = max(1, min(*page, total))
	return ui.Row(c).Gap(4).AlignItems(ui.Center).Children(func() {
		Button(c, "‹ Previous", ButtonOptions{Variant: Ghost, Disabled: *page == 1}).OnClick(func() { *page-- })
		start, end := max(1, *page-2), min(total, *page+2)
		if start > 1 {
			Button(c, "1", ButtonOptions{Variant: Ghost}).OnClick(func() { *page = 1 })
			if start > 2 {
				ui.Text(c, "…")
			}
		}
		for i := start; i <= end; i++ {
			v := Ghost
			if i == *page {
				v = Outline
			}
			Button(c, fmt.Sprint(i), ButtonOptions{Variant: v, Size: Icon}).OnClick(func() { *page = i })
		}
		if end < total {
			if end < total-1 {
				ui.Text(c, "…")
			}
			Button(c, fmt.Sprint(total), ButtonOptions{Variant: Ghost}).OnClick(func() { *page = total })
		}
		Button(c, "Next ›", ButtonOptions{Variant: Ghost, Disabled: *page == total}).OnClick(func() { *page++ })
	})
}
func Sidebar(c *ui.Context, selected *string, build func()) ui.Element {
	outer := ui.Box(c).Width(256)
	original := c.Theme()
	local := ui.Local(outer, "sidebar-theme", func() ui.Theme { return *original })
	*local = *original
	local.Surface = token(c, "sidebar", original.Sidebar)
	local.Sidebar = local.Surface
	local.Text = token(c, "sidebar-foreground", original.Text)
	local.Accent = token(c, "sidebar-primary", original.Accent)
	local.AccentText = token(c, "sidebar-primary-foreground", original.AccentText)
	local.SurfaceHover = token(c, "sidebar-accent", original.SurfaceHover)
	local.Border = token(c, "sidebar-border", original.Border)
	local.Focus = token(c, "sidebar-ring", original.Focus)
	return outer.TextColor(local.Text).Children(func() {
		c.SetTheme(local)
		defer c.SetTheme(original)
		ui.Sidebar(c, selected, build).Fill().Padding(8).Gap(4).Background(local.Sidebar).TextColor(local.Text).FontSize(fontSize(c, 14))
	})
}
func SidebarItem(c *ui.Context, id, label string) ui.Element {
	return ui.SidebarItem(c, id, nil, label).Height(32).Radius(radius(c, 12)).Padding(8)
}
func NavigationMenu(c *ui.Context, build func()) ui.Element {
	return ui.Row(c).Gap(4).FocusGroup(ui.Horizontal).Children(build)
}
func Menubar(c *ui.Context, build func()) ui.Element {
	return ui.Row(c).Gap(4).Padding(4).Radius(radius(c, 16)).Border(1, c.Theme().Border).FocusGroup(ui.Horizontal).Children(build)
}

// CommandItem is a command-palette action. ID must be stable and unique.
type CommandItem struct {
	ID, Label, Group, Shortcut string
	Disabled                   bool
	Run                        func()
}

// Command provides keyboard-searchable actions. Selection runs only once.
func Command(c *ui.Context, query *string, items []CommandItem) ui.Element {
	var chosen func()
	e := ui.Column(c).Radius(radius(c, 16)).Background(c.Theme().Surface).Border(1, c.Theme().Border).FocusGroup(ui.Vertical)
	e.Children(func() {
		Input(c, query, InputOptions{Placeholder: "Type a command or search…"}).Label("Command search")
		Separator(c)
		ui.Scroll(c).MaxHeight(300).Children(func() {
			group := ""
			count := 0
			for _, item := range items {
				if !strings.Contains(strings.ToLower(item.Label), strings.ToLower(*query)) {
					continue
				}
				count++
				if item.Group != group {
					group = item.Group
					CardDescription(c, group).FontSize(fontSize(c, 12)).Padding(8)
				}
				button := Button(c, item.Label, ButtonOptions{Variant: Ghost, Disabled: item.Disabled}).Radius(radius(c, 12)).Justify(ui.Start).Label(item.ID)
				if button.Clicked() {
					chosen = item.Run
				}
			}
			if count == 0 {
				CardDescription(c, "No results found.").Padding(16)
			}
		})
	})
	if chosen != nil {
		chosen()
	}
	return e
}
