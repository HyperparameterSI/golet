// Author: HyperparameterSI
package forma

import (
	"github.com/HyperparameterSI/golet/backend/ui"
	"sync"
)

type messageIcons struct{ good, bad, retry, edit *ui.SVG }

// Message-only assets are initialized on first use, never on the app path.
var messageActionIcons = sync.OnceValue(func() messageIcons {
	icon := func(body string) *ui.SVG {
		return ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.65" stroke-linecap="round" stroke-linejoin="round">` + body + `</svg>`))
	}
	return messageIcons{
		good:  icon(`<path d="M7 10v11H3V10h4Zm0 0 5-8c2 0 3 1 2 4l-1 4h6a2 2 0 0 1 2 2l-2 7a3 3 0 0 1-3 2H7"/>`),
		bad:   icon(`<path d="M7 14V3H3v11h4Zm0 0 5 8c2 0 3-1 2-4l-1-4h6a2 2 0 0 0 2-2l-2-7a3 3 0 0 0-3-2H7"/>`),
		retry: icon(`<path d="M3 4v6h6m12 10v-6h-6M4 9a8 8 0 0 1 14-4l3 5M3 14l3 5a8 8 0 0 0 14-4"/>`),
		edit:  icon(`<path d="m15 5 4 4M4 16l-1 5 5-1L21 7a2.8 2.8 0 0 0-4-4L4 16Z"/>`),
	}
})
