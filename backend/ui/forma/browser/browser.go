// Author: HyperparameterSI

// Package browser connects Forma themes to the optional standalone CEF host.
// Import it explicitly; ordinary Forma components do not depend on Chromium.
// This is not an embedded ui.Element.
package browser

import (
	"context"
	"fmt"

	"github.com/HyperparameterSI/golet/backend/plugins/chromium"
	"github.com/HyperparameterSI/golet/backend/ui"
)

// Theme converts a Forma (or core UI) theme to a browser shell snapshot.
// The shell limits typography to 10–20 DIP to keep its address bar usable.
func Theme(t *ui.Theme) chromium.Theme {
	if t == nil {
		t = ui.FormaTheme(false)
	}
	color := func(c ui.Color) string { return fmt.Sprintf("%02x%02x%02x%02x", c.R, c.G, c.B, c.A) }
	return chromium.Theme{
		Dark:       t.Dark,
		Background: color(t.Background), Surface: color(t.Surface), Muted: color(t.SurfaceHover),
		Foreground: color(t.Text), MutedForeground: color(t.TextMuted), Border: color(t.Border), Primary: color(t.Accent),
		Radius: min(64, max(0, t.Radius)), FontSize: min(20, max(10, t.FontSize)), Font: t.Font,
	}
}

// Launch opens a standalone browser using a snapshot of the selected theme.
// Custom shell themes currently require macOS and chromium.StyleAlloy.
func Launch(ctx context.Context, opts chromium.Options, theme *ui.Theme) (*chromium.Browser, error) {
	snapshot := Theme(theme)
	opts.Theme = &snapshot
	return chromium.Launch(ctx, opts)
}

// SetTheme applies a changed Forma theme without relaunching the browser.
// Call when the theme changes, outside the rendering loop.
func SetTheme(ctx context.Context, b *chromium.Browser, theme *ui.Theme) error {
	return b.SetTheme(ctx, Theme(theme))
}
