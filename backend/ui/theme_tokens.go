// Author: HyperparameterSI
package ui

// ThemeTokens is an immutable set of namespaced colors for opt-in component
// libraries. An ordinary theme can leave it nil without allocating a map.
type ThemeTokens struct{ colors map[string]Color }

// NewThemeTokens takes a snapshot, so later changes to colors cannot affect
// windows already using this set. Replace the pointer to publish a new set.
func NewThemeTokens(colors map[string]Color) *ThemeTokens {
	t := &ThemeTokens{colors: make(map[string]Color, len(colors))}
	for name, color := range colors {
		t.colors[name] = color
	}
	return t
}

// Lookup distinguishes an explicitly transparent token from an absent one.
// It is safe on a nil receiver and safe for concurrent readers.
func (t *ThemeTokens) Lookup(name string) (Color, bool) {
	if t == nil {
		return Transparent, false
	}
	color, ok := t.colors[name]
	return color, ok
}
