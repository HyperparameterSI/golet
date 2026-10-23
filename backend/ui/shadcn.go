// Author: HyperparameterSI

package ui

// ShadcnTheme returns the neutral shadcn palette with 14 DIP typography and
// soft, Maia-inspired corners. It shares LightTheme/DarkTheme's semantic colors.
// The returned value belongs to the caller: create it once, customize it and
// pass it to Context.SetTheme. No CSS, JavaScript or webview is required.
func ShadcnTheme(dark bool) *Theme {
	t := LightTheme()
	if dark {
		t = DarkTheme()
	}
	t.FontSize = 14
	return t
}
