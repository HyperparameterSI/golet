// Author: HyperparameterSI
package shadcn

import (
	"math"
	"strings"
	"time"

	"github.com/HyperparameterSI/golet/backend/ui"
)

// InputOptions controls validation and availability without changing the binding.
type InputOptions struct {
	Placeholder                 string
	Disabled, Invalid, ReadOnly bool
}

func inputStyle(c *ui.Context, e ui.Element, radius float32, o InputOptions) ui.Element {
	t := c.Theme()
	border := t.InputBorder
	if border == ui.Transparent {
		border = t.Border
	}
	e.Radius(radius).Border(1, border).Background(border.Alpha(.3)).FontSize(14).TextColor(t.Text).Disabled(o.Disabled).ReadOnly(o.ReadOnly).Placeholder(o.Placeholder)
	if o.Disabled {
		e.Opacity(.5)
	}
	focus(c, e, radius, o.Invalid)
	return e
}
func Input(c *ui.Context, value *string, options ...InputOptions) ui.Element {
	o := InputOptions{}
	if len(options) > 0 {
		o = options[0]
	}
	return inputStyle(c, ui.TextInputBase(c, value).Height(36).Padding(4, 12).MinWidth(0), 32, o)
}
func Textarea(c *ui.Context, value *string, options ...InputOptions) ui.Element {
	o := InputOptions{}
	if len(options) > 0 {
		o = options[0]
	}
	return inputStyle(c, ui.TextAreaBase(c, value).MinHeight(64).Padding(12).MinWidth(0), 12, o)
}

// InputGroup shares one surface and focus outline across its input and addons.
// Use InputGroupInput inside it to avoid nested input borders.
func InputGroup(c *ui.Context, build func()) ui.Element {
	e := ui.Row(c).AlignItems(ui.Center).Gap(6).MinHeight(36).Radius(32).Border(1, c.Theme().Border).Background(c.Theme().Border.Alpha(.3)).PaddingX(12).Role(ui.RoleGroup)
	e.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if e.FocusWithin() {
			p.Stroke(r, c.Theme().Focus, 32, 1)
			p.Stroke(ui.Rect{X: r.X - 1.5, Y: r.Y - 1.5, W: r.W + 3, H: r.H + 3}, c.Theme().Focus.Alpha(.5), 32, 3)
		}
	})
	return e.Children(build)
}

// InputGroupInput is the editable, borderless child of InputGroup.
func InputGroupInput(c *ui.Context, value *string, options ...InputOptions) ui.Element {
	o := InputOptions{}
	if len(options) > 0 {
		o = options[0]
	}
	e := ui.TextInputBase(c, value).Height(34).Grow(1).MinWidth(0).Padding(4, 0).FontSize(14).TextColor(c.Theme().Text).Placeholder(o.Placeholder).Disabled(o.Disabled).ReadOnly(o.ReadOnly).FocusRing(false)
	if o.Disabled {
		e.Opacity(.5)
	}
	if o.Invalid {
		e.Border(1, c.Theme().Danger).Radius(4)
	}
	return e
}
func Field(c *ui.Context, label string, build func()) ui.Element {
	return ui.Field(c, label, build).Gap(8).FontSize(14)
}
func Form(c *ui.Context, build func()) ui.Element {
	return ui.Column(c).Gap(24).Role(ui.RoleGroup).Children(build)
}

func Checkbox(c *ui.Context, on *bool, label string) ui.Element {
	t := c.Theme()
	e := ui.CheckboxBase(c, on).Gap(8).AlignItems(ui.Center).Label(label)
	e.Children(func() {
		box := ui.Box(c).Size(16, 16).Shrink(0)
		path := ui.Local(box, "check-path", func() ui.Path { return ui.Path{} })
		box.Draw(func(p *ui.Painter, r ui.Rect) {
			bg, border := t.Background, t.Border
			if *on {
				bg, border = t.Accent, t.Accent
			}
			p.Fill(r, bg, 6)
			p.Stroke(r, border, 6, 1)
			if *on {
				path.Reset().MoveTo(r.X+3, r.Y+8).LineTo(r.X+6.5, r.Y+11).LineTo(r.X+12.5, r.Y+5)
				p.StrokePath(path, 1.5, t.AccentText)
			}
			if e.FocusVisible() {
				p.Stroke(ui.Rect{X: r.X - 3, Y: r.Y - 3, W: r.W + 6, H: r.H + 6}, t.Focus.Alpha(.5), 9, 3)
			}
		})
		if label != "" {
			Label(c, label)
		}
	})
	e.FocusRing(false)
	return e
}
func Switch(c *ui.Context, on *bool) ui.Element {
	t := c.Theme()
	e := ui.SwitchBase(c, on).Size(32, 18.4).Radius(32)
	e.Draw(func(p *ui.Painter, r ui.Rect) {
		bg := t.Border
		if *on {
			bg = t.Accent
		}
		p.Fill(r, bg, 32)
		x := r.X + 1.2
		if *on {
			x += 13.6
		}
		thumb := t.Background
		if t.Dark {
			thumb = t.Text
			if *on {
				thumb = t.AccentText
			}
		}
		p.Fill(ui.Rect{X: x, Y: r.Y + 1.2, W: 16, H: 16}, thumb, 16)
	})
	focus(c, e, 32, false)
	return e
}
func RadioGroup(c *ui.Context, build func()) ui.Element { return ui.RadioGroup(c, build).Gap(12) }
func Radio[T comparable](c *ui.Context, selected *T, value T, label string) ui.Element {
	t := c.Theme()
	e := ui.RadioBase(c, selected, value).Gap(8).AlignItems(ui.Center).Label(label)
	e.Children(func() {
		ui.Box(c).Size(16, 16).Draw(func(p *ui.Painter, r ui.Rect) {
			p.Stroke(r, t.Border, 16, 1)
			if *selected == value {
				p.Fill(ui.Rect{X: r.X + 4, Y: r.Y + 4, W: 8, H: 8}, t.Accent, 8)
			}
		})
		Label(c, label)
	})
	focus(c, e, 8, false)
	return e
}
func Toggle(c *ui.Context, on *bool, label string) ui.Element {
	e := ui.ToggleBase(c, on).Height(36).MinWidth(36).PaddingX(12).Radius(32).Gap(4).Label(label)
	e.Draw(func(p *ui.Painter, r ui.Rect) {
		if *on || e.Hovered() {
			p.Fill(r, c.Theme().SurfaceHover, 32)
		}
	})
	focus(c, e, 32, false)
	e.Children(func() { Label(c, label) })
	return e
}
func ToggleGroup(c *ui.Context, build func()) ui.Element { return ui.ToggleGroup(c, build).Gap(4) }
func Slider(c *ui.Context, value *float64, lo, hi float64) ui.Element {
	t := c.Theme()
	e := ui.SliderBase(c, value, lo, hi).Height(20).PaddingX(8)
	e.Draw(func(p *ui.Painter, r ui.Rect) {
		fraction := float32(0)
		if hi > lo {
			fraction = float32(math.Max(0, math.Min(1, (*value-lo)/(hi-lo))))
		}
		track := ui.Rect{X: r.X + 8, Y: r.Y + (r.H-12)/2, W: max(0, r.W-16), H: 12}
		p.Fill(track, t.SurfaceHover, 32)
		p.Fill(ui.Rect{X: track.X, Y: track.Y, W: track.W * fraction, H: 12}, t.Accent, 32)
		thumb := ui.Rect{X: track.X + track.W*fraction - 8, Y: r.Y + (r.H-16)/2, W: 16, H: 16}
		p.Fill(thumb, ui.RGB(255, 255, 255), 32)
		p.Stroke(thumb, t.Accent, 32, 1)
		if e.FocusVisible() {
			p.Stroke(ui.Rect{X: thumb.X - 4, Y: thumb.Y - 4, W: 24, H: 24}, t.Focus.Alpha(.5), 32, 4)
		}
	})
	return e.FocusRing(false)
}
func Progress(c *ui.Context, value float64) ui.Element {
	return ui.Progress(c, value).Height(8).Radius(32)
}

// DatePicker composes Golet's native date editing and calendar popover.
func DatePicker(c *ui.Context, date *time.Time) ui.Element {
	return ui.DateInput(c, date).Height(36).Radius(32).PaddingX(12)
}

// InputOTP accepts ASCII digits, supports paste and preserves a single native
// editing/IME surface for keyboard navigation and accessibility.
func InputOTP(c *ui.Context, value *string, length int) ui.Element {
	length = max(1, min(length, 32))
	clean := func(s string) string {
		var b strings.Builder
		for _, r := range s {
			if r >= '0' && r <= '9' && b.Len() < length {
				b.WriteRune(r)
			}
		}
		return b.String()
	}
	glyphWidth, _ := c.MeasureText(100, ui.Span{Text: "0", Font: "monospace", Size: 14})
	e := Input(c, value).Width(float32(length * 36)).Font("monospace").LetterSpacing(36 - glyphWidth).PaddingX((36 - glyphWidth) / 2).Label("One-time password").OnChange(func() { *value = clean(*value) })
	e.DrawOver(func(p *ui.Painter, r ui.Rect) {
		for i := 1; i < length; i++ {
			x := r.X + float32(i*36)
			p.Line(x, r.Y, x, r.Y+r.H, 1, c.Theme().Border)
		}
		if e.FocusVisible() {
			_, caret := e.TextAnchorCaret()
			caret = min(length-1, caret)
			p.Stroke(ui.Rect{X: r.X + float32(caret*36), Y: r.Y, W: 36, H: r.H}, c.Theme().Focus, caretRadius(caret, length), 1)
		}
	})
	return e
}

func caretRadius(index, length int) float32 {
	if index == 0 || index == length-1 {
		return 12
	}
	return 0
}
