// Author: HyperparameterSI
package shadcn

import "github.com/HyperparameterSI/golet/backend/ui"

// Card is a Maia card. Compose its Header, Content and Footer inside Children.
func Card(c *ui.Context) ui.Element {
	return ui.Column(c).Gap(24).PaddingY(24).Radius(16).Background(c.Theme().Surface).Border(1, c.Theme().Text.Alpha(.1)).FontSize(14)
}
func CardHeader(c *ui.Context) ui.Element { return ui.Column(c).Gap(8).PaddingX(24) }
func CardTitle(c *ui.Context, title string) ui.Element {
	return ui.Text(c, title).FontSize(16).FontWeight(500)
}
func CardDescription(c *ui.Context, text string) ui.Element {
	return ui.Text(c, text).FontSize(14).TextColor(c.Theme().TextMuted)
}
func CardContent(c *ui.Context) ui.Element { return ui.Column(c).PaddingX(24).Gap(16) }
func CardFooter(c *ui.Context) ui.Element  { return ui.Row(c).PaddingX(24).Gap(8) }
func Label(c *ui.Context, text string) ui.Element {
	return ui.Text(c, text).FontSize(14).FontWeight(500)
}

// Separator separates adjacent groups in either orientation.
func Separator(c *ui.Context, vertical ...bool) ui.Element {
	e := ui.Box(c).Background(c.Theme().Border).Shrink(0)
	if len(vertical) > 0 && vertical[0] {
		return e.Width(1)
	}
	return e.Height(1)
}
func AspectRatio(c *ui.Context, ratio float32) ui.Element {
	if ratio <= 0 {
		ratio = 1
	}
	return ui.Box(c).AspectRatio(ratio).Clip()
}
func ScrollArea(c *ui.Context) ui.Element { return ui.Scroll(c) }
func Resizable(c *ui.Context, size *float32, first, second func()) ui.Element {
	return ui.Split(c, size, first, second)
}
func Direction(c *ui.Context, d ui.Direction, build func()) ui.Element {
	return ui.Box(c).Direction(d).Children(build)
}
func ButtonGroup(c *ui.Context, build func()) ui.Element {
	return ui.Row(c).Gap(0).Role(ui.RoleGroup).Children(build)
}

// GroupItem joins an element to its neighbors without overlapping borders.
// Apply this after styling each child of ButtonGroup, using zero-based indexes.
// vertical selects a column group; call Column on the container as well.
func GroupItem(e ui.Element, index, count int, vertical ...bool) ui.Element {
	if count <= 1 {
		return e
	}
	first, last := float32(0), float32(0)
	if index == 0 {
		first = 32
	}
	if index == count-1 {
		last = 32
	}
	if len(vertical) > 0 && vertical[0] {
		e.Radius(first, first, last, last)
		if index > 0 {
			e.BorderWidth(0, 1, 1, 1)
		}
	} else {
		e.Radius(first, last, last, first)
		if index > 0 {
			e.BorderWidth(1, 1, 1, 0)
		}
	}
	return e
}
func Kbd(c *ui.Context, text string) ui.Element {
	return ui.Text(c, text).FontSize(12).FontWeight(500).Padding(0, 4).Height(20).Radius(4).Background(c.Theme().SurfaceHover).TextColor(c.Theme().TextMuted)
}
func Alert(c *ui.Context, title, description string, destructive ...bool) ui.Element {
	fg := c.Theme().Text
	if len(destructive) > 0 && destructive[0] {
		fg = c.Theme().Danger
	}
	return ui.Column(c).Role(ui.RoleStatus).Padding(16).Gap(4).Radius(16).Border(1, c.Theme().Border).TextColor(fg).Children(func() { Label(c, title); ui.Text(c, description).FontSize(14) })
}
func Empty(c *ui.Context, title, description string, actions func()) ui.Element {
	return ui.Column(c).Center().Padding(48).Gap(16).Radius(8).Children(func() {
		ui.Text(c, title).FontSize(18).FontWeight(500)
		CardDescription(c, description)
		if actions != nil {
			actions()
		}
	})
}
func Item(c *ui.Context, title, description string, actions func()) ui.Element {
	return ui.Row(c).Padding(16).Gap(16).Radius(16).Border(1, c.Theme().Border).AlignItems(ui.Center).Children(func() {
		ui.Column(c).Grow(1).Gap(4).Children(func() { Label(c, title); CardDescription(c, description) })
		if actions != nil {
			actions()
		}
	})
}

// Avatar displays an image or initials on Maia's neutral fallback surface.
func Avatar(c *ui.Context, name string, bitmap *ui.Bitmap) ui.Element {
	return ui.Avatar(c, name, bitmap).Size(32, 32).Radius(32).Background(c.Theme().SurfaceHover).TextColor(c.Theme().TextMuted)
}

// Breadcrumb reports an activated ancestor in selected; the last item is the
// current page and is deliberately not interactive.
func Breadcrumb(c *ui.Context, items []string, selected *int) ui.Element {
	return ui.Row(c).AlignItems(ui.Center).Gap(6).Role(ui.RoleGroup).Label("Breadcrumb").Children(func() {
		for i, item := range items {
			if i > 0 {
				separator := ui.Box(c).Size(14, 14).Shrink(0)
				path := ui.Local(separator, "chevron", func() ui.Path { return ui.Path{} })
				separator.Draw(func(p *ui.Painter, r ui.Rect) {
					path.Reset().MoveTo(r.X+5, r.Y+3).LineTo(r.X+9, r.Y+7).LineTo(r.X+5, r.Y+11)
					p.StrokePath(path, 1.5, c.Theme().TextMuted)
				})
			}
			if i == len(items)-1 {
				ui.Text(c, item).FontSize(14)
				continue
			}
			link := ui.ButtonBase(c).Label(item).Role(ui.RoleLink).Padding(0).FontSize(14).TextColor(c.Theme().TextMuted)
			if link.Hovered() {
				link.TextColor(c.Theme().Text)
			}
			link.Children(func() { ui.Text(c, item).SingleLine() })
			link.OnClick(func() {
				if selected != nil {
					*selected = i
				}
			})
		}
	})
}
func Spinner(c *ui.Context) ui.Element { return ui.Spinner(c).Size(16, 16) }
