// Author: HyperparameterSI
package forma

import (
	"time"

	"github.com/HyperparameterSI/golet/backend/ui"
)

func Popover(c *ui.Context, anchor ui.Element, open *bool, build func()) ui.Element {
	return ui.PopoverBase(c, anchor, open, func(p ui.Element) { panel(c, p).Width(288).Padding(16).Gap(16); build() })
}
func Tooltip(c *ui.Context, anchor ui.Element, text string) ui.Element {
	return ui.TooltipBase(c, anchor, func(tip ui.Element) {
		tip.Padding(6, 12).Radius(radius(c, 12)).FontSize(fontSize(c, 12)).Background(c.Theme().Text).TextColor(c.Theme().Background).Children(func() { ui.Text(c, text) })
	})
}
func Dialog(c *ui.Context, open *bool, title string, build func()) ui.Element {
	return ui.DialogBase(c, open, func(back, p ui.Element) {
		back.Background(ui.RGBA(0, 0, 0, .8))
		w, _ := c.Size()
		panel(c, p).Width(min(448, max(0, w-32))).Radius(radius(c, 32)).Padding(24).Gap(24).Label(title)
		ui.Row(c).AlignItems(ui.Center).Children(func() {
			CardTitle(c, title).Grow(1)
			Button(c, "×", ButtonOptions{Variant: Ghost, Size: IconSmall}).Label("Close dialog").OnClick(func() { *open = false })
		})
		build()
	})
}
func AlertDialog(c *ui.Context, open *bool, title, description string, confirm func()) ui.Element {
	return ui.DialogBase(c, open, func(back, p ui.Element) {
		// An alert must not be dismissed by an accidental backdrop click.
		if back.Clicked() {
			*open = true
		}
		back.Background(ui.RGBA(0, 0, 0, .8))
		w, _ := c.Size()
		panel(c, p).Width(min(448, max(0, w-32))).Radius(radius(c, 32)).Padding(24).Gap(24).Label(title)
		CardTitle(c, title)
		CardDescription(c, description)
		ui.Row(c).Gap(8).Justify(ui.End).Children(func() {
			Button(c, "Cancel", ButtonOptions{Variant: Outline}).AutoFocus().OnClick(func() { *open = false })
			Button(c, "Continue").OnClick(func() {
				*open = false
				if confirm != nil {
					confirm()
				}
			})
		})
	})
}

// Side locates a sheet at the corresponding window edge.
type Side uint8

const (
	Right Side = iota
	Left
	Top
	Bottom
)

func Sheet(c *ui.Context, open *bool, side Side, title string, build func()) ui.Element {
	return ui.DialogBase(c, open, func(back, p ui.Element) {
		back.Background(ui.RGBA(0, 0, 0, .8))
		w, h := c.Size()
		p.Absolute().Background(token(c, "popover", c.Theme().Surface)).TextColor(token(c, "popover-foreground", c.Theme().Text)).Padding(24).Gap(24).Label(title)
		switch side {
		case Left:
			p.Left(0).Top(0).Bottom(0).Width(min(384, w*.75))
		case Top:
			p.Left(0).Right(0).Top(0).MaxHeight(h * .8)
		case Bottom:
			p.Left(0).Right(0).Bottom(0).MaxHeight(h * .8)
		default:
			p.Right(0).Top(0).Bottom(0).Width(min(384, w*.75))
		}
		CardTitle(c, title)
		Button(c, "×", ButtonOptions{Variant: Ghost, Size: IconSmall}).Absolute().Top(16).Right(16).Label("Close sheet").OnClick(func() { *open = false })
		build()
	})
}

// Drawer provides a bottom panel with native drag-to-dismiss and keyboard close.
func Drawer(c *ui.Context, open *bool, title string, build func()) ui.Element {
	return ui.DialogBase(c, open, func(back, p ui.Element) {
		back.Background(ui.RGBA(0, 0, 0, .8))
		_, h := c.Size()
		p.Absolute().Left(8).Right(8).Bottom(8).MaxHeight(h - 16).Radius(radius(c, 32)).Padding(24).Gap(16).Background(token(c, "popover", c.Theme().Surface)).TextColor(token(c, "popover-foreground", c.Theme().Text)).Label(title)
		offset := ui.Local(p, "drawer-offset", func() ui.Point { return ui.Point{} })
		p.Draggable(offset, ui.DraggableOptions{Axis: ui.Vertical, Bounds: &ui.Rect{H: h}, Snap: func(v ui.Point) ui.Point {
			if v.Y > h*.25 {
				*open = false
			}
			return ui.Point{}
		}})
		ui.Box(c).Size(100, 6).Radius(radius(c, 32)).Background(c.Theme().SurfaceHover).AlignSelf(ui.Center)
		CardTitle(c, title)
		build()
	})
}

// HoverCard keeps its content open while either the anchor or panel is hovered.
func HoverCard(c *ui.Context, anchor ui.Element, build func()) ui.Element {
	state := ui.Local(anchor, "maia-hover-card", func() struct {
		open bool
		left time.Time
	} {
		return struct {
			open bool
			left time.Time
		}{}
	})
	if anchor.Hovered() || anchor.FocusWithin() {
		state.open = true
		state.left = time.Time{}
	}
	p := Popover(c, anchor, &state.open, build)
	if state.open {
		if p.Hovered() {
			state.left = time.Time{}
		} else if !anchor.Hovered() && !anchor.FocusWithin() {
			if state.left.IsZero() {
				state.left = c.Now()
			}
			if c.Now().Sub(state.left) > 200*time.Millisecond {
				state.open = false
			} else {
				c.After(200 * time.Millisecond)
			}
		}
	}
	return p
}

// MenuItem describes a native-drawn popup row, checkbox, or submenu.
type MenuItem struct {
	Label, Shortcut     string
	Disabled, Separator bool
	Checked             *bool
	Children            []MenuItem
	Action              func()
}

func menuItems(c *ui.Context, open *bool, items []MenuItem) {
	first := true
	for _, item := range items {
		if item.Separator {
			Separator(c).MarginY(4)
			continue
		}
		e := Button(c, item.Label, ButtonOptions{Variant: Ghost, Disabled: item.Disabled}).Radius(radius(c, 12)).Justify(ui.Start).Padding(8, 12)
		if !e.Hovered() || item.Disabled {
			e.TextColor(token(c, "popover-foreground", c.Theme().Text))
		}
		if first && !item.Disabled {
			e.AutoFocus()
			first = false
		}
		if item.Checked != nil {
			e.Children(func() {
				if *item.Checked {
					ui.Text(c, "✓")
				}
			})
		}
		if item.Shortcut != "" {
			e.Children(func() { Kbd(c, item.Shortcut) })
		}
		if len(item.Children) > 0 {
			sub := ui.Local(e, "submenu", func() bool { return false })
			e.OnClick(func() { *sub = !*sub })
			Popover(c, e, sub, func() { menuItems(c, open, item.Children) })
		} else {
			e.OnClick(func() {
				if item.Checked != nil {
					*item.Checked = !*item.Checked
				}
				*open = false
				if item.Action != nil {
					item.Action()
				}
			})
		}
	}
}
func DropdownMenu(c *ui.Context, label string, items []MenuItem) ui.Element {
	e := Button(c, label, ButtonOptions{Variant: Outline})
	open := ui.Local(e, "menu-open", func() bool { return false })
	e.OnClick(func() { *open = !*open })
	ui.PopoverBase(c, e, open, func(p ui.Element) { panel(c, p).MinWidth(192).FocusGroup(ui.Vertical); menuItems(c, open, items) })
	return e
}

// ContextMenu draws its popup in Go, keeping Maia styling on every OS.
// Register it on an anchor that does not already own a HandleInput callback.
func ContextMenu(c *ui.Context, anchor ui.Element, items []MenuItem) ui.Element {
	state := ui.Local(anchor, "maia-context-menu", func() struct {
		open  bool
		point ui.Point
	} {
		return struct {
			open  bool
			point ui.Point
		}{}
	})
	anchor.Focusable().HandleInput(func(ev ui.InputEvent) bool {
		if ev.Kind == ui.InputPointerDown && ev.Button == 1 {
			state.open = true
			state.point = ui.Point{X: ev.X, Y: ev.Y}
			return true
		}
		if ev.Kind == ui.InputKeyDown && ev.Key == ui.KeyF10 && ev.Mods == ui.Shift {
			state.open = true
			state.point = ui.Point{}
			return true
		}
		return false
	})
	ui.PopoverBase(c, anchor, &state.open, func(p ui.Element) {
		panel(c, p).MinWidth(192).AttachTo(anchor, ui.AnchorTopStart, ui.AnchorTopStart).Margin(state.point.Y, 0, 0, state.point.X).FocusGroup(ui.Vertical)
		menuItems(c, &state.open, items)
	})
	return anchor
}
