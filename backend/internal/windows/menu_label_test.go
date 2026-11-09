// Author: HyperparameterSI

//go:build windows && (amd64 || arm64)

package windows

import (
	"testing"

	"github.com/HyperparameterSI/golet/backend/internal/platform"
)

// TestMenuLabel keeps the mnemonic of the label, and shows the subtitle's
// ampersands as they are.
func TestMenuLabel(t *testing.T) {
	for _, tt := range []struct{ label, subtitle, want string }{
		{"&Open", "", "&Open"},
		{"&Open", "Bed & Bath", "&Open — Bed && Bath"},
	} {
		if got := menuLabel(&platform.MenuItem{Label: tt.label, Subtitle: tt.subtitle}); got != tt.want {
			t.Errorf("menuLabel(%q, %q) = %q, want %q", tt.label, tt.subtitle, got, tt.want)
		}
	}
}
