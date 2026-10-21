// Author: HyperparameterSI

package ui

// Tabs creates a row of tabs showing labels, of which *selected is the
// index of the one chosen. A click chooses a tab, as do the arrows, Home
// and End while one has the keyboard focus, which follows the choice:
//
//	ui.Tabs(c, &app.tab, "General", "Appearance", "Advanced")
//	switch app.tab {
//	case 0:
//		app.general(c)
//	…
//	}
//
// Changed reports a new choice.
func coreTabs(c *context, selected *int, labels ...string) *node {
	t := c.theme
	tabs := coreTabsBase(c, selected, len(labels))
	list := tabs.List.Gap(t.Space(1)).Padding(t.Space(1)).Radius(t.Radius + 2).Background(t.SurfaceHover)
	list.Children(func() {
		for i, label := range labels {
			tab := tabs.Tab(i).Padding(t.Space(1.5), t.Space(3)).Radius(t.Radius).FocusRing(false)
			on := i == *selected
			tab.TextColor(t.TextMuted)
			if on {
				tab.TextColor(t.Text).FontWeight(500).Background(t.Surface).Border(1, t.Border)
			}
			tab.styleFn = func(tab *node) {
				if !on && tab.Hovered() {
					tab.ts.color = t.Text
				}
			}
			tab.DrawOver(func(p *Painter, r Rect) {
				if tab.FocusVisible() {
					p.FocusRing(r, [4]float32{t.Radius, t.Radius, t.Radius, t.Radius})
				}
			})
			tab.Children(func() { coreText(c, label).SingleLine() })
		}
	})
	return list
}
