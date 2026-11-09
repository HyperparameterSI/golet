// Author: HyperparameterSI
package forma

import (
	"github.com/HyperparameterSI/golet/backend/ui"
	"sync"
)

type codeIcons struct{ terminal, wrap, download, copy, check, failed *ui.SVG }

// Parsed only when a CodeBlock is displayed, never on the common app path.
var codeBlockIcons = sync.OnceValue(func() codeIcons {
	icon := func(body string) *ui.SVG {
		return ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.65" stroke-linecap="round" stroke-linejoin="round">` + body + `</svg>`))
	}
	return codeIcons{
		terminal: icon(`<path d="m4 5 7 7-7 7m10 0h6"/>`),
		wrap:     icon(`<path d="M3 6h18M3 11h14a3 3 0 0 1 0 6h-5m2-2-2 2 2 2M3 17h4"/>`),
		download: icon(`<path d="M12 3v12m-4-4 4 4 4-4M4 15v5h16v-5"/>`),
		copy:     icon(`<rect x="8" y="8" width="12" height="13" rx="2"/><path d="M16 5V4a2 2 0 0 0-2-2H5a2 2 0 0 0-2 2v10a2 2 0 0 0 2 2"/>`),
		check:    icon(`<path d="m5 12 4 4L19 6"/>`),
		failed:   icon(`<path d="m6 6 12 12M18 6 6 18"/>`),
	}
})

func codeBlockAction(c *ui.Context, label string, icon *ui.SVG) ui.Element {
	b := ui.ButtonBase(c).Size(28, 28).Radius(radius(c, 6)).Label(label).TextColor(c.Theme().TextMuted).AlignItems(ui.Center).Justify(ui.Center)
	if b.Hovered() {
		b.Background(c.Theme().Text.Alpha(.06)).TextColor(c.Theme().Text)
	}
	b.Children(func() { ui.Icon(c, icon).Size(16, 16) })
	Tooltip(c, b, label)
	return b
}
