// Author: HyperparameterSI

package ui

// ShadcnTheme returns a neutral theme matching Golet's shadcn web template.
// The returned value belongs to the caller and can be customized. Create it
// once and pass it to Context.SetTheme in the view. This is a native Go theme;
// it does not load CSS, JavaScript or a webview.
func ShadcnTheme(dark bool) *Theme {
	t := LightTheme()
	if dark {
		t = DarkTheme()
		t.Background = Hex("#0a0a0a")
		t.Surface = Hex("#171717")
		t.SurfaceHover = Hex("#262626")
		t.SurfacePressed = Hex("#404040")
		t.Border = Hex("#404040")
		t.Text, t.TextMuted = Hex("#fafafa"), Hex("#a3a3a3")
		t.Accent, t.AccentHover, t.AccentPressed = Hex("#fafafa"), Hex("#e5e5e5"), Hex("#d4d4d4")
		t.AccentText = Hex("#171717")
		t.Selection = RGBA(163, 163, 163, 0.3)
		t.Focus = Hex("#a3a3a3")
	} else {
		t.Background = Hex("#fafafa")
		t.Surface = Hex("#ffffff")
		t.SurfaceHover = Hex("#f5f5f5")
		t.SurfacePressed = Hex("#e5e5e5")
		t.Border = Hex("#e5e5e5")
		t.Text, t.TextMuted = Hex("#171717"), Hex("#737373")
		t.Accent, t.AccentHover, t.AccentPressed = Hex("#171717"), Hex("#262626"), Hex("#404040")
		t.AccentText = Hex("#fafafa")
		t.Selection = RGBA(115, 115, 115, 0.25)
		t.Focus = Hex("#737373")
	}
	t.Radius, t.Spacing, t.FontSize = 8, 4, 14
	return t
}
