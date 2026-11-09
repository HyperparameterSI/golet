// Author: HyperparameterSI
package forma

import (
	"github.com/HyperparameterSI/golet/backend/ui"
	"testing"
)

func TestPromptKeyboardAndSnapshot(t *testing.T) {
	s := PromptInputState{Text: "hello", Attachments: []PromptAttachment{{Name: "a.txt"}}}
	var sent []PromptSubmission
	tt := ui.NewTester(func(c *ui.Context) {
		PromptInput(c, &s, PromptInputOptions{OnSubmit: func(p PromptSubmission) { sent = append(sent, p) }})
	}, 560, 300)
	tt.Click("Prompt text")
	tt.Key(ui.Shift, ui.KeyEnter)
	if len(sent) != 0 {
		t.Fatal("Shift+Enter submitted")
	}
	tt.Key(0, ui.KeyEnter)
	if len(sent) != 1 || sent[0].Text != s.Text {
		t.Fatalf("snapshot %+v input %q", sent, s.Text)
	}
	s.Attachments[0].Name = "changed"
	if sent[0].Attachments[0].Name != "a.txt" {
		t.Fatal("snapshot aliased draft")
	}
	tt.Frame()
	tt.Click("Remove attachment changed")
	if len(s.Attachments) != 0 {
		t.Fatal("attachment not removed")
	}
}
func TestPromptBusyAndEmpty(t *testing.T) {
	s := PromptInputState{}
	calls, stops := 0, 0
	status := PromptReady
	disabled := false
	tt := ui.NewTester(func(c *ui.Context) {
		PromptInput(c, &s, PromptInputOptions{Disabled: disabled, Status: status, OnSubmit: func(PromptSubmission) { calls++ }, OnStop: func() { stops++ }})
	}, 560, 300)
	tt.Click("Send prompt")
	if calls != 0 {
		t.Fatal("empty submit")
	}
	s.Attachments = []PromptAttachment{{Name: "file.txt"}}
	tt.Frame()
	tt.Click("Send prompt")
	if calls != 1 {
		t.Fatal("attachment-only submission failed")
	}
	status = PromptStreaming
	tt.Frame()
	tt.Click("Stop response")
	if stops != 1 {
		t.Fatal("missing stop")
	}
	tt.Click("Remove attachment file.txt")
	if len(s.Attachments) != 1 {
		t.Fatal("busy draft mutated")
	}
	disabled = true
	tt.Frame()
	tt.Click("Stop response")
	if stops != 1 {
		t.Fatal("disabled stop invoked")
	}
}
func TestPromptDrops(t *testing.T) {
	s := PromptInputState{}
	var files []string
	tt := ui.NewTester(func(c *ui.Context) {
		PromptInput(c, &s, PromptInputOptions{OnDropFiles: func(p []string) { files = p }})
	}, 560, 300)
	r, _ := tt.Find("Prompt input")
	tt.DropFiles(r.X+4, r.Y+4, "/tmp/example.txt")
	if len(files) != 1 || files[0] != "/tmp/example.txt" {
		t.Fatal(files)
	}
}
