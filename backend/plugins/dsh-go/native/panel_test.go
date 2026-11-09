// Author: HyperparameterSI
package native

import (
	"testing"

	"github.com/HyperparameterSI/golet/backend/ui"
	"github.com/HyperparameterSI/golet/backend/ui/forma"
)

// A protected invocation must remain an explicit single-call decision in both
// themes, and its parameters must be available before the user approves it.
func TestApprovalPanel(t *testing.T) {
	for _, dark := range []bool{false, true} {
		answer := make(chan bool, 1)
		p := &Panel{busy: true, approval: &approval{name: "current_time", arguments: "{}", answer: answer}}
		tt := ui.NewTester(func(c *ui.Context) {
			c.SetTheme(forma.Theme(dark))
			p.View(c)
		}, 800, 600)
		tt.Frame()
		for _, text := range []string{"允许执行工具：current_time？", "{}", "仅允许本次调用", "允许", "拒绝"} {
			if !tt.HasText(text) {
				t.Fatalf("dark=%v missing %q", dark, text)
			}
		}
		if err := tt.Click("允许"); err != nil {
			t.Fatal(err)
		}
		select {
		case allowed := <-answer:
			if !allowed {
				t.Fatal("approval denied")
			}
		default:
			t.Fatal("approval not delivered")
		}
		if p.approval != nil {
			t.Fatal("approval remained active")
		}
	}
}
