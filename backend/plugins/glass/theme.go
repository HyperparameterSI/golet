// Author: HyperparameterSI

package glass

import "github.com/HyperparameterSI/golet/backend/ui"

// Theme returns an opt-in Liquid Glass theme for native controls. Cache one
// for each appearance and set it on the UI context before building controls.
// Only apps importing glass carry its shaders. Content keeps a solid surface;
// buttons, editable fields, navigation tracks and popovers use vivid glass.
// High-contrast preferences automatically keep the controls' solid fills.
func Theme(dark bool) *ui.Theme {
	t := ui.ShadcnTheme(dark)
	// Stronger secondary text stays legible over tinted backdrops.
	t.TextMuted = ui.Hex("#48484a")
	t.Accent = ui.Hex("#007aff")
	t.AccentHover = ui.Hex("#006ee6")
	t.AccentPressed = ui.Hex("#005fcc")
	t.AccentText = ui.Hex("#ffffff")
	if dark {
		t.TextMuted = ui.Hex("#d1d1d6")
		t.Accent = ui.Hex("#0a84ff")
		t.AccentHover = ui.Hex("#409cff")
		t.AccentPressed = ui.Hex("#0070e0")
	}
	t.Focus = t.Accent
	t.ControlMaterial = Glass{Style: Vivid, Interactive: true}
	t.AccentMaterial = Glass{Style: Vivid, Tint: t.Accent, Interactive: true}
	t.PanelMaterial = Glass{Style: Vivid}
	return t
}
