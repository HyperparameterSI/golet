// Author: HyperparameterSI
package forma

import (
	"fmt"
	"math"
	"time"

	"github.com/HyperparameterSI/golet/backend/ui"
	"github.com/HyperparameterSI/golet/backend/ui/chart"
)

func Skeleton(c *ui.Context) ui.Element {
	t := c.Theme()
	reduced := c.Preferences().ReduceMotion
	e := ui.Box(c).Radius(radius(c, 12))
	e.Draw(func(p *ui.Painter, r ui.Rect) {
		alpha := float32(1)
		if !reduced {
			alpha = .75 + .25*float32(math.Cos(float64(p.Now().UnixNano())/float64(time.Second)*math.Pi))
			p.After(time.Second / 30)
		}
		p.Fill(r, t.SurfaceHover.Alpha(alpha), radius(c, 12))
	})
	return e
}

// Carousel presents a page with previous/next controls, arrow keys and drag.
func Carousel(c *ui.Context, index *int, count int, slide func(int)) ui.Element {
	count = max(0, count)
	*index = max(0, min(*index, max(0, count-1)))
	e := ui.Row(c).Gap(12).AlignItems(ui.Center)
	e.Children(func() {
		Button(c, "‹", ButtonOptions{Variant: Outline, Size: IconSmall, Disabled: *index == 0}).Label("Previous slide").OnClick(func() { *index-- })
		body := ui.Box(c).Grow(1).MinWidth(0).Clip().Focusable().Role(ui.RoleGroup).Label(fmt.Sprintf("Slide %d of %d", *index+1, count))
		offset := ui.Local(body, "carousel-drag", func() ui.Point { return ui.Point{} })
		body.Draggable(offset, ui.DraggableOptions{Axis: ui.Horizontal, Bounds: &ui.Rect{X: -120, W: 240}, Snap: func(point ui.Point) ui.Point {
			if point.X > 40 && *index > 0 {
				*index--
			}
			if point.X < -40 && *index < count-1 {
				*index++
			}
			return ui.Point{}
		}})
		body.Children(func() {
			if count > 0 {
				slide(*index)
			}
		})
		body.HandleInput(func(ev ui.InputEvent) bool {
			if ev.Kind == ui.InputKeyDown {
				if ev.Key == ui.KeyRight && *index < count-1 {
					*index++
					return true
				}
				if ev.Key == ui.KeyLeft && *index > 0 {
					*index--
					return true
				}
			}
			return false
		})
		Button(c, "›", ButtonOptions{Variant: Outline, Size: IconSmall, Disabled: count == 0 || *index >= count-1}).Label("Next slide").OnClick(func() { *index++ })
	})
	return e
}

// Table draws Maia headers and cells over Golet's virtualized, selectable list.
// Column resizing and cell-block selection remain available in ui.Table.
func Table(c *ui.Context, state *ui.ListState, columns []ui.TableColumn, count int, cell func(int, int)) ui.Element {
	root := ui.Column(c).Role(ui.RoleTable).FontSize(fontSize(c, 14))
	sizeCell := func(e ui.Element, col ui.TableColumn) {
		if col.Width > 0 {
			e.Width(col.Width).Shrink(0)
		} else {
			e.Grow(1).MinWidth(max(60, col.MinWidth))
		}
	}
	root.Children(func() {
		ui.Row(c).Height(48).BorderWidth(0, 0, 1, 0).BorderColor(c.Theme().Border).Role(ui.RoleRow).Children(func() {
			for _, col := range columns {
				e := ui.Column(c).Justify(ui.Center).PaddingX(12).Role(ui.RoleColumnHeader)
				sizeCell(e, col)
				e.Children(func() {
					if col.Sortable && state != nil && state.Sort != nil {
						Button(c, col.Title, ButtonOptions{Variant: Ghost, Size: Small}).Padding(0).OnClick(func() {
							id := col.ID
							if id == "" {
								id = col.Title
							}
							descending := state.Sort.Column == id && !state.Sort.Descending
							*state.Sort = ui.SortOrder{Column: id, Descending: descending}
						})
					} else {
						Label(c, col.Title)
					}
				})
			}
		})
		ui.List(c, state, count).Grow(1).Rows(func(row ui.ListRow) {
			r := ui.Row(c).Role(ui.RoleRow).MinHeight(44)
			if row.Index < count-1 {
				r.BorderWidth(0, 0, 1, 0).BorderColor(c.Theme().Border)
			}
			if row.Selected() {
				r.Background(c.Theme().SurfaceHover)
			} else if r.Hovered() {
				r.Background(c.Theme().SurfaceHover.Alpha(.5))
			}
			r.Children(func() {
				for i, col := range columns {
					e := ui.Column(c).Padding(12).Role(ui.RoleCell)
					sizeCell(e, col)
					e.Children(func() { cell(row.Index, i) })
				}
			})
		})
	})
	return root
}

// ChartContainer is Maia's themed composition surface for ui/chart plots.
// Importing other shadcn controls does not link chart geometry into an app.
func ChartContainer(c *ui.Context, axes chart.Axes, build func(ui.Element)) ui.Element {
	return chart.Plot(c, axes, build).FontSize(fontSize(c, 12)).TextColor(c.Theme().TextMuted).AspectRatio(16.0 / 9)
}
func ChartLegend(c *ui.Context, labels []string, colors []ui.Color) ui.Element {
	defaults := ChartColors(c)
	return ui.Row(c).Gap(16).Center().Children(func() {
		for i, label := range labels {
			color := defaults[i%5]
			if i < len(colors) {
				color = colors[i]
			}
			ui.Row(c).Gap(6).AlignItems(ui.Center).Children(func() {
				ui.Box(c).Size(8, 8).Radius(radius(c, 2)).Background(color)
				ui.Text(c, label).FontSize(fontSize(c, 12))
			})
		}
	})
}
func ChartTooltip(c *ui.Context, anchor ui.Element, open *bool, title string, build func()) ui.Element {
	return Popover(c, anchor, open, func() { Label(c, title); build() }).Width(160).Padding(10).Gap(6)
}

// Toast adds a notification to Golet's managed queue. Build Sonner once in the view.
func Toast(c *ui.Context, t ui.Toast) string { return c.AddToast(t) }
func Sonner(c *ui.Context) {
	ui.ToastViewportBase(c, func(viewport ui.Element, toasts []ui.Toast) {
		viewport.Gap(8)
		for _, toast := range toasts {
			parts := ui.ToastBase(c, toast)
			parts.Root.Width(356).Padding(16).Gap(8).Radius(radius(c, 16)).Background(c.Theme().Surface).Border(1, c.Theme().Border).Shadow(0, 4, 12, 0, ui.RGBA(0, 0, 0, .1)).Children(func() {
				Label(c, toast.Title)
				if toast.Description != "" {
					CardDescription(c, toast.Description)
				}
				ui.Row(c).Gap(8).Children(func() {
					if toast.Action != "" {
						parts.ActionButton().Padding(6, 12).Radius(radius(c, 32)).Background(c.Theme().Accent).TextColor(c.Theme().AccentText).Children(func() { ui.Text(c, toast.Action) })
					}
					parts.CloseButton().Padding(6).Children(func() { ui.Text(c, "×") })
				})
			})
		}
	})
}
