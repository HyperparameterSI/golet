// Author: HyperparameterSI

package main

import (
	"encoding/json"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestShadcnTemplate(t *testing.T) {
	dir := t.TempDir()
	data := templateData{Name: `Demo "App"`, Slug: "demo", Runtime: "^0.1.0", CLI: "^0.1.0", Golet: "golet", ConfigImport: "golet-cli"}
	if err := writeTemplate(dir, shadcnTemplate, data); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"package.json", "frontend/tsconfig.json", "frontend/components.json"} {
		b, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil || !json.Valid(b) {
			t.Fatalf("invalid %s: %v", name, err)
		}
	}
	for _, name := range []string{"frontend/src/components/app.tsx", "frontend/src/routes/__root.tsx", "frontend/src/routes/index.tsx", "frontend/src/router.tsx", "frontend/src/components/ui/button.tsx", "frontend/src/components/ui/input.tsx", "frontend/src/components/ui/card.tsx", "frontend/src/lib/utils.ts", "frontend/src/style.css"} {
		if _, err := os.Stat(filepath.Join(dir, name)); err != nil {
			t.Fatal(err)
		}
	}
	b, _ := os.ReadFile(filepath.Join(dir, "backend", "main.go"))
	if _, err := parser.ParseFile(token.NewFileSet(), "main.go", b, parser.AllErrors); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(b), "time.Tick") {
		t.Fatal("idle template should not poll")
	}
}
