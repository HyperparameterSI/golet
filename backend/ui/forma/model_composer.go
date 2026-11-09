// Author: HyperparameterSI
package forma

import (
	"github.com/HyperparameterSI/golet/backend/ui"
	"strings"
	"sync"
)

func modelPill(c *ui.Context, e ui.Element) ui.Element {
	return e.Height(32).Padding(5, 12).Gap(8).Radius(16).Border(1, c.Theme().Border).Background(c.Theme().Background).TextColor(c.Theme().TextMuted).FontSize(fontSize(c, 14)).AlignItems(ui.Center)
}

type modelIcons struct{ chevron, add, send *ui.SVG }

var modelControlIcons = sync.OnceValue(func() modelIcons {
	icon := func(body string) *ui.SVG {
		return ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">` + body + `</svg>`))
	}
	return modelIcons{icon(`<path d="m6 9 6 6 6-6"/>`), icon(`<path d="M12 5v14M5 12h14"/>`), icon(`<path d="M12 19V5m-6 6 6-6 6 6"/>`)}
})

// ModelComposerOptions supplies effects explicitly. The component never sends
// requests or opens file dialogs itself. OnSubmit receives an immutable snapshot
// of the prompt, model and effort; the application owns clearing the text.
type ModelComposerOptions struct {
	Placeholder   string
	Disabled      bool
	OnAdd         func()
	OnSubmit      func(text, model, effort string)
	OnModelChange func(model, effort string)
}

// ModelComposer combines a multiline prompt and ModelSelector in one muted
// panel. Enter inserts a newline; the send button submits nonblank text only.
func ModelComposer(c *ui.Context, text, model, effort *string, models []ModelOption, options ...ModelComposerOptions) ui.Element {
	o := ModelComposerOptions{}
	if len(options) > 0 {
		o = options[0]
	}
	if o.Placeholder == "" {
		o.Placeholder = "Ask the agent to change your code…"
	}
	root := ui.Column(c).FillWidth().MinWidth(0).Padding(12).Gap(12).Radius(radius(c, 20)).Background(token(c, "muted", c.Theme().SurfaceHover))
	return root.Children(func() {
		ui.TextAreaBase(c, text).FillWidth().MinWidth(0).Height(48).Padding(2).FontSize(fontSize(c, 14)).TextColor(c.Theme().Text).Placeholder(o.Placeholder).Label("Message input").Disabled(o.Disabled).FocusRing(false)
		ui.Row(c).FillWidth().Gap(8).AlignItems(ui.Center).Children(func() {
			add := modelCircle(c, "Add attachment", modelControlIcons().add, false).Disabled(o.Disabled || o.OnAdd == nil)
			if add.Clicked() && o.OnAdd != nil {
				o.OnAdd()
			}
			ui.Spacer(c)
			ModelSelector(c, model, effort, models, ModelSelectorOptions{Disabled: o.Disabled, OnChange: o.OnModelChange})
			valid := false
			for _, m := range models {
				if m.ID == *model && !m.Disabled {
					valid = true
					break
				}
			}
			send := modelCircle(c, "Send message", modelControlIcons().send, true).Disabled(o.Disabled || o.OnSubmit == nil || strings.TrimSpace(*text) == "" || !valid)
			if send.Clicked() && o.OnSubmit != nil {
				o.OnSubmit(*text, *model, *effort)
			}
		})
	})
}
func modelCircle(c *ui.Context, label string, icon *ui.SVG, primary bool) ui.Element {
	b := ui.ButtonBase(c).Size(32, 32).Radius(16).Label(label).AlignItems(ui.Center).Justify(ui.Center).Background(c.Theme().Background).Border(1, c.Theme().Border).TextColor(c.Theme().TextMuted)
	if primary {
		b.Background(c.Theme().TextMuted).TextColor(c.Theme().Background).Border(0, c.Theme().Border)
	} else if b.Hovered() {
		b.Background(c.Theme().SurfaceHover)
	}
	b.Children(func() { ui.Icon(c, icon).Size(18, 18) })
	Tooltip(c, b, label)
	return b
}
