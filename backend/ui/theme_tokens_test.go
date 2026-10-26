// Author: HyperparameterSI
package ui

import "testing"

func TestThemeTokensSnapshot(t *testing.T) {
	var empty *ThemeTokens
	if _, ok := empty.Lookup("x"); ok {
		t.Fatal("nil tokens should be empty")
	}
	source := map[string]Color{"x": Transparent}
	tokens := NewThemeTokens(source)
	source["x"] = Hex("#ff0000")
	value, ok := tokens.Lookup("x")
	if !ok || value != Transparent {
		t.Fatal("token set did not snapshot caller data")
	}
	if _, ok := tokens.Lookup("missing"); ok {
		t.Fatal("unexpected token")
	}
}
