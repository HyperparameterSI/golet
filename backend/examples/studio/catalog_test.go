// Author: HyperparameterSI
package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCatalogSourcesAndIDs(t *testing.T) {
	seen := map[string]bool{}
	for _, demo := range demos {
		if seen[demo.ID] || demo.ID == "" || len(demo.Controls) == 0 {
			t.Fatalf("invalid demo: %+v", demo)
		}
		seen[demo.ID] = true
		if _, err := os.Stat(filepath.Join("../../..", demo.Source)); err != nil {
			t.Errorf("%s: %v", demo.ID, err)
		}
	}
	if !seen["Browser"] {
		t.Fatal("Browser demo missing")
	}
	if !seen["Scene3D"] {
		t.Fatal("Scene3D demo missing")
	}
}

// Keep visual widget coverage current when the generated public API grows.
// Base variants share their styled counterpart's demo; nonvisual helpers are excluded.
func TestNativeWidgetCoverage(t *testing.T) {
	used := map[string]bool{}
	for _, dir := range []string{".", "../internal/gallerydemo", "../internal/editordemo", "../internal/loadingdemo"} {
		files, err := filepath.Glob(dir + "/*.go")
		if err != nil {
			t.Fatal(err)
		}
		for _, name := range files {
			if strings.HasSuffix(name, "_test.go") {
				continue
			}
			file, err := parser.ParseFile(token.NewFileSet(), name, nil, 0)
			if err != nil {
				t.Fatal(err)
			}
			ast.Inspect(file, func(n ast.Node) bool {
				if selector, ok := n.(*ast.SelectorExpr); ok {
					if pkg, ok := selector.X.(*ast.Ident); ok && pkg.Name == "ui" {
						used[selector.Sel.Name] = true
					}
				}
				return true
			})
		}
	}
	api, err := parser.ParseFile(token.NewFileSet(), "../../ui/api_gen.go", nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	helpers := map[string]bool{"DataDragOver": true, "DragOver": true, "Drop": true, "DropData": true, "Local": true, "NewTester": true, "Render": true, "Overlay": true}
	for _, decl := range api.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Recv != nil || !ast.IsExported(fn.Name.Name) {
			continue
		}
		name := fn.Name.Name
		if strings.HasSuffix(name, "Base") || helpers[name] {
			continue
		}
		if !used[name] {
			t.Errorf("visual widget %s has no Studio demo", name)
		}
	}
}
