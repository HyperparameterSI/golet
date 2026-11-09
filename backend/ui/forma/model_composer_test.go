// Author: HyperparameterSI
package forma

import (
	"github.com/HyperparameterSI/golet/backend/ui"
	"testing"
)

func TestModelComposerActions(t *testing.T) {
	text, model, effort := "", "demo", "medium"
	sent := ""
	adds := 0
	disabled := false
	models := []ModelOption{{ID: "demo", Name: "Demo", Efforts: []ModelEffort{{ID: "medium", Name: "Medium"}}}}
	tt := ui.NewTester(func(c *ui.Context) {
		ModelComposer(c, &text, &model, &effort, models, ModelComposerOptions{Disabled: disabled, OnAdd: func() { adds++ }, OnSubmit: func(s, m, e string) { sent = s + "/" + m + "/" + e }})
	}, 600, 240)
	tt.Click("Send message")
	if sent != "" {
		t.Fatal("empty submission")
	}
	text = "hello"
	tt.Frame()
	tt.Click("Send message")
	if sent != "hello/demo/medium" || text != "hello" {
		t.Fatalf("snapshot %q, input %q", sent, text)
	}
	tt.Click("Add attachment")
	if adds != 1 {
		t.Fatal(adds)
	}
	disabled = true
	sent = ""
	tt.Frame()
	tt.Click("Send message")
	tt.Click("Add attachment")
	if sent != "" || adds != 1 {
		t.Fatal("disabled callbacks ran")
	}
}
