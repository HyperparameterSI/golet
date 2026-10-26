// Author: HyperparameterSI
package forma

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// TestPinnedRegistry prevents silent reference changes or catalog omissions.
// Reference sources are read only by tests; none are embedded in applications.
func TestPinnedRegistry(t *testing.T) {
	data, err := os.ReadFile("../../../docs/references/shadcn-base-maia/manifest.json")
	if err != nil {
		t.Fatal(err)
	}
	var manifest struct {
		Style      string `json:"style"`
		Components []struct {
			Name  string                          `json:"name"`
			Files []struct{ File, SHA256 string } `json:"files"`
		} `json:"components"`
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatal(err)
	}
	if manifest.Style != "base-maia" {
		t.Fatalf("unexpected style %q", manifest.Style)
	}
	names := make(map[string]bool)
	for _, component := range manifest.Components {
		if names[component.Name] {
			t.Fatalf("duplicate component %q", component.Name)
		}
		names[component.Name] = true
		// Upstream's form entry is intentionally metadata-only; native Form
		// provides composition rather than claiming a source-level port.
		if len(component.Files) == 0 && component.Name != "form" {
			t.Fatalf("missing reference for %s", component.Name)
		}
		for _, file := range component.Files {
			data, err := os.ReadFile(filepath.Join("../../../docs/references/shadcn-base-maia", file.File))
			if err != nil {
				t.Fatal(err)
			}
			hash := sha256.Sum256(data)
			if hex.EncodeToString(hash[:]) != file.SHA256 {
				t.Errorf("reference changed: %s", file.File)
			}
		}
	}
	for _, name := range Components() {
		if !names[name] {
			t.Errorf("catalog entry absent or repeated in reference: %s", name)
		}
		delete(names, name)
	}
	for name := range names {
		t.Errorf("missing native catalog entry: %s", name)
	}
}
