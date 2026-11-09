// Author: HyperparameterSI
package forma

import (
	"github.com/HyperparameterSI/golet/backend/ui"
	"testing"
)

func TestConfirmationDecisionIsSingleUse(t *testing.T) {
	for _, allow := range []bool{true, false} {
		state := ConfirmationPending
		calls := 0
		tt := ui.NewTester(func(c *ui.Context) {
			Confirmation(c, "Read project?", "main.go", &state, func(v bool) {
				calls++
				if v != allow {
					t.Error("wrong decision")
				}
			})
		}, 600, 300)
		button := "拒绝"
		want := ConfirmationDenied
		if allow {
			button = "允许"
			want = ConfirmationApproved
		}
		if err := tt.Click(button); err != nil {
			t.Fatal(err)
		}
		tt.Frame()
		tt.Frame()
		if state != want || calls != 1 {
			t.Fatalf("state %s, callbacks %d", state, calls)
		}
		if _, found := tt.Find(button); found {
			t.Fatal("decision button remains after decision")
		}
	}
}
func TestContextUsageBounds(t *testing.T) {
	for _, tc := range []struct {
		used, limit int64
		want        float64
	}{{5, 0, 0}, {-1, 10, 0}, {20, 10, 100}, {5, 20, 25}} {
		if got := ContextUsagePercent(ContextUsage{Used: tc.used, Limit: tc.limit}); got != tc.want {
			t.Fatalf("got %v want %v", got, tc.want)
		}
	}
}
func TestSourceDeduplicationAndURLValidation(t *testing.T) {
	sources := []Source{{URL: "https://go.dev"}, {URL: "https://go.dev"}, {URL: "javascript:alert(1)"}, {}}
	if got := len(uniqueSources(sources)); got != 2 {
		t.Fatalf("sources %d", got)
	}
	if sourceAllowed(sources[2]) || sourceAllowed(sources[3]) || !sourceAllowed(sources[0]) {
		t.Fatal("URL validation")
	}
	if sourceName(sources[0]) != "go.dev" {
		t.Fatal("hostname fallback")
	}
}
func TestFileTreeSelectionCallback(t *testing.T) {
	state := FileTreeState{Open: map[string]bool{"root": true}}
	selected := ""
	tt := ui.NewTester(func(c *ui.Context) {
		FileTree(c, &state, []FileNode{{ID: "root", Name: "backend", Directory: true, Children: []FileNode{{ID: "main", Name: "main.go", Path: "backend/main.go"}}}}, func(n FileNode) { selected = n.Path })
	}, 600, 400)
	if err := tt.Click("File main"); err != nil {
		t.Fatal(err)
	}
	if state.Selected != "main" || selected != "backend/main.go" {
		t.Fatal("selection not delivered")
	}
}
func TestFailedTestSuiteStartsExpanded(t *testing.T) {
	state := TestResultsState{Open: true}
	tt := ui.NewTester(func(c *ui.Context) {
		TestResults(c, &state, []AITestSuite{{ID: "one", Name: "suite", Tests: []AITestCase{{Name: "bad", Status: AIComplete, Error: "missing configuration"}}}})
	}, 600, 400)
	tt.Frame()
	if !state.Suites["one"] || !tt.HasText("missing configuration") || !tt.HasText("Tests · 0 passed · 1 failed") {
		t.Fatal("failure not exposed")
	}
	if err := tt.Click("suite"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if state.Suites["one"] {
		t.Fatal("collapsed state not preserved")
	}
}

func TestStepsStatusPrecedence(t *testing.T) {
	for _, tc := range []struct {
		steps []AIStep
		want  AIStatus
	}{
		{nil, AIPending},
		{[]AIStep{{Status: AIComplete}, {Status: AICancelled}}, AICancelled},
		{[]AIStep{{Status: AIPending}, {Status: AIRunning}}, AIRunning},
		{[]AIStep{{Status: AIRunning}, {Status: AIError}}, AIError},
	} {
		if got := stepsStatus(tc.steps); got != tc.want {
			t.Fatalf("got %s want %s", got, tc.want)
		}
	}
}
